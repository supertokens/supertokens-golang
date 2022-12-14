/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
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

package multitenancymodels

import "github.com/supertokens/supertokens-golang/supertokens"

type APIInterface struct {
	LoginMethodsGET func(tenantId *string, userContext supertokens.UserContext) (LoginMethodsGETResponse, error)
}

type LoginMethodsGETResponse struct {
	OK *struct {
		Emailpassword struct {
			Enabled bool
		}
		Passwordless struct {
			Enabled bool
		}
		ThirdParty struct {
			Enabled   bool
			Providers []struct {
				Id   string
				Name string
			}
		}
	}
	GeneralError *supertokens.GeneralErrorResponse
}