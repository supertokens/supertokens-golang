package session

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestThatThereIsNoMemoryLeakWithJWKSCache(t *testing.T) {
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
			Init(&sessmodels.TypeInput{
				JWKSRefreshIntervalSec: &[]uint64{0}[0],
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

	testServer := GetTestServer(t)
	defer func() {
		testServer.Close()
	}()

	sess, err := CreateNewSessionWithoutRequestResponse("public", "testuser", map[string]interface{}{}, map[string]interface{}{}, nil)
	assert.NoError(t, err)

	accessToken := sess.GetAccessToken()

	_, err = GetSessionWithoutRequestResponse(accessToken, nil, nil)
	assert.NoError(t, err)

	numGoroutinesBeforeJWKSRefresh := runtime.NumGoroutine()

	for i := 0; i < 100; i++ {
		_, err = GetSessionWithoutRequestResponse(accessToken, nil, nil)
		assert.NoError(t, err)

		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(1 * time.Second)
	assert.Equal(t, numGoroutinesBeforeJWKSRefresh, runtime.NumGoroutine())
}
