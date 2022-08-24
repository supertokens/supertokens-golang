package claims

func BooleanClaim(key string, fetchValue FetchValueFunc, defaultMaxAgeInSeconds *int64) *TypeBooleanClaim {
	primitiveClaim := PrimitiveClaim(key, fetchValue, defaultMaxAgeInSeconds)

	booleanClaim := &TypeBooleanClaim{
		TypePrimitiveClaim: *primitiveClaim,
	}

	booleanClaim.Validators = &BooleanClaimValidators{
		PrimitiveClaimValidators: primitiveClaim.Validators,

		IsTrue: func(maxAgeInSeconds *int64) *SessionClaimValidator {
			return primitiveClaim.Validators.HasValue(true, maxAgeInSeconds, nil)
		},

		IsFalse: func(maxAgeInSeconds *int64) *SessionClaimValidator {
			return primitiveClaim.Validators.HasValue(false, maxAgeInSeconds, nil)
		},
	}

	return booleanClaim
}

type TypeBooleanClaim struct {
	TypePrimitiveClaim
	Validators *BooleanClaimValidators
}

type BooleanClaimValidators struct {
	*PrimitiveClaimValidators
	IsTrue  func(maxAgeInSeconds *int64) *SessionClaimValidator
	IsFalse func(maxAgeInSeconds *int64) *SessionClaimValidator
}
