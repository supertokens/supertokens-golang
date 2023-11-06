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

func makeRecipeImplementation(querier supertokens.Querier, config userrolesmodels.TypeNormalisedInput, appInfo supertokens.NormalisedAppinfo) userrolesmodels.RecipeInterface {

	addRoleToUser := func(userID string, role string, tenantId string, userContext supertokens.UserContext) (userrolesmodels.AddRoleToUserResponse, error) {
		response, err := querier.SendPutRequest(tenantId+"/recipe/user/role", map[string]interface{}{
			"userId": userID,
			"role":   role,
		}, userContext)
		if err != nil {
			return userrolesmodels.AddRoleToUserResponse{}, err
		}

		if response["status"] == "OK" {
			return userrolesmodels.AddRoleToUserResponse{
				OK: &struct{ DidUserAlreadyHaveRole bool }{
					DidUserAlreadyHaveRole: response["didUserAlreadyHaveRole"].(bool),
				},
			}, nil
		}

		return userrolesmodels.AddRoleToUserResponse{
			UnknownRoleError: &userrolesmodels.UnknownRoleError{},
		}, nil
	}

	removeUserRole := func(userID string, role string, tenantId string, userContext supertokens.UserContext) (userrolesmodels.RemoveUserRoleResponse, error) {
		response, err := querier.SendPostRequest(tenantId+"/recipe/user/role/remove", map[string]interface{}{
			"userId": userID,
			"role":   role,
		}, userContext)
		if err != nil {
			return userrolesmodels.RemoveUserRoleResponse{}, err
		}

		if response["status"] == "OK" {
			return userrolesmodels.RemoveUserRoleResponse{
				OK: &struct{ DidUserHaveRole bool }{
					DidUserHaveRole: response["didUserHaveRole"].(bool),
				},
			}, nil
		}

		return userrolesmodels.RemoveUserRoleResponse{
			UnknownRoleError: &userrolesmodels.UnknownRoleError{},
		}, nil
	}

	getRolesForUser := func(userID string, tenantId string, userContext supertokens.UserContext) (userrolesmodels.GetRolesForUserResponse, error) {
		response, err := querier.SendGetRequest(tenantId+"/recipe/user/roles", map[string]string{
			"userId": userID,
		}, userContext)
		if err != nil {
			return userrolesmodels.GetRolesForUserResponse{}, err
		}

		return userrolesmodels.GetRolesForUserResponse{
			OK: &struct{ Roles []string }{
				Roles: convertToStringArray(response["roles"].([]interface{})),
			},
		}, nil

	}

	getUsersThatHaveRole := func(role string, tenantId string, userContext supertokens.UserContext) (userrolesmodels.GetUsersThatHaveRoleResponse, error) {
		response, err := querier.SendGetRequest(tenantId+"/recipe/role/users", map[string]string{
			"role": role,
		}, userContext)
		if err != nil {
			return userrolesmodels.GetUsersThatHaveRoleResponse{}, err
		}

		if response["status"] == "OK" {
			return userrolesmodels.GetUsersThatHaveRoleResponse{
				OK: &struct{ Users []string }{
					Users: convertToStringArray(response["users"].([]interface{})),
				},
			}, nil
		}

		return userrolesmodels.GetUsersThatHaveRoleResponse{
			UnknownRoleError: &userrolesmodels.UnknownRoleError{},
		}, nil
	}

	createNewRoleOrAddPermissions := func(role string, permissions []string, userContext supertokens.UserContext) (userrolesmodels.CreateNewRoleOrAddPermissionsResponse, error) {
		response, err := querier.SendPutRequest("/recipe/role", map[string]interface{}{
			"role":        role,
			"permissions": permissions,
		}, userContext)
		if err != nil {
			return userrolesmodels.CreateNewRoleOrAddPermissionsResponse{}, err
		}

		return userrolesmodels.CreateNewRoleOrAddPermissionsResponse{
			OK: &struct{ CreatedNewRole bool }{
				CreatedNewRole: response["createdNewRole"].(bool),
			},
		}, nil
	}

	getPermissionsForRole := func(role string, userContext supertokens.UserContext) (userrolesmodels.GetPermissionsForRoleResponse, error) {
		response, err := querier.SendGetRequest("/recipe/role/permissions", map[string]string{
			"role": role,
		}, userContext)
		if err != nil {
			return userrolesmodels.GetPermissionsForRoleResponse{}, err
		}

		if response["status"] == "OK" {
			return userrolesmodels.GetPermissionsForRoleResponse{
				OK: &struct{ Permissions []string }{
					Permissions: convertToStringArray(response["permissions"].([]interface{})),
				},
			}, nil
		}

		return userrolesmodels.GetPermissionsForRoleResponse{
			UnknownRoleError: &userrolesmodels.UnknownRoleError{},
		}, nil
	}

	removePermissionsFromRole := func(role string, permissions []string, userContext supertokens.UserContext) (userrolesmodels.RemovePermissionsFromRoleResponse, error) {
		response, err := querier.SendPostRequest("/recipe/role/permissions/remove", map[string]interface{}{
			"role":        role,
			"permissions": permissions,
		}, userContext)
		if err != nil {
			return userrolesmodels.RemovePermissionsFromRoleResponse{}, err
		}

		if response["status"] == "OK" {
			return userrolesmodels.RemovePermissionsFromRoleResponse{
				OK: &struct{}{},
			}, nil
		}

		return userrolesmodels.RemovePermissionsFromRoleResponse{
			UnknownRoleError: &userrolesmodels.UnknownRoleError{},
		}, nil
	}

	getRolesThatHavePermission := func(permission string, userContext supertokens.UserContext) (userrolesmodels.GetRolesThatHavePermissionResponse, error) {
		response, err := querier.SendGetRequest("/recipe/permission/roles", map[string]string{
			"permission": permission,
		}, userContext)
		if err != nil {
			return userrolesmodels.GetRolesThatHavePermissionResponse{}, err
		}

		return userrolesmodels.GetRolesThatHavePermissionResponse{
			OK: &struct{ Roles []string }{
				Roles: convertToStringArray(response["roles"].([]interface{})),
			},
		}, nil
	}

	deleteRole := func(role string, userContext supertokens.UserContext) (userrolesmodels.DeleteRoleResponse, error) {
		response, err := querier.SendPostRequest("/recipe/role/remove", map[string]interface{}{
			"role": role,
		}, userContext)
		if err != nil {
			return userrolesmodels.DeleteRoleResponse{}, err
		}

		return userrolesmodels.DeleteRoleResponse{
			OK: &struct{ DidRoleExist bool }{
				DidRoleExist: response["didRoleExist"].(bool),
			},
		}, nil
	}

	getAllRoles := func(userContext supertokens.UserContext) (userrolesmodels.GetAllRolesResponse, error) {
		response, err := querier.SendGetRequest("/recipe/roles", map[string]string{}, userContext)
		if err != nil {
			return userrolesmodels.GetAllRolesResponse{}, err
		}

		return userrolesmodels.GetAllRolesResponse{
			OK: &struct{ Roles []string }{
				Roles: convertToStringArray(response["roles"].([]interface{})),
			},
		}, nil
	}

	return userrolesmodels.RecipeInterface{
		AddRoleToUser:                 &addRoleToUser,
		RemoveUserRole:                &removeUserRole,
		GetRolesForUser:               &getRolesForUser,
		GetUsersThatHaveRole:          &getUsersThatHaveRole,
		CreateNewRoleOrAddPermissions: &createNewRoleOrAddPermissions,
		GetPermissionsForRole:         &getPermissionsForRole,
		RemovePermissionsFromRole:     &removePermissionsFromRole,
		GetRolesThatHavePermission:    &getRolesThatHavePermission,
		DeleteRole:                    &deleteRole,
		GetAllRoles:                   &getAllRoles,
	}
}
