package testing

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func startingHelper() {
	CleanST()
	SetUpST()
	StartUpST("localhost", "8080")
}

func endingHelper() {
	ResetAll()
	KillAllST()
}

func normalizeCookieSameSite(input string) string {
	str := strings.TrimSpace(input)
	val := strings.ToLower(str)
	return val
}

func TestSuperTokensInitWithAbsentOptionalFieldsInAppInfo(t *testing.T) {
	apiBasePath := "test/"
	websiteBasePath := "test1/"
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
		{
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
				session.Init(nil),
			},
		},
	}
	for _, configValue := range configValues {
		startingHelper()
		err := supertokens.Init(configValue)
		if err != nil {
			fmt.Println("Failed to get a supertokens instance")
		}
		supertokensInstance, err := supertokens.GetInstanceOrThrowError()

		if err != nil {
			fmt.Println("could not find a supertokens instance")
		}

		if configValue.AppInfo.APIBasePath != nil {
			assert.Equal(t, "/test", supertokensInstance.AppInfo.APIBasePath.GetAsStringDangerous())
			assert.Equal(t, "/test1", supertokensInstance.AppInfo.WebsiteBasePath.GetAsStringDangerous())
		} else {
			assert.Equal(t, "/auth", supertokensInstance.AppInfo.APIBasePath.GetAsStringDangerous())
			assert.Equal(t, "/auth", supertokensInstance.AppInfo.WebsiteBasePath.GetAsStringDangerous())
		}
		endingHelper()
	}
}

func TestSuperTokensInitWithAbsenceOfCompulsoryInputInAppInfo(t *testing.T) {
	configValues := []supertokens.TypeInput{
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				AppName:   "SuperTokens",
				APIDomain: "api.supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
	}
	for _, configValue := range configValues {
		startingHelper()
		err := supertokens.Init(configValue)
		if err != nil {
			errMessage := err.Error()
			if configValue.AppInfo.AppName != "SuperTokens" {
				assert.Equal(t, errMessage, "Please provide your appName inside the appInfo object when calling supertokens.init")
			} else if configValue.AppInfo.APIDomain != "api.supertokens.io" {
				assert.Equal(t, errMessage, "Please provide your apiDomain inside the appInfo object when calling supertokens.init")
			} else if configValue.AppInfo.WebsiteDomain != "supertokens.io" {
				assert.Equal(t, errMessage, "Please provide your websiteDomain inside the appInfo object when calling supertokens.init")
			}
		}
		endingHelper()
	}
}

func TestSuperTokensInitWithDifferentLengthOfRecipeModules(t *testing.T) {
	configValues := []supertokens.TypeInput{
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
				APIDomain:     "api.supertokens.io",
			},
			RecipeList: []supertokens.Recipe{},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
				APIDomain:     "api.supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
				APIDomain:     "api.supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
				emailpassword.Init(nil),
			},
		},
	}
	for _, configValue := range configValues {
		startingHelper()
		err := supertokens.Init(configValue)
		if err != nil {
			errorMessage := err.Error()
			if errorMessage != "please provide at least one recipe to the supertokens.init function call" {
				fmt.Println(errorMessage)
				log.Fatalf(err.Error())
			} else {
				assert.Equal(t, errorMessage, "please provide at least one recipe to the supertokens.init function call")
			}
			continue
		}
		supertokensInstance, err := supertokens.GetInstanceOrThrowError()

		if err != nil {
			fmt.Println("could not find a supertokens instance")
			log.Fatalf(err.Error())
		}

		assert.Equal(t, len(configValue.RecipeList), len(supertokensInstance.RecipeModules))

		endingHelper()
	}
}

func TestSuperTokensInitWithConfigForSessionModules(t *testing.T) {
	cookieDomain := "testDomain"
	sessionExpiredStatusCode := 111
	cookieSecure := true
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
			session.Init(
				&sessmodels.TypeInput{
					CookieDomain:             &cookieDomain,
					SessionExpiredStatusCode: &sessionExpiredStatusCode,
					CookieSecure:             &cookieSecure,
				},
			),
		},
	}
	startingHelper()
	err := supertokens.Init(configValue)
	if err != nil {
		log.Fatalf(err.Error())
	}

	sessionSingletonInstance, err := session.GetRecipeInstanceOrThrowError()
	if err != nil {
		log.Fatalf(err.Error())
	}
	assert.Equal(t, *sessionSingletonInstance.Config.CookieDomain, "testdomain")
	assert.Equal(t, sessionSingletonInstance.Config.SessionExpiredStatusCode, 111)
	assert.Equal(t, sessionSingletonInstance.Config.CookieSecure, true)
	endingHelper()
}

func TestSuperTokensInitWithConfigForSessionModulesWithVariousSameSiteValues(t *testing.T) {
	cookieSameSite0 := " Lax "
	cookieSameSite1 := "None "
	cookieSameSite2 := " STRICT "
	cookieSameSite3 := "random "
	cookieSameSite4 := " "
	cookieSameSite5 := "lax"
	cookieSameSite6 := "none"
	cookieSameSite7 := "strict"

	configValues := []supertokens.TypeInput{
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(&sessmodels.TypeInput{
					CookieSameSite: &cookieSameSite0,
				}),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(&sessmodels.TypeInput{
					CookieSameSite: &cookieSameSite1,
				}),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(&sessmodels.TypeInput{
					CookieSameSite: &cookieSameSite2,
				}),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(&sessmodels.TypeInput{
					CookieSameSite: &cookieSameSite3,
				}),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(&sessmodels.TypeInput{
					CookieSameSite: &cookieSameSite4,
				}),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(&sessmodels.TypeInput{
					CookieSameSite: &cookieSameSite5,
				}),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(&sessmodels.TypeInput{
					CookieSameSite: &cookieSameSite6,
				}),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(&sessmodels.TypeInput{
					CookieSameSite: &cookieSameSite7,
				}),
			},
		},
		{
			Supertokens: &supertokens.ConnectionInfo{
				ConnectionURI: "http://localhost:8080",
			},
			AppInfo: supertokens.AppInfo{
				APIDomain:     "api.supertokens.io",
				AppName:       "SuperTokens",
				WebsiteDomain: "supertokens.io",
			},
			RecipeList: []supertokens.Recipe{
				session.Init(nil),
			},
		},
	}

	for _, configValue := range configValues {
		startingHelper()
		err := supertokens.Init(configValue)
		if err != nil {
			errorMessage := err.Error()
			if errorMessage != `cookie same site must be one of "strict", "lax", or "none"` {
				log.Fatalf(err.Error())
			} else {
				assert.Equal(t, errorMessage, `cookie same site must be one of "strict", "lax", or "none"`)
				continue
			}
		}
		sessionSingletonInstance, err := session.GetRecipeInstanceOrThrowError()
		if err != nil {
			log.Fatalf(err.Error())
		}
		assert.Contains(t, []string{"lax", "strict", "none"}, sessionSingletonInstance.Config.CookieSameSite)
		endingHelper()
	}
}
