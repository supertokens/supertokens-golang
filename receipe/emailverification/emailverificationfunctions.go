package emailverification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/supertokens/supertokens-golang/supertokens"
)

func getEmailVerificationURL(appInfo supertokens.NormalisedAppinfo) string {
	return appInfo.WebsiteDomain.GetAsStringDangerous() + appInfo.WebsiteBasePath.GetAsStringDangerous() + "/verify-email"
}

func createAndSendCustomEmail(appInfo supertokens.NormalisedAppinfo) func(user User, emailVerifyURLWithToken string) {
	return func(user User, emailVerifyURLWithToken string) {
		url := "https://api.supertokens.io/0/st/auth/email/verify"
		data := map[string]interface{}{
			"email":          user.email,
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
