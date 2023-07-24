package supertokens

import (
	"context"
	"net/http"
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

func TestMakeDefaultUserContextFromAPI(t *testing.T) {
	var validRID = "valid-request-id"
	superTokensInstance = &superTokens{RequestIDKey: RequestIDKey}

	reqCtx, _ := http.NewRequestWithContext(context.WithValue(context.TODO(), RequestIDKey, validRID),
		"GET", "", nil)

	reqHeader, _ := http.NewRequest("", "", nil)
	reqHeader.Header.Set("X-Request-ID", validRID)

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			"set valid requestID from header",
			args{r: reqHeader},
			validRID,
		},
		{
			"set valid requestID from context",
			args{r: reqCtx},
			validRID,
		},
		{
			"failure - no requestID present in context or header",
			args{r: &http.Request{}},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, (*MakeDefaultUserContextFromAPI(tt.args.r))[RequestIDKey], "MakeDefaultUserContextFromAPI(%v)", tt.args.r)
		})
	}
}
