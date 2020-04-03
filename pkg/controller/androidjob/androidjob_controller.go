package androidjob

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"sync"
	"time"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/android"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_androidjob")

// Add creates a new AndroidJob Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAndroidJob{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("androidjob-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AndroidJob
	err = c.Watch(&source.Kind{Type: &androidv1alpha1.AndroidJob{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// start the TTL controller
	go runTTLController(mgr.GetClient())

	return nil
}

// blank assignment to verify that ReconcileAndroidJob implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAndroidJob{}

// ReconcileAndroidJob reconciles a AndroidJob object
type ReconcileAndroidJob struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// status is used to feed back status updates when running jobs in parallel
type status struct {
	name   string
	status androidv1alpha1.DeviceJobStatus
}

// Reconcile reads that state of the cluster for a AndroidJob object and makes changes based on the state read
// and what is in the AndroidJob.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAndroidJob) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AndroidJob")

	// Fetch the AndroidJob instance
	instance := &androidv1alpha1.AndroidJob{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// fetch the job template
	jobTemplate := &androidv1alpha1.AndroidJobTemplate{}
	if err := r.client.Get(context.TODO(), instance.TemplateNamespacedName(), jobTemplate); err != nil {
		return reconcile.Result{}, err
	}

	// fetch the target devices
	reqLogger.Info("Fetching target devices for job")
	targetDevices, err := getTargetDevices(r.client, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// instantiate the status map if it is nil
	if instance.Status.JobStatus == nil {
		reqLogger.Info("Initializing status entries for job")
		instance.Status.JobStatus = make(map[string]androidv1alpha1.DeviceJobStatus)
	}

	// setup channels for running jobs concurrently
	errChan := make(chan error)
	statusChan := make(chan status)
	var wg sync.WaitGroup

	// run the jobs
	for _, device := range targetDevices {
		wg.Add(1)
		reqLogger.Info("Starting job worker for device", "DeviceName", device.Name)
		go runJobWorker(reqLogger, instance, device, jobTemplate, statusChan, errChan, &wg)
	}

	// wait and close the channels
	go func() {
		wg.Wait()
		reqLogger.Info("Jobs are complete, closing channels")
		close(statusChan)
		close(errChan)
	}()

	// retrieve results from the channels
	instance, errOcurred := watchJobChannels(reqLogger, instance, statusChan, errChan)

	// push status updates
	reqLogger.Info("Publishing status updates for job")
	if err := r.client.Status().Update(context.TODO(), instance); err != nil {
		return reconcile.Result{}, err
	}

	// determine if any errors happened (we need to requeue)
	if errOcurred {
		reqLogger.Info("One or more errors ocurred while processing the job, requeing until completion")
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Duration(3) * time.Second,
		}, nil
	}

	return reconcile.Result{}, nil
}

func watchJobChannels(reqLogger logr.Logger, instance *androidv1alpha1.AndroidJob, statusChan chan status, errChan chan error) (*androidv1alpha1.AndroidJob, bool) {
	var errOcurred bool
	for {
		select {
		case jobStatus, ok := <-statusChan:
			if ok && jobStatus.name != "" {
				reqLogger.Info("Received status update for job", "DeviceName", jobStatus.name, "Status", jobStatus.status)
				instance.Status.JobStatus[jobStatus.name] = jobStatus.status
			}
			if !ok {
				// status channel has closed
				statusChan = nil
			}
		case err, ok := <-errChan:
			if ok {
				reqLogger.Error(err, "Received error on jobs channel")
				errOcurred = true
			}
			if !ok {
				// error channel has closed
				errChan = nil
			}
		}
		if statusChan == nil && errChan == nil {
			break
		}
	}
	return instance, errOcurred
}

func runJobWorker(reqLogger logr.Logger, instance *androidv1alpha1.AndroidJob, device corev1.Pod, jobTemplate *androidv1alpha1.AndroidJobTemplate, statusChan chan status, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	// check if job has already been run
	if status, ok := instance.Status.JobStatus[device.Name]; ok {
		if status.Status == androidv1alpha1.StatusComplete || status.Status == androidv1alpha1.StatusFailed {
			return
		}
	}

	// run the jobs
	jobStatus, err := runDeviceJobs(reqLogger, instance, device, jobTemplate)
	if err != nil {
		errChan <- err
	}
	if jobStatus.Status != "" {
		statusChan <- status{name: device.Name, status: jobStatus}
	}
}

func getTargetDevices(c client.Client, instance *androidv1alpha1.AndroidJob) ([]corev1.Pod, error) {
	targetDevices := make([]corev1.Pod, 0)

	// lookup the device(s)
	if instance.Spec.DeviceName != "" {
		device := &corev1.Pod{}
		if err := c.Get(context.TODO(), instance.DeviceNamespacedName(), device); err != nil {
			return nil, err
		}
		targetDevices = append(targetDevices, *device)
	} else if instance.Spec.DeviceSelector != nil {
		deviceList := &corev1.PodList{}
		if err := c.List(context.TODO(), deviceList, client.InNamespace(instance.Namespace), client.MatchingLabels(instance.Spec.DeviceSelector)); err != nil {
			return nil, err
		}
		targetDevices = append(targetDevices, deviceList.Items...)
	}

	return targetDevices, nil
}

func runDeviceJobs(reqLogger logr.Logger, instance *androidv1alpha1.AndroidJob, device corev1.Pod, jobTemplate *androidv1alpha1.AndroidJobTemplate) (androidv1alpha1.DeviceJobStatus, error) {
	// lookup the adb port for the device
	adbPort, err := util.GetPodADBPort(device)
	if err != nil {
		return androidv1alpha1.DeviceJobStatus{
			Status:  androidv1alpha1.StatusFailed,
			Message: "Could not determine ADB port for device",
		}, nil
	}

	// connect to the device and run the activities
	sess, err := android.NewSession(reqLogger, device.Status.PodIP, adbPort)
	if err != nil {
		return androidv1alpha1.DeviceJobStatus{}, fmt.Errorf("%s: %s", device.Name, err.Error())
	}
	defer sess.Close()

	for _, job := range jobTemplate.Spec.Actions {
		if job.Activity == androidv1alpha1.CommandActivity {
			status, err := runCommandActivity(sess, instance, device, job)
			if err != nil {
				return androidv1alpha1.DeviceJobStatus{}, err
			}
			if status.Status != "" {
				return status, nil
			}
		}
	}

	return androidv1alpha1.DeviceJobStatus{
		Status:  androidv1alpha1.StatusComplete,
		Message: "The job completed successfully",
	}, nil

}

func runCommandActivity(sess android.DeviceSession, instance *androidv1alpha1.AndroidJob, device corev1.Pod, job androidv1alpha1.Action) (androidv1alpha1.DeviceJobStatus, error) {
	for _, cmd := range job.Commands {
		tmplCmd, err := templateCommand(device, cmd)
		if err != nil {
			err = fmt.Errorf("Failed to template command: %s", err.Error())
			return androidv1alpha1.DeviceJobStatus{
				Status:  androidv1alpha1.StatusFailed,
				Message: err.Error(),
			}, nil
		}
		if _, err = sess.RunCommand(job.RunAsRoot, tmplCmd); err != nil {
			return androidv1alpha1.DeviceJobStatus{}, fmt.Errorf("%s: %s", device.Name, err.Error())
		}
	}
	return androidv1alpha1.DeviceJobStatus{}, nil
}

func templateCommand(pod corev1.Pod, cmd string) (string, error) {
	tmpl, err := template.New(pod.Name).Parse(cmd)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	if err := tmpl.Execute(&out, pod); err != nil {
		return "", err
	}
	return out.String(), nil
}
