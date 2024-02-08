/*
 * Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

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

func SendSMTPEmail(settings SMTPSettings, content EmailContent) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", settings.From.Name, settings.From.Email))
	m.SetHeader("To", content.ToEmail)
	m.SetHeader("Subject", content.Subject)

	if content.IsHtml {
		m.SetBody("text/html", content.Body)
	} else {
		m.SetBody("text/plain", content.Body)
	}

	username := settings.From.Email
	if settings.Username != nil {
		username = *settings.Username
	}

	d := gomail.NewDialer(settings.Host, settings.Port, username, settings.Password)

	if settings.TLSConfig != nil {
		d.TLSConfig = settings.TLSConfig
	} else {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true, ServerName: settings.Host}
	}

	if settings.Secure {
		d.SSL = true
	}
	return d.DialAndSend(m)
}
