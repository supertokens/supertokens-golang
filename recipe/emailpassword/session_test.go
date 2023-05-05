package emailpassword

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShouldMakeSignInUpReturn500WhenUsingProtectedProp(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			session.Init(&sessmodels.TypeInput{
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						originalCreateNewSession := *originalImplementation.CreateNewSession
						newCreateNewSession := func(userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, disableAntiCsrf *bool, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							accessTokenPayload["sub"] = "asdf"

							return originalCreateNewSession(userID, accessTokenPayload, sessionDataInDatabase, disableAntiCsrf, userContext)
						}

						*originalImplementation.CreateNewSession = newCreateNewSession

						return originalImplementation
					},
				},
			}),
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()

	passwordVal := "validPass123"

	emailVal := "random@email.com"

	formFields := map[string][]map[string]string{
		"formFields": {
			{
				"id":    "email",
				"value": emailVal,
			},
			{
				"id":    "password",
				"value": passwordVal,
			},
		},
	}

	postBody, err := json.Marshal(formFields)
	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))
	assert.Equal(t, 500, resp.StatusCode)
	cookies := unittesting.ExtractInfoFromResponse(resp)
	assert.True(t, cookies["accessTokenFromAny"] == "")
	assert.True(t, cookies["refreshTokenFromAny"] == "")
	assert.True(t, cookies["frontToken"] == "")
}

func checkResponse(t *testing.T, res *http.Response, exposed bool) {
	info := unittesting.ExtractInfoFromResponse(res)

	if exposed {
		assert.Equal(t, info["sAccessToken"], info["accessTokenFromHeader"])
	} else {
		assert.Equal(t, info["accessTokenFromHeader"], "")
		assert.NotEqual(t, info["sAccessToken"], "")
	}
}

func GetTestServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()
	checkDBTrue := true

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		dataInBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		var result map[string]interface{}
		err = json.Unmarshal(dataInBytes, &result)

		var payload map[string]interface{}

		if result["payload"] != nil {
			payload = result["payload"].(map[string]interface{})
		}

		_, err2 := session.CreateNewSession(r, rw, "uniqueId", payload, map[string]interface{}{})

		if err2 != nil {
			http.Error(rw, fmt.Sprint(err2), 400)
		}
	})

	mux.HandleFunc("/verify", verifySession2(true, &checkDBTrue, func(rw http.ResponseWriter, r *http.Request) {
		session := session.GetSessionFromRequestContext(r.Context())
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"message":       true,
			"sessionHandle": session.GetHandle(),
			"sessionExists": session != nil,
		})
	}))

	mux.HandleFunc("/merge-into-payload", verifySession2(true, nil, func(rw http.ResponseWriter, r *http.Request) {
		session := session.GetSessionFromRequestContext(r.Context())

		dataInBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err.Error())
		}
		var result map[string]interface{}
		err = json.Unmarshal(dataInBytes, &result)

		err = session.MergeIntoAccessTokenPayload(result["payload"].(map[string]interface{}))
		assert.NoError(t, err)

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"message":       true,
			"sessionHandle": session.GetHandle(),
			"sessionExists": session != nil,
			"newPayload":    session.GetAccessTokenPayload(),
		})
	}))

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	return testServer
}

func verifySession2(sessionRequired bool, checkDatabase *bool, otherHandler http.HandlerFunc) http.HandlerFunc {
	return session.VerifySession(&sessmodels.VerifySessionOptions{
		SessionRequired: &sessionRequired,
		CheckDatabase:   checkDatabase,
	}, otherHandler)
}
