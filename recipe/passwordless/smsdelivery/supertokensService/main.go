package supertokensService

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const SUPERTOKENS_SMS_SERVICE_URL = "https://api.supertokens.com/0/services/sms"

func MakeSupertokensService(config smsdelivery.SupertokensServiceConfig) smsdelivery.SmsDeliveryInterface {
	sendPasswordlessLoginSms := func(input smsdelivery.PasswordlessLoginType, userContext supertokens.UserContext) error {
		instance, err := supertokens.GetInstanceOrThrowError()
		if err != nil {
			return err
		}

		data := map[string]interface{}{
			"apiKey": config.ApiKey,
			"smsInput": map[string]interface{}{
				"type":         "PASSWORDLESS_LOGIN",
				"phoneNumber":  input.PhoneNumber,
				"codeLifetime": input.CodeLifetime,
				"appName":      instance.AppInfo.AppName,
			},
		}
		if input.UrlWithLinkCode != nil {
			data["smsInput"].(map[string]interface{})["urlWithLinkCode"] = *input.UrlWithLinkCode
		}
		if input.UserInputCode != nil {
			data["smsInput"].(map[string]interface{})["userInputCode"] = *input.UserInputCode
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
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 300 {
			return errors.New(fmt.Sprintf("Could not send SMS. The API returned %d status.", resp.StatusCode))
		}
		return nil
	}

	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		if input.PasswordlessLogin != nil {
			return sendPasswordlessLoginSms(*input.PasswordlessLogin, userContext)
		} else {
			return errors.New("should never come here")
		}
	}

	return smsdelivery.SmsDeliveryInterface{
		SendSms: &sendSms,
	}
}
