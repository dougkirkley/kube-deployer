package controllers

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	appsv1 "github.com/dougkirkley/kube-deployer/pkg/controllers/apps/v1"
)

var clientset = CreateClient()

// CreateClient returns in cluster config clientset
func CreateClient() *kubernetes.Clientset {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return clientset
}

// ListPods lists all Pods or one Pod
func ListPods(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]

	namespace := r.URL.Query().Get("namespace")
	PodClient := clientset.CoreV1().Pods(namespace)

	if id == "" {
		list, err := PodClient.List(metav1.ListOptions{})
		if err != nil {
			log.Print(err.Error())
		}
		json.NewEncoder(w).Encode(list)
	} else {
		list, err := PodClient.Get(id, metav1.GetOptions{})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			errResp := &appsv1.ErrorResponse{Status: 404, Err: fmt.Sprintf("Error, %v", err.Error())}
			json.NewEncoder(w).Encode(errResp)
			log.Print(err.Error())
		}
		json.NewEncoder(w).Encode(list)
	}
}