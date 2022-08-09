package claims

type BooleanClaim struct {
	PrimitiveClaim
}

func (claim *BooleanClaim) GetValidators() BooleanClaimValidators {
	primitiveClaimValidators := claim.PrimitiveClaim.GetValidators()
	return BooleanClaimValidators{
		PrimitiveClaimValidators: primitiveClaimValidators,
		IsTrue: func(maxAgeInSeconds *int64) SessionClaimValidator {
			if maxAgeInSeconds != nil {
				return &HasFreshValueImpl{
					id:              claim.Key + "-freshVal",
					claim:           claim,
					maxAgeInSeconds: *maxAgeInSeconds,
					val:             true,
				}
			}
			return &HasValueImpl{
				id:    claim.Key,
				claim: claim,
				val:   true,
			}
		},
		IsFalse: func(maxAgeInSeconds *int64) SessionClaimValidator {
			if maxAgeInSeconds != nil {
				return &HasFreshValueImpl{
					id:              claim.Key + "-freshVal",
					claim:           claim,
					maxAgeInSeconds: *maxAgeInSeconds,
					val:             false,
				}
			}
			return &HasValueImpl{
				id:    claim.Key,
				claim: claim,
				val:   false,
			}
		},
	}
}

type BooleanClaimValidators struct {
	PrimitiveClaimValidators
	IsTrue  func(maxAgeInSeconds *int64) SessionClaimValidator
	IsFalse func(maxAgeInSeconds *int64) SessionClaimValidator
}
