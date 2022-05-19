package smtpService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	epsmtpService "github.com/supertokens/supertokens-golang/recipe/emailpassword/emaildelivery/smtpService"
	evsmtpService "github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeSmtpService(config emaildelivery.SMTPTypeInput) emaildelivery.EmailDeliveryInterface {
	serviceImpl := makeServiceImplementation(config.SMTPSettings)

	if config.Override != nil {
		serviceImpl = config.Override(serviceImpl)
	}

	emailVerificationServiceImpl := evsmtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: config.SMTPSettings,
		Override:     makeEmailverificationServiceImplementation(serviceImpl),
	})
	emailPasswordServiceImpl := epsmtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: config.SMTPSettings,
		Override:     makeEmailpasswordServiceImplementation(serviceImpl),
	})

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.EmailVerification != nil {
			return (*emailVerificationServiceImpl.SendEmail)(input, userContext)

		} else if input.PasswordReset != nil {
			return (*emailPasswordServiceImpl.SendEmail)(input, userContext)

		} else {
			return errors.New("should never come here")
		}
	}

	return emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
