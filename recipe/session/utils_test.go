package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type GetTopLevelDomainForSameSiteResolutionTest struct {
	Input  string
	Output string
}

func TestGetTopLevelDomainForSameSiteResolution(t *testing.T) {
	input := []GetTopLevelDomainForSameSiteResolutionTest{{
		Input:  "http://a.b.test.com",
		Output: "test.com",
	}, {
		Input:  "https://a.b.test.com",
		Output: "test.com",
	}, {
		Input:  "http://a.b.test.co.uk",
		Output: "test.co.uk",
	}, {
		Input:  "http://test.com",
		Output: "test.com",
	}, {
		Input:  "https://test.com",
		Output: "test.com",
	}, {
		Input:  "http://localhost",
		Output: "localhost",
	}, {
		Input:  "http://localhost.org",
		Output: "localhost",
	}, {
		Input:  "http://8.8.8.8",
		Output: "localhost",
	}, {
		Input:  "http://8.8.8.8:8080",
		Output: "localhost",
	}, {
		Input:  "http://localhost:3000",
		Output: "localhost",
	}, {
		Input:  "http://test.com:3567",
		Output: "test.com",
	}, {
		Input:  "https://test.com:3567",
		Output: "test.com",
	}}
	for _, val := range input {
		domain, _ := GetTopLevelDomainForSameSiteResolution(val.Input)
		assert.Equal(t, val.Output, domain, val.Input)
	}
}
