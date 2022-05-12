package emaildelivery

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type SMTPServiceConfig struct {
	Host     string
	From     SMTPServiceFromConfig
	Port     int
	Password string
	Secure   *bool
}

type SMTPServiceFromConfig struct {
	Name  string
	Email string
}

type SMTPGetContentResult struct {
	Body    string
	IsHtml  bool
	Subject string
	ToEmail string
}

type SMTPServiceInterface struct {
	SendRawEmail *func(input SMTPGetContentResult, userContext supertokens.UserContext) error
	GetContent   *func(input EmailType, userContext supertokens.UserContext) (SMTPGetContentResult, error)
}

type SMTPTypeInput struct {
	SMTPSettings SMTPServiceConfig
	Override     func(originalImplementation SMTPServiceInterface) SMTPServiceInterface
}
