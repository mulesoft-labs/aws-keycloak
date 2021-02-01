package cmd

import (
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:     "check",
	Short:   "Check will authenticate you through keycloak and store session.",
	Example: "  aws-keycloak -p power-devx check",
	Args:    cobra.MaximumNArgs(0),
	RunE:    runCheck,
}

func init() {
	RootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) error {
	return runWithAwsEnv(true, "aws", "sts", "get-caller-identity")
}
