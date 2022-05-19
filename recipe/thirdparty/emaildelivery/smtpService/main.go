package smtpService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeSmtpService(config emaildelivery.SMTPTypeInput) emaildelivery.EmailDeliveryInterface {
	serviceImpl := makeServiceImplementation(config.SMTPSettings)

	if config.Override != nil {
		serviceImpl = config.Override(serviceImpl)
	}

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.EmailVerification != nil {
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
