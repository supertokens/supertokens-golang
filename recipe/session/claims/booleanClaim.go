package claims

func BooleanClaim(key string, fetchValue FetchValueFunc) *TypeBooleanClaim {
	primitiveClaim := PrimitiveClaim(key, fetchValue)

	booleanClaim := &TypeBooleanClaim{
		TypePrimitiveClaim: *primitiveClaim,
	}

	booleanClaim.Validators = &BooleanClaimValidators{
		PrimitiveClaimValidators: primitiveClaim.Validators,

		IsTrue: func(maxAgeInSeconds *int64) *SessionClaimValidator {
			if maxAgeInSeconds != nil {
				return primitiveClaim.Validators.HasFreshValue(true, *maxAgeInSeconds, nil)
			}
			return primitiveClaim.Validators.HasValue(true, nil)
		},

		IsFalse: func(maxAgeInSeconds *int64) *SessionClaimValidator {
			if maxAgeInSeconds != nil {
				return primitiveClaim.Validators.HasFreshValue(false, *maxAgeInSeconds, nil)
			}
			return primitiveClaim.Validators.HasValue(false, nil)
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
