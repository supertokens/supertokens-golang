package supertokens

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

type CreateOrUpdateTenantIDConfigResponse struct {
	OK *struct {
		Created bool
		Updated bool
	}
}

func CreateOrUpdateTenantIDConfigMapping(thirdPartyProviderID string, tenantID string, config interface{}) (CreateOrUpdateTenantIDConfigResponse, error) {
	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return CreateOrUpdateTenantIDConfigResponse{}, err
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		return CreateOrUpdateTenantIDConfigResponse{}, err
	}

	// TODO: Update version
	if maxVersion(cdiVersion, "2.15") != cdiVersion {
		return CreateOrUpdateTenantIDConfigResponse{}, errors.New("Please upgrade the SuperTokens core to >= 3.15.0")
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return CreateOrUpdateTenantIDConfigResponse{}, err
	}

	configStr := base64.RawStdEncoding.EncodeToString(configBytes)

	data := map[string]interface{}{
		"thirdPartyProviderId": thirdPartyProviderID,
		"tenantId":             tenantID,
		"config":               configStr,
	}

	resp, err := querier.SendPostRequest("/recipe/multitenant-config/map", data)
	if err != nil {
		return CreateOrUpdateTenantIDConfigResponse{}, err
	}
	return CreateOrUpdateTenantIDConfigResponse{
		OK: &struct {
			Created bool
			Updated bool
		}{
			Created: resp["created"].(bool),
			Updated: resp["updated"].(bool),
		},
	}, nil
}

type FetchTenantIDConfigResponse struct {
	OK                  *struct{}
	UnknownMappingError *struct{}
}

func FetchTenantIDConfigMapping(thirdPartyProviderID string, tenantID string, config interface{}) (FetchTenantIDConfigResponse, error) {
	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return FetchTenantIDConfigResponse{}, err
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		return FetchTenantIDConfigResponse{}, err
	}

	// TODO: Update version
	if maxVersion(cdiVersion, "2.15") != cdiVersion {
		return FetchTenantIDConfigResponse{}, errors.New("Please upgrade the SuperTokens core to >= 3.15.0")
	}

	data := map[string]string{
		"thirdPartyProviderId": thirdPartyProviderID,
		"tenantId":             tenantID,
	}

	resp, err := querier.SendGetRequest("/recipe/multitenant-config/map", data)
	if err != nil {
		return FetchTenantIDConfigResponse{}, err
	}
	if resp["status"] == "OK" {
		configStr := resp["config"].(string)
		configBytes, err := base64.RawStdEncoding.DecodeString(configStr)
		if err != nil {
			return FetchTenantIDConfigResponse{}, err
		}
		err = json.Unmarshal(configBytes, config)
		if err != nil {
			return FetchTenantIDConfigResponse{}, err
		}

		return FetchTenantIDConfigResponse{
			OK: &struct{}{},
		}, nil

	} else {
		return FetchTenantIDConfigResponse{
			UnknownMappingError: &struct{}{},
		}, nil
	}
}

type DeleteTenantIDConfigResponse struct {
	OK *struct {
		DidMappingExist bool
	}
}

func DeleteTenantIDConfigMapping(thirdPartyProviderID string, tenantID string) (DeleteTenantIDConfigResponse, error) {
	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return DeleteTenantIDConfigResponse{}, err
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		return DeleteTenantIDConfigResponse{}, err
	}

	// TODO: Update version
	if maxVersion(cdiVersion, "2.15") != cdiVersion {
		return DeleteTenantIDConfigResponse{}, errors.New("Please upgrade the SuperTokens core to >= 3.15.0")
	}

	data := map[string]interface{}{
		"thirdPartyProviderId": thirdPartyProviderID,
		"tenantId":             tenantID,
	}
	resp, err := querier.SendPostRequest("/recipe/multitenant-config/map/remove", data)
	if err != nil {
		return DeleteTenantIDConfigResponse{}, err
	}
	return DeleteTenantIDConfigResponse{
		OK: &struct{ DidMappingExist bool }{
			DidMappingExist: resp["didMappingExist"].(bool),
		},
	}, nil
}

type ListMultitenantConfigResponse struct {
	OK *struct {
		Configs []struct {
			ThirdPartyProviderId string
			TenantId             string
			Config               interface{}
		}
	}
}

func ListMultitenantConfigMapping(thirdPartyProviderId *string) (ListMultitenantConfigResponse, error) {
	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return ListMultitenantConfigResponse{}, err
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		return ListMultitenantConfigResponse{}, err
	}

	// TODO: Update version
	if maxVersion(cdiVersion, "2.15") != cdiVersion {
		return ListMultitenantConfigResponse{}, errors.New("Please upgrade the SuperTokens core to >= 3.15.0")
	}

	data := map[string]string{}

	if thirdPartyProviderId != nil {
		data["thirdPartyProviderId"] = *thirdPartyProviderId
	}

	resp, err := querier.SendGetRequest("/recipe/multitenant-config/map/list", data)
	if err != nil {
		return ListMultitenantConfigResponse{}, err
	}

	configs := resp["configs"].([]interface{})
	configObjs := make([]struct {
		ThirdPartyProviderId string
		TenantId             string
		Config               interface{}
	}, len(configs))

	for idx, config := range configs {
		configObjs[idx].ThirdPartyProviderId = config.(map[string]interface{})["thirdPartyProviderId"].(string)
		configObjs[idx].TenantId = config.(map[string]interface{})["tenantId"].(string)
		configStr := config.(map[string]interface{})["config"].(string)
		configBytes, err := base64.RawStdEncoding.DecodeString(configStr)
		if err != nil {
			return ListMultitenantConfigResponse{}, err
		}
		err = json.Unmarshal(configBytes, &configObjs[idx].Config)
		if err != nil {
			return ListMultitenantConfigResponse{}, err
		}
	}

	return ListMultitenantConfigResponse{
		OK: &struct {
			Configs []struct {
				ThirdPartyProviderId string
				TenantId             string
				Config               interface{}
			}
		}{
			Configs: configObjs,
		},
	}, nil
}
