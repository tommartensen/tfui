package system

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/tommartensen/tfui/pkg/client"
	"github.com/tommartensen/tfui/pkg/util"
)

func Reset(cmd *cobra.Command, args []string) {
	if util.AskForConfirmation("Execute system reset?") {
		resp, err := client.Request(http.MethodDelete, "reset", nil)
		if err != nil {
			resp.Body.Close()
			log.Fatalln(fmt.Errorf("system reset failed: %v", err))
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalln(fmt.Errorf("system reset failed with %d", resp.StatusCode))
		} else {
			log.Println("System reset successfully")
		}
	} else {
		log.Println("System reset aborted")
	}
}
