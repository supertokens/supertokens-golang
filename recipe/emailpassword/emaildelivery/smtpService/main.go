package smtpService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	evsmtpService "github.com/supertokens/supertokens-golang/recipe/emailverification/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeSmtpService(config emaildelivery.SMTPTypeInput) emaildelivery.EmailDeliveryInterface {
	serviceImpl := MakeServiceImplementation(config.SMTPSettings)

	if config.Override != nil {
		serviceImpl = config.Override(serviceImpl)
	}

	emailVerificationServiceImpl := evsmtpService.MakeSmtpService(emaildelivery.SMTPTypeInput{
		SMTPSettings: config.SMTPSettings,
		Override:     makeEmailverificationServiceImplementation(serviceImpl),
	})

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.EmailVerification != nil {
			return (*emailVerificationServiceImpl.SendEmail)(input, userContext)

		} else if input.PasswordReset != nil {
			content, err := (*serviceImpl.GetContent)(input, userContext)
			if err != nil {
				return err
			}
			return (*serviceImpl.SendRawEmail)(content, userContext)

		} else {
			return errors.New("should never come here")
		}
	}

	return emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
