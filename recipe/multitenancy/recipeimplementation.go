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

package multitenancy

import (
	"errors"
	"fmt"

	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeRecipeImplementation(querier supertokens.Querier, config multitenancymodels.TypeNormalisedInput, appInfo supertokens.NormalisedAppinfo) multitenancymodels.RecipeInterface {
	getTenantId := func(tenantIdFromFrontend string, userContext supertokens.UserContext) (string, error) {
		return tenantIdFromFrontend, nil
	}

	createOrUpdateTenant := func(tenantId string, config multitenancymodels.TenantConfig, userContext supertokens.UserContext) (multitenancymodels.CreateOrUpdateTenantResponse, error) {
		requestBody := map[string]interface{}{
			"tenantId": tenantId,
		}
		if config.EmailPasswordEnabled != nil {
			requestBody["emailPasswordEnabled"] = *config.EmailPasswordEnabled
		}
		if config.PasswordlessEnabled != nil {
			requestBody["passwordlessEnabled"] = *config.PasswordlessEnabled
		}
		if config.ThirdPartyEnabled != nil {
			requestBody["thirdPartyEnabled"] = *config.ThirdPartyEnabled
		}
		if config.CoreConfig != nil {
			requestBody["coreConfig"] = config.CoreConfig
		}
		createOrUpdateResponse, err := querier.SendPutRequest("/recipe/multitenancy/tenant", requestBody)
		if err != nil {
			return multitenancymodels.CreateOrUpdateTenantResponse{}, err
		}

		_, ok := createOrUpdateResponse["status"].(string)
		if ok {
			return multitenancymodels.CreateOrUpdateTenantResponse{
				OK: &struct{ CreatedNew bool }{
					CreatedNew: createOrUpdateResponse["createdNew"].(bool),
				},
			}, nil
		}

		return multitenancymodels.CreateOrUpdateTenantResponse{}, errors.New("should not come here")
	}

	deleteTenant := func(tenantId string, userContext supertokens.UserContext) (multitenancymodels.DeleteTenantResponse, error) {
		deleteTenantResponse, err := querier.SendPostRequest("/recipe/multitenancy/tenant/remove", map[string]interface{}{
			"tenantId": tenantId,
		})
		if err != nil {
			return multitenancymodels.DeleteTenantResponse{}, err
		}
		_, ok := deleteTenantResponse["status"].(string)
		if ok {
			return multitenancymodels.DeleteTenantResponse{
				OK: &struct{ DidExist bool }{
					DidExist: deleteTenantResponse["didExist"].(bool),
				},
			}, nil
		}

		return multitenancymodels.DeleteTenantResponse{}, errors.New("should not come here")
	}

	getTenant := func(tenantId string, userContext supertokens.UserContext) (*multitenancymodels.Tenant, error) {
		tenantResponse, err := querier.SendGetRequest(fmt.Sprintf("/%s/recipe/multitenancy/tenant", tenantId), map[string]string{})
		if err != nil {
			return nil, err
		}
		status, ok := tenantResponse["status"].(string)
		if ok {
			if status == "TENANT_NOT_FOUND_ERROR" {
				return nil, nil
			}

			result := &multitenancymodels.Tenant{}
			err = supertokens.MapToStruct(tenantResponse, result)
			if err != nil {
				return nil, err
			}

			if status == "OK" {
				return result, nil
			}
		}

		return nil, errors.New("should not come here")
	}

	listAllTenants := func(userContext supertokens.UserContext) (multitenancymodels.ListAllTenantsResponse, error) {
		tenantsResponse, err := querier.SendGetRequest("/recipe/multitenancy/tenant/list", map[string]string{})
		if err != nil {
			return multitenancymodels.ListAllTenantsResponse{}, err
		}
		result := multitenancymodels.ListAllTenantsResponse{
			OK: &struct {
				Tenants []multitenancymodels.Tenant `json:"tenants"`
			}{},
		}
		err = supertokens.MapToStruct(tenantsResponse, result.OK)
		if err != nil {
			return multitenancymodels.ListAllTenantsResponse{}, err
		}
		return result, nil
	}

	createOrUpdateThirdPartyConfig := func(tenantId string, config tpmodels.ProviderConfig, skipValidation *bool, userContext supertokens.UserContext) (multitenancymodels.CreateOrUpdateThirdPartyConfigResponse, error) {
		configMap, err := supertokens.StructToMap(config)
		if err != nil {
			return multitenancymodels.CreateOrUpdateThirdPartyConfigResponse{}, err
		}

		requestBody := map[string]interface{}{
			"config": configMap,
		}
		if skipValidation != nil {
			requestBody["skipApiCleanUp"] = *skipValidation
		}
		response, err := querier.SendPutRequest(fmt.Sprintf("/%s/recipe/multitenancy/config/thirdparty", tenantId), requestBody)
		if err != nil {
			return multitenancymodels.CreateOrUpdateThirdPartyConfigResponse{}, err
		}
		return multitenancymodels.CreateOrUpdateThirdPartyConfigResponse{
			OK: &struct {
				CreatedNew bool
			}{
				CreatedNew: response["createdNew"].(bool),
			},
		}, nil
	}

	deleteThirdPartyConfig := func(tenantId string, thirdPartyId string, userContext supertokens.UserContext) (multitenancymodels.DeleteThirdPartyConfigResponse, error) {
		response, err := querier.SendPostRequest(fmt.Sprintf("/%s/recipe/multitenancy/config/thirdparty/remove", tenantId), map[string]interface{}{
			"thirdPartyId": thirdPartyId,
		})
		if err != nil {
			return multitenancymodels.DeleteThirdPartyConfigResponse{}, err
		}

		return multitenancymodels.DeleteThirdPartyConfigResponse{
			OK: &struct{ DidConfigExist bool }{
				DidConfigExist: response["didConfigExist"].(bool),
			},
		}, nil
	}

	associateUserToTenant := func(tenantId string, userId string, userContext supertokens.UserContext) (multitenancymodels.AssociateUserToTenantResponse, error) {
		response, err := querier.SendPostRequest(fmt.Sprintf("/%s/recipe/multitenancy/tenant/user", tenantId), map[string]interface{}{
			"userId": userId,
		})
		if err != nil {
			return multitenancymodels.AssociateUserToTenantResponse{}, err
		}
		return multitenancymodels.AssociateUserToTenantResponse{
			OK: &struct{ WasAlreadyAssociated bool }{
				WasAlreadyAssociated: response["wasAlreadyAssociated"].(bool),
			},
		}, nil
	}

	disassociateUserFromTenant := func(tenantId string, userId string, userContext supertokens.UserContext) (multitenancymodels.DisassociateUserFromTenantResponse, error) {
		response, err := querier.SendPostRequest(fmt.Sprintf("/%s/recipe/multitenancy/tenant/user/remove", tenantId), map[string]interface{}{
			"userId": userId,
		})
		if err != nil {
			return multitenancymodels.DisassociateUserFromTenantResponse{}, err
		}
		return multitenancymodels.DisassociateUserFromTenantResponse{
			OK: &struct{ WasAssociated bool }{
				WasAssociated: response["wasAssociated"].(bool),
			},
		}, nil
	}

	return multitenancymodels.RecipeInterface{
		GetTenantId: &getTenantId,

		CreateOrUpdateTenant: &createOrUpdateTenant,
		DeleteTenant:         &deleteTenant,
		GetTenant:            &getTenant,
		ListAllTenants:       &listAllTenants,

		CreateOrUpdateThirdPartyConfig: &createOrUpdateThirdPartyConfig,
		DeleteThirdPartyConfig:         &deleteThirdPartyConfig,

		AssociateUserToTenant:      &associateUserToTenant,
		DisassociateUserFromTenant: &disassociateUserFromTenant,
	}
}
