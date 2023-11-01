package session

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

// Added the logger tests here because supertokens/logger_test.go causes cyclic import errors due to imports in test/unittesting/testingUtils.go

func resetLogger() {
	supertokens.Logger = log.New(os.Stdout, "com.supertokens", 0)
	os.Unsetenv("SUPERTOKENS_DEBUG")
}

func TestLogDebugMessageWhenDebugTrue(t *testing.T) {
	var logMessage = "test log message"
	var buf bytes.Buffer

	supertokens.Logger = log.New(&buf, "test", 0)

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
		Debug: true,
	}
	BeforeEach()

	unittesting.StartUpST("localhost", "8080")

	defer AfterEach()
	defer resetLogger()

	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	supertokens.LogDebugMessage(logMessage)
	assert.Contains(t, buf.String(), logMessage, "checking log message in logs")
}

func TestLogDebugMessageWhenDebugFalse(t *testing.T) {
	var logMessage = "test log message"
	var buf bytes.Buffer

	supertokens.Logger = log.New(&buf, "test", 0)

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
		Debug: false,
	}
	BeforeEach()

	unittesting.StartUpST("localhost", "8080")

	defer AfterEach()
	defer resetLogger()

	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	supertokens.LogDebugMessage(logMessage)
	assert.NotContains(t, buf.String(), logMessage, "checking log message in logs")
}

func TestLogDebugMessageWhenDebugNotSet(t *testing.T) {
	var logMessage = "test log message"
	var buf bytes.Buffer

	supertokens.Logger = log.New(&buf, "test", 0)

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}
	BeforeEach()

	unittesting.StartUpST("localhost", "8080")

	defer AfterEach()
	defer resetLogger()

	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	supertokens.LogDebugMessage(logMessage)
	assert.NotContains(t, buf.String(), logMessage, "checking log message in logs")
}

func TestLogDebugMessageWithEnvVar(t *testing.T) {
	var logMessage = "test log message"
	var buf bytes.Buffer

	supertokens.Logger = log.New(&buf, "test", 0)
	os.Setenv("SUPERTOKENS_DEBUG", "1")

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "api.supertokens.io",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(nil),
		},
	}
	BeforeEach()

	unittesting.StartUpST("localhost", "8080")

	defer AfterEach()
	defer resetLogger()

	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	supertokens.LogDebugMessage(logMessage)
	assert.Contains(t, buf.String(), logMessage, "checking log message in logs")
}
