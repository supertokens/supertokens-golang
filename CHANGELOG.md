# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2021-10-21

### Breaking changes:
- Makes recipe and API interface have pointers to functions to fix https://github.com/supertokens/supertokens-node/issues/199
-   Support for FDI 1.10:
    -   Allow thirdparty `/signinup POST` API to take `authCodeResponse` XOR `code` so that it can supprt OAuth via PKCE
    -   Adds apple sign in callback API
-   Optional `getRedirectURI` function added to social providers in case we set the `redirect_uri` on the backend.
-   Adds optional `IsDefault` param to auth providers so that they can be reused with different credentials.

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