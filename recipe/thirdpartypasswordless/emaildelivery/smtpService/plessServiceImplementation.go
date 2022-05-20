package smtpService

import (
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makePasswordlessServiceImplementation(serviceImpl emaildelivery.SMTPServiceInterface) func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
	return func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
		sendRawEmail := func(input emaildelivery.SMTPGetContentResult, userContext supertokens.UserContext) error {
			return (*serviceImpl.SendRawEmail)(input, userContext)
		}

		getContent := func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPGetContentResult, error) {
			return (*serviceImpl.GetContent)(input, userContext)
		}

		return emaildelivery.SMTPServiceInterface{
			SendRawEmail: &sendRawEmail,
			GetContent:   &getContent,
		}
	}
}
