package mBinding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateZoneDomain(t *testing.T) {
	type testCase struct {
		Name   string
		Domain string
		Valid  bool
	}

	testCases := []testCase{
		testCase{
			Name:   "invalid ending 1",
			Domain: "cn2apollo.backup2-prod.chinaslb.com.lucky.test.error.cdn.cname.test2",
			Valid:  false,
		},
		testCase{
			Name:   "invalid ending 2",
			Domain: "cn2apollo.backup2-prod.chinaslb.com.lucky.test.error.cdn.cname.test-",
			Valid:  false,
		},
		testCase{
			Name:   "invalid ending 3",
			Domain: "cn2apollo.backup2-prod.chinaslb.com.lucky.test.error.cdn.cname.te-st",
			Valid:  false,
		},
		testCase{
			Name:   "normal 1",
			Domain: "cn2apollo.backup2-prod.chinaslb.com.lucky.test.error.cdn.cname.test",
			Valid:  true,
		},
		testCase{
			Name:   "normal 2",
			Domain: "cn2-apollo.backup2-prod.chinaslb.com.lucky.test.error.cdn.cname.test",
			Valid:  true,
		},
		testCase{
			Name:   "duplicate dash",
			Domain: "cn2apollo.backup2--prod.chinaslb.com.lucky.test.error.cdn.cname.test",
			Valid:  true,
		},
		testCase{
			Name:   "wildcard",
			Domain: "*.chinaslb.com.lucky.test.error.cdn.cname.test",
			Valid:  true,
		},
		testCase{
			Name:   "invalid wildcard",
			Domain: "*.chinaslb.com.lucky.*.error.cdn.cname.test",
			Valid:  false,
		},
	}

	for _, tt := range testCases {
		assert.Equal(t, tt.Valid, hostnameRegexRFC1123AndWildcard.MatchString(tt.Domain), tt.Name)
	}
}

func TestValidateFQDNDomain(t *testing.T) {
	type testCase struct {
		Name   string
		Domain string
		Valid  bool
	}

	testCases := []testCase{
		testCase{
			Name:   "invalid ending 1",
			Domain: "cn2apollo.backup2-prod.chinaslb.com.lucky.test.error.cdn.cname.test2",
			Valid:  false,
		},
		testCase{
			Name:   "invalid ending 2",
			Domain: "cn2apollo.backup2-prod.chinaslb.com.lucky.test.error.cdn.cname.test-",
			Valid:  false,
		},
		testCase{
			Name:   "invalid ending 3",
			Domain: "cn2apollo.backup2-prod.chinaslb.com.lucky.test.error.cdn.cname.te-st",
			Valid:  false,
		},
		testCase{
			Name:   "normal 1",
			Domain: "cn2apollo.backup2-prod.chinaslb.com.lucky.test.error.cdn.cname.test",
			Valid:  true,
		},
		testCase{
			Name:   "normal 2",
			Domain: "cn2-apollo.backup2-prod.chinaslb.com.lucky.test.error.cdn.cname.test",
			Valid:  true,
		},
		testCase{
			Name:   "duplicate dash",
			Domain: "cn2apollo.backup2--prod.chinaslb.com.lucky.test.error.cdn.cname.test",
			Valid:  true,
		},
		testCase{
			Name:   "not support wildcard",
			Domain: "*.chinaslb.com.lucky.test.error.cdn.cname.test",
			Valid:  false,
		},
		testCase{
			Name:   "invalid wildcard",
			Domain: "*.chinaslb.com.lucky.*.error.cdn.cname.test",
			Valid:  false,
		},
	}

	for _, tt := range testCases {
		assert.Equal(t, tt.Valid, fqdnDomainRegex.MatchString(tt.Domain), tt.Name)
	}
}
