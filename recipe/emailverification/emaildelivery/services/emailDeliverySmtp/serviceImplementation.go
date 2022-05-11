package emailDeliverySmtp

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery/emaildeliverymodels"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery/services/smtpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeServiceImplementation(config smtpmodels.SMTPServiceConfig) smtpmodels.ServiceInterface {
	sendRawEmail := func(input smtpmodels.GetContentResult, userContext supertokens.UserContext) error {
		return emaildelivery.SendSMTPEmail(config, input)
	}

	getContent := func(input emaildeliverymodels.EmailType, userContext supertokens.UserContext) (smtpmodels.GetContentResult, error) {
		if input.EmailVerification != nil {
			return getEmailVerifyEmailContent(*input.EmailVerification)
		} else {
			return smtpmodels.GetContentResult{}, errors.New("should never come here")
		}
	}

	return smtpmodels.ServiceInterface{
		SendRawEmail: &sendRawEmail,
		GetContent:   &getContent,
	}
}
