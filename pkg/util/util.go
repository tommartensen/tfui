package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func AskForConfirmation(msg string) bool {
	var response string
	fmt.Printf("%s (y/n): ", msg)
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	}
	if containsString(nokayResponses, response) {
		return false
	}
	fmt.Println("Please type yes or no and then press enter.")
	return AskForConfirmation(msg)
}

func containsString(slice []string, element string) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetCurrentDatetime() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02 15:04:05")
}

func GetCommitID() string {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}
