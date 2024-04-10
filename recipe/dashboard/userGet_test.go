package dashboard

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/api"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/api/userdetails"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

/*
- Get user from the user get API
- Check that user has public tenant
*/
func TestThatUserGetReturnsTenantIDsCorrectly(t *testing.T) {
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
			emailpassword.Init(nil),
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

	signupResponse, err := emailpassword.SignUp("public", "testing@supertokens.com", "abcd1234")
	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, signupResponse.OK)

	userId := signupResponse.OK.User.ID

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/dashboard/api/user?userId="+userId+"&recipeId=emailpassword", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer testapikey")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err.Error())
	}

	var response userdetails.UserGetResponse
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &response)

	assert.True(t, len(response.User.TenantIds) > 0)
	assert.Equal(t, response.User.TenantIds[0], "public")
}

var customProvider1 = tpmodels.ProviderInput{
	Config: tpmodels.ProviderConfig{
		ThirdPartyId:          "custom",
		AuthorizationEndpoint: "https://test.com/oauth/auth",
		TokenEndpoint:         "https://test.com/oauth/token",

		Clients: []tpmodels.ProviderClientConfig{
			{
				ClientID: "supertokens",
			},
		},
	},

	Override: func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			return tpmodels.TypeUserInfo{
				ThirdPartyUserId: "user",
				Email: &tpmodels.EmailStruct{
					ID:         "email@test.com",
					IsVerified: true,
				},
			}, nil
		}
		return originalImplementation
	},
}

func TestThatUserGetReturnsValidUserForThirdPartyUserWhenUsingThirdPartyPasswordless(t *testing.T) {
	config := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.ProviderInput{
						customProvider1,
					},
				},
			}),
			passwordless.Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
				},
			}),
			Init(&dashboardmodels.TypeInput{
				ApiKey: "testapikey",
			}),
			session.Init(nil),
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

	_, err = unittesting.SigninupCustomRequest(testServer.URL, "test@gmail.com", "testPass0")

	if err != nil {
		t.Error(err.Error())
	}

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/auth/dashboard/api/users?limit=10", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer testapikey")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err.Error())
	}

	var listResponse api.UsersGetResponse
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &listResponse)

	user := listResponse.Users[0].User

	req, err = http.NewRequest(http.MethodGet, testServer.URL+"/auth/dashboard/api/user?userId="+user.Id+"&recipeId=thirdparty", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer testapikey")
	res, err = http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err.Error())
	}

	var response userdetails.UserGetResponse
	body, _ = io.ReadAll(res.Body)
	json.Unmarshal(body, &response)

	assert.Equal(t, response.Status, "OK")
}
