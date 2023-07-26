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

package userrolesmodels

import "github.com/supertokens/supertokens-golang/supertokens"

type UnknownRoleError struct {
}

type AddRoleToUserResponse struct {
	OK *struct {
		DidUserAlreadyHaveRole bool
	}
	UnknownRoleError *UnknownRoleError
}

type RemoveUserRoleResponse struct {
	OK *struct {
		DidUserHaveRole bool
	}
	UnknownRoleError *UnknownRoleError
}

type GetRolesForUserResponse struct {
	OK *struct {
		Roles []string
	}
}

type GetUsersThatHaveRoleResponse struct {
	OK *struct {
		Users []string
	}
	UnknownRoleError *UnknownRoleError
}

type CreateNewRoleOrAddPermissionsResponse struct {
	OK *struct {
		CreatedNewRole bool
	}
}

type GetPermissionsForRoleResponse struct {
	OK *struct {
		Permissions []string
	}
	UnknownRoleError *UnknownRoleError
}

type RemovePermissionsFromRoleResponse struct {
	OK               *struct{}
	UnknownRoleError *UnknownRoleError
}

type GetRolesThatHavePermissionResponse struct {
	OK *struct {
		Roles []string
	}
}

type DeleteRoleResponse struct {
	OK *struct {
		DidRoleExist bool
	}
}

type GetAllRolesResponse struct {
	OK *struct {
		Roles []string
	}
}

type RecipeInterface struct {
	AddRoleToUser                 *func(userID string, role string, tenantId string, userContext supertokens.UserContext) (AddRoleToUserResponse, error)
	RemoveUserRole                *func(userID string, role string, tenantId string, userContext supertokens.UserContext) (RemoveUserRoleResponse, error)
	GetRolesForUser               *func(userID string, tenantId string, userContext supertokens.UserContext) (GetRolesForUserResponse, error)
	GetUsersThatHaveRole          *func(role string, tenantId string, userContext supertokens.UserContext) (GetUsersThatHaveRoleResponse, error)
	CreateNewRoleOrAddPermissions *func(role string, permissions []string, userContext supertokens.UserContext) (CreateNewRoleOrAddPermissionsResponse, error)
	GetPermissionsForRole         *func(role string, userContext supertokens.UserContext) (GetPermissionsForRoleResponse, error)
	RemovePermissionsFromRole     *func(role string, permissions []string, userContext supertokens.UserContext) (RemovePermissionsFromRoleResponse, error)
	GetRolesThatHavePermission    *func(permission string, userContext supertokens.UserContext) (GetRolesThatHavePermissionResponse, error)
	DeleteRole                    *func(role string, userContext supertokens.UserContext) (DeleteRoleResponse, error)
	GetAllRoles                   *func(userContext supertokens.UserContext) (GetAllRolesResponse, error)
}
