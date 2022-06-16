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
