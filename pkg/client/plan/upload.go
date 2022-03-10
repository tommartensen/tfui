package plan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/tommartensen/tfui/pkg/client"
	"github.com/tommartensen/tfui/pkg/models"
	"github.com/tommartensen/tfui/pkg/util"
)

func buildPlanMeta() models.TfPlanMeta {
	return models.TfPlanMeta{
		Region:    os.Getenv("REGION"),
		Project:   os.Getenv("PROJECT"),
		Workspace: os.Getenv("WORKSPACE"),
		Date:      util.GetCurrentDatetime(),
		CommitID:  util.GetCommitID(),
	}
}

func readFile(filepath string) []byte {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Could not read file: %v", err)
	}
	return data
}

func UploadPlan(cmd *cobra.Command, args []string) {
	file := cmd.Flags().Lookup("file").Value.String()
	log.Printf("Uploading plan: %s", file)
	fileContents := readFile(file)

	tfPlan := models.TfPlan{}
	err := json.Unmarshal(fileContents, &tfPlan)
	if err != nil {
		log.Fatalf("Could not upload plan because JSON parsing failed: %v", err)
	}

	tfPlan.Meta = buildPlanMeta()
	body, err := json.Marshal(tfPlan)
	if err != nil {
		log.Fatalf("Could not upload plan because JSON marshalling failed: %v", err)
	}

	resp, err := client.Request(http.MethodPost, "plan", body)
	if err != nil {
		resp.Body.Close()
		log.Fatalln(fmt.Errorf("plan upload failed: %v", err))
	}
	if resp.StatusCode != http.StatusCreated {
		log.Fatalln(fmt.Errorf("plan upload failed (HTTP %d)", resp.StatusCode))
	} else {
		log.Println("Plan uploaded")
	}
}
