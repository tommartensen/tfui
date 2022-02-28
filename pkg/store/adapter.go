package store

import (
	"github.com/tommartensen/tfui/pkg/models"
)

func New() error {
	return createFileStore()
}

func WritePlan(plan models.TfPlan) error {
	return storeFilePlan(plan)
}

func ReadAllPlans() ([]models.TfPlan, error) {
	return readAllFilePlans()
}

func ReadPlan(region string, project string, workspace string) (models.TfPlan, error) {
	return readFilePlan(region, project, workspace)
}

func DeleteAllPlans() error {
	return deleteAllFilePlans()
}

func DeletePlan(region string, project string, workspace string) error {
	return deleteFilePlan(region, project, workspace)
}
