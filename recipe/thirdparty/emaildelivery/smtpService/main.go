package smtpService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	evsmtpService "github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeSmtpService(config emaildelivery.SMTPTypeInput) emaildelivery.EmailDeliveryInterface {
	emailVerificationServiceImpl := evsmtpService.MakeSmtpService(config)

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.EmailVerification != nil {
			return (*emailVerificationServiceImpl.SendEmail)(input, userContext)

		} else {
			return errors.New("should never come here")
		}
	}

	return emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
