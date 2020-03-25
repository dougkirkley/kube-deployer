package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"

	routes "github.com/dougkirkley/kube-deployer/pkg/routes/v1"
)

func main() {
	var Port string = "9090"
	// Handle routes
	router := routes.Handlers()
	http.Handle("/api/v1/", router)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe(":"+Port, handlers.LoggingHandler(os.Stdout, router)))
}