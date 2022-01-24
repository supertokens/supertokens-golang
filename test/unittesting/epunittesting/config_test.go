/*
 * Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
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

package epunittesting

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDefaultConfigForEmailPasswordModule(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
		},
	}
	unittesting.StartingHelper()
	err := supertokens.Init(configValue)
	if err != nil {
		log.Fatal(err.Error())
	}
	singletonEmailPasswordInstance, err := emailpassword.GetRecipeInstanceOrThrowError()
	if err != nil {
		log.Fatal(err.Error())
	}
	signupFeature := singletonEmailPasswordInstance.Config.SignUpFeature
	assert.Equal(t, len(signupFeature.FormFields), 2)
	for i := 0; i < len(signupFeature.FormFields); i++ {
		assert.Equal(t, signupFeature.FormFields[i].Optional, false)
		//*to add test for validate function
	}

	singInFeature := singletonEmailPasswordInstance.Config.SignInFeature
	assert.Equal(t, len(singInFeature.FormFields), 2)
	for i := 0; i < len(singInFeature.FormFields); i++ {
		assert.Equal(t, singInFeature.FormFields[i].Optional, false)
		//*to add test for validate function
	}

	resetPasswordUsingTokenFeature := singletonEmailPasswordInstance.Config.ResetPasswordUsingTokenFeature
	assert.Equal(t, len(resetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm), 1)

	unittesting.EndingHelper()

}
