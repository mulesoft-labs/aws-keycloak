package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/99designs/keyring"
	log "github.com/Sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"

	"github.com/mulesoft-labs/aws-keycloak/provider"
)

// Errors returned from frontend commands
var (
	ErrCommandMissing         = errors.New("must specify command to run")
	ErrTooManyArguments       = errors.New("too many arguments")
	ErrTooFewArguments        = errors.New("too few arguments")
	ErrFailedToSetCredentials = errors.New("Failed to set credentials in your keyring")
)

const (
	KeycloakConfigUrl = "https://wiki.corp.mulesoft.com/download/attachments/53909517/keycloak-config?api=v2"
)

// global flags
var (
	backend    string
	kr         keyring.Keyring
	debug      bool
	quiet      bool
	configFile string
	kcprofile  string
	awsrole    string
	region     string
	kcConf     map[string]string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:               "aws-keycloak [flags] -- <command>",
	Short:             "aws-keycloak allows you to authenticate with AWS using your keycloak credentials",
	Example:           "  aws-keycloak -p power-devx -- aws sts get-caller-identity",
	SilenceUsage:      true,
	SilenceErrors:     true,
	PersistentPreRunE: prerun,
	RunE:              runCommand,
	Version:           "1.4.1",
}

func runCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return ErrTooFewArguments
	}
	return runWithAwsEnv(true, args[0], args[1:]...)
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		switch err {
		case ErrTooFewArguments, ErrTooManyArguments:
			RootCmd.Usage()
		}
		os.Exit(1)
	}
}

func prerun(cmd *cobra.Command, args []string) error {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else if quiet {
		log.SetLevel(log.ErrorLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if cmd.Name() == "help" {
		return nil
	}

	// Load backend from env var if not set as a flag
	if !cmd.Flags().Lookup("backend").Changed {
		if backendFromEnv, ok := os.LookupEnv("AWS_KEYCLOAK_BACKEND"); ok {
			backend = backendFromEnv
		}
	}

	var allowedBackends []keyring.BackendType
	if backend != "" {
		allowedBackends = append(allowedBackends, keyring.BackendType(backend))
	}
	ring, err := keyring.Open(keyring.Config{
		AllowedBackends:          allowedBackends,
		KeychainTrustApplication: true,
		ServiceName:              "keycloak-login",
		LibSecretCollectionName:  "awsvault",
	})
	kr = ring
	if err != nil {
		return err
	}

	if !cmd.Flags().Lookup("config").Changed {
		configFile, err = provider.EnvFileOrDefault()
		if err != nil {
			return err
		}
	}

	config, err := provider.NewConfigFromFile(configFile)
	if err != nil {
		log.Errorf("No configuration found at %s.", configFile)
		fetchConfig()
		return fmt.Errorf("Please install configuration file and try again.")
	}

	sections, err := config.Parse()
	if err != nil {
		return err
	}

	// So hacky!
	if cmd == openCmd && len(args) == 1 {
		awsrole = args[0]
	}

	aliases := provider.Aliases(sections["aliases"])
	if aliases.Exists(awsrole) {
		alias := awsrole
		kcprofile, awsrole, region = aliases.Lookup(alias)
		log.Debugf("Found alias for %s: %s %s %s", alias, kcprofile, awsrole, region)
	}

	kcConf = sections[kcprofile]
	if len(kcConf) == 0 {
		return fmt.Errorf("No keycloak profile found at %s", kcprofile)
	}

	return nil
}

func init() {
	backendsAvailable := []string{}
	for _, backendType := range keyring.AvailableBackends() {
		backendsAvailable = append(backendsAvailable, string(backendType))
	}
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", provider.DefaultConf, "Keycloak provider configuration")
	RootCmd.PersistentFlags().StringVarP(&backend, "backend", "b", "", fmt.Sprintf("Secret backend to use %s", backendsAvailable))
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")
	RootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Minimize output")
	RootCmd.PersistentFlags().StringVarP(&kcprofile, "keycloak-profile", "k", provider.DefaultKeycloak, "Keycloak system to auth to")
	RootCmd.PersistentFlags().StringVarP(&awsrole, "profile", "p", "", "AWS profile to run against (recommended)")
}

func fetchConfig() {
	fmt.Print("You need to put a keycloak-config file in your ~/.aws/ directory.\nWould you like to download one? [y/N] ")
	var choice string
	fmt.Scanln(&choice)
	if choice == "y" || choice == "Y" {
		open.Run(KeycloakConfigUrl)
	}
	fmt.Println(KeycloakConfigUrl)
}
