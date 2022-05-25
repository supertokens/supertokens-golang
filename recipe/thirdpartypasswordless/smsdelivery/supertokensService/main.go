package supertokensService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/smsdelivery/supertokensService"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeSupertokensService(config smsdelivery.SupertokensServiceConfig) smsdelivery.SmsDeliveryInterface {
	plessServiceImpl := supertokensService.MakeServiceImplementation(config)

	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		if input.PasswordlessLogin != nil {
			return (*plessServiceImpl.SendSms)(*input.PasswordlessLogin, userContext)
		} else {
			return errors.New("should never come here")
		}
	}

	return smsdelivery.SmsDeliveryInterface{
		SendSms: &sendSms,
	}
}
