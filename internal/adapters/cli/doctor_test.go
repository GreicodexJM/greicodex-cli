package cli

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestAddDoctorCommand(t *testing.T) {
	rootCmd := &cobra.Command{}
	AddDoctorCommand(rootCmd)

	// Check that the command was added
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "doctor" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Doctor command was not added to the root command")
	}
}
