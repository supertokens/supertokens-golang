package emailpassword

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestLoginMethodsReturnEmptyProviderArray(t *testing.T) {
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
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/auth/loginmethods")
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respObj := map[string]interface{}{}
	respBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(respBytes, &respObj)
	assert.Nil(t, err)
	providersArr, ok := respObj["thirdParty"].(map[string]interface{})["providers"]
	assert.True(t, ok)
	assert.NotNil(t, providersArr)
}
