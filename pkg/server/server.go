package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/tommartensen/tfui/pkg/config"
	"github.com/tommartensen/tfui/pkg/handler"
	"github.com/tommartensen/tfui/pkg/middleware"
	"github.com/tommartensen/tfui/pkg/store"
)

type Server struct {
	Config *config.Configuration
	Router *mux.Router
}

func (a *Server) Initialize() {
	a.Config = config.New()
	a.Router = mux.NewRouter()
	if err := store.New(); err != nil {
		log.Fatalf("Could not create store: %v", err)
	}

	api := a.Router.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/health", handler.HealthCheck).Methods("GET")
	api.HandleFunc("/reset", handler.ResetSystem).Methods("DELETE")
	api.HandleFunc("/plan", handler.GetPlans).Methods("GET")
	api.HandleFunc("/plan", handler.UploadPlan).Methods("POST")
	api.HandleFunc("/plan/by-params", handler.GetPlan).
		Methods("GET").
		Queries("region", "{region}", "project", "{project}", "workspace", "{workspace}")
	api.HandleFunc("/plan/by-params", handler.DeletePlan).
		Methods("DELETE").
		Queries("region", "{region}", "project", "{project}", "workspace", "{workspace}")
	api.HandleFunc("/plan/summary", handler.GetAllPlansSummary).Methods("GET")
	api.HandleFunc("/plan/summary/by-params", handler.GetPlanSummary).
		Methods("GET").
		Queries("region", "{region}", "project", "{project}", "workspace", "{workspace}")
	if len(a.Config.ApplicationToken) > 0 {
		api.Use(middleware.Middleware)
	}
	a.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist")))
}

func (a *Server) Run() {
	log.Fatal(http.ListenAndServe(":"+a.Config.Port, a.Router))
}

func RunServer(cmd *cobra.Command, args []string) {
	s := Server{}
	s.Initialize()
	s.Run()
}
