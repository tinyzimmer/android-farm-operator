package api

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/tinyzimmer/android-farm-operator/pkg/util"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/android"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var logger = log.Log.WithName("api-server")

type FarmAPI interface {
	PostCommand(namespace, device, command string) (out []byte, err error)
	GetFile(namespace, device, path string, writer io.Writer) (err error)
}

type farmAPI struct {
	client client.Client
}

func NewFarmAPI(c client.Client) FarmAPI {
	return &farmAPI{client: c}
}

func (f *farmAPI) getDevice(namespace, device string) (*corev1.Pod, error) {
	nn := types.NamespacedName{Name: device, Namespace: namespace}
	pod := &corev1.Pod{}
	return pod, f.client.Get(context.TODO(), nn, pod)
}

func (f *farmAPI) getSession(pod *corev1.Pod) (android.DeviceSession, error) {
	port, err := util.GetPodADBPort(*pod)
	if err != nil {
		return nil, err
	}
	return android.NewSession(logger, pod.Status.PodIP, port)
}

func (f *farmAPI) PostCommand(namespace, device, command string) (out []byte, err error) {
	pod, err := f.getDevice(namespace, device)
	if err != nil {
		return nil, errors.NewAPIError(err.Error())
	}
	sess, err := f.getSession(pod)
	if err != nil {
		return nil, errors.NewAPIError(err.Error())
	}
	defer sess.Close()
	out, err = sess.RunCommand(true, command)
	if err != nil {
		return nil, errors.NewAPIError(err.Error())
	}
	return out, nil
}

func (f *farmAPI) GetFile(namespace, device, fpath string, writer io.Writer) (err error) {
	pod, err := f.getDevice(namespace, device)
	if err != nil {
		return errors.NewAPIError(err.Error())
	}
	sess, err := f.getSession(pod)
	if err != nil {
		return errors.NewAPIError(err.Error())
	}
	defer sess.Close()
	if err := sess.DownloadFile(path.Clean(fpath), writer); err != nil {
		return errors.NewAPIError(fmt.Sprintf("Could not retrieve %s from device", fpath))
	}
	return nil
}
