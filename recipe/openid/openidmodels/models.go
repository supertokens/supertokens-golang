package openidmodels

import (
	"github.com/supertokens/supertokens-golang/recipe/jwt/jwtmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TypeInput struct {
	Issuer             *string
	JwtValiditySeconds *uint64
	Override           *OverrideStruct
}

type TypeNormalisedInput struct {
	IssuerDomain       supertokens.NormalisedURLDomain
	IssuerPath         supertokens.NormalisedURLPath
	JwtValiditySeconds *uint64
	Override           OverrideStruct
}

type OverrideStruct struct {
	Functions  func(originalImplementation RecipeInterface) RecipeInterface
	APIs       func(originalImplementation APIInterface) APIInterface
	JwtFeature *jwtmodels.OverrideStruct
}
