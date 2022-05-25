package smsdelivery

import "github.com/supertokens/supertokens-golang/supertokens"

type SupertokensServiceConfig struct {
	ApiKey string

	Override func(originalImplementation SmsDeliveryInterface) SmsDeliveryInterface
}

type SupertokensService struct {
	SendSms *func(input PasswordlessLoginType, userContext supertokens.UserContext) error
}
