module github.com/mulesoft-labs/aws-keycloak

go 1.14

require (
	github.com/99designs/keyring v1.1.5
	github.com/aws/aws-sdk-go v1.19.18
	github.com/golang/mock v1.5.0
	github.com/lordbyron/oauth2-auth-cli v0.0.0-20190425203937-c9a64b4ef0b3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nmrshll/rndm-go v0.0.0-20170430161430-8da3024e53de // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/skratchdot/open-golang v0.0.0-20160302144031-75fb7ed4208c
	github.com/spf13/cobra v1.2.1
	github.com/vaughan0/go-ini v0.0.0-20130923145212-a98ad7ee00ec
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
	golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602
)

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
