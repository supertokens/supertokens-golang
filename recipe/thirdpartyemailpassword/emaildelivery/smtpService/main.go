package smtpService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	epsmtpService "github.com/supertokens/supertokens-golang/recipe/emailpassword/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeSmtpService(config emaildelivery.SMTPTypeInput) emaildelivery.EmailDeliveryInterface {
	emailPasswordServiceImpl := epsmtpService.MakeSmtpService(config)

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.EmailVerification != nil || input.PasswordReset != nil {
			return (*emailPasswordServiceImpl.SendEmail)(input, userContext)

		} else {
			return errors.New("should never come here")
		}
	}

	return emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
