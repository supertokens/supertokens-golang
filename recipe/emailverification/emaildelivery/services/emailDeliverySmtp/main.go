package emailDeliverySmtp

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery/emaildeliverymodels"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery/services/smtpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeSmtpService(config smtpmodels.TypeInput) emaildeliverymodels.EmailDeliveryInterface {
	serviceImpl := makeServiceImplementation(config.SMTPSettings)

	if config.Override != nil {
		serviceImpl = config.Override(serviceImpl)
	}

	sendEmail := func(input emaildeliverymodels.EmailType, userContext supertokens.UserContext) error {
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

	return emaildeliverymodels.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
