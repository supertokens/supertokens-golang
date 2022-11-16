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

func findProvider(options tpmodels.APIOptions, thirdPartyId string, tenantId *string) (tpmodels.TypeProvider, error) {

	for _, provider := range options.Providers {
		if provider.ID == thirdPartyId {
			return provider, nil
		}
	}

	if tenantId == nil {
		return tpmodels.TypeProvider{}, supertokens.BadInputError{Msg: "The third party provider " + thirdPartyId + " seems to be missing from the backend configs."}
	}

	// If tenantId is not nil, we need to create the provider on the fly,
	// so that the GetConfig function will make use of the core to fetch the config
	return createProvider(thirdPartyId), nil
}

func createProvider(thirdPartyId string) tpmodels.TypeProvider {
	switch thirdPartyId {
	case "active-directory":
		return providers.ActiveDirectory(tpmodels.ProviderInput{})
	case "apple":
		return providers.Apple(tpmodels.ProviderInput{})
	case "discord":
		return providers.Discord(tpmodels.ProviderInput{})
	case "facebook":
		return providers.Facebook(tpmodels.ProviderInput{})
	case "github":
		return providers.Github(tpmodels.ProviderInput{})
	case "google":
		return providers.Google(tpmodels.ProviderInput{})
	case "google-workspaces":
		return providers.GoogleWorkspaces(tpmodels.ProviderInput{})
	case "okta":
		return providers.Okta(tpmodels.ProviderInput{})
	case "linkedin":
		return providers.Linkedin(tpmodels.ProviderInput{})
	case "boxyhq":
		// TODO
	}

	return providers.NewProvider(tpmodels.ProviderInput{
		ThirdPartyID: thirdPartyId,
	})
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
