package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tommartensen/tfui/pkg/models"
	"github.com/tommartensen/tfui/pkg/server"
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s := server.Server{}
	os.Setenv("BASE_DIR", "../../test/plans")
	defer os.Unsetenv("BASE_DIR")
	s.Initialize()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, actual, expected int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func setup() {
	os.MkdirAll("../../test/plans", 0744)
	_, err := os.Stat("../../test/plans/us-west-2_tfui_default.tfplan")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.Link("../../test/plan.json", "../../test/plans/us-west-2_tfui_default.tfplan")
			if err != nil {
				log.Fatal(err.Error())
			}
		} else {
			log.Fatal(err.Error())
		}
	}
}

func teardown() {
	if err := os.RemoveAll("../../test/plans"); err != nil {
		log.Fatal(err.Error())
	}
}

func TestGetPlan(t *testing.T) {
	setup()

	// Not existent plan
	req, err := http.NewRequest("GET", "/api/plan/by-params?region=not-exists&workspace=default&project=tfui", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	response := executeRequest(req)
	checkResponseCode(t, response.Code, http.StatusNotFound)

	// Existing plan
	req, err = http.NewRequest("GET", "/api/plan/by-params?region=us-west-2&workspace=default&project=tfui", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	response = executeRequest(req)
	checkResponseCode(t, response.Code, http.StatusOK)

	plan := models.TfPlan{}
	err = json.NewDecoder(response.Body).Decode(&plan)
	if err != nil {
		log.Fatal(err.Error())
	}
	assert.Equal(t, "1.0.10", plan.TfVersion)
	assert.Equal(t, 18, len(plan.ResourceChanges))
	assert.Equal(t, "tfui", plan.Meta.Project)
	assert.Equal(t, "us-west-2", plan.Meta.Region)
	teardown()
}

func TestGetAllPlansSummary(t *testing.T) {
	setup()
	req, err := http.NewRequest("GET", "/api/plan/summary", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	response := executeRequest(req)
	checkResponseCode(t, response.Code, http.StatusOK)
	assert.Equal(t, response.Header().Get("Content-Type"), "application/json")

	summaries := []models.TfPlanSummary{}
	err = json.NewDecoder(response.Body).Decode(&summaries)
	if err != nil {
		log.Fatal(err.Error())
	}
	assert.Equal(t, summaries[0].PlannedChanges.Create, 18)
	assert.Equal(t, summaries[0].PlannedChanges.Delete, 0)
	assert.Equal(t, summaries[0].PlannedChanges.NoOp, 0)
	assert.Equal(t, summaries[0].PlannedChanges.Update, 0)
	teardown()
}

func TestDeletePlan(t *testing.T) {
	setup()
	req, err := http.NewRequest("DELETE", "/api/plan/by-params?region=us-west-2&workspace=default&project=tfui", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	response := executeRequest(req)
	checkResponseCode(t, response.Code, http.StatusOK)
	teardown()

	// TODO: delete of not existing plan fails with 404
}

func TestGetPlanSummary(t *testing.T) {
	setup()
	req, err := http.NewRequest("GET", "/api/plan/summary/by-params?region=us-west-2&workspace=default&project=tfui", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	response := executeRequest(req)
	checkResponseCode(t, response.Code, http.StatusOK)
	assert.Equal(t, response.Header().Get("Content-Type"), "application/json")

	summary := models.TfPlanSummary{}
	err = json.NewDecoder(response.Body).Decode(&summary)
	if err != nil {
		log.Fatal(err.Error())
	}
	assert.Equal(t, summary.Meta.Project, "tfui")
	assert.Equal(t, summary.Meta.Region, "us-west-2")
	assert.Equal(t, summary.State, "create")
	assert.Equal(t, summary.PlannedChanges.Create, 18)
	assert.Equal(t, summary.PlannedChanges.Delete, 0)
	assert.Equal(t, summary.PlannedChanges.NoOp, 0)
	assert.Equal(t, summary.PlannedChanges.Update, 0)
	teardown()
}

func TestGetPlans(t *testing.T) {
	setup()
	req, err := http.NewRequest("GET", "/api/plan", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	response := executeRequest(req)
	checkResponseCode(t, response.Code, http.StatusOK)
	assert.Equal(t, response.Header().Get("Content-Type"), "application/json")

	plans := []models.TfPlan{}
	err = json.NewDecoder(response.Body).Decode(&plans)
	if err != nil {
		log.Fatal(err.Error())
	}
	assert.Equal(t, len(plans), 1)
	assert.Equal(t, plans[0].TfVersion, "1.0.10")
	assert.Equal(t, len(plans[0].ResourceChanges), 18)
	teardown()
}

func TestUploadPlan(t *testing.T) {
	file, _ := ioutil.ReadFile("../../test/plan.json")
	body := bytes.NewBuffer(file)
	req, err := http.NewRequest("POST", "/api/plan", body)
	if err != nil {
		log.Fatal(err.Error())
	}
	response := executeRequest(req)
	checkResponseCode(t, response.Code, http.StatusCreated)
	teardown()
}
