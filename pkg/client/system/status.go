package system

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/tommartensen/tfui/pkg/client"
)

func Status(cmd *cobra.Command, args []string) {
	resp, err := client.Request(http.MethodGet, "health", nil)
	if err != nil {
		resp.Body.Close()
		log.Fatalln(fmt.Errorf("system status request failed: %v", err))
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalln(fmt.Errorf("system status not healthy (HTTP %d)", resp.StatusCode))
	} else {
		log.Println("System is healthy")
	}
}
