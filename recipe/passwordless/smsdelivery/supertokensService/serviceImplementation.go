package supertokensService

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const SUPERTOKENS_SMS_SERVICE_URL = "https://api.supertokens.com/0/services/sms"

func MakeServiceImplementation(config smsdelivery.SupertokensServiceConfig) smsdelivery.SupertokensService {
	sendSms := func(input smsdelivery.PasswordlessLoginType, userContext supertokens.UserContext) error {
		instance, err := supertokens.GetInstanceOrThrowError()
		if err != nil {
			return err
		}

		data := map[string]map[string]interface{}{
			"smsInput": {
				"type":         "PASSWORDLESS_LOGIN",
				"phoneNumber":  input.PhoneNumber,
				"codeLifetime": input.CodeLifetime,
				"appName":      instance.AppInfo.AppName,
			},
		}
		if input.UrlWithLinkCode != nil {
			data["smsInput"]["urlWithLinkCode"] = *input.UrlWithLinkCode
		}
		if input.UserInputCode != nil {
			data["smsInput"]["userInputCode"] = *input.UserInputCode
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		req, err := http.NewRequest("POST", SUPERTOKENS_SMS_SERVICE_URL, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		req.Header.Set("content-type", "application/json")
		req.Header.Set("api-version", "0")
		client := &http.Client{}
		_, err = client.Do(req)
		return err
	}

	return smsdelivery.SupertokensService{
		SendSms: &sendSms,
	}
}
