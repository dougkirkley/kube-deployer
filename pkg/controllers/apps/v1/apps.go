package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	chartutil "github.com/dougkirkley/kube-deployer/pkg/chartutil/v1"
	"helm.sh/helm/v3/pkg/action"
)

// ErrorResponse struct type
type ErrorResponse struct {
	Status int
	Err    string
}

var config = CreateConfig("")

// CreateConfig returns in cluster config
func CreateConfig(namespace string) *action.Configuration {
	// creates the cluster config
	config := new(action.Configuration)
	if err := config.Init(nil, namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		log.Printf(format, v)
	}); err != nil {
		panic(err)
	}
	return config

}

// Health function tests API
func Health(w http.ResponseWriter, r *http.Request) {
	if err := config.KubeClient.IsReachable(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("API is healthy")
	}
}

// List lists charts or one chart
func List(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	namespace := r.URL.Query().Get("namespace")
	if namespace != "" {
		config = CreateConfig(namespace)
	}
	if id != "" {
		chartGet := action.NewGet(config)
		get, err := chartGet.Run(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err.Error())
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(get)
	} else {
		chartlist := action.NewList(config)
		list, err := chartlist.Run()
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err.Error())
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(list)
	}
}

// Install creates new charts
func Install(w http.ResponseWriter, r *http.Request) {
	// Create the file
	var chartFile = "/tmp/chart_install.tgz"
	log.Print("Creating chart file")
	out, err := os.Create(chartFile)
	if err != nil {
		log.Print(err.Error())
	}
	defer out.Close()

	// Write the body to file
	log.Print("Writing chart to file")
	_, copyerr := io.Copy(out, r.Body)
	if copyerr != nil {
		var response = ErrorResponse{
			Status: 400,
			Err:    copyerr.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	}
	log.Print("Loading chart")
	chart, err := chartutil.ChartLoader(chartFile)
	if err != nil {
		var response = ErrorResponse{
			Status: 400,
			Err:    err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	}
	log.Print("Loaded chart...")
	newInstall := action.NewInstall(config)
	install, err := newInstall.Run(chart, nil)
	if err != nil {
		var response = ErrorResponse{
			Status: 400,
			Err:    err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	}
	if err := os.Remove(chartFile); err != nil {
		log.Print("failed to delete chart file")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fmt.Sprintf("Successfully installed release: ", install.Name))
}

// Export downloads charts
func Export(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	namespace := r.URL.Query().Get("namespace")
	if namespace != "" {
		config = CreateConfig(namespace)
	}
	if id == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode("Must specify name to export")
	}
	exportChart := action.NewChartExport(config)
	var out io.Writer
	err := exportChart.Run(out, id)
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode("Not Implemented")
	}
}

// Upgrade upgrades charts
func Upgrade(w http.ResponseWriter, r *http.Request) {
	//TODO
	namespace := r.URL.Query().Get("namespace")
	if namespace != "" {
		config = CreateConfig(namespace)
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode("Not Implemented")
}

// Remove deletes charts
func Remove(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	namespace := r.URL.Query().Get("namespace")
	if namespace != "" {
		config = CreateConfig(namespace)
	}
	if id == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode("Must specify name to delete")
	}
	chartRemove := action.NewChartRemove(config)
	var out io.Writer
	err := chartRemove.Run(out, id)
	if err != nil {
		log.Print(err.Error())
		json.NewEncoder(w).Encode(err.Error())
	}
	var response = fmt.Sprintf("Deleted chart: %s", id)
	json.NewEncoder(w).Encode(response)
}
