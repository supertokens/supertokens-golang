package backwardCompatibilityService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeBackwardCompatibilityService(createAndSendCustomSms func(phoneNumber string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error) smsdelivery.SmsDeliveryInterface {
	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		if input.PasswordlessLogin != nil {
			return createAndSendCustomSms(
				input.PasswordlessLogin.PhoneNumber,
				input.PasswordlessLogin.UserInputCode,
				input.PasswordlessLogin.UrlWithLinkCode,
				input.PasswordlessLogin.CodeLifetime,
				input.PasswordlessLogin.PreAuthSessionId,
				userContext,
			)
		} else {
			return errors.New("should never come here")
		}
	}

	return smsdelivery.SmsDeliveryInterface{
		SendSms: &sendSms,
	}
}
