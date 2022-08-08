package claims

type PrimitiveArrayClaim struct {
	PrimitiveClaim
}

func (claim *PrimitiveArrayClaim) GetValidators() PrimitiveArrayClaimValidators {
	return PrimitiveArrayClaimValidators{
		PrimitiveClaimValidators: claim.PrimitiveClaim.GetValidators(),
	}
}

type PrimitiveArrayClaimValidators struct {
	PrimitiveClaimValidators
	Includes    func(val Any, maxAgeInSeconds *int, id *string)
	Excludes    func(val Any, maxAgeInSeconds *int, id *string)
	IncludesAll func(vals []interface{}, maxAgeInSeconds *int, id *string)
	ExcludesAll func(vals []interface{}, maxAgeInSeconds *int, id *string)
}
