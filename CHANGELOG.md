# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [unreleased]

- `session.CreateNewSession` now defaults to the value of the `st-auth-mode` header (if available) if the configured `config.GetTokenTransferMethod` returns `any`.

## [0.17.5] - 2024-03-14
- Adds a type uint64 to the `accessTokenCookiesExpiryDurationMillis` local variable in `recipe/session/utils.go`. It also removes the redundant `uint64` type forcing needed because of the untyped variable.
- Fixes the passing of `tenantId` in `getAllSessionHandlesForUser` and `revokeAllSessionsForUser` based on `fetchAcrossAllTenants` and `revokeAcrossAllTenants` inputs respectively.
- Updated fake email generation

## [0.17.4] - 2024-02-08

- Adds `TLSConfig` to SMTP settings.
- `TLSConfig` is always passed to gomail so that it can be used when gomail uses `STARTTLS` to upgrade the connection to TLS. - https://github.com/supertokens/supertokens-golang/issues/392
- Not setting `InsecureSkipVerify` to `true` in the SMTP settings because it is not recommended to use it in production.

## [0.17.3] - 2023-12-12

- CI/CD changes

## [0.17.2] - 2023-12-06

- Updates LinkedIn OAuth implementation as per the latest [changes](https://learn.microsoft.com/en-us/linkedin/consumer/integrations/self-serve/sign-in-with-linkedin-v2?context=linkedin%2Fconsumer%2Fcontext#authenticating-members).

## [0.17.1] - 2023-11-24

### Added

-   Adds support for configuring multiple frontend domains to be used with the same backend
-   Added new `Origin` and `GetOrigin` properties to `AppInfo`, this can be configured to allow you to conditionally return the value of the frontend domain. This property will replace `WebsiteDomain` in a future release of `supertokens-golang`
-   `WebsiteDomain` inside `AppInfo` is now optional. Using `Origin` or `GetOrigin` is recommended over using `WebsiteDomain`. This is not a breaking change and using `WebsiteDomain` will continue to work.

## [0.17.0] - 2023-11-14

### Breaking change

-   Fixes github user id conversion to be consistent with other SDKs.

### Migration

If you were using the SDK Versions >= `0.13.0` and < `0.17.0`, use the following override function for github:

```go
{
	Config: tpmodels.ProviderConfig{
		ThirdPartyId: "github",
		// other config
	},
	Override: func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		originalGetUserInfo := originalImplementation.GetUserInfo
		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			userInfo, err := originalGetUserInfo(oAuthTokens, userContext)
			if err != nil {
				return userInfo, err
			}
			number, err := strconv.ParseFloat(userInfo.ThirdPartyUserId, 64)
			if err != nil {
				return userInfo, err
			}
			userInfo.ThirdPartyUserId = fmt.Sprint(number)
			return userInfo, nil
		}

		return originalImplementation
	},
},
```

If you were using the SDK Versions < `0.13.0`, use the following override function for github:

```go
{
	Config: tpmodels.ProviderConfig{
		ThirdPartyId: "github",
		// other config
	},
	Override: func(originalImplementation *tpmodels.TypeProvider) *tpmodels.TypeProvider {
		originalGetUserInfo := originalImplementation.GetUserInfo
		originalImplementation.GetUserInfo = func(oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
			userInfo, err := originalGetUserInfo(oAuthTokens, userContext)
			if err != nil {
				return userInfo, err
			}
			number, err := strconv.ParseFloat(userInfo.ThirdPartyUserId, 64)
			if err != nil {
				return userInfo, err
			}
			userInfo.ThirdPartyUserId = fmt.Sprintf("%f", number)
			return userInfo, nil
		}

		return originalImplementation
	},
},
```

## [0.16.6] - 2023-11-3

-   Added `NetworkInterceptor` to the `ConnectionInfo` config.
  - This can be used to capture/modify all the HTTP requests sent to the core.
  - Solves the issue - https://github.com/supertokens/supertokens-core/issues/865

## [0.16.5] - 2023-11-1

-   Adds `debug` flag to the `TypeInput`. If set to `true`, debug logs will be printed.

## [0.16.4] - 2023-10-20

-   Fixes an issue where sometimes the `Access-Control-Expose-Headers` header value would contain duplicates

## [0.16.3] - 2023-10-19

-   Fixes an issue where trying to view details of a third party user would show an error for the user not being found when using the thirdpartypasswordless recipe

## [0.16.2] - 2023-10-17

-   Fixes an issue where tenant ids returned for a user from the user get API of the dashboard recipe would always be nil for thirdpartyemailpassword and thirdpartypasswordless recipes

## [0.16.1] - 2023-10-03

### Changes

-   Added `ValidateAccessToken` to the configuration for social login providers, this function allows you to verify the access token returned by the social provider. If you are using Github as a provider, there is a default implementation provided for this function.

### Fixes

-   Fixes `timeJoined` casing in emailpassword and passwordless user objects.

## [0.16.0] - 2023-09-27

### Fixes

-   Solves issue with clock skew during third party sign-in/up - https://github.com/supertokens/supertokens-golang/issues/362
    -   Bumped `github.com/golang-jtw/jwt` version from v4 to v5.
    -   Bumped `github.com/MicahParks/keyfunc` version from v1 to v2.

### Breaking Changes

-   Minimum golang version supported is 1.18

## [0.15.0] - 2023-09-26

-   Adds Twitter/X as a default provider to the third party recipe
-   Added a `Cache-Control` header to `/jwt/jwks.json` (`GetJWKSGET`)
-   Added `ValidityInSeconds` to the return value of the overrideable `GetJWKS` function.
    -   This can be used to control the `Cache-Control` header mentioned above.
    -   It defaults to `60` or the value set in the cache-control header returned by the core
    -   This is optional (so you are not required to update your overrides). Returning undefined means that the header is not set.
-   Handle AWS Public URLs (ending with `.amazonaws.com`) separately while extracting TLDs for SameSite attribute.
-   Return `500` status instead of panic when `supertokens.Middleware` is used without initializing the SDK.
-   Updates fiber adaptor package in the fiber example.

## [0.14.0] - 2023-09-11

### Added

-   The Dashboard recipe now accepts a new `Admins` property which can be used to give Dashboard Users write privileges for the user dashboard.

### Changes

-   Dashboard APIs now return a status code `403` for all non-GET requests if the currently logged in Dashboard User is not listed in the `admins` array
-   Now ignoring protected props in the payload in `CreateNewSession` and `CreateNewSessionWithoutRequestResponse`

## [0.13.2] - 2023-08-28

-   Adds logic to retry network calls if the core returns status 429

## [0.13.1] - 2023-08-24

-   Fixes login methods API to return empty provider array instead of `null`
-   Fixes thirdpartypasswordless initialization when there are no static providers configured

## [0.13.0] - 2023-08-07

### Added

-   Added Multitenancy Recipe & always initialized by default.
-   Adds Multitenancy support to all the recipes
-   Added new Social login providers - LinkedIn
-   Added new Multi-tenant SSO providers - Okta, Active Directory, Boxy SAML
-   All APIs handled by Supertokens middleware can have an optional `tenantId` prefixed in the path. e.g. <basePath>/<tenantId>/signinup
-   Following recipe functions have been added:
    -   `emailpassword.CreateResetPasswordLink`
    -   `emailpassword.SendResetPasswordEmail`
    -   `emailverification.CreateEmailVerificationLink`
    -   `emailverification.SendEmailVerificationEmail`
    -   `thirdparty.GetProvider`
    -   `thirdpartyemailpassword.ThirdPartyGetProvider`
    -   `thirdpartyemailpassword.CreateResetPasswordLink`
    -   `thirdpartyemailpassword.SendResetPasswordEmail`
    -   `thirdpartypasswordless.ThirdPartyGetProvider`

### Breaking changes

-   Only supporting FDI 1.17
-   Core must be upgraded to 6.0
-   For consistency, all `UnknownUserIDError` have been renamed to `UnknownUserIdError`
-   `getUsersOldestFirst` & `getUsersNewestFirst` has mandatory parameter `tenantId`. Pass `'public'` if not using multitenancy.
-   Added mandatory field `tenantId` to `EmailDeliveryInterface` and `SmsDeliveryInterface`. Pass `'public'` if not using multitenancy.
-   Removed deprecated config `createAndSendCustomEmail` and `createAndSendCustomTextMessage`.
-   EmailPassword recipe changes:
    -   Added mandatory `tenantId` field to `TypeEmailPasswordPasswordResetEmailDeliveryInput`
    -   Removed `resetPasswordUsingTokenFeature` from `TypeInput`
    -   Added `tenantId` param to `validate` function in `TypeInputFormField`
    -   Added mandatory `tenantId` as first parameter to the following recipe index functions:
        -   `SignUp`
        -   `SignIn`
        -   `GetUserByEmail`
        -   `CreateResetPasswordToken`
        -   `ResetPasswordUsingToken`
    -   Added mandatory `tenantId` in the input for the following recipe interface functions. If any of these functions are overridden, they need to be updated accordingly:
        -   `SignUp`
        -   `SignIn`
        -   `GetUserByEmail`
        -   `CreateResetPasswordToken`
        -   `ResetPasswordUsingToken`
        -   `UpdateEmailOrPassword`
    -   Added mandatory `tenantId` in the input for the following API interface functions. If any of these functions are overridden, they need to be updated accordingly:
        -   `EmailExistsGET`
        -   `GeneratePasswordResetTokenPOST`
        -   `PasswordResetPOST`
        -   `SignInPOST`
        -   `SignUpPOST`
-   EmailVerification recipe changes:
    -   Added mandatory `TenantId` field to `EmailVerificationType`
    -   Added mandatory `tenantId` as first parameter to the following recipe index functions:
        -   `CreateEmailVerificationToken`
        -   `VerifyEmailUsingToken`
        -   `RevokeEmailVerificationTokens`
    -   Added mandatory `tenantId` in the input for the following recipe interface functions. If any of these functions are overridden, they need to be updated accordingly:
        -   `CreateEmailVerificationToken`
        -   `VerifyEmailUsingToken`
        -   `RevokeEmailVerificationTokens`
    -   Added mandatory `tenantId` in the input for the following API interface functions. If any of these functions are overridden, they need to be updated accordingly:
        -   `VerifyEmailPOST`
-   Passwordless recipe changes:
    -   Added `tenantId` param to `ValidateEmailAddress`, `ValidatePhoneNumber` and `GetCustomUserInputCode` functions in `TypeInput`
    -   Added mandatory `TenantId` field to `emaildelivery.PasswordlessLoginType` and `smsdelivery.PasswordlessLoginType`
    -   Added mandatory `tenantId` in the input to the following recipe index functions:
        -   `CreateCodeWithEmail`
        -   `CreateCodeWithPhoneNumber`
        -   `CreateNewCodeForDevice`
        -   `ConsumeCodeWithUserInputCode`
        -   `ConsumeCodeWithLinkCode`
        -   `GetUserByEmail`
        -   `GetUserByPhoneNumber`
        -   `RevokeAllCodesByEmail`
        -   `RevokeAllCodesByPhoneNumber`
        -   `RevokeCode`
        -   `ListCodesByEmail`
        -   `ListCodesByPhoneNumber`
        -   `ListCodesByDeviceID`
        -   `ListCodesByPreAuthSessionID`
        -   `CreateMagicLinkByEmail`
        -   `CreateMagicLinkByPhoneNumber`
        -   `SignInUpByEmail`
        -   `SignInUpByPhoneNumber`
    -   Added mandatory `tenantId` in the input for the following recipe interface functions. If any of these functions are overridden, they need to be updated accordingly:
        -   `CreateCode`
        -   `CreateNewCodeForDevice`
        -   `ConsumeCode`
        -   `GetUserByEmail`
        -   `GetUserByPhoneNumber`
        -   `RevokeAllCodes`
        -   `RevokeCode`
        -   `ListCodesByEmail`
        -   `ListCodesByPhoneNumber`
        -   `ListCodesByDeviceID`
        -   `ListCodesByPreAuthSessionID`
    -   Added mandatory `tenantId` in the input for the following API interface functions. If any of these functions are overridden, they need to be updated accordingly:
        -   `CreateCodePOST`
        -   `ResendCodePOST`
        -   `ConsumeCodePOST`
        -   `EmailExistsGET`
        -   `PhoneNumberExistsGET`
-   ThirdParty recipe changes
    -   The providers array in `SignInUpFeature` accepts `[]ProviderInput` instead of `[]TypeProvider`. TypeProvider interface is re-written. Refer migration section for more info.
    -   Removed `SignInUp` and added `ManuallyCreateOrUpdateUser` instead in the recipe index functions.
    -   Added `ManuallyCreateOrUpdateUser` to recipe interface which is being called by the function mentioned above.
        -   `ManuallyCreateOrUpdateUser` recipe interface function should not be overridden as it is not going to be called by the SDK in the sign in/up flow.
        -   `SignInUp` recipe interface functions is not removed and is being used by the sign in/up flow.
    -   Added mandatory `tenantId` as first parameter to the following recipe index functions:
        -   `GetUsersByEmail`
        -   `GetUserByThirdPartyInfo`
    -   Added mandatory `tenantId` in the input for the following recipe interface functions. If any of these functions are overridden, they need to be updated accordingly:
        -   `GetUsersByEmail`
        -   `GetUserByThirdPartyInfo`
        -   `SignInUp`
    -   Added mandatory `tenantId` in the input for the following API interface functions. If any of these functions are overridden, they need to be updated accordingly:
        -   `AuthorisationUrlGET`
        -   `SignInUpPOST`
    -   Updated `SignInUp` recipe interface function in thirdparty with new parameters:
        -   `oAuthTokens` - contains all the tokens (access_token, id_token, etc.) as returned by the provider
        -   `rawUserInfoFromProvider` - contains all the user profile info as returned by the provider
    -   Updated `AuthorisationUrlGET` API
        -   Changed: Doesn't accept `clientId` anymore and accepts `clientType` instead to determine the matching config
        -   Added: optional `pkceCodeVerifier` in the response, to support PKCE
    -   Updated `SignInUpPOST` API
        -   Removed: `clientId`, `redirectURI`, `authCodeResponse` and `code` from the input
        -   Instead,
            -   accepts `clientType` to determine the matching config
            -   One of redirectURIInfo (for code flow) or oAuthTokens (for token flow) is required
    -   Updated `AppleRedirectHandlerPOST`
        -   to accept all the form fields instead of just the code
        -   to use redirect URI encoded in the `state` parameter instead of using the websiteDomain config.
        -   to use HTTP 303 instead of javascript based redirection.
-   Session recipe changes
    -   Added mandatory `tenantId` parameter to the following recipe index functions:
        -   `CreateNewSession`
        -   `CreateNewSessionWithoutRequestResponse`
        -   `ValidateClaimsInJWTPayload`
    -   Added mandatory `tenantId` in the input for the following recipe interface functions. If any of these functions are overridden, they need to be updated accordingly:
        -   `CreateNewSession`
        -   `GetGlobalClaimValidators`
    -   Added `tenantId` and `revokeAcrossAllTenants` params to `RevokeAllSessionsForUser` in the recipe interface.
    -   Added `tenantId` and `fetchAcrossAllTenants` params to `GetAllSessionHandlesForUser` in the recipe interface.
    -   Added `GetTenantId` function to `TypeSessionContainer`
    -   Added `tenantId` to `FetchValue` function in `PrimitiveClaim`, `PrimitiveArrayClaim`.
-   UserRoles recipe changes
    -   Added mandatory `tenantId` as first parameter to the following recipe index functions:
        -   `AddRoleToUser`
        -   `RemoveUserRole`
        -   `GetRolesForUser`
        -   `GetUsersThatHaveRole`
    -   Added mandatory `tenantId` in the input for the following recipe interface functions. If any of these functions are overridden, they need to be updated accordingly:
        -   `AddRoleToUser`
        -   `RemoveUserRole`
        -   `GetRolesForUser`
        -   `GetUsersThatHaveRole`
-   Similar changes in combination recipes (thirdpartyemailpassword and thirdpartypasswordless) have been made
-   Even if thirdpartyemailpassword and thirdpartpasswordless recipes do not have a providers array as an input, they will still expose the third party recipe routes to the frontend.
-   Returns 400 status code in emailpassword APIs if the input email or password are not of type string.

### Changes

-   Recipe function changes:
    -   Added optional `tenantIdForPasswordPolicy` param to `emailpassword.UpdateEmailOrPassword`, `thirdpartyemailpassword.UpdateEmailOrPassword`
    -   Added optional param `tenantId` to `session.RevokeAllSessionsForUser`. If tenantId is nil, sessions are revoked across all tenants
    -   Added optional param `tenantId` to `session.getAllSessionHandlesForUser`. If tenantId is nil, sessions handles across all tenants are returned
-   Adds optional param `tenantId` to `GetUserCount` which returns total count across all tenants if not passed.
-   Adds protected prop `tId` to the accessToken payload
-   Adds `IncludesAny` claim validator to `PrimitiveArrayClaim`

### Migration

-   To call any recipe function that has `tenantId` added to it, pass `'public`'

    Before:

    ```ts
    emailpassword.SignUp("test@example.com", "password");
    ```

    After:

    ```ts
    emailpassword.SignUp("public", "test@example.com", "password");
    ```

-   Input for provider array change as follows:

    Before:

    ```go
    thirdparty.Google(tpmodels.GoogleConfig{
        ClientID: "...",
        ClientSecret: "...",
    })
    ```

    After:

    ```go
    tpmodels.ProviderInput{
        ThirdPartyId: "google"
        Config: tpmodels.ProviderConfig{
            Clients: []tpmodels.ProviderClientConfig{
                {
                    ClientID:     "...",
                    ClientSecret: "...",
                },
            },
        },
    }
    ```

-   Single instance with multiple clients of each provider instead of multiple instances of them. Also use `clientType` to differentiate them. `clientType` passed from the frontend will be used to determine the right config. `IsDefault` option has been removed and `clientType` is expected to be passed when there are more than one client. If there is only one client, `clientType` is optional and will be used by default.

    Before:

    ```go
    []tpmodels.TypeProvider{
        thirdparty.Google(tpmodels.GoogleConfig{
            IsDefault: true,
            ClientID: "clientid1",
            ClientSecret: "...",
        }),
        thirdparty.Google(tpmodels.GoogleConfig{
            ClientID: "clientid2",
            ClientSecret: "...",
        }),
    }
    ```

    After:

    ```go
    []tpmodels.ProviderInput{
        {
            ThirdPartyId: "google",
            Config: tpmodels.ProviderConfig{
                Clients: []tpmodels.ProviderClientConfig{
                    {
                        ClientType: "web",
                        ClientID:     "clientid1",
                        ClientSecret: "...",
                    },
                    {
                        ClientType: "android",
                        ClientID:     "clientid2",
                        ClientSecret: "...",
                    },
                },
            },
        },
    }
    ```

-   Change in the implementation of custom providers

    -   All config is part of `ProviderInput`
    -   To provide implementation for `GetProfileInfo`
        -   either use `UserInfoEndpoint`, `UserInfoEndpointQueryParams` and `UserInfoMap` to fetch the user info from the provider
        -   or specify custom implementation in an override for `GetUserInfo` (override example in the next section)

    Before:

    ```go
    tpmodels.TypeProvider{
        ID: "custom",
        Get: func(redirectURI, authCodeFromRequest *string, userContext supertokens.UserContext) tpmodels.TypeProviderGetResponse {
            return tpmodels.TypeProviderGetResponse{
                AccessTokenAPI: "...",
                AuthorisationRedirect: tpmodels.AuthorisationRedirect{
                    URL:    "...",
                    Params: map[string]interface{}{},
                },
                GetClientId: func(userContext supertokens.UserContext) string {
                    return "..."
                },
                GetRedirectURI: func(userContext supertokens.UserContext) (string, error) {
                    return "...", nil
                },
                GetProfileInfo: func(authCodeResponse interface{}, userContext supertokens.UserContext) (tpmodels.UserInfo, error) {
                    return tpmodels.UserInfo{
                        ID: "...",
                        Email: &tpmodels.EmailStruct{
                            ID:         "...",
                            IsVerified: true,
                        },
                    }, nil
                },
            }
        },
    }
    ```

    After:

    ```go
    tpmodels.ProviderInput{
        ThirdPartyID: "custom",
        Config: tpmodels.ProviderConfig{
            Clients: []tpmodels.ProviderClientConfig{
                {
                    ClientID:     "...",
                    ClientSecret: "...",
                },
            },
            AuthorizationEndpoint:            "...",
            AuthorizationEndpointQueryParams: map[string]interface{}{},
            TokenEndpoint:                    "...",
            TokenEndpointBodyParams:          map[string]interface{}{},
            UserInfoEndpoint:                 "...",
            UserInfoEndpointQueryParams:      map[string]interface{}{},
            UserInfoMap: tpmodels.TypeUserInfoMap{
                FromUserInfoAPI: tpmodels.TypeUserInfoFields{
                    UserId:        "id",
                    Email:         "email",
                    EmailVerified: "email_verified",
                },
            },
        },
    }
    ```

    Also, if the custom provider supports openid, it can automatically discover the endpoints

    ```go
    tpmodels.ProviderInput{
        ThirdPartyID: "custom",
        Config: tpmodels.ProviderConfig{
            Clients: []tpmodels.ProviderClientConfig{
                {
                    ClientID:     "...",
                    ClientSecret: "...",
                },
            },
            OIDCDiscoveryEndpoint: "...",
            UserInfoMap: tpmodels.TypeUserInfoMap{
                FromUserInfoAPI: tpmodels.TypeUserInfoFields{
                    UserId:        "id",
                    Email:         "email",
                    EmailVerified: "emailVerified",
                },
            },
        },
    }
    ```

    Note: The SDK will fetch the oauth2 endpoints from the provider's OIDC discovery endpoint. No need to `/.well-known/openid-configuration` to the `oidcDiscoveryEndpoint` config. For eg. if `oidcDiscoveryEndpoint` is set to `"https://accounts.google.com/"`, the SDK will fetch the endpoints from `"https://accounts.google.com/.well-known/openid-configuration"`

-   Any of the functions in the TypeProvider can be overridden for custom implementation

    -   Overrides can do the following:
        -   update params, headers dynamically for the authorization redirect url or in the exchange of code to tokens
        -   add custom logic to exchange code to tokens
        -   add custom logic to get the user info

    ```go
    tpmodels.ProviderInput{
        ThirdPartyID: "custom",
        Config: tpmodels.ProviderConfig{
            Clients: []tpmodels.ProviderClientConfig{
                {
                    ClientID:     "...",
                    ClientSecret: "...",
                },
            },
            OIDCDiscoveryEndpoint: "...",
            UserInfoMap: tpmodels.TypeUserInfoMap{
                FromUserInfoAPI: tpmodels.TypeUserInfoFields{
                    UserId:        "id",
                    Email:         "email",
                    EmailVerified: "emailVerified",
                },
            },
        },
        Override: func(provider *tpmodels.TypeProvider) *tpmodels.TypeProvider {
            oGetAuthorisationRedirectURL := provider.GetAuthorisationRedirectURL
            provider.GetAuthorisationRedirectURL = func(config tpmodels.ProviderConfigForClientType, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeAuthorisationRedirect, error) {
                resp, err := oGetAuthorisationRedirectURL(config, redirectURIOnProviderDashboard, userContext)
                // your logic here
                return resp, err
            }

            oExchangeAuthCodeForOAuthTokens := provider.ExchangeAuthCodeForOAuthTokens
            provider.ExchangeAuthCodeForOAuthTokens = func(config tpmodels.ProviderConfigForClientType, code string, redirectURIOnProviderDashboard string, userContext supertokens.UserContext) (tpmodels.TypeOAuthResponse, error) {
                resp, err := oExchangeAuthCodeForOAuthTokens(config, code, redirectURIOnProviderDashboard, userContext)
                // your logic here
                return resp, err
            }

            oGetUserInfo := provider.GetUserInfo
            provider.GetUserInfo = func(config tpmodels.ProviderConfigForClientType, oAuthTokens tpmodels.TypeOAuthTokens, userContext supertokens.UserContext) (tpmodels.TypeUserInfo, error) {
                resp, err := oGetUserInfo(config, oAuthTokens, userContext)
                // your logic here
                return resp, err
            }

            return provider
        },
    }
    ```

-   To get access token and raw user info from the provider, override the signInUp function

    ```go
    thirdparty.Init(&tpmodels.TypeInput{
        Override: &tpmodels.OverrideStruct{
            Functions: func(originalImplementation tpmodels.RecipeInterface) tpmodels.RecipeInterface {
                oSignInUp := *originalImplementation.SignInUp
                nSignInUp := func(thirdPartyID string, thirdPartyUserID string, email string, oAuthTokens tpmodels.TypeOAuthTokens, rawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider, tenantId string, userContext supertokens.UserContext) (tpmodels.SignInUpResponse, error) {
                    resp, err := oSignInUp(thirdPartyID, thirdPartyUserID, email, oAuthTokens, rawUserInfoFromProvider, tenantId, userContext)
                    // resp.OK.OAuthTokens["access_token"]
                    // resp.OK.OAuthTokens["id_token"]
                    // resp.OK.RawUserInfoFromProvider.FromUserInfoAPI
                    // resp.OK.RawUserInfoFromProvider.FromIdTokenPayload

                    return resp, err
                }
                *originalImplementation.SignInUp = nSignInUp

                return originalImplementation
            },
        },
    })
    ```

-   Request body of thirdparty signinup API has changed

    -   If using auth code:

        Before:

        ```json
        {
            "thirdPartyId": "...",
            "clientId": "...",
            "redirectURI": "...", // optional
            "code": "..."
        }
        ```

        After:

        ```json
        {
            "thirdPartyId": "...",
            "clientType": "...",
            "redirectURIInfo": {
                "redirectURIOnProviderDashboard": "...", // required
                "redirectURIQueryParams": {
                    "code": "...",
                    "state": "..."
                    // ... all callback query params
                },
                "pkceCodeVerifier": "..." // optional, use this if using PKCE flow
            }
        }
        ```

    -   If using tokens:

        Before:

        ```json
        {
            "thirdPartyId": "...",
            "clientId": "...",
            "redirectURI": "...",
            "authCodeResponse": {
                "access_token": "...", // required
                "id_token": "..."
            }
        }
        ```

        After:

        ```json
        {
            "thirdPartyId": "...",
            "clientType": "...",
            "oAuthTokens": {
                "access_token": "...", // now optional
                "id_token": "..."
                // rest of the oAuthTokens as returned by the provider
            }
        }
        ```

### SDK and core compatibility

-   Compatible with Core>=6.0.0 (CDI 4.0)
-   Compatible with frontend SDKs:
    -   supertokens-auth-react@0.34.0
    -   supertokens-web-js@0.7.0
    -   supertokens-website@17.0.2

## [0.12.10] - 2023-07-31

- Fixes error handling with regenerate access token when the access token of the session is revoked.
- Fixes payload in get session when the  access token version <= 2

## [0.12.9] - 2023-07-26

- Fixes an issue where updating the user's password from the user management dashboard would result in a crash when using the thirdpartyemailpassword recipe (https://github.com/supertokens/supertokens-golang/issues/311)

## [0.12.8] - 2023-07-10

- Adds additional tests for session verification

### Fixes

- Now properly ignoring missing anti-csrf tokens in optional session validation

## [0.12.7] - 2023-06-05

### Fixes

- Update email templates to fix an issue with styling on some email clients
- Fixes an issue where session verification would fail when using JWTs created by the JWT recipe (and not the session recipe)

## [0.12.6] - 2023-06-01

### Fixes

- Fixes a bug in the session recipe where the SDK would try to fetch the JWKs from the core multiple times per minute

## [0.12.5] - 2023-05-26

### Fixes

-   Fixes bug in debug logging where line number was being printed incorrectly.

## [0.12.4] - 2023-05-23

### Changes

-   Added a new `GetRequestFromUserContext` function that can be used to read the original network request from the user context in overridden APIs and recipe functions

## [0.12.3] - 2023-05-22

### Added

-   Adds additional debug logs whenever the SDK returns a `TryRefreshTokenError` or `UnauthorizedError` to make debugging easier

## [0.12.2] - 2023-05-19

-   Adds additional tests for the session recipe
-   Fixes https://github.com/supertokens/supertokens-golang/issues/284

## [0.12.1] - 2023-05-12

### Changes

-   Made the access token string optional in the overrideable `GetSession` function
-   Moved checking if the access token is nil into the overrideable `GetSession` function

## [0.12.0] - 2023-05-05

### Added

- added optional password policy check in `updateEmailOrPassword`

### Breaking Changes

-   Changed the interface and configuration of the Session recipe, see below for details. If you do not use the Session recipe directly and do not provide custom configuration, then no migration is necessary.
-   Renamed `GetSessionData` to `GetSessionDataInDatabase` to clarify that it always hits the DB
-   Renamed `GetSessionDataWithContext` to `GetSessionDataInDatabaseWithContext` to clarify that it always hits the DB
-   Renamed `UpdateSessionData` to `UpdateSessionDataInDatabase`
-   Renamed `UpdateSessionDataWithContext` to `UpdateSessionDataInDatabaseWithContext` to clarify that it always hits the DB
-   Renamed `SessionData` to `SessionDataInDatabase` in `SessionInformation`
-   Renamed `sessionData` to `sessionDataInDatabase` in the input to `CreateNewSession`
-   Added `useStaticSigningKey` to `CreateJWT` and `CreateJWTWithContext`
-   Added support for CDI version `2.21`
-   Dropped support for CDI version `2.8`-`2.20`
-   `GetAccessTokenPayload` will now return standard (`sub`, `iat`, `exp`) claims and some SuperTokens specific claims along the user defined ones in `GetAccessTokenPayload`.
-   Some claim names are now prohibited in the root level of the access token payload
    -   They are: `sub`, `iat`, `exp`, `sessionHandle`, `parentRefreshTokenHash1`, `refreshTokenHash1`, `antiCsrfToken`
    -   If you used these in the root level of the access token payload, then you'll need to migrate your sessions or they will be logged out during the next refresh
    -   These props should be renamed (e.g., by adding a prefix) or moved inside an object in the access token payload
    -   You can migrate these sessions by updating their payload to match your new structure, by calling `MergeIntoAccessTokenPayload`
-   New access tokens are valid JWTs now
    -   They can be used directly (i.e.: by calling `GetAccessToken` on the session) if you need a JWT
    -   The `jwt` prop in the access token payload is removed
-   JWT and OpenId related configuration has been removed from the Session recipe config. If necessary, they can be added by initializing the OpenId recipe before the Session recipe.
-   Changed the Session recipe interface - CreateNewSession, GetSession and RefreshSession overrides now do not take response and request and return status instead of throwing
-   Renamed `AccessTokenPayload` to `CustomClaimsInAccessTokenPayload` in `SessionInformation` (the return value of `GetSessionInformation`). This reflects the fact that it doesn't contain some default claims (`sub`, `iat`, etc.)

### Changed

-   Refactors the URL for the JWKS endpoint exposed by SuperTokens core
-   Added new optional `useStaticSigningKey` param to `CreateJWT`
-   The Session recipe now always initializes the OpenID recipe if it hasn't been initialized.
-   Refactored how access token validation is done
-   Added support for new access token version
-   Removed the handshake call to improve start-up times
-   Removed `GetAccessTokenLifeTimeMS` and `GetRefreshTokenLifeTimeMS` functions
-   Added `ExposeAccessTokenToFrontendInCookieBasedAuth` (defaults to `false`) option to the Session recipe config
-   Added new `checkDatabase` param to `VerifySession` and `GetSession`
-   Removed deprecated `UpdateAccessTokenPayload`, `UpdateAccessTokenPayloadWithContext`, `RegenerateAccessToken` and `RegenerateAccessTokenWithContext` from the Session recipe interface
-   Added `CreateNewSessionWithoutRequestResponse`, `CreateNewSessionWithContextWithoutRequestResponse`, `GetSessionWithoutRequestResponse`, `GetSessionWithContextWithoutRequestResponse`, `RefreshSession`, `RefreshSessionWithContextWithoutRequestResponse` to the Session recipe.
-   Added `GetAllSessionTokensDangerously` to session objects (`SessionContainer`)
-   Added `AttachToRequestResponse` to session objects (`SessionContainer`)

### Migration

#### If self-hosting core

1. You need to update the core version
2. There are manual migration steps needed. Check out the core changelogs for more details.

#### If you used the jwt feature of the session recipe

1. Add `ExposeAccessTokenToFrontendInCookieBasedAuth: true` to the Session recipe config on the backend if you need to access the JWT on the frontend.
2. Choose a prop from the following list. We'll use `sub` in the code below, but you can replace it with another from the list if you used it in a custom access token payload.
    - `sub`
    - `iat`
    - `exp`
    - `sessionHandle`
3. On the frontend where you accessed the JWT before by: `(await Session.getAccessTokenPayloadSecurely()).jwt` update to:

```tsx
let jwt = null;
const accessTokenPayload = await Session.getAccessTokenPayloadSecurely();
if (accessTokenPayload.sub !== undefined) {
    jwt = await Session.getAccessToken();
} else {
    // This branch is only required if there are valid access tokens created before the update
    // It can be removed after the validity period ends
    jwt = accessTokenPayload.jwt;
}
```

4. On the backend if you accessed the JWT before by `session.GetAccessTokenPayload()["jwt"]` please update to:

```go
var jwt string
accessTokenPayload := session.GetAccessTokenPayload();
if (accessTokenPayload["sub"] != nil) {
    jwt =  session.GetAccessToken();
} else {
    // This branch is only required if there are valid access tokens created before the update
    // It can be removed after the validity period ends
    jwt = accessTokenPayload["jwt"].(string);
}
```

#### If you used to set an issuer in the session recipe `Jwt` configuration

-   You can add an issuer claim to access tokens by overriding the `CreateNewSession` function in the session recipe init.
    -   Check out https://supertokens.com/docs/passwordless/common-customizations/sessions/claims/access-token-payload#during-session-creation for more information
-   You can add an issuer claim to JWTs created by the JWT recipe by passing the `iss` claim as part of the payload.
-   You can set the OpenId discovery configuration as follows:

Before:

```go
func main() {
    supertokens.Init(supertokens.TypeInput{
        RecipeList: []supertokens.Recipe{
            session.Init(&sessmodels.TypeInput{
                Jwt: &sessmodels.JWTInputConfig{
                    Enable: true,
					Issuer: "...",
                },
            }),
        },
    })
}
```

After:

```go
func main() {
    supertokens.Init(supertokens.TypeInput{
	RecipeList: []supertokens.Recipe{
		session.Init(&sessmodels.TypeInput{
			GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
				return sessmodels.HeaderTransferMethod
			},
			Override: &sessmodels.OverrideStruct{
				OpenIdFeature: &openidmodels.OverrideStruct{
					Functions: func(originalImplementation openidmodels.RecipeInterface) openidmodels.RecipeInterface {
						(*originalImplementation.GetOpenIdDiscoveryConfiguration) = func(userContext *map[string]interface{}) (openidmodels.GetOpenIdDiscoveryConfigurationResponse, error) {
							return openidmodels.GetOpenIdDiscoveryConfigurationResponse{
								OK: &struct{Issuer string; Jwks_uri string}{
									Issuer:   "your issuer",
									Jwks_uri: "https://your.api.domain/auth/jwt/jwks.json",
								},
							}, nil
						}

						return originalImplementation
					},
				},
			},
		}),
	},
})
}
```

#### If you used `sessionData` (not `accessTokenPayload`)

Related functions/prop names have changes (`sessionData` became `sessionDataFromDatabase`):

-   Renamed `GetSessionData` to `GetSessionDataFromDatabase` to clarify that it always hits the DB
-   Renamed `UpdateSessionData` to `UpdateSessionDataInDatabase`
-   Renamed `sessionData` to `sessionDataInDatabase` in `SessionInformation` and the input to `CreateNewSession`

#### If you used to set `access_token_blacklisting` in the core config

-   You should now set `CheckDatabase` to true in the verifySession params.

#### If you used to set `access_token_signing_key_dynamic` in the core config

-   You should now set `useDynamicAccessTokenSigningKey` in the Session recipe config.

#### If you used to use standard/protected props in the access token payload root:

1. Update you application logic to rename those props (e.g., by adding a prefix)
2. Update the session recipe config (in this example `sub` is the protected property we are updating by adding the `app` prefix):

Before:

```go
func main() {
    supertokens.Init(supertokens.TypeInput{
	RecipeList: []supertokens.Recipe{
		session.Init(&sessmodels.TypeInput{
			Override: &sessmodels.OverrideStruct{
				Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
					originalCreateNewSession := *originalImplementation.CreateNewSession
					(*originalImplementation.CreateNewSession) = func(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload, sessionData map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
						if accessTokenPayload == nil {
							accessTokenPayload = map[string]interface{}{}
						}

						accessTokenPayload["sub"] = userID + "!!!"

						return originalCreateNewSession(req, res, userID, accessTokenPayload, sessionData, userContext)
					}

					return originalImplementation
				},
			},
		}),
	},
})
}
```

After:

```go
func main() {
    supertokens.Init(supertokens.TypeInput{
	RecipeList: []supertokens.Recipe{
		session.Init(&sessmodels.TypeInput{
			Override: &sessmodels.OverrideStruct{
				Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
					originalGetSession := *originalImplementation.GetSession

					(*originalImplementation.GetSession) = func(req *http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions, userContext *map[string]interface{}) (*sessmodels.TypeSessionContainer, error) {
						result, err := originalGetSession(req, res, options, userContext)

						if result != nil {
							originalPayload := result.GetAccessTokenPayload()

							if originalPayload["appSub"] == nil {
								result.MergeIntoAccessTokenPayload(map[string]interface{}{
									"appSub": originalPayload["sub"],
									"sub": nil,
								})
							}
						}

						return result, err
					}

					originalCreateNewSession := *originalImplementation.CreateNewSession


					(*originalImplementation.CreateNewSession) = func(req *http.Request, res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, userContext *map[string]interface{}) (*sessmodels.TypeSessionContainer, error) {
						if accessTokenPayload == nil {
							accessTokenPayload = map[string]interface{}{}
						}

						accessTokenPayload["sub"] = userID + "!!!"

						return originalCreateNewSession(req, res, userID, accessTokenPayload, sessionData, userContext)
					}

					return originalImplementation
				},
			},
		}),
	},
})
}
```

#### If you added an override for `CreateNewSession`/`RefreshSession`/`GetSession`:

This example uses `GetSession`, but the changes required for the other ones are very similar. Before:

```go
session.Init(&sessmodels.TypeInput{
    Override: &sessmodels.OverrideStruct{
        Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
            originalGetSession := *originalImplementation.GetSession

            newGetSession := func(req *http.Request, res http.ResponseWriter, options *sessmodels.VerifySessionOptions, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
                response, err := originalGetSession(req, res, options, userContext)

                if err != nil {
                    return nil, err
                }

                return response, nil
            }
            
            *originalImplementation.GetSession = newGetSession
            return originalImplementation
        },
    },
})
```

After:

```go
session.Init(&sessmodels.TypeInput{
    Override: &sessmodels.OverrideStruct{
        Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
            originalGetSession := *originalImplementation.GetSession

            newGetSession := func(accessToken string, antiCSRFToken *string, options *sessmodels.VerifySessionOptions, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
                defaultUserContext := (*userContext)["_default"].(map[string]interface{})
                request := defaultUserContext["request"]
                
                print(request)
                
                response, err := originalGetSession(accessToken, antiCSRFToken, options, userContext)
                if err != nil {
                    return nil, err
                }
                
                return response, nil
            }
            
            *originalImplementation.GetSession = newGetSession

            return originalImplementation
        },
    },
}),
```

## [0.11.0] - 2023-04-28
- Added missing arguments in `GetUsersNewestFirst` and `GetUsersOldestFirst`

## [0.10.8] - 2023-04-18
- Email template for verify email updated 

## [0.10.7] - 2023-04-11
- Changed email template to render correctly in outlook

## [0.10.6]

-   Fixes panic issue in input validation for emailpassword APIs - https://github.com/supertokens/supertokens-golang/issues/254

## [0.10.5] - 2023-03-31

-   Adds search APIs to the dashboard recipe

## [0.10.4] - 2023-03-30

-   Adds a telemetry API to the dashboard recipe

## [0.10.3] - 2023-03-29
-   Adds unit test for Apple callback form post
-   Updates all example apps to also initialise dashboard recipe
-   Adds login with gitlab (for single tenant only) and bitbucket

## [0.10.2] - 2023-02-24
-   Adds APIs and logic to the dashboard recipe to enable email password based login

## [0.10.1] - 2023-02-06

-   Email template updates

## [0.10.0] - 2023-02-02

### Fixes
-   Fixes issue with go-fiber example, where updating accessTokenPayload from user defined endpoint doesn't reflect in the response cookies.

### Added
-   Added support for authorizing requests using the `Authorization` header instead of cookies
-   Optional `GetTokenTransferMethod` config is Session recipe input, which determines the token transfer method.
-   Check out https://supertokens.com/docs/thirdpartyemailpassword/common-customizations/sessions/token-transfer-method for more information

### Removed
-   ID Refresh token is removed from the SDK

### Breaking changes
-   The frontend SDK should be updated to a version supporting the header-based sessions!
    -   supertokens-auth-react: >= 0.31.0
    -   supertokens-web-js: >= 0.5.0
    -   supertokens-website: >= 16.0.0
    -   supertokens-react-native: >= 4.0.0
    -   supertokens-ios >= 0.2.0
    -   supertokens-android >= 0.3.0
    -   supertokens-flutter >= 0.1.0
-   `CreateNewSession` now requires passing the request as well as the response.
    -   This only requires a change if you manually created sessions (e.g.: during testing)
    -   Check the migration example below
-   `CreateNewSessionWithContext` and `CreateNewSession` in the session recipe accepts new 
-   Only supporting FDI 1.16
parameter `req` of type `*http.Request`

### Migration

Before:

```go
func httpHandler(w http.ResponseWriter, r *http.Request,) {
    sessionContainer, err := session.CreateNewSession(w, "userId", map[string]interface{}{}, map[string]interface{}{})
    if err != nil {
        // handle error
    }
    // ...
}
```

After:

```go
func httpHandler(w http.ResponseWriter, r *http.Request,) {
    sessionContainer, err := session.CreateNewSession(r, w, "userId", map[string]interface{}{}, map[string]interface{}{})
    if err != nil {
        // handle error
    }
    // ...
}
```

## [0.9.14] - 2022-12-26

-   Fixes an issue in the dashboard recipe when fetching user details for passwordless users that don't have an email associated with their accounts
-   Updates dashboard version
-   Updates user GET API for dashboard recipe

## [0.9.13] - 2022-12-26
-   Adds optional `Username` to `SMTPSettings`, which can be used for SMTP login if username is different from `From.Email`.

## [0.9.12] - 2022-12-26
-   Fixes `newPassword` validation in Dashboard API

## [0.9.11]
-   Fixes panic issue with dashboard usersGet API

## [0.9.10]
-   Fixes issue where if SendEmail is overridden with a different email, it will reset that email.

## [0.9.9] - 2022-11-24

### Added:
- Adds APIs for user details to the dashboard recipe

### Changed:
- Updates dashboard version to 0.2

## [0.9.8] - 2022-11-16
- Fixes go fiber to handle handler chaining correctly with verifySession.
- Added test to check JWT contains updated value when MergeIntoAccessTokenPayload is called.
- Adds updating of session claims in email verification token generation API in case the session claims are outdated.

## [0.9.7] - 2022-10-20
- Updated Frontend integration test server for angular tests

### Fixes
- Fixes Apple secret key computation

## [0.9.6] - 2022-10-17
- Updated google token endpoint

## [0.9.5] - 2022-10-14
### Fixes
- Fixes crash in findRightProvider

### Changed:

-   Removed default defaultMaxAge from session claim base classes
-   Added a 5 minute defaultMaxAge to UserRoleClaim, PermissionClaim and EmailVerificationClaim

## [0.9.4] - 2022-09-30
### Fixes
- Using UnixNano instead of UnixMilli to support go version < 1.17

## [0.9.3] - 2022-09-29
### Fixes
- Clears cookies before calling onUnauthorizedError handler if ClearCookies is nil or set to true
- Email verification endpoints will now clear the session if called by a deleted/unknown user

## [0.9.2] - 2022-09-22
### Changed

- Email verification endpoints will now clear the session if called by a deleted/unknown user


## [0.9.1] - 2022-09-20

### Adds:

- Adds Dashboard recipe


## [0.9.0] - 2022-09-14

### Added

- Added support for session claims with related interfaces and classes.
- Added `EmailVerificationClaim`.
- `Mode` config is added to `evmodels.TypeInput`
- `GetEmailForUserID` config is added to `evmodels.TypeInput`
- Added `OnInvalidClaim` optional error handler to send InvalidClaim error responses.
- Added `INVALID_CLAIMS` to `SessionErrors`.
- Added `InvalidClaimStatusCode` optional config to set the status code of InvalidClaim errors.
- Added `OverrideGlobalClaimValidators` to options of `getSession` and `verifySession`.
- Added `MergeIntoAccessTokenPayload` to the Session recipe and session objects which should be preferred to the now deprecated `UpdateAccessTokenPayload`.
- Added `AssertClaims`, `ValidateClaimsForSessionHandle`, `ValidateClaimsInJWTPayload` to the Session recipe to support validation of the newly added `EmailVerificationClaim`.
- Added `FetchAndSetClaim`, `GetClaimValue`, `SetClaimValue` and `RemoveClaim` to the Session recipe to manage claims.
- Added `AssertClaims`, `FetchAndSetClaim`, `GetClaimValue`, `SetClaimValue` and `RemoveClaim` to session objects to manage claims.
- Added sessionContainer to the input of `GenerateEmailVerifyTokenPOST`, `VerifyEmailPOST`, `IsEmailVerifiedGET`.
- Adds default UserContext for verifySession calls that contains the request object.
- Added `UserRoleClaim` and `PermissionClaim` to user roles recipe.

### Breaking changes
- Removes support for FDI < 1.15
- Changed `SignInUp` third party recipe function to accept an email string instead of an object that takes `{ID: string, IsVerified: boolean}`.
- The frontend SDK should be updated to a version supporting session claims!
  - supertokens-auth-react: >= 0.25.0
  - supertokens-web-js: >= 0.2.0
- `EmailVerification` recipe is now not initialized as part of auth recipes, it should be added to the `recipeList` directly instead.
- Email verification related overrides (`EmailVerificationFeature` prop of `Override`) moved from auth recipes into the `EmailVerification` recipe config.
- ThirdParty recipe no longer takes EmailDelivery config -> use Emailverification recipe's EmailDelivery instead.
- Moved email verification related configs from the `EmailDelivery` config of auth recipes into a separate `EmailVerification` email delivery config.
- Updated return type of `GetEmailForUserId` in the `EmailVerification` recipe config. It should now return `OK`, `EmailDoesNotExistError` or `UnknownUserIDError` as response.
- Removed `GetResetPasswordURL`, `GetEmailVerificationURL`, `GetLinkDomainAndPath`. Changing these urls can be done in the email delivery configs instead.
- Removed `UnverifyEmail`, `RevokeEmailVerificationTokens`, `IsEmailVerified`, `VerifyEmailUsingToken` and `CreateEmailVerificationToken` from auth recipes. These should be called on the `EmailVerification` recipe instead.
- Changed function signature for email verification APIs to accept a sessionContainer as an input.
- Changed Session API interface functions:
  - `RefreshPOST` now returns a Session container object.
  - `SignOutPOST` now takes in an optional session object as a parameter.
- `SessionContainer` is renamed to `TypeSessionContainer` and `SessionContainer` is now an alias for `*TypeSessionContainer`. All `*SessionContainer` is now replaced with `SessionContainer`.
- Removed unused parameter `email` from `thirdpartyemailpassword.GetUserByThirdPartyInfoWithContext` function.

### Migration

Before:

```go

supertokens.Init(supertokens.TypeInput{
    AppInfo: supertokens.AppInfo{
        AppName:       "...",
        APIDomain:     "...",
        WebsiteDomain: "...",
    },

    RecipeList: []supertokens.Recipe{
        emailpassword.Init(&epmodels.TypeInput{
            EmailVerificationFeature: evmodels.TypeInput{
                // ...
            },
            Override: &epmodels.OverrideStruct{
                EmailVerificationFeature: &evmodels.OverrideStruct{
                    // ...
                },
            },
        }),
    },
})

```

After the update:

```go

supertokens.Init(supertokens.TypeInput{
    AppInfo: supertokens.AppInfo{
        AppName:       "...",
        APIDomain:     "...",
        WebsiteDomain: "...",
    },

    RecipeList: []supertokens.Recipe{
        emailverification.Init(evmodels.TypeInput{
            Mode: evmodels.ModeRequired, // or evmodels.ModeOptional
            // all config should be moved here from the emailVerificationFeature prop of the EmailPassword recipe config
            Override: &evmodels.OverrideStruct{
                // move the overrides from the emailVerificationFeature prop of the override config in the EmailPassword init here
            },
        }),
        emailpassword.Init(nil),
    },
})

```

### Passwordless users and email verification

If you turn on email verification your email-based passwordless users may be redirected to an email verification screen in their existing session.
Logging out and logging in again will solve this problem or they could click the link in the email to verify themselves.

You can avoid this by running a script that will:

1. list all users of passwordless
2. create an emailverification token for each of them if they have email addresses
3. user the token to verify their address

Something similar to this script:

```go
package main

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func main() {
	supertokens.Init(supertokens.TypeInput{
		AppInfo: supertokens.AppInfo{
			AppName:       "...",
			APIDomain:     "...",
			WebsiteDomain: "...",
		},

		RecipeList: []supertokens.Recipe{
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeRequired,
			}),
			passwordless.Init(plessmodels.TypeInput{
				FlowType: "USER_INPUT_CODE_AND_MAGIC_LINK",
				ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
					Enabled: true,
				},
			}),
			session.Init(nil),
		},
	})

	var paginationToken *string
	recipeList := []string{"passwordless"}
	limit := 100
	done := false

	for !done {
		userList, err := supertokens.GetUsersNewestFirst(paginationToken, &limit, &recipeList)
		if err != nil {
			panic(err)
		}

		for _, user := range userList.Users {
			if user.RecipeId == "passwordless" && user.User["email"] != nil {
				token, err := emailverification.CreateEmailVerificationToken(user.User["id"].(string), nil)
				if err != nil {
					panic(err)
				}
				if token.OK != nil {
					_, err := emailverification.VerifyEmailUsingToken(token.OK.Token)
					if err != nil {
						panic(err)
					}
				}
			}

			done = (userList.NextPaginationToken == nil)
			paginationToken = userList.NextPaginationToken
		}
	}
}
```

#### User roles

The UserRoles recipe now adds role and permission information into the access token payload by default. If you are already doing this manually, this will result in duplicate data in the access token.

- You can disable this behaviour by setting `SkipAddingRolesToAccessToken` and `SkipAddingPermissionsToAccessToken` to true in the recipe init.
- Check how to use the new claims in the updated guide: https://supertokens.com/docs/userroles/protecting-routes


## [0.8.3] - 2022-07-30
### Added
- Adds test to verify that session container uses overridden functions
- Adds with-go-zero example: https://github.com/supertokens/supertokens-golang/issues/157
- UserId Mapping functionality and compatibility with CDI 2.15
- Adds `CreateUserIdMapping`, `GetUserIdMapping`, `DeleteUserIdMapping`, `UpdateOrDeleteUserIdMappingInfo` functions to supertokens package


## [0.8.2] - 2022-07-18

### Fixes:
- Fixes JWKS Keyfunc call that resulted in a goroutine leak: https://github.com/supertokens/supertokens-golang/issues/155

## [0.8.1] - 2022-07-12

### Fixes:
- Fixes issue with 404 status being sent for apple redirect callback route.

## [0.8.0] - 2022-07-08

### Breaking change:
-   Changes session recipe interfaces to not return an `UNAUTHORISED` error when the input is a sessionHandle: https://github.com/supertokens/backend/issues/83
-   `GetSessionInformation` now returns `nil` is the session does not exist
-   `UpdateSessionData` now returns `nil` if the input `sessionHandle` does not exist.
-   `UpdateAccessTokenPayload` now returns `false` if the input `sessionHandle` does not exist.
-   `RegenerateAccessToken` now returns `nil` if the input access token's `sessionHandle` does not exist.
-   The session container functions have not changed in behaviour and return errors if `sessionHandle` does not exist. This works on the current session.

### Fixes
-   Clears cookies when RevokeSession is called using the session container, even if the session did not exist from before: https://github.com/supertokens/supertokens-node/issues/343

### Adds:
-   Adds default userContext for API calls that contains the request object. It can be used in APIs / functions override like so:

```golang
SignIn: func (..., userContext supertokens.UserContext) {
    if _default, ok := (*userContext)["_default"].(map[string]interface{}); ok {
        if req, ok := _default["request"].(*http.Request); ok {
            // do something here with the request object
        }
    }
}
```

## [0.7.2] - 2022-06-29
-   Adds unit tests for resend email & sms services for passwordless and thirdpartypasswordless recipes
-   Adds User Roles recipe and compatibility with CDI 2.14

## [0.7.1] - 2022-06-27
-   Fixes panic while returning empty result object with nil error in the API overrides. Related to https://github.com/supertokens/supertokens-golang/issues/107

## [0.7.0] - 2022-06-23
### Breaking change
-   Renamed `SMTPServiceConfig` to `SMTPSettings`
-   Changed type of `Secure` in `SMTPSettings` from `*bool` to `bool`
-   Renamed `SMTPServiceFromConfig` to `SMTPFrom`
-   Renamed `SMTPGetContentResult` to `EmailContent`
-   Renamed `SMTPTypeInput` to `SMTPServiceConfig`
-   Renamed field `SMTPSettings` to `Settings` in `SMTPServiceConfig`
-   Renamed `SMTPServiceInterface` to `SMTPInterface`
-   Renamed all instances of `MakeSmtpService` to `MakeSMTPService`
-   All instances of `MakeSMTPService` returns `*EmailDeliveryInterface` instead of `EmailDeliveryInterface`
-   Renamed `TwilioServiceConfig` to `TwilioSettings`
-   Renamed `TwilioGetContentResult` to `SMSContent`
-   Renamed `TwilioTypeInput` to `TwilioServiceConfig`
-   Renamed field `TwilioSettings` to `Settings` in `TwilioServiceConfig`
-   Changed types of fields `From` and `MessagingServiceSid` in `TwilioSettings` from `*string` to `string`
-   Renamed `MakeSupertokensService` to `MakeSupertokensSMSService`
-   All instances of `MakeSupertokensSMSService` and `MakeTwilioService` returns `*SmsDeliveryInterface` instead of `SmsDeliveryInterface`
-   Removed `SupertokensServiceConfig` and `MakeSupertokensSMSService` accepts `apiKey` directly instead of `SupertokensServiceConfig`
-   Renamed `TwilioServiceInterface` to `TwilioInterface`
- Removes support for FDIs that are < 1.14

### Added
-   Exposed `MakeSMTPService` from emailverification, emailpassword, passwordless, thirdparty, thirdpartyemailpassword and thirdpartypasswordless recipes
-   Exposed `MakeSupertokensSMSService` and `MakeTwilioService` from passwordless and thirdpartypasswordless recipes

### Fixes
- Fixes Cookie SameSite config validation.
- Changes `getEmailForUserIdForEmailVerification` function inside thirdpartypasswordless to take into account passwordless emails and return an empty string in case a passwordless email doesn't exist. This helps situations where the dev wants to customise the email verification functions in the thirdpartypasswordless recipe.

## [0.6.8] - 2022-06-17
### Added
- `EmailDelivery` user config for Emailpassword, Thirdparty, ThirdpartyEmailpassword, Passwordless and ThirdpartyPasswordless recipes.
- `SmsDelivery` user config for Passwordless and ThirdpartyPasswordless recipes.
- `Twilio` service integration for SmsDelivery ingredient.
- `SMTP` service integration for EmailDelivery ingredient.
- `Supertokens` service integration for SmsDelivery ingredient.

### Deprecated
- For Emailpassword recipe input config, `ResetPasswordUsingTokenFeature.CreateAndSendCustomEmail` and `EmailVerificationFeature.CreateAndSendCustomEmail` have been deprecated.
- For Thirdparty recipe input config, `EmailVerificationFeature.CreateAndSendCustomEmail` has been deprecated.
- For ThirdpartyEmailpassword recipe input config, `ResetPasswordUsingTokenFeature.CreateAndSendCustomEmail` and `EmailVerificationFeature.CreateAndSendCustomEmail` have been deprecated.
- For Passwordless recipe input config, `CreateAndSendCustomEmail` and `CreateAndSendCustomTextMessage` have been deprecated.
- For ThirdpartyPasswordless recipe input config, `CreateAndSendCustomEmail`, `CreateAndSendCustomTextMessage` and `EmailVerificationFeature.CreateAndSendCustomEmail` have been deprecated.

### Migration

Following is an example of ThirdpartyPasswordless recipe migration. If your existing code looks like

```go
func passwordlessLoginEmail(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
	// some custom logic
}

func passwordlessLoginSms(phoneNumber string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
	// some custom logic
}

func verifyEmail(user tplmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
	// some custom logic
}

supertokens.Init(supertokens.TypeInput{
    AppInfo: supertokens.AppInfo{
        AppName:       "...",
        APIDomain:     "...",
        WebsiteDomain: "...",
    },
    RecipeList: []supertokens.Recipe{
        thirdpartypasswordless.Init(tplmodels.TypeInput{
            FlowType: "...",
            ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
                Enabled: true,
                CreateAndSendCustomEmail: passwordlessLoginEmail,
                CreateAndSendCustomTextMessage: passwordlessLoginSms,
            },
            EmailVerificationFeature: &tplmodels.TypeInputEmailVerificationFeature{
                CreateAndSendCustomEmail: verifyEmail,
            },
        }),
    },
})
```

After migration to using new `EmailDelivery` and `SmsDelivery` config, your code would look like:
```go
func passwordlessLoginEmail(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
	// some custom logic
	return nil
}

func passwordlessLoginSms(phoneNumber string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
	// some custom logic
	return nil
}

func verifyEmail(user tplmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
	// some custom logic
}

var sendEmail = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
	if input.EmailVerification != nil {
		verifyEmail(tplmodels.User{ID: input.EmailVerification.User.ID, Email: &input.EmailVerification.User.Email}, input.EmailVerification.EmailVerifyLink, userContext)
	} else if input.PasswordlessLogin != nil {
		return passwordlessLoginEmail(input.PasswordlessLogin.Email, input.PasswordlessLogin.UserInputCode, input.PasswordlessLogin.UrlWithLinkCode, input.PasswordlessLogin.CodeLifetime, input.PasswordlessLogin.PreAuthSessionId, userContext)
	}
	return nil
}

var sendSms = func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
	if input.PasswordlessLogin != nil {
		return passwordlessLoginSms(input.PasswordlessLogin.PhoneNumber, input.PasswordlessLogin.UserInputCode, input.PasswordlessLogin.UrlWithLinkCode, input.PasswordlessLogin.CodeLifetime, input.PasswordlessLogin.PreAuthSessionId, userContext)
	}
	return nil
}

supertokens.Init(supertokens.TypeInput{
    AppInfo: supertokens.AppInfo{
        AppName:       "...",
        APIDomain:     "...",
        WebsiteDomain: "...",
    },
    RecipeList: []supertokens.Recipe{
        thirdpartypasswordless.Init(tplmodels.TypeInput{
            FlowType: "...",
            ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
                Enabled: true,
            },
            EmailDelivery: &emaildelivery.TypeInput{
                Service: &emaildelivery.EmailDeliveryInterface{
                    SendEmail: &sendEmail,
                },
            },
            SmsDelivery: &smsdelivery.TypeInput{
                Service: &smsdelivery.SmsDeliveryInterface{
                    SendSms: &sendSms,
                },
            },
        }),
    },
})
```

## [0.6.7]
- Fixes panic when call to thirdparty provider API returns a non 2xx status.

### Breaking change
-   https://github.com/supertokens/supertokens-node/issues/220
    -   Adds `{status: "GENERAL_ERROR", message: string}` as a possible output to all the APIs.
    -   Changes `FIELD_ERROR` output status in third party recipe API to be `GENERAL_ERROR`.
    -   Replaced `FIELD_ERROR` status type in third party signinup API with `GENERAL_ERROR`.
    -   Removed `FIELD_ERROR` status type from third party signinup recipe function.
- Changes output of `VerifyEmailPOST` to `VerifyEmailPOSTResponse`
- Changes output of `PasswordResetPOST` to `ResetPasswordPOSTResponse`
- `SignInUp` recipe function doesn't return `FIELD_ERROR` anymore in thirdparty, thirdpartypasswordless and thirdpartyemailpassword recipe.
- `SignInUpPOST` api function returns `GENERAL_ERROR` instead of `FIELD_ERROR` in thirdparty, thirdpartypasswordless and thirdpartyemailpassword recipe.
- If there is an error in sending SMS or email in passwordless based recipes, then we no longer return a GENERAL_ERROR, but instead, we return a regular golang error.
- Changes `GetJWKSGET` in JWT recipe to return `GetJWKSAPIResponse` (that also contains a General Error response)
- Changes `GetOpenIdDiscoveryConfigurationGET` in Open ID recipe to return `GetOpenIdDiscoveryConfigurationAPIResponse` (that also contains a General Error response)
- Renames `OnGeneralError` callback (that's in user input) to `OnSuperTokensAPIError`
- If there is an error in the `errorHandler`, we no longer call `OnSuperTokensAPIError` in that, but instead, we return an error back.

## [0.6.6]
- Fixes facebook login

## [0.6.5]
- Fixes issue in reading request body in API override: https://github.com/supertokens/supertokens-golang/issues/116

## [0.6.4]
- Fixes issue in writing custom response in API override with general error
### Added
- Adds unit tests to thirdpartypasswordless recipe

## [0.6.3] - 2022-05-19
### Fixes
- Fixes the function signature of the `GetUserByThirdPartyInfo` function in the `thirdpartypasswordless` recipe.

## [0.6.2] - 2022-05-18
### Fixes
- Fixes issue in writing custom response in API Override

## [0.6.1] - 2022-05-17
### Fixes
- https://github.com/supertokens/supertokens-golang/issues/102. Sending `preAuthSessionID` instead of `preAuthSessionId` to the core.
- Fixes the error message in AuthorizationUrlAPI function in the `api` module of the thirdparty recipe in case when providers is nil

## [0.6.0] - 2022-05-13
### Breaking Change

- Adds both with context and without context functions to thirdparty passwordless recipe, Like all other recipes. Where we expose both WithContext functions and without context functions, which are basically the same as WithContext ones with an emtpy map[string]interface{} passed as context

### Added
- Adds unit tests to passwordless recipe 

### Fixes
- Fixes existing action to run go mod tidy in the examples folder
- Fixes stopSt function in testing utils

## [0.5.9] - 2022-05-10
### Fixes
- Fixes bug in the revokeCode function of the recipeimplementation in passwordless recipe 

## [0.5.8] - 2022-05-05
### Added
- Adds Github Actions for testing and pre-commit hooks.
- Adds more unit tests for thirdpary email password recipe
- Adds test to jwt recipe
- Adds test to opendID recipe


### Fixes
- Third party sign in up API response correction.

## [0.5.7] - 2022-04-23
- Adds functions to delete passwordless user info in recipes that have passwordless users.
- Fixes bug in signinup helper function exposed by passwordless recipe

## [0.5.6] - 2022-04-18

- Adds UserMetadata recipe

## [0.5.5] - 2022-04-11
### Added 
-   Adds functions for debug logging

## [0.5.4] - 2022-03-30

### Added
 - workflow to enforce go mod tidy is run when issuing a PR. 

## [0.5.3] - 2022-03-24

### Fixes
- Checks if discord returned email before setting it in the profile info obj.

## [0.5.2] - 2022-03-17
- Adds thirdpartypasswordless recipe: https://github.com/supertokens/supertokens-core/issues/331

## [0.5.1] - 2022-02-07

-   Adds testing framework along with unit tests for the recipes
-   Adds unit tests for thirdparty recipe and thirdpartyemailpassword recipe
-   Adds example implementation with go fiber

## [0.5.0] - 2022-02-20
### Breaking Change

-   Adds user context to all functions exposed to the user, and to API and Recipe interface functions. This is a non breaking change for User exposed function calls, but a breaking change if you are using the Recipe or APIs override feature
-   Returns session from API interface functions that create a session
-   Renames functions in ThirdPartyEmailPassword recipe (https://github.com/supertokens/supertokens-node/issues/219):
    -   Recipe Interface:
        -   `SignInUp` -> `ThirdPartySignInUp`
        -   `SignUp` -> `EmailPasswordSignUp`
        -   `SignIn` -> `EmailPasswordSignIn`
    -   API Interface:
        -   `EmailExistsGET` -> `EmailPasswordEmailExistsGET`
    -   User exposed functions (in `recipe/thirdpartyemailpassword/main.go`)
        -   `SignInUp` -> `ThirdPartySignInUp`
        -   `SignUp` -> `EmailPasswordSignUp`
        -   `SignIn` -> `EmailPasswordSignIn`

### Change:

-   Uses recipe interface inside session class so that any modification to those get reflected in the session class functions too.

## [0.4.2] - 2022-01-31
- Adds ability to give a path for each of the hostnames in the connectionURI: https://github.com/supertokens/supertokens-node/issues/252
- Adds workflow to verify if pr title follows conventional commits
- Added userId as an optional property to the response of `recipe/user/password/reset` (Compatibility with CDI 2.12).

### Added

-   Added `regenerateAccessToken` as a new recipe function for the session recipe.
-   Added a bunch of new functions inside the session container which gives user the ability to either call a       function with userContext or just call the function without it (for example: `RevokeSession` and `RevokeSessionWithContext`)
 
### Breaking changes:

-   Allows passing of custom user context everywhere: https://github.com/supertokens/supertokens-golang/issues/64


## [0.4.1] - 2022-01-27
-   Fixes https://github.com/supertokens/supertokens-node/issues/244 - throws an error if a user tries to update email / password of a third party login user.
-   Adds check to see if user has provided empty connectionInfo
-   Adds fixes to solve casting of data in session-functions

## [0.4.0] - 2022-01-14

-   Adds passwordless recipe
-   Adds compatibility with FDI 1.11 and CDI 2.11

## [0.3.5] - 2022-01-08

### Fixes
- Fixes issue of methods getting hidden due to DoneWriter wrapper around ResponseWriter: https://github.com/supertokens/supertokens-golang/issues/55

## [0.3.4] - 2022-01-06

### Fixes
- Sends application/json content-type in `SendNon200Response` function: https://github.com/supertokens/supertokens-golang/issues/53

## [0.3.3] - 2021-12-20

### Added
- Add DeleteUser function

## [0.3.2] - 2021-12-06
### Added
-   The ability to enable JWT creation with session management, this allows easier integration with services that require JWT based authentication: https://github.com/supertokens/supertokens-core/issues/250

## [0.3.1] - 2021-12-06
### Changes
- Upgrade `keyfunc` dependency to stable version.

### Fixes
- Removes use of apiGatewayPath from apple's redirect URI since that is already there in the apiBasePath


## [0.3.0] - 2021-11-23

### Breaking changes:
- Changes `FIELD_ERROR` type in sign in up response from `Error` to `ErrorMsg`

### Addition
- Sign in with google workspaces and discord

### Changes
- If getting profile info from third party provider throws an error, that is propagated a `FIELD_ERROR` to the client.

## [0.2.2] - 2021-11-15

### Changes
- Does not send a response if the user has already sent the response: https://github.com/supertokens/supertokens-node/issues/197

## [0.2.1] - 2021-11-08

### Changes
-   When routing, ignores `rid` value `"anti-csrf"`: https://github.com/supertokens/supertokens-node/issues/202

## [0.2.0] - 2021-10-21

### Breaking changes:
- Makes recipe and API interface have pointers to functions to fix https://github.com/supertokens/supertokens-node/issues/199
-   Support for FDI 1.10:
    -   Allow thirdparty `/signinup POST` API to take `authCodeResponse` XOR `code` so that it can supprt OAuth via PKCE

### Added:
- Makes recipe and API interface have pointers to functions to fix https://github.com/supertokens/supertokens-node/issues/199
-   Support for FDI 1.10:
    -   Adds apple sign in callback API
-   Optional `getRedirectURI` function added to social providers in case we set the `redirect_uri` on the backend.
-   Adds optional `IsDefault` param to auth providers so that they can be reused with different credentials.
- Adds sign in with apple support: https://github.com/supertokens/supertokens-golang/issues/20

## [0.1.0] - 2021-10-21

### Breaking change:

- Removes `SignInUpPost` from thirdpartyemailpassword API interface and replaces it with three APIs: `EmailPasswordSignInPOST`, `EmailPasswordSignUpPOST` and `ThirdPartySignInUpPOST`: https://github.com/supertokens/supertokens-node/issues/192
- Renames all JWT function names to use AccessToken instead for clarity

## [0.0.6] - 2021-10-18

### Changed

-  Changes implementation such that actual client IDs are not in the SDK, removes imports for OAuth dev related code.

## [0.0.5] - 2021-10-18

### Fixed

- URL protocol is being taken into account when determining the value of cookie same site: https://github.com/supertokens/supertokens-golang/issues/36

## [0.0.4] - 2021-10-12

### Added

- Adds OAuth development keys for Google and Github for faster recipe implementation.

## [0.0.3] - 2021-09-25

### Added

- Support for FDI 1.9
- JWT Recipe

### Fixed
- Sets response content-type as JSON

## [0.0.2] - 2021-09-22

### Added

-   Support for multiple access token signing keys: https://github.com/supertokens/supertokens-core/issues/305
-   Supporting CDI 2.9

## [0.0.1] - 2021-09-18

### Added
- Initial version of the repo
