package smsdelivery

import "github.com/supertokens/supertokens-golang/supertokens"

type SmsDeliveryInterface struct {
	SendSms *func(input SmsType, userContext supertokens.UserContext) error
}

type TypeInput struct {
	Service  *SmsDeliveryInterface
	Override func(originalImplementation SmsDeliveryInterface) SmsDeliveryInterface
}

type TypeInputWithService struct {
	Service  SmsDeliveryInterface
	Override func(originalImplementation SmsDeliveryInterface) SmsDeliveryInterface
}

type SmsType struct {
	PasswordlessLogin *PasswordlessLoginType
}

type PasswordlessLoginType struct {
	PhoneNumber      string
	UserInputCode    *string
	UrlWithLinkCode  *string
	CodeLifetime     uint64
	PreAuthSessionId string
}

type User struct {
	ID          string
	PhoneNumber string
}
