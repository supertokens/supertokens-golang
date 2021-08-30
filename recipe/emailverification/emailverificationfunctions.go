package emailverification

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func DefaultGetEmailVerificationURL(appInfo supertokens.NormalisedAppinfo) func(models.User) (string, error) {
	return func(user models.User) (string, error) {
		return appInfo.WebsiteDomain.GetAsStringDangerous() + appInfo.WebsiteBasePath.GetAsStringDangerous() + "/verify-email", nil
	}
}

// TODO: send only if not testing
func DefaultCreateAndSendCustomEmail(appInfo supertokens.NormalisedAppinfo) func(user models.User, emailVerifyURLWithToken string) {
	return func(user models.User, emailVerifyURLWithToken string) {
		const url = "https://api.supertokens.io/0/st/auth/email/verify"

		data := map[string]string{
			"email":          user.Email,
			"appName":        appInfo.AppName,
			"emailVerifyURL": emailVerifyURLWithToken,
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
		return
	}
}
