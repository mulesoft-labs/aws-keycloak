package provider_test

import (
	"fmt"
	"testing"

	"github.com/mulesoft-labs/aws-keycloak/provider"
)

func testAlias(a provider.Aliases, alias, _kcprofile, _awsrole, _region string, _duration uint64) error {
	kcprofile, awsrole, region, duration := a.Lookup(alias)
	if kcprofile != _kcprofile {
		return fmt.Errorf("kcprofile does not match expected. actual: %s - expected: %s", kcprofile, _kcprofile)
	}
	if awsrole != _awsrole {
		return fmt.Errorf("awsrole does not match expected. actual: %s - expected: %s", awsrole, _awsrole)
	}
	if region != _region {
		return fmt.Errorf("region does not match expected. actual: %s - expected: %s", region, _region)
	}
	if duration != _duration {
		return fmt.Errorf("duration does not match expected. actual: %d - expected: %d", duration, _duration)
	}
	return nil
}

func TestAliasMulti(t *testing.T) {
	a := provider.Aliases{
		"t2": "id:power-devx",
		"t3": "id:power-devx:us-east-1",
		"t4": "id:power-devx:us-east-1:7200",
	}
	if err := testAlias(a, "t2", "id", "power-devx", "", 0); err != nil {
		t.Error(err)
	}
	if err := testAlias(a, "t3", "id", "power-devx", "us-east-1", 0); err != nil {
		t.Error(err)
	}
	if err := testAlias(a, "t4", "id", "power-devx", "us-east-1", 7200); err != nil {
		t.Error(err)
	}
}
