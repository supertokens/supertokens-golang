package emailpassword

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TestAPIWithSupertokensMiddlewareButNotInitialized(t *testing.T) {
	BeforeEach()
	defer AfterEach()

	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, 500, resp.StatusCode)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	bodyStr := string(bodyBytes)
	assert.Equal(t, "initialisation not done. Did you forget to call the SuperTokens.init function?\n", bodyStr)
}
