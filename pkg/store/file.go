package store

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tommartensen/tfui/pkg/config"
	"github.com/tommartensen/tfui/pkg/models"
)

var extension = "tfplan"

func buildFilename(region, project, workspace string) string {
	var baseDir = config.New().BaseDir
	return fmt.Sprintf("%s/%s_%s_%s.%s",
		baseDir, region, project, workspace, extension,
	)
}

func buildFilenameFromPlan(plan models.TfPlan) string {
	return buildFilename(plan.Meta.Region, plan.Meta.Project, plan.Meta.Workspace)
}

func createFileStore() error {
	var baseDir = config.New().BaseDir
	return os.MkdirAll(baseDir, 0744)
}

func deleteAllFilePlans() error {
	filenames, err := findPlansInDir()
	if err != nil {
		return err
	}
	for _, path := range filenames {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

func deleteFilePlan(region, project, workspace string) error {
	filename := buildFilename(region, project, workspace)
	return os.Remove(filename)
}

func findPlansInDir() ([]string, error) {
	var baseDir = config.New().BaseDir
	filenames := []string{}
	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		return filenames, err
	}
	for _, file := range files {
		if isPlanFile(file) {
			filenames = append(
				filenames,
				fmt.Sprintf("%s/%s", baseDir, file.Name()),
			)
		}
	}
	return filenames, nil
}

func isPlanFile(file fs.FileInfo) bool {
	return !file.IsDir() && strings.HasSuffix(file.Name(), "."+extension)
}

func readAllFilePlans() ([]models.TfPlan, error) {
	plans := []models.TfPlan{}
	files, err := findPlansInDir()
	if err != nil {
		return plans, err
	}
	for _, filename := range files {
		plan := models.TfPlan{}
		if err := readPlanFromDir(filename, &plan); err != nil {
			return plans, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

func readFilePlan(region, project, workspace string) (models.TfPlan, error) {
	plan := models.TfPlan{}
	filename := buildFilename(region, project, workspace)
	if err := readPlanFromDir(filename, &plan); err != nil {
		return plan, err
	}
	return plan, nil
}

func readPlanFromDir(filename string, plan *models.TfPlan) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &plan)
}

func storeFilePlan(plan models.TfPlan) error {
	filename := buildFilenameFromPlan(plan)
	file, err := json.Marshal(plan)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, file, 0600)
	if err != nil {
		return err
	}
	return nil
}
