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
	"github.com/supertokens/supertokens-golang/supertokens"
)

type SMTPSettings struct {
	Host     string
	From     SMTPFrom
	Port     int
	Password string
	Secure   bool
}

type SMTPFrom struct {
	Name  string
	Email string
}

type SMTPContent struct {
	Body    string
	IsHtml  bool
	Subject string
	ToEmail string
}

type SMTPServiceInterface struct {
	SendRawEmail *func(input SMTPContent, userContext supertokens.UserContext) error
	GetContent   *func(input EmailType, userContext supertokens.UserContext) (SMTPContent, error)
}

type SMTPServiceConfig struct {
	Settings SMTPSettings
	Override func(originalImplementation SMTPServiceInterface) SMTPServiceInterface
}
