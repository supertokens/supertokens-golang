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

package thirdparty

const (
	AuthorisationAPI = "/authorisationurl"
	SignInUpAPI      = "/signinup"

	// If Third Party login is used with one of the following development keys, then the dev authorization url and the redirect url will be used.
	// When adding or changing client id's they should be in the following order: Google and Facebook

	DevOauthAuthorisationUrl = "https://supertokens.io/dev/oauth/redirect-to-provider"
	DevOauthRedirectUrl      = "https://supertokens.io/dev/oauth/redirect-to-app"
)

var DevOauthClientIds = map[string]bool{
	"1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com": true, // google
	"467101b197249757c71f": true, // github
}
