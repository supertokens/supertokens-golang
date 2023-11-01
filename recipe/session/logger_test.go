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

func TestLogDebugMessageWhenDebugTrue(t *testing.T) {
	var logMessage = "test log message"
	var buf bytes.Buffer

	debug := true
	supertokens.Logger = log.New(&buf, "", 0)

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
		Debug: &debug,
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

	supertokens.LogDebugMessage(logMessage)
	assert.Equal(t, &debug, supertokensInstance.Debug)
	assert.Contains(t, buf.String(), logMessage, "checking log message in logs")
}

func TestLogDebugMessageWhenDebugFalse(t *testing.T) {
	var logMessage = "test log message"
	var buf bytes.Buffer

	debug := false
	supertokens.Logger = log.New(&buf, "", 0)

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
		Debug: &debug,
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

	supertokens.LogDebugMessage(logMessage)
	assert.Equal(t, &debug, supertokensInstance.Debug)
	assert.NotContains(t, buf.String(), logMessage, "checking log message in logs")
}

func TestLogDebugMessageWithEnvVar(t *testing.T) {
	var logMessage = "test log message"
	var buf bytes.Buffer

	supertokens.Logger = log.New(&buf, "", 0)
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

	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	supertokens.LogDebugMessage(logMessage)
	assert.Contains(t, buf.String(), logMessage, "checking log message in logs")
}
