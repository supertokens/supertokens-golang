package testing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TestSuperTokensInit(t *testing.T) {
	configValues := []supertokens.TypeInput{
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				AppName:       "SuperTokens",
				APIDomain:     "api.supertokens.io",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
	}
	for _, configValue := range configValues {
		CleanST()
		// KillAllST()
		SetUpST()
		StartUpST("localhost", "8080")
		err := supertokens.Init(configValue)
		if err != nil {
			fmt.Println("Failed to get a supertokens instance")
		}
		supertokensInstance, err := supertokens.GetInstanceOrThrowError()

		if err != nil {
			fmt.Println("could not find a supertokens instance")
		}

		assert.Equal(t, "/auth", supertokensInstance.AppInfo.APIBasePath.GetAsStringDangerous())
		ResetAll()
		KillAllST()
	}
}
