# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [unreleased]

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