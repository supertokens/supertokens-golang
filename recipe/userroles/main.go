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

func AddRoleToUser(userID string, role string, userContext ...supertokens.UserContext) (userrolesmodels.AddRoleToUserResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.AddRoleToUserResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.AddRoleToUser)(userID, role, userContext[0])
}

func RemoveUserRole(userID string, role string, userContext ...supertokens.UserContext) (userrolesmodels.RemoveUserRoleResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.RemoveUserRoleResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.RemoveUserRole)(userID, role, userContext[0])
}

func GetRolesForUser(userID string, userContext ...supertokens.UserContext) (userrolesmodels.GetRolesForUserResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.GetRolesForUserResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetRolesForUser)(userID, userContext[0])
}

func GetUsersThatHaveRole(role string, userContext ...supertokens.UserContext) (userrolesmodels.GetUsersThatHaveRoleResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.GetUsersThatHaveRoleResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetUsersThatHaveRole)(role, userContext[0])
}

func CreateNewRoleOrAddPermissions(role string, permissions []string, userContext ...supertokens.UserContext) (userrolesmodels.CreateNewRoleOrAddPermissionsResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.CreateNewRoleOrAddPermissionsResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.CreateNewRoleOrAddPermissions)(role, permissions, userContext[0])
}

func GetPermissionsForRole(role string, userContext ...supertokens.UserContext) (userrolesmodels.GetPermissionsForRoleResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.GetPermissionsForRoleResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetPermissionsForRole)(role, userContext[0])
}

func RemovePermissionsFromRole(role string, permissions []string, userContext ...supertokens.UserContext) (userrolesmodels.RemovePermissionsFromRoleResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.RemovePermissionsFromRoleResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.RemovePermissionsFromRole)(role, permissions, userContext[0])
}

func GetRolesThatHavePermission(permission string, userContext ...supertokens.UserContext) (userrolesmodels.GetRolesThatHavePermissionResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.GetRolesThatHavePermissionResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetRolesThatHavePermission)(permission, userContext[0])
}

func DeleteRole(role string, userContext ...supertokens.UserContext) (userrolesmodels.DeleteRoleResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.DeleteRoleResponse{}, err
	}
	return (*instance.RecipeImpl.DeleteRole)(role, userContext[0])
}

func GetAllRoles(userContext ...supertokens.UserContext) (userrolesmodels.GetAllRolesResponse, error) {
	instance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		return userrolesmodels.GetAllRolesResponse{}, err
	}
	if len(userContext) == 0 {
		userContext = append(userContext, &map[string]interface{}{})
	}
	return (*instance.RecipeImpl.GetAllRoles)(userContext[0])
}
