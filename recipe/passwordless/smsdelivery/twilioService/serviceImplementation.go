package twilioService

import (
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeServiceImplementation(config smsdelivery.TwilioServiceConfig) smsdelivery.TwilioServiceInterface {
	sendRawSms := func(input smsdelivery.GetContentResult, userContext supertokens.UserContext) error {
		return smsdelivery.SendTwilioSms(config, input)
	}

	getContent := func(input smsdelivery.SmsType, userContext supertokens.UserContext) (smsdelivery.GetContentResult, error) {
		result := getPasswordlessLoginSmsContent(*input.PasswordlessLogin)
		return result, nil
	}

	return smsdelivery.TwilioServiceInterface{
		SendRawSms: &sendRawSms,
		GetContent: &getContent,
	}
}
