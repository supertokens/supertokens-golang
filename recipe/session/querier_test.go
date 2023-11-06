package session

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func resetQuerier() {
	supertokens.SetQuerierApiVersionForTests("")
}

func TestThatNetworkCallIsRetried(t *testing.T) {
	resetAll()
	mux := http.NewServeMux()

	numberOfTimesCalled := 0
	numberOfTimesSecondCalled := 0
	numberOfTimesThirdCalled := 0

	mux.HandleFunc("/testing", func(rw http.ResponseWriter, r *http.Request) {
		numberOfTimesCalled++
		rw.WriteHeader(supertokens.RateLimitStatusCode)
		rw.Header().Set("Content-Type", "application/json")
		response, err := json.Marshal(map[string]interface{}{})
		if err != nil {
			t.Error(err.Error())
		}
		rw.Write(response)
	})

	mux.HandleFunc("/testing2", func(rw http.ResponseWriter, r *http.Request) {
		numberOfTimesSecondCalled++
		rw.Header().Set("Content-Type", "application/json")

		if numberOfTimesSecondCalled == 3 {
			rw.WriteHeader(200)
		} else {
			rw.WriteHeader(supertokens.RateLimitStatusCode)
		}

		response, err := json.Marshal(map[string]interface{}{})
		if err != nil {
			t.Error(err.Error())
		}
		rw.Write(response)
	})

	mux.HandleFunc("/testing3", func(rw http.ResponseWriter, r *http.Request) {
		numberOfTimesThirdCalled++
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(200)
		response, err := json.Marshal(map[string]interface{}{})
		if err != nil {
			t.Error(err.Error())
		}
		rw.Write(response)
	})

	testServer := httptest.NewServer(mux)

	defer func() {
		testServer.Close()
	}()

	config := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			// We need the querier to call the test server and not the core
			ConnectionURI: testServer.URL,
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}

	err := supertokens.Init(config)

	if err != nil {
		t.Error(err.Error())
	}

	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	supertokens.SetQuerierApiVersionForTests("3.0")
	defer resetQuerier()

	if err != nil {
		t.Error(err.Error())
	}

	_, err = q.SendGetRequest("/testing", map[string]string{}, nil)
	if err == nil {
		t.Error(errors.New("request should have failed but didnt").Error())
	} else {
		if !strings.Contains(err.Error(), "with status code: 429") {
			t.Error(errors.New("request failed with an unexpected error").Error())
		}
	}

	_, err = q.SendGetRequest("/testing2", map[string]string{}, nil)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = q.SendGetRequest("/testing3", map[string]string{}, nil)
	if err != nil {
		t.Error(err.Error())
	}

	// One initial call + 5 retries
	assert.Equal(t, numberOfTimesCalled, 6)
	assert.Equal(t, numberOfTimesSecondCalled, 3)
	assert.Equal(t, numberOfTimesThirdCalled, 1)
}

func TestThatRateLimitErrorsAreThrownBackToTheUser(t *testing.T) {
	resetAll()
	mux := http.NewServeMux()

	mux.HandleFunc("/testing", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(supertokens.RateLimitStatusCode)
		rw.Header().Set("Content-Type", "application/json")
		response, err := json.Marshal(map[string]interface{}{
			"status": "RATE_LIMIT_ERROR",
		})
		if err != nil {
			t.Error(err.Error())
		}
		rw.Write(response)
	})

	testServer := httptest.NewServer(mux)

	defer func() {
		testServer.Close()
	}()

	config := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			// We need the querier to call the test server and not the core
			ConnectionURI: testServer.URL,
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}

	err := supertokens.Init(config)

	if err != nil {
		t.Error(err.Error())
	}

	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	supertokens.SetQuerierApiVersionForTests("3.0")
	defer resetQuerier()

	if err != nil {
		t.Error(err.Error())
	}

	_, err = q.SendGetRequest("/testing", map[string]string{}, nil)
	if err == nil {
		t.Error(errors.New("request should have failed but didnt").Error())
	} else {
		if !strings.Contains(err.Error(), "with status code: 429") {
			t.Error(errors.New("request failed with an unexpected error" + err.Error()).Error())
		}

		assert.True(t, strings.Contains(err.Error(), "message: {\"status\":\"RATE_LIMIT_ERROR\"}"))
	}
}

func TestThatParallelCallsHaveIndependentRetryCounters(t *testing.T) {
	resetAll()
	mux := http.NewServeMux()

	numberOfTimesFirstCalled := 0
	numberOfTimesSecondCalled := 0

	mux.HandleFunc("/testing", func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("id") == "1" {
			numberOfTimesFirstCalled++
		} else {
			numberOfTimesSecondCalled++
		}

		rw.WriteHeader(supertokens.RateLimitStatusCode)
		rw.Header().Set("Content-Type", "application/json")
		response, err := json.Marshal(map[string]interface{}{})
		if err != nil {
			t.Error(err.Error())
		}
		rw.Write(response)
	})

	testServer := httptest.NewServer(mux)

	defer func() {
		testServer.Close()
	}()

	config := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			// We need the querier to call the test server and not the core
			ConnectionURI: testServer.URL,
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}

	err := supertokens.Init(config)

	if err != nil {
		t.Error(err.Error())
	}

	q, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	supertokens.SetQuerierApiVersionForTests("3.0")
	defer resetQuerier()

	if err != nil {
		t.Error(err.Error())
	}

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		_, err = q.SendGetRequest("/testing", map[string]string{
			"id": "1",
		}, nil)
		if err == nil {
			t.Error(errors.New("request should have failed but didnt").Error())
		} else {
			if !strings.Contains(err.Error(), "with status code: 429") {
				t.Error(errors.New("request failed with an unexpected error" + err.Error()).Error())
			}
		}

		wg.Done()
	}()

	go func() {
		_, err = q.SendGetRequest("/testing", map[string]string{
			"id": "2",
		}, nil)
		if err == nil {
			t.Error(errors.New("request should have failed but didnt").Error())
		} else {
			if !strings.Contains(err.Error(), "with status code: 429") {
				t.Error(errors.New("request failed with an unexpected error" + err.Error()).Error())
			}
		}

		wg.Done()
	}()

	wg.Wait()

	assert.Equal(t, numberOfTimesFirstCalled, 6)
	assert.Equal(t, numberOfTimesSecondCalled, 6)
}
