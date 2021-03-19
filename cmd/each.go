package cmd

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	eachFilter string
	printEnv   bool
)

var eachCmd = &cobra.Command{
	Use:     "each",
	Short:   "Run the command for each matching profile available to you",
	Example: " aws-keycloak each -- aws iam list-account-aliases",
	Args:    cobra.MinimumNArgs(1),
	RunE:    runEach,
}

func init() {
	eachCmd.PersistentFlags().StringVarP(&eachFilter, "filter", "f", "", "Regex to filter listed roles (eg. 'admin').")
	eachCmd.PersistentFlags().BoolVarP(&printEnv, "print-env", "", false, "Print the name of each env to stdout before running.")
	RootCmd.AddCommand(eachCmd)
}

func runEach(cmd *cobra.Command, args []string) error {
	filter, err := regexp.Compile(eachFilter)
	if err != nil {
		return err
	}

	roles, err := listRoles()
	if err != nil {
		return err
	}

	re := regexp.MustCompile("role/keycloak-([^/]+)$")
	for _, role := range roles {
		if !filter.MatchString(role) {
			continue
		}
		p := re.FindStringSubmatch(role)
		awsrole = p[1]
		log.Infof("role %s\n", awsrole)
		if printEnv {
			fmt.Printf("# aws-keycloak profile: %s\n", awsrole)
		}
		err = runWithAwsEnv(true, args[0], args[1:]...)
		if err != nil {
			return err
		}
	}

	return nil
}
