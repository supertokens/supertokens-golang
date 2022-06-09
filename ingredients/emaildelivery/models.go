package emaildelivery

import "github.com/supertokens/supertokens-golang/supertokens"

type EmailDeliveryInterface struct {
	SendEmail *func(input EmailType, userContext supertokens.UserContext) error
}

type TypeInput struct {
	Service  *EmailDeliveryInterface
	Override func(originalImplementation EmailDeliveryInterface) EmailDeliveryInterface
}

type TypeInputWithService struct {
	Service  EmailDeliveryInterface
	Override func(originalImplementation EmailDeliveryInterface) EmailDeliveryInterface
}

type EmailType struct {
	EmailVerification *EmailVerificationType
	PasswordReset     *PasswordResetType
	PasswordlessLogin *PasswordlessLoginType
}

type EmailVerificationType struct {
	User            User
	EmailVerifyLink string
}

type PasswordResetType struct {
	User              User
	PasswordResetLink string
}

type PasswordlessLoginType struct {
	Email            string
	UserInputCode    *string
	UrlWithLinkCode  *string
	CodeLifetime     uint64
	PreAuthSessionId string
}

type User struct {
	ID    string
	Email string
}
