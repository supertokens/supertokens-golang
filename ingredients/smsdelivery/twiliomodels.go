package smsdelivery

import (
	"errors"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type TwilioServiceConfig struct {
	AccountSid          string
	AuthToken           string
	From                *string
	MessagingServiceSid *string
}

type TwilioGetContentResult struct {
	Body          string
	ToPhoneNumber string
}

type TwilioServiceInterface struct {
	SendRawSms *func(input TwilioGetContentResult, userContext supertokens.UserContext) error
	GetContent *func(input SmsType, userContext supertokens.UserContext) (TwilioGetContentResult, error)
}

type TwilioTypeInput struct {
	TwilioSettings TwilioServiceConfig
	Override       func(originalImplementation TwilioServiceInterface) TwilioServiceInterface
}

func NormaliseTwilioTypeInput(input TwilioTypeInput) (TwilioTypeInput, error) {
	if input.TwilioSettings.From == nil && input.TwilioSettings.MessagingServiceSid == nil {
		return TwilioTypeInput{}, errors.New("either 'From' or 'MessagingServiceSid' must be set")
	}
	if input.TwilioSettings.From != nil && input.TwilioSettings.MessagingServiceSid != nil {
		return TwilioTypeInput{}, errors.New("only one of 'From' or 'MessagingServiceSid' must be set")
	}
	return input, nil
}
