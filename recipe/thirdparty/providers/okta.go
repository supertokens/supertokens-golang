package providers

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Okta(input tpmodels.ProviderInput) *tpmodels.TypeProvider {
	if input.Config.Name == "" {
		input.Config.Name = "Okta"
	}

	oOverride := input.Override

	input.Override = func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		oGetConfig := originalImplementation.GetConfigForClientType
		originalImplementation.GetConfigForClientType = func(clientType *string, userContext supertokens.UserContext) (tpmodels.ProviderConfigForClientType, error) {
			config, err := oGetConfig(clientType, userContext)
			if err != nil {
				return tpmodels.ProviderConfigForClientType{}, err
			}

			if config.OIDCDiscoveryEndpoint == "" {
				config.OIDCDiscoveryEndpoint = fmt.Sprintf("https://%s.okta.com", config.AdditionalConfig["oktaDomain"])
			}

			if len(config.Scope) == 0 {
				config.Scope = []string{"openid", "email"}
			}

			if config.ClientSecret == "" && config.AdditionalConfig["privateKey"] != nil {
				if config.TokenEndpointBodyParams == nil {
					config.TokenEndpointBodyParams = map[string]interface{}{}
				}
				config.TokenEndpointBodyParams["client_assertion_type"] = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
				ca, err := getOktaClientAssertion(config)
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

func getOktaClientAssertion(config tpmodels.ProviderConfigForClientType) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 3600,
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		Audience:  fmt.Sprintf("https://%s.okta.com/oauth2/v1/token", config.AdditionalConfig["oktaDomain"]),
		Subject:   getActualClientIdFromDevelopmentClientId(config.ClientID),
		Issuer:    getActualClientIdFromDevelopmentClientId(config.ClientID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["alg"] = "RS256"

	privateKey := config.AdditionalConfig["privateKey"].(string)

	block, _ := pem.Decode([]byte(privateKey))
	// Check if it's a private key
	if block == nil || block.Type != "PRIVATE KEY" {
		return "", errors.New("failed to decode PEM block containing private key")
	}
	// Get the encoded bytes
	x509Encoded := block.Bytes

	// Now you need an instance of *ecdsa.PrivateKey
	parsedKey, err := x509.ParsePKCS8PrivateKey(x509Encoded) // EDIT to x509Encoded from p8bytes
	if err != nil {
		return "", err
	}

	return token.SignedString(parsedKey)
}
