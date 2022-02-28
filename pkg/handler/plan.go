package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/tommartensen/tfui/pkg/models"
	"github.com/tommartensen/tfui/pkg/store"
)

func getStateForPlannedChanges(plannedChanges models.PlannedChanges) string {
	if plannedChanges.Create+plannedChanges.Delete+plannedChanges.NoOp+plannedChanges.Update == 0 {
		return "error"
	}
	if plannedChanges.Create > 0 {
		return "create"
	}
	if plannedChanges.Delete > 0 {
		return "delete"
	}
	if plannedChanges.Update > 0 {
		return "update"
	}
	return "no change"
}

func calculatePlannedChanges(plan *models.TfPlan) models.TfPlanSummary {
	plannedChanges := models.PlannedChanges{}
	for _, resource := range plan.ResourceChanges {
		for _, change := range resource.Change.Actions {
			switch change {
			case "create":
				plannedChanges.Create++
			case "delete":
				plannedChanges.Delete++
			case "no-op":
				plannedChanges.NoOp++
			case "update":
				plannedChanges.Update++
			}
		}
	}
	state := getStateForPlannedChanges(plannedChanges)
	return models.TfPlanSummary{
		Meta:           plan.Meta,
		PlannedChanges: plannedChanges,
		State:          state,
	}
}

func writeJSON(w http.ResponseWriter, content interface{}) {
	response, err := json.Marshal(content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func GetAllPlansSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	plans, err := store.ReadAllPlans()
	if err != nil {
		log.Fatalf(err.Error())
	}
	planSummaries := []models.TfPlanSummary{}
	for i := range plans {
		plannedChanges := calculatePlannedChanges(&plans[i])
		planSummaries = append(planSummaries, plannedChanges)
	}
	writeJSON(w, planSummaries)
}

func DeletePlan(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := store.DeletePlan(params["region"], params["project"], params["workspace"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Ok")
}

func readPlanFromStore(w http.ResponseWriter, r *http.Request) (models.TfPlan, error) {
	params := mux.Vars(r)
	plan, err := store.ReadPlan(params["region"], params["project"], params["workspace"])
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return plan, err
	}
	return plan, nil
}

func GetPlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	plan, err := readPlanFromStore(w, r)
	if err != nil {
		return
	}
	writeJSON(w, plan)
}

func GetPlanSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	plan, err := readPlanFromStore(w, r)
	if err != nil {
		return
	}
	plannedChanges := calculatePlannedChanges(&plan)
	writeJSON(w, plannedChanges)
}

func GetPlans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	plans, err := store.ReadAllPlans()
	if err != nil {
		log.Fatalf(err.Error())
	}
	writeJSON(w, plans)
}

func UploadPlan(w http.ResponseWriter, r *http.Request) {
	var plan models.TfPlan
	err := json.NewDecoder(r.Body).Decode(&plan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = store.WritePlan(plan); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Ok")
}
