package twilioService

import (
	"github.com/supertokens/supertokens-golang/ingredients/smsdelivery"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/smsdelivery/twilioService"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeTwilioService(config smsdelivery.TwilioTypeInput) smsdelivery.SmsDeliveryInterface {
	plessServiceImpl := twilioService.MakeTwilioService(config)

	sendSms := func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
		return (*plessServiceImpl.SendSms)(input, userContext)
	}

	return smsdelivery.SmsDeliveryInterface{
		SendSms: &sendSms,
	}
}
