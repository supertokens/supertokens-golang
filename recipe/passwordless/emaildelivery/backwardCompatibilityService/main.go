package backwardCompatibilityService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeBackwardCompatibilityService(appInfo supertokens.NormalisedAppinfo, createAndSendCustomEmail func(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error) emaildelivery.EmailDeliveryInterface {
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.PasswordlessLogin != nil {
			return createAndSendCustomEmail(
				input.PasswordlessLogin.Email,
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

	return emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
