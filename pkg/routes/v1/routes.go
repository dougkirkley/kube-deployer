package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	appsv1 "github.com/dougkirkley/kube-deployer/pkg/controllers/apps/v1"
	containersv1 "github.com/dougkirkley/kube-deployer/pkg/controllers/containers/v1"
)

// Handlers handles v1 routes
func Handlers() *mux.Router {
	var r = mux.NewRouter()
	r = mux.NewRouter().StrictSlash(true)
	r.Use(CommonMiddleware)
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/health", appsv1.Health).Methods("GET")
	s.HandleFunc("/apps", appsv1.List).Methods("GET")
	s.HandleFunc("/apps/{id}", appsv1.List).Methods("GET")
	s.HandleFunc("/apps/{id}", appsv1.Remove).Methods("DELETE")
	s.HandleFunc("/upload", appsv1.Install).Methods("POST")
	s.HandleFunc("/upgrade/{id}", appsv1.Upgrade).Methods("POST")
	s.HandleFunc("/export/{id}", appsv1.Export).Methods("PUT")
	//container handlers
	s.HandleFunc("/containers", containersv1.ListPods).Methods("GET")
	return r
}

// CommonMiddleware --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "OPTIONS" {
			if origin := w.Header().Get("Origin"); origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			w.Header().Add("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
			next.ServeHTTP(w, r)
		}
	})
}
