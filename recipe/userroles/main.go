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

package userroles

import (
	"github.com/supertokens/supertokens-golang/recipe/userroles/userrolesmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init(config *userrolesmodels.TypeInput) supertokens.Recipe {
	return recipeInit(config)
}

func AddRoleToUser(userID string, role string, tenantId *string, userContext supertokens.UserContext) (userrolesmodels.AddRoleToUserResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.AddRoleToUserResponse{}, err
	}
	return (*instance.RecipeImpl.AddRoleToUser)(userID, role, tenantId, userContext)
}

func RemoveUserRole(userID string, role string, tenantId *string, userContext supertokens.UserContext) (userrolesmodels.RemoveUserRoleResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.RemoveUserRoleResponse{}, err
	}
	return (*instance.RecipeImpl.RemoveUserRole)(userID, role, tenantId, userContext)
}

func GetRolesForUser(userID string, tenantId *string, userContext supertokens.UserContext) (userrolesmodels.GetRolesForUserResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.GetRolesForUserResponse{}, err
	}
	return (*instance.RecipeImpl.GetRolesForUser)(userID, tenantId, userContext)
}

func GetUsersThatHaveRole(role string, tenantId *string, userContext supertokens.UserContext) (userrolesmodels.GetUsersThatHaveRoleResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.GetUsersThatHaveRoleResponse{}, err
	}
	return (*instance.RecipeImpl.GetUsersThatHaveRole)(role, tenantId, userContext)
}

func CreateNewRoleOrAddPermissions(role string, permissions []string, tenantId *string, userContext supertokens.UserContext) (userrolesmodels.CreateNewRoleOrAddPermissionsResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.CreateNewRoleOrAddPermissionsResponse{}, err
	}
	return (*instance.RecipeImpl.CreateNewRoleOrAddPermissions)(role, permissions, tenantId, userContext)
}

func GetPermissionsForRole(role string, tenantId *string, userContext supertokens.UserContext) (userrolesmodels.GetPermissionsForRoleResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.GetPermissionsForRoleResponse{}, err
	}
	return (*instance.RecipeImpl.GetPermissionsForRole)(role, tenantId, userContext)
}

func RemovePermissionsFromRole(role string, permissions []string, tenantId *string, userContext supertokens.UserContext) (userrolesmodels.RemovePermissionsFromRoleResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.RemovePermissionsFromRoleResponse{}, err
	}
	return (*instance.RecipeImpl.RemovePermissionsFromRole)(role, permissions, tenantId, userContext)
}

func GetRolesThatHavePermission(permission string, tenantId *string, userContext supertokens.UserContext) (userrolesmodels.GetRolesThatHavePermissionResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.GetRolesThatHavePermissionResponse{}, err
	}
	return (*instance.RecipeImpl.GetRolesThatHavePermission)(permission, tenantId, userContext)
}

func DeleteRole(role string, tenantId *string, userContext supertokens.UserContext) (userrolesmodels.DeleteRoleResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.DeleteRoleResponse{}, err
	}
	return (*instance.RecipeImpl.DeleteRole)(role, tenantId, userContext)
}

func GetAllRoles(tenantId *string, userContext supertokens.UserContext) (userrolesmodels.GetAllRolesResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.GetAllRolesResponse{}, err
	}
	return (*instance.RecipeImpl.GetAllRoles)(tenantId, userContext)
}
