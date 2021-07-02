package emailpassword

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func defaultGetResetPasswordURL(appInfo supertokens.NormalisedAppinfo) func(_ models.User) string {
	return func(_ models.User) string {
		return appInfo.WebsiteDomain.GetAsStringDangerous() + appInfo.WebsiteBasePath.GetAsStringDangerous() + "/reset-password"
	}
}

func defaultCreateAndSendCustomPasswordResetEmail(appInfo supertokens.NormalisedAppinfo) func(user models.User, passwordResetURLWithToken string) error {
	return func(user models.User, passwordResetURLWithToken string) error {
		url := "https://api.supertokens.io/0/st/auth/password/reset"
		data := map[string]string{
			"email":            user.Email,
			"appName":          appInfo.AppName,
			"passwordResetURL": passwordResetURLWithToken,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		req.Header.Set("api-version", "0")
		client := &http.Client{}
		client.Do(req)
		return nil
	}
}
