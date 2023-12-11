/* Copyright (c) 2023, VRAI Labs and/or its affiliates. All rights reserved.
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

type UserPaginationResult struct {
	Users               []User
	NextPaginationToken *string
}

type CanCreatePrimaryUserResponse struct {
	OK *struct {
		WasAlreadyAPrimaryUser bool
	}
	RecipeUserIdAlreadyLinkedWithPrimaryUserIdError *struct {
		PrimaryUserId string
		Description   string
	}
	AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError *struct {
		PrimaryUserId string
		Description   string
	}
}

type CreatePrimaryUserResponse struct {
	OK *struct {
		User                   User
		WasAlreadyAPrimaryUser bool
	}
	RecipeUserIdAlreadyLinkedWithPrimaryUserIdError *struct {
		PrimaryUserId string
		Description   string
	}
	AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError *struct {
		PrimaryUserId string
		Description   string
	}
}

type CanLinkAccountResponse struct {
	OK *struct {
		AccountsAlreadyLinked bool
	}
	RecipeUserIdAlreadyLinkedWithAnotherPrimaryUserIdError *struct {
		PrimaryUserId string
		Description   string
	}
	AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError *struct {
		PrimaryUserId string
		Description   string
	}
	InputUserIsNotAPrimaryUserError *struct{}
}

type LinkAccountResponse struct {
	OK *struct {
		AccountsAlreadyLinked bool
		User                  User
	}
	RecipeUserIdAlreadyLinkedWithAnotherPrimaryUserIdError *struct {
		PrimaryUserId string
		User          User
	}
	AccountInfoAlreadyAssociatedWithAnotherPrimaryUserIdError *struct {
		PrimaryUserId string
		Description   string
	}
	InputUserIsNotAPrimaryUserError *struct{}
}

type UnlinkAccountsResponse struct {
	WasRecipeUserDeleted bool
	WasLinked            bool
}

type AccountLinkingRecipeInterface struct {
	GetUsersWithSearchParams *func(tenantID string, timeJoinedOrder string, paginationToken *string, limit *int, includeRecipeIds *[]string, searchParams map[string]string, userContext UserContext) (UserPaginationResult, error)

	CanCreatePrimaryUser *func(recipeUserId RecipeUserID, userContext UserContext) (CanCreatePrimaryUserResponse, error)

	CreatePrimaryUser *func(recipeUserId RecipeUserID, userContext UserContext) (CreatePrimaryUserResponse, error)

	CanLinkAccounts *func(recipeUserId RecipeUserID, primaryUserId string, userContext UserContext) (CanLinkAccountResponse, error)

	LinkAccounts *func(recipeUserId RecipeUserID, primaryUserId string, userContext UserContext) (LinkAccountResponse, error)

	UnlinkAccounts *func(recipeUserId RecipeUserID, userContext UserContext) (UnlinkAccountsResponse, error)

	GetUser *func(userId string, userContext UserContext) (*User, error)

	ListUsersByAccountInfo *func(tenantID string, accountInfo AccountInfo, doUnionOfAccountInfo bool, userContext UserContext) ([]User, error)

	DeleteUser *func(userId string, removeAllLinkedAccounts bool, userContext UserContext) error
}
