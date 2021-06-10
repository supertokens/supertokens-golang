package emailverification

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func DefaultGetEmailVerificationURL(appInfo supertokens.NormalisedAppinfo) func(schema.User) string {
	return func(user schema.User) string {
		return appInfo.WebsiteDomain.GetAsStringDangerous() + appInfo.WebsiteBasePath.GetAsStringDangerous() + "/verify-email"
	}
}

func DefaultCreateAndSendCustomEmail(appInfo supertokens.NormalisedAppinfo) func(user schema.User, emailVerifyURLWithToken string) {
	return func(user schema.User, emailVerifyURLWithToken string) {
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

		req.Header.Set("api-version", "0")
		client := &http.Client{}
		client.Do(req)
	}
}
