package emailpassword

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestEVGenerateUpdatesSessionClaims(t *testing.T) {
	antiCsrfConf := "VIA_TOKEN"
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
			Init(nil),
			emailverification.Init(evmodels.TypeInput{
				Mode:                     evmodels.ModeOptional,
				CreateAndSendCustomEmail: func(user evmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {},
			}),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &antiCsrfConf,
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
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
	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if unittesting.MaxVersion("2.10", cdiVersion) != cdiVersion {
		return
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	assert.NoError(t, err)

	assert.Equal(t, 200, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	bodyObj := map[string]interface{}{}
	err = json.Unmarshal(bodyBytes, &bodyObj)
	assert.NoError(t, err)

	userId := bodyObj["user"].(map[string]interface{})["id"].(string)
	infoFromResponse := unittesting.ExtractInfoFromResponse(resp)
	antiCsrf := infoFromResponse["antiCsrf"]

	token, err := emailverification.CreateEmailVerificationToken(userId, nil, nil)
	assert.NoError(t, err)
	_, err = emailverification.VerifyEmailUsingToken(token.OK.Token, nil)
	assert.NoError(t, err)

	resp, err = unittesting.EmailVerifyTokenRequest(
		testServer.URL,
		userId,
		infoFromResponse["sAccessToken"],
		antiCsrf,
	)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	infoFromResponse = unittesting.ExtractInfoFromResponse(resp)
	b64Bytes, err := base64.StdEncoding.DecodeString(infoFromResponse["frontToken"])
	assert.NoError(t, err)
	frontendInfo := map[string]interface{}{}
	err = json.Unmarshal(b64Bytes, &frontendInfo)
	assert.NoError(t, err)

	val := frontendInfo["up"].(map[string]interface{})["st-ev"].(map[string]interface{})["v"].(bool)
	fmt.Println(val)
	assert.True(t, frontendInfo["up"].(map[string]interface{})["st-ev"].(map[string]interface{})["v"].(bool))

	// Calling again should not modify access token
	resp, err = unittesting.EmailVerifyTokenRequest(
		testServer.URL,
		userId,
		infoFromResponse["sAccessToken"],
		antiCsrf,
	)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	infoFromResponse2 := unittesting.ExtractInfoFromResponse(resp)
	assert.Empty(t, infoFromResponse2["frontToken"])

	// now we mark the email as unverified and try again
	emailverification.UnverifyEmail(userId, nil, nil)
	resp, err = unittesting.EmailVerifyTokenRequest(
		testServer.URL,
		userId,
		infoFromResponse["sAccessToken"],
		antiCsrf,
	)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	infoFromResponse = unittesting.ExtractInfoFromResponse(resp)
	b64Bytes, err = base64.StdEncoding.DecodeString(infoFromResponse["frontToken"])
	assert.NoError(t, err)
	frontendInfo = map[string]interface{}{}
	err = json.Unmarshal(b64Bytes, &frontendInfo)
	assert.NoError(t, err)
	assert.False(t, frontendInfo["up"].(map[string]interface{})["st-ev"].(map[string]interface{})["v"].(bool))

	// calling the API again should not modify the access token again
	resp, err = unittesting.EmailVerifyTokenRequest(
		testServer.URL,
		userId,
		infoFromResponse["sAccessToken"],
		antiCsrf,
	)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	infoFromResponse2 = unittesting.ExtractInfoFromResponse(resp)
	assert.Empty(t, infoFromResponse2["frontToken"])
}
