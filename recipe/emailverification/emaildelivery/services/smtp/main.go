package smtp

import (
	"errors"
	"net/smtp"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery/emaildeliverymodels"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery/services/smtpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeSmtpService(config smtpmodels.TypeInput) emaildeliverymodels.EmailDeliveryInterface {
	// TODO: check sending email..
	smtpAuth := smtp.PlainAuth(config.SMTPSettings.From.Name, config.SMTPSettings.From.Email, config.SMTPSettings.Password, config.SMTPSettings.Host)

	serviceImpl := makeServiceImplementation(smtpAuth, config.SMTPSettings.Host, config.SMTPSettings.Port, config.SMTPSettings.From)
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
