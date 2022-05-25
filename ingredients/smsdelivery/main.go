package smsdelivery

import (
	"errors"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Ingredient struct {
	IngredientInterfaceImpl SmsDeliveryInterface
}

func MakeIngredient(config TypeInputWithService) Ingredient {

	result := Ingredient{
		IngredientInterfaceImpl: config.Service,
	}

	if config.Override != nil {
		result.IngredientInterfaceImpl = config.Override(result.IngredientInterfaceImpl)
	}

	return result
}

func SendTwilioSms(config TwilioServiceConfig, content GetContentResult) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.AccountSid,
		Password: config.AuthToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(content.ToPhoneNumber)
	params.SetBody(content.Body)

	if config.From != nil {
		params.SetFrom(*config.From)
	} else if config.MessagingServiceSid != nil {
		params.SetMessagingServiceSid(*config.MessagingServiceSid)
	} else {
		return errors.New("should not come here")
	}

	_, err := client.Api.CreateMessage(params)

	return err
}
