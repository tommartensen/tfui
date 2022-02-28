package handler

import (
	"fmt"
	"net/http"

	"github.com/tommartensen/tfui/pkg/store"
)

func ResetSystem(w http.ResponseWriter, r *http.Request) {
	if err := store.DeleteAllPlans(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Ok")
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok")
}
