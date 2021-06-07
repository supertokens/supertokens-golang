package emailverification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func GetEmailVerificationURL(appInfo supertokens.NormalisedAppinfo) func(schema.User) string {
	return func(userId schema.User) string {
		return appInfo.WebsiteDomain.GetAsStringDangerous() + appInfo.WebsiteBasePath.GetAsStringDangerous() + "/verify-email"
	}
}

func CreateAndSendCustomEmail(appInfo supertokens.NormalisedAppinfo) func(user schema.User, emailVerifyURLWithToken string) {
	return func(user schema.User, emailVerifyURLWithToken string) {
		const url = "https://api.supertokens.io/0/st/auth/email/verify"

		data := map[string]string{
			"email":          user.Email,
			"appName":        appInfo.AppName,
			"emailVerifyURL": emailVerifyURLWithToken,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err) // todo: handle error
			return
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println(err) // todo: handle error
			return
		}

		req.Header.Set("api-version", "0")
		client := &http.Client{}
		client.Do(req)
	}
}
