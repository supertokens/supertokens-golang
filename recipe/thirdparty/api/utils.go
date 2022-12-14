package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/derekstavis/go-qs"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/providers"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func findAndCreateProviderInstance(options tpmodels.APIOptions, thirdPartyId string, tenantId *string) (tpmodels.TypeProvider, error) {
	for _, provider := range options.Providers {
		if provider.Config.ThirdPartyId == thirdPartyId {
			useForDefault := true
			if provider.UseForDefaultTenant != nil {
				useForDefault = *provider.UseForDefaultTenant
			}
			if (tenantId == nil || *tenantId == tpmodels.DefaultTenantId) && !useForDefault {
				return tpmodels.TypeProvider{}, fmt.Errorf("the provider %s is disabled for default tenant", thirdPartyId)
			}

			providerInstance := createProvider(provider)
			return *providerInstance, nil
		}
	}
	return tpmodels.TypeProvider{}, fmt.Errorf("the provider %s could not be found in the configuration", thirdPartyId)
}

func createProvider(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	var providerInstance *tpmodels.TypeProvider

	switch input.Config.ThirdPartyId {
	case "active-directory":
		providerInstance = providers.ActiveDirectory(input)
	case "apple":
		providerInstance = providers.Apple(input)
	case "discord":
		providerInstance = providers.Discord(input)
	case "facebook":
		providerInstance = providers.Facebook(input)
	case "github":
		providerInstance = providers.Github(input)
	case "google":
		providerInstance = providers.Google(input)
	case "google-workspaces":
		providerInstance = providers.GoogleWorkspaces(input)
	case "okta":
		providerInstance = providers.Okta(input)
	case "linkedin":
		providerInstance = providers.Linkedin(input)
	case "boxy-saml":
		providerInstance = providers.BoxySaml(input)
	default:
		providerInstance = providers.NewProvider(input)
	}

	if input.Override != nil {
		providerInstance = input.Override(providerInstance)
	}

	return providerInstance
}

func discoverOIDCEndpoints(config tpmodels.ProviderConfigForClientType) (tpmodels.ProviderConfigForClientType, error) {
	if config.OIDCDiscoveryEndpoint != "" {
		oidcInfo, err := getOIDCDiscoveryInfo(config.OIDCDiscoveryEndpoint)
		if err != nil {
			return tpmodels.ProviderConfigForClientType{}, err
		}

		if authURL, ok := oidcInfo["authorization_endpoint"].(string); ok {
			if config.AuthorizationEndpoint == "" {
				config.AuthorizationEndpoint = authURL
			}
		}

		if tokenURL, ok := oidcInfo["token_endpoint"].(string); ok {
			if config.TokenEndpoint == "" {
				config.TokenEndpoint = tokenURL
			}
		}

		if userInfoURL, ok := oidcInfo["userinfo_endpoint"].(string); ok {
			if config.UserInfoEndpoint == "" {
				config.UserInfoEndpoint = userInfoURL
			}
		}

		if jwksUri, ok := oidcInfo["jwks_uri"].(string); ok {
			config.JwksURI = jwksUri
		}
	}
	return config, nil
}

// OIDC utils
var oidcInfoMap = map[string]map[string]interface{}{}
var oidcInfoMapLock = sync.Mutex{}

func getOIDCDiscoveryInfo(issuer string) (map[string]interface{}, error) {
	normalizedDomain, err := supertokens.NewNormalisedURLDomain(issuer)
	if err != nil {
		return nil, err
	}
	normalizedPath, err := supertokens.NewNormalisedURLPath(issuer)
	if err != nil {
		return nil, err
	}

	openIdConfigPath, err := supertokens.NewNormalisedURLPath("/.well-known/openid-configuration")
	if err != nil {
		return nil, err
	}

	normalizedPath = normalizedPath.AppendPath(openIdConfigPath)

	if oidcInfo, ok := oidcInfoMap[issuer]; ok {
		return oidcInfo, nil
	}

	oidcInfoMapLock.Lock()
	defer oidcInfoMapLock.Unlock()

	// Check again to see if it was added while we were waiting for the lock
	if oidcInfo, ok := oidcInfoMap[issuer]; ok {
		return oidcInfo, nil
	}

	oidcInfo, err := doGetRequest(normalizedDomain.GetAsStringDangerous()+normalizedPath.GetAsStringDangerous(), nil, nil)
	if err != nil {
		return nil, err
	}
	oidcInfoMap[issuer] = oidcInfo.(map[string]interface{})
	return oidcInfoMap[issuer], nil
}

func doGetRequest(url string, queryParams map[string]interface{}, headers map[string]string) (interface{}, error) {
	if queryParams != nil {
		querystring, err := qs.Marshal(queryParams)
		if err != nil {
			return nil, err
		}
		url = url + "?" + querystring
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("GET request to %s resulted in %d status with body %s", url, resp.StatusCode, string(body))
	}
	return result, nil
}
