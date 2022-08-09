package claims

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TestBooleanClaim(t *testing.T) {
	boolClaim := BooleanClaim{
		PrimitiveClaim: PrimitiveClaim{
			Key: "test",
			fetchValue: func(userId string, userContext supertokens.UserContext) interface{} {
				return map[string]interface{}{}
			},
		},
	}
	payload := map[string]interface{}{}
	payload = boolClaim.AddToPayload_internal(payload, true, nil)

	validators := boolClaim.GetValidators()
	assert.True(t, validators.IsTrue(nil).Validate(payload, nil).IsValid)
	assert.False(t, validators.IsFalse(nil).Validate(payload, nil).IsValid)

	maxAge := int64(1)
	assert.True(t, validators.IsTrue(&maxAge).Validate(payload, nil).IsValid)
	time.Sleep(2 * time.Second)
	assert.False(t, validators.IsTrue(&maxAge).Validate(payload, nil).IsValid)
}
