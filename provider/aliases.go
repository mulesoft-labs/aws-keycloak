package provider

import (
	"strconv"
	"strings"
)

const (
	DefaultRegion              = "us-east-1"
	DefaultKeycloak            = "id"
	DefaultSAMLSessionDuration = 3600
)

type Aliases map[string]string

func (as Aliases) Exists(alias string) bool {
	_, exists := as[alias]
	return exists
}

func (as Aliases) Lookup(alias string) (kcprofile, awsrole, region string, duration uint64) {
	s := strings.Split(as[alias], ":")
	kcprofile = s[0]
	awsrole = s[1]
	if len(s) >= 3 {
		region = s[2]
	}
	// else region is empty
	if len(s) >= 4 {
		duration, _ = strconv.ParseUint(s[3], 10, 64)
	}
	return
}
