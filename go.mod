module github.com/mulesoft-labs/aws-keycloak

go 1.14

require (
	github.com/99designs/keyring v1.1.5
	github.com/aws/aws-sdk-go v1.19.18
	github.com/golang/mock v1.2.0
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/lordbyron/oauth2-auth-cli v0.0.0-20190425203937-c9a64b4ef0b3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nmrshll/rndm-go v0.0.0-20170430161430-8da3024e53de // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/skratchdot/open-golang v0.0.0-20160302144031-75fb7ed4208c
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/vaughan0/go-ini v0.0.0-20130923145212-a98ad7ee00ec
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/net v0.0.0-20190603091049-60506f45cf65
	golang.org/x/oauth2 v0.0.0-20181120190819-8f65e3013eba
	google.golang.org/appengine v1.6.6 // indirect
)

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
