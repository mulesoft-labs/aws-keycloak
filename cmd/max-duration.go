package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var maxDCmd = &cobra.Command{
	Use:     "max-duration",
	Short:   "Figure out the maximum duration you can request for a session with this role.",
	Example: "  aws-keycloak -p power-devx max-duration",
	RunE:    runMaxDCmd,
}

func init() {
	RootCmd.AddCommand(maxDCmd)
}

func runMaxDCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("%s does not take any arguments", cmd.Use)
	}
	if !regexp.MustCompile("^[a-z-]+$").MatchString(awsrole) {
		return fmt.Errorf("%s is not a valid role to request", awsrole)
	}
	args = strings.Split(fmt.Sprintf("iam get-role --role-name keycloak-%s --query Role.MaxSessionDuration", awsrole), " ")
	return runWithAwsEnv(false, "aws", args...)
}
