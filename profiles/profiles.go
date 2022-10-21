package profiles

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// Uses exec to run AWS CLI command to list profiles
func GetProfiles() []string {
	cmd := exec.Command("aws", "configure", "list-profiles")
	cmd.Stderr = os.Stderr
	data, err := cmd.Output()

	if err != nil {
		log.Fatalf("Failed to call cmd.Output(): %v", err)
	}

	profiles := strings.Split(string(data), "\n")

	return profiles
}
