package smsdelivery

import "github.com/supertokens/supertokens-golang/supertokens"

type TwilioServiceConfig struct {
	AccountSid          string
	AuthToken           string
	From                *string
	MessagingServiceSid *string
}

type GetContentResult struct {
	Body          string
	ToPhoneNumber string
}

type TwilioServiceInterface struct {
	SendRawSms *func(input GetContentResult, userContext supertokens.UserContext) error
	GetContent *func(input SmsType, userContext supertokens.UserContext) (GetContentResult, error)
}

type TwilioTypeInput struct {
	TwilioSettings TwilioServiceConfig
	Override       func(originalImplementation TwilioServiceInterface) TwilioServiceInterface
}
