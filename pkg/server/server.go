package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/tinyzimmer/android-farm-operator/pkg/server/api"
	"github.com/tinyzimmer/android-farm-operator/pkg/util/errors"
	"github.com/gorilla/mux"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var defaultWait = time.Duration(15)
var gracefulWait time.Duration

func init() {
	wait := os.Getenv("API_GRACEFUL_SHUTDOWN")
	if wait != "" {
		val, err := strconv.Atoi(wait)
		if err != nil {
			panic(err)
		}
		gracefulWait = time.Second * time.Duration(val)
	} else {
		gracefulWait = time.Second * defaultWait
	}
}

type webServer struct {
	api api.FarmAPI
}

type commandRequest struct {
	Command string `json:"command"`
}

func (s *webServer) runDeviceCommand(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req commandRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse([]byte(`{"error": "Could not read request body"}`), w)
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeResponse([]byte(fmt.Sprintf(`{"error": "Could not decode request body: %s"}`, err.Error())), w)
		return
	}
	namespace, device, _ := getVars(r)
	out, err := s.api.PostCommand(namespace, device, req.Command)
	if err != nil {
		if apierr, ok := errors.IsAPIError(err); ok {
			writeResponse(apierr.ErrorJSON(), w)
			return
		} else {
			writeResponse([]byte(err.Error()), w)
			return
		}
	}
	res, err := json.MarshalIndent(map[string]string{
		"stdout": string(out),
	}, "  ", "")
	if err != nil {
		writeResponse([]byte(`{"error": "Could not write response body"}`), w)
		return
	}
	writeResponse(append(res, []byte("\n")...), w)
}

func (s *webServer) getDeviceFile(w http.ResponseWriter, r *http.Request) {
	namespace, device, path := getVars(r)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(path)))
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	if err := s.api.GetFile(namespace, device, fmt.Sprintf("/%s", path), w); err != nil {
		w.Header().Set("Content-Disposition", "")
		w.Header().Set("Content-Type", "application/json")
		if apierr, ok := errors.IsAPIError(err); ok {
			writeResponse(apierr.ErrorJSON(), w)
		} else {
			writeResponse([]byte(err.Error()), w)
		}
	}
}

func writeResponse(res []byte, w http.ResponseWriter) {
	if _, err := w.Write(res); err != nil {
		log.Println("Error writing response:", err)
	}
}

func getVars(r *http.Request) (namespace, device, path string) {
	vars := mux.Vars(r)
	namespace = vars["namespace"]
	device = vars["device"]
	path = vars["path"]
	return
}

func RunServer(stopChan <-chan struct{}, client client.Client) {

	// create a new router
	r := mux.NewRouter()

	// Add routes
	websrv := &webServer{api: api.NewFarmAPI(client)}
	r.HandleFunc("/{namespace}/{device}/command", websrv.runDeviceCommand).
		Methods("POST")
	r.PathPrefix("/{namespace}/{device}/{path:.*}").
		HandlerFunc(websrv.getDeviceFile).
		Methods("GET")

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	// Block until we receive our signal.
	<-stopChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), gracefulWait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Error while shutting down API server:", err)
		return
	}
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("Shut down API server")
}
