package dashboard

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

/*
- Initialise with thirdpartyemailpassword and provide no custom form fields
- Create an emailpassword user using the thirdpartyemailpassword recipe
- Try to change the password of the user
- Should result in no errors
- Sign in with new password
- Should result in no errors and same user should be returned
*/
func TestThatUpdatingPasswordWithNoSignUpFeatureInTPEPWorks(t *testing.T) {
	config := supertokens.TypeInput{
		OnSuperTokensAPIError: func(err error, req *http.Request, res http.ResponseWriter) {
			print(err)
		},
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			thirdpartyemailpassword.Init(nil),
			Init(&dashboardmodels.TypeInput{
				ApiKey: "testapikey",
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(config)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	signupResponse, err := thirdpartyemailpassword.EmailPasswordSignUp("public", "testing@supertokens.com", "abcd1234")
	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, signupResponse.OK)

	userId := signupResponse.OK.User.ID

	req, err := http.NewRequest(http.MethodPut, testServer.URL+"/auth/dashboard/api/user/password", strings.NewReader(`{"userId": "`+userId+`", "newPassword": "newabcd1234"}`))

	if err != nil {
		t.Error(err.Error())
	}

	req.Header.Set("Authorization", "Bearer testapikey")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)

	signInResponse, err := thirdpartyemailpassword.EmailPasswordSignIn("public", "testing@supertokens.com", "newabcd1234")

	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, signInResponse.OK)
	assert.Equal(t, signInResponse.OK.User.ID, userId)
}
