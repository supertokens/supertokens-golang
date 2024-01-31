/* Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
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

package supertokens

func validateAndNormaliseAccountLinkingUserInput(appInfo NormalisedAppinfo, config *AccountLinkingTypeInput) AccountLinkingTypeNormalisedInput {

	typeNormalisedInput := makeTypeNormalisedInput(appInfo)

	if config != nil && config.ShouldDoAutomaticAccountLinking != nil {
		typeNormalisedInput.ShouldDoAutomaticAccountLinking = config.ShouldDoAutomaticAccountLinking
	}

	if config != nil && config.OnAccountLinked != nil {
		typeNormalisedInput.OnAccountLinked = config.OnAccountLinked
	}

	if config != nil && config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
	}

	return typeNormalisedInput
}

func makeTypeNormalisedInput(appInfo NormalisedAppinfo) AccountLinkingTypeNormalisedInput {
	return AccountLinkingTypeNormalisedInput{
		OnAccountLinked: func(user User, newAccountUser RecipeLevelUser, userContext UserContext) error {
			return nil
		},
		ShouldDoAutomaticAccountLinking: func(newAccountInfo AccountInfoWithRecipeIdAndWithRecipeUserId, user *User, tenantID string, userContext UserContext) (ShouldDoAutomaticAccountLinkingResponse, error) {
			return ShouldDoAutomaticAccountLinkingResponse{
				ShouldAutomaticallyLink:   false,
				ShouldRequireVerification: false,
			}, nil
		},
		Override: AccountLinkingOverrideStruct{
			Functions: func(originalImplementation AccountLinkingRecipeInterface) AccountLinkingRecipeInterface {
				return originalImplementation
			},
		},
	}
}
