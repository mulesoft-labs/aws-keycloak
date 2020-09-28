package cmd_test

import (
	"testing"

	"github.com/mulesoft-labs/aws-keycloak/cmd"
)

func TestValidateProfile(t *testing.T) {
	profiles := []string{
		"dev-rtf",
		"ms-rtf01",
		"wws-sandbox",
	}
	for _, p := range profiles {
		if !cmd.ValidateProfile(p) {
			t.Errorf("Profile `%s` should validate, but it failed", p)
		}
	}
}

func TestDontValidateProfile(t *testing.T) {
	profiles := []string{
		"profile with spaces",
		"devx;reboot",
		"devx|cat",
	}
	for _, p := range profiles {
		if cmd.ValidateProfile(p) {
			t.Errorf("Profile `%s` should NOT validate, but it did", p)
		}
	}
}
