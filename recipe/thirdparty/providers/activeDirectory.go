package providers

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"golang.org/x/crypto/pkcs12"
)

func ActiveDirectory(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	if input.Config.Name == "" {
		input.Config.Name = "Active Directory"
	}

	oOverride := input.Override

	input.Override = func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := originalImplementation.GetConfigForClientType
		originalImplementation.GetConfigForClientType = func(clientType *string, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			if config.AdditionalConfig == nil || config.AdditionalConfig["directoryId"] == nil {
				if config.OIDCDiscoveryEndpoint == "" {
					return tpmodels.ProviderConfigForClientType{}, fmt.Errorf("Please provide the directoryId in the additionalConfig of the Active Directory provider.")
				}
			} else {
				config.OIDCDiscoveryEndpoint = fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0/.well-known/openid-configuration", config.AdditionalConfig["directoryId"])
			}

			// The config could be coming from core where we didn't add the well-known previously
			config.OIDCDiscoveryEndpoint = normaliseOIDCEndpointToIncludeWellKnown(config.OIDCDiscoveryEndpoint)

			if len(config.Scope) == 0 {
				config.Scope = []string{"openid", "email"}
			}

			if config.ClientSecret == "" && config.AdditionalConfig["certificate"] != nil {
				if config.TokenEndpointBodyParams == nil {
					config.TokenEndpointBodyParams = map[string]interface{}{}
				}
				config.TokenEndpointBodyParams["client_assertion_type"] = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
				ca, err := getADClientAssertion(config)
				if err != nil {
					return tpmodels.ProviderConfigForClientType{}, err
				}
				config.TokenEndpointBodyParams["client_assertion"] = ca
			}

			return config, nil
		}

		if oOverride != nil {
			originalImplementation = oOverride(originalImplementation)
		}
		return originalImplementation
	}

	return NewProvider(input)
}

func getADClientAssertion(config tpmodels.ProviderConfigForClientType) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Audience:  jwt.ClaimStrings{fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", config.AdditionalConfig["directoryId"])},
		Subject:   getActualClientIdFromDevelopmentClientId(config.ClientID),
		Issuer:    getActualClientIdFromDevelopmentClientId(config.ClientID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	thumbBytes, err := hex.DecodeString(config.AdditionalConfig["certificateThumbprint"].(string))
	if err != nil {
		return "", err
	}
	token.Header["x5t"] = base64.StdEncoding.EncodeToString(thumbBytes)
	token.Header["alg"] = "RS256"

	pfxbytes, err := base64.StdEncoding.DecodeString(config.AdditionalConfig["certificate"].(string))
	if err != nil {
		return "", err
	}
	pk, _, err := pkcs12.Decode(pfxbytes, "")
	if err != nil {
		return "", err
	}
	if pk == nil {
		return "", errors.New("private key not found")
	}

	return token.SignedString(pk)
}
