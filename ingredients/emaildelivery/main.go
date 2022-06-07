package emaildelivery

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

type Ingredient struct {
	IngredientInterfaceImpl EmailDeliveryInterface
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

func SendSMTPEmail(config SMTPServiceConfig, content SMTPGetContentResult) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", config.From.Name, config.From.Email))
	m.SetHeader("To", content.ToEmail)
	m.SetHeader("Subject", content.Subject)

	if content.IsHtml {
		m.SetBody("text/html", content.Body)
	} else {
		m.SetBody("text/plain", content.Body)
	}

	secure := false
	if config.Secure != nil {
		secure = *config.Secure
	}

	d := gomail.NewDialer(config.Host, config.Port, config.From.Email, config.Password)

	if secure {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true, ServerName: config.Host}
		d.SSL = true
	}
	return d.DialAndSend(m)
}
