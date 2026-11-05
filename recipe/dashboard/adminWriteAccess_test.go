package dashboard

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

const dashboardTestUserEmail = "dashboard-user@supertokens.com"
const dashboardTestUserPassword = "password1234"

// setUpDashboardWithAdmins initialises the dashboard (without an api key) with the given
// admins list, creates a dashboard user and returns the server, a dashboard session id and
// the id of an email password user that write operations can target.
func setUpDashboardWithAdmins(t *testing.T, admins *[]string) (*httptest.Server, string, string) {
	connectionURI := unittesting.StartUpST("localhost", "8080")
	config := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: connectionURI,
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			Init(&dashboardmodels.TypeInput{
				Admins: admins,
			}),
		},
	}
	err := supertokens.Init(config)
	assert.NoError(t, err)

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("dashboard")
	assert.NoError(t, err)
	createResp, err := querier.SendPostRequest("/recipe/dashboard/user", map[string]interface{}{
		"email":    dashboardTestUserEmail,
		"password": dashboardTestUserPassword,
	}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "OK", createResp["status"])

	testServer := httptest.NewServer(supertokens.Middleware(http.NewServeMux()))

	signInRes, err := http.Post(testServer.URL+"/auth/dashboard/api/signin", "application/json", strings.NewReader(`{"email": "`+dashboardTestUserEmail+`", "password": "`+dashboardTestUserPassword+`"}`))
	assert.NoError(t, err)
	signInBody := map[string]interface{}{}
	assert.NoError(t, json.NewDecoder(signInRes.Body).Decode(&signInBody))
	signInRes.Body.Close()
	assert.Equal(t, "OK", signInBody["status"])
	sessionId, _ := signInBody["sessionId"].(string)
	assert.NotEmpty(t, sessionId)

	signUpResponse, err := emailpassword.SignUp("public", "testing@supertokens.com", "abcd1234")
	assert.NoError(t, err)
	assert.NotNil(t, signUpResponse.OK)

	return testServer, sessionId, signUpResponse.OK.User.ID
}

func sendDashboardRequest(t *testing.T, method string, url string, sessionId string, body string) (int, string) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+sessionId)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	return res.StatusCode, string(resBody)
}

func passwordUpdateBody(userId string) string {
	return `{"userId": "` + userId + `", "newPassword": "newabcd1234"}`
}

func TestNonAdminCannotWriteWithExemptPathInQueryString(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer, sessionId, userId := setUpDashboardWithAdmins(t, &[]string{"admin@supertokens.com"})
	defer testServer.Close()

	for _, query := range []string{"?x=/api/signout", "?x=/api/analytics"} {
		status, body := sendDashboardRequest(t, http.MethodPut, testServer.URL+"/auth/dashboard/api/user/password"+query, sessionId, passwordUpdateBody(userId))
		assert.Equal(t, http.StatusForbidden, status, query)
		assert.Contains(t, body, "You are not permitted to perform this operation", query)
	}

	signInResponse, err := emailpassword.SignIn("public", "testing@supertokens.com", "abcd1234")
	assert.NoError(t, err)
	assert.NotNil(t, signInResponse.OK)
}

func TestReadOnlyUserCannotWriteWithExemptPathInQueryString(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer, sessionId, userId := setUpDashboardWithAdmins(t, &[]string{})
	defer testServer.Close()

	for _, query := range []string{"?x=/api/signout", "?x=/api/analytics"} {
		status, body := sendDashboardRequest(t, http.MethodPut, testServer.URL+"/auth/dashboard/api/user/password"+query, sessionId, passwordUpdateBody(userId))
		assert.Equal(t, http.StatusForbidden, status, query)
		assert.Contains(t, body, "You are not permitted to perform this operation", query)
	}
}

func TestNonAdminCanPostAnalyticsAndSignOut(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer, sessionId, _ := setUpDashboardWithAdmins(t, &[]string{"admin@supertokens.com"})
	defer testServer.Close()

	status, _ := sendDashboardRequest(t, http.MethodPost, testServer.URL+"/auth/dashboard/api/analytics", sessionId, `{"email": "`+dashboardTestUserEmail+`", "dashboardVersion": "0.7.0"}`)
	assert.Equal(t, http.StatusOK, status)

	status, _ = sendDashboardRequest(t, http.MethodPost, testServer.URL+"/auth/dashboard/api/signout", sessionId, "")
	assert.Equal(t, http.StatusOK, status)
}

func TestReadOnlyUserCanPostAnalyticsAndSignOut(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer, sessionId, _ := setUpDashboardWithAdmins(t, &[]string{})
	defer testServer.Close()

	status, _ := sendDashboardRequest(t, http.MethodPost, testServer.URL+"/auth/dashboard/api/analytics", sessionId, `{"email": "`+dashboardTestUserEmail+`", "dashboardVersion": "0.7.0"}`)
	assert.Equal(t, http.StatusOK, status)

	status, _ = sendDashboardRequest(t, http.MethodPost, testServer.URL+"/auth/dashboard/api/signout", sessionId, "")
	assert.Equal(t, http.StatusOK, status)
}

func TestAdminCanWrite(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer, sessionId, userId := setUpDashboardWithAdmins(t, &[]string{dashboardTestUserEmail})
	defer testServer.Close()

	status, _ := sendDashboardRequest(t, http.MethodPut, testServer.URL+"/auth/dashboard/api/user/password", sessionId, passwordUpdateBody(userId))
	assert.Equal(t, http.StatusOK, status)

	signInResponse, err := emailpassword.SignIn("public", "testing@supertokens.com", "newabcd1234")
	assert.NoError(t, err)
	assert.NotNil(t, signInResponse.OK)
}

func TestEveryDashboardUserCanWriteWhenAdminsIsNotSet(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	testServer, sessionId, userId := setUpDashboardWithAdmins(t, nil)
	defer testServer.Close()

	status, _ := sendDashboardRequest(t, http.MethodPut, testServer.URL+"/auth/dashboard/api/user/password", sessionId, passwordUpdateBody(userId))
	assert.Equal(t, http.StatusOK, status)

	signInResponse, err := emailpassword.SignIn("public", "testing@supertokens.com", "newabcd1234")
	assert.NoError(t, err)
	assert.NotNil(t, signInResponse.OK)
}
