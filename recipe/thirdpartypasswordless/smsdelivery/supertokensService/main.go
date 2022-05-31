package supertokensService

import (
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/smsdelivery/supertokensService"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeSupertokensService(config smsdelivery.SupertokensServiceConfig) smsdelivery.SmsDeliveryInterface {
	plessService := supertokensService.MakeSupertokensService(config)

	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		return (*plessService.SendSms)(input, userContext)
	}

	return smsdelivery.SmsDeliveryInterface{
		SendSms: &sendSms,
	}
}
