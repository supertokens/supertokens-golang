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

package smtpService

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	plesssmtpService "github.com/supertokens/supertokens-golang/recipe/passwordless/emaildelivery/smtpService"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeSMTPService(config emaildelivery.SMTPServiceConfig) *emaildelivery.EmailDeliveryInterface {
	serviceImpl := makeServiceImplementation(config.Settings)

	if config.Override != nil {
		serviceImpl = config.Override(serviceImpl)
	}

	passwordlessServiceImpl := plesssmtpService.MakeSMTPService(emaildelivery.SMTPServiceConfig{
		Settings: config.Settings,
		Override: makePasswordlessServiceImplementation(serviceImpl),
	})

	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.PasswordlessLogin != nil {
			return (*passwordlessServiceImpl.SendEmail)(input, userContext)

		} else {
			return errors.New("should never come here")
		}
	}

	return &emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}
