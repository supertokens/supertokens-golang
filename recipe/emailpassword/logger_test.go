package emailpassword

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestLogMessageRequestIDInHeader(t *testing.T) {
	os.Setenv("SUPERTOKENS_DEBUG", "true")
	var logMessage = "test log message"
	var validRequestID = "valid-request-ID"
	var buf bytes.Buffer
	customLogger := &mockLogger{
		log.New(&buf, "", 0),
	}
	var user *epmodels.User
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					Functions: func(originalImplementation epmodels.RecipeInterface) epmodels.RecipeInterface {
						originalSignIn := *originalImplementation.SignIn
						*originalImplementation.SignIn = func(email, password string, userContext supertokens.UserContext) (epmodels.SignInResponse, error) {
							supertokens.LogNewDebugMessage(userContext, logMessage)
							res, err := originalSignIn(email, password, userContext)
							return res, err
						}
						return originalImplementation
					},
				},
			}),
		},
		Log: customLogger,
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/user", func(rw http.ResponseWriter, r *http.Request) {
		fetchedUser := map[string]interface{}{}
		jsonResp, err := json.Marshal(fetchedUser)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	})
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	user = nil
	assert.Nil(t, user)

	_, _ = unittesting.SignInRequestWithRequestIDHeader("random@gmail.com", "validpass123", testServer.URL, validRequestID)

	assert.NoError(t, err)

	assert.Containsf(t, buf.String(), validRequestID, "Checking request id in log")
	assert.Containsf(t, buf.String(), supertokens.RequestIDKey, "Checking request id key in log")
	assert.Containsf(t, buf.String(), logMessage, "Checking log message in log")

	user = nil
	buf.Reset()
	assert.Nil(t, user)

	_, _ = unittesting.SignInRequestWithRequestIDHeader("random@gmail.com", "validpass123", testServer.URL, "new-request-id")

	assert.NoError(t, err)

	assert.Containsf(t, buf.String(), "new-request-id", "Checking request id in log")
	assert.Containsf(t, buf.String(), supertokens.RequestIDKey, "Checking request id key in log")
	assert.Containsf(t, buf.String(), logMessage, "Checking for the log message in log")

	os.Clearenv()
}

func TestLogMessageRequestIDInContext(t *testing.T) {
	os.Setenv("SUPERTOKENS_DEBUG", "true")
	var logMessage = "test log message"
	var buf bytes.Buffer
	customLogger := &mockLogger{
		log.New(&buf, "", 0),
	}
	var user *epmodels.User
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					APIs: func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
						originalSignInPOST := *originalImplementation.SignInPOST
						*originalImplementation.SignInPOST = func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignInPOSTResponse, error) {
							res, err := originalSignInPOST(formFields, options, userContext)
							supertokens.LogNewDebugMessage(userContext, logMessage)
							return res, err
						}
						return originalImplementation
					},
				},
			}),
		},
		Log:          customLogger,
		RequestIDKey: "rID",
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/user", func(rw http.ResponseWriter, r *http.Request) {
		fetchedUser := map[string]interface{}{}
		jsonResp, err := json.Marshal(fetchedUser)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	})
	testServer := httptest.NewServer(loggingMiddleware(supertokens.Middleware(mux)))
	defer func() {
		testServer.Close()
	}()

	user = nil
	assert.Nil(t, user)

	_, _ = unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)

	assert.NoError(t, err)

	assert.Containsf(t, buf.String(), strconv.FormatUint(rid, 10), "Checking request id in log")
	assert.Containsf(t, buf.String(), supertokens.RequestIDKey, "Checking request id key in log")
	assert.Containsf(t, buf.String(), logMessage, "Checking log message in log")

	user = nil
	buf.Reset()
	assert.Nil(t, user)

	_, _ = unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)

	assert.NoError(t, err)

	assert.Containsf(t, buf.String(), strconv.FormatUint(rid, 10), "Checking request id in log")
	assert.Containsf(t, buf.String(), supertokens.RequestIDKey, "Checking request id key in log")
	assert.Containsf(t, buf.String(), logMessage, "Checking for the log message in log")

	os.Clearenv()
}

func TestLogMessageWithoutCustomLogger(t *testing.T) {
	os.Setenv("SUPERTOKENS_DEBUG", "true")
	var logMessage = "test log message"
	var buf bytes.Buffer
	customLogger := &mockLogger{
		log.New(&buf, "", 0),
	}
	customLogger.Log("")
	var user *epmodels.User
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					APIs: func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
						originalSignInPOST := *originalImplementation.SignInPOST
						*originalImplementation.SignInPOST = func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignInPOSTResponse, error) {
							res, err := originalSignInPOST(formFields, options, userContext)
							supertokens.LogNewDebugMessage(userContext, logMessage)
							return res, err
						}
						return originalImplementation
					},
				},
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/user", func(rw http.ResponseWriter, r *http.Request) {
		fetchedUser := map[string]interface{}{}
		jsonResp, err := json.Marshal(fetchedUser)
		if err != nil {
			t.Errorf("Error happened in JSON marshal. Err: %s", err)
		}
		rw.WriteHeader(200)
		rw.Write(jsonResp)
	})
	testServer := httptest.NewServer(loggingMiddleware(supertokens.Middleware(mux)))
	defer func() {
		testServer.Close()
	}()

	user = nil
	assert.Nil(t, user)

	_, _ = unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)

	assert.NoError(t, err)

	assert.NotContainsf(t, buf.String(), strconv.FormatUint(rid, 10), "Checking request id in log")
	assert.NotContainsf(t, buf.String(), supertokens.RequestIDKey, "Checking request id key in log")
	assert.NotContainsf(t, buf.String(), logMessage, "Checking log message in log")

	user = nil
	buf.Reset()
	assert.Nil(t, user)

	_, _ = unittesting.SignInRequest("random@gmail.com", "validpass123", testServer.URL)

	assert.NoError(t, err)

	assert.NotContainsf(t, buf.String(), strconv.FormatUint(rid, 10), "Checking request id in log")
	assert.NotContainsf(t, buf.String(), supertokens.RequestIDKey, "Checking request id key in log")
	assert.NotContainsf(t, buf.String(), logMessage, "Checking for the log message in log")

	os.Clearenv()
}

func TestSuperTokensInitWithCustomLogger(t *testing.T) {
	var buf bytes.Buffer
	customLogger := &mockLogger{
		log.New(&buf, "", 0),
	}
	apiBasePath := "/"
	websiteBasePath := "/"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "SuperTokens",
			APIDomain:       "api.supertokens.io",
			WebsiteDomain:   "supertokens.io",
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			Init(&epmodels.TypeInput{}),
		},
		Log: customLogger,
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")

	defer AfterEach()

	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	supertokensInstance, err := supertokens.GetInstanceOrThrowError()

	if err != nil {
		t.Error(err.Error())
	}

	logMsg := "successful log\n"
	supertokensInstance.Log.Log(logMsg)
	assert.Equal(t, logMsg, buf.String())
	assert.Equal(t, true, reflect.DeepEqual(customLogger, supertokensInstance.Log))
}

func TestSuperTokensInitWithRequestIDKey(t *testing.T) {
	var buf bytes.Buffer
	customLogger := &mockLogger{
		log.New(&buf, "", 0),
	}
	apiBasePath := "/"
	websiteBasePath := "/"
	ridKey := "requestID"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "SuperTokens",
			APIDomain:       "api.supertokens.io",
			WebsiteDomain:   "supertokens.io",
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			Init(&epmodels.TypeInput{}),
		},
		Log:          customLogger,
		RequestIDKey: ridKey,
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")

	defer AfterEach()

	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	supertokensInstance, err := supertokens.GetInstanceOrThrowError()

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "requestID", supertokensInstance.RequestIDKey)
}

type mockLogger struct {
	*log.Logger
}

func (ml *mockLogger) Log(msg string) {
	ml.Print(msg)
}

var rid uint64

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "rID", strconv.FormatUint(atomic.AddUint64(&rid, 1), 10))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
