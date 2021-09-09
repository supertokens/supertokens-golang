package emailpassword

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func defaultGetResetPasswordURL(appInfo supertokens.NormalisedAppinfo) func(_ models.User) string {
	return func(_ models.User) string {
		return appInfo.WebsiteDomain.GetAsStringDangerous() + appInfo.WebsiteBasePath.GetAsStringDangerous() + "/reset-password"
	}
}

// TODO: add test to see query
func defaultCreateAndSendCustomPasswordResetEmail(appInfo supertokens.NormalisedAppinfo) func(user models.User, passwordResetURLWithToken string) {
	return func(user models.User, passwordResetURLWithToken string) {
		if flag.Lookup("test.v") != nil {
			// if running in test mode, we do not want to send this.
			return
		}
		url := "https://api.supertokens.io/0/st/auth/password/reset"
		data := map[string]string{
			"email":            user.Email,
			"appName":          appInfo.AppName,
			"passwordResetURL": passwordResetURLWithToken,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			return
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set("api-version", "0")

		client := &http.Client{}
		_, err = client.Do(req)
		if err != nil {
			return
		}
	}
}
