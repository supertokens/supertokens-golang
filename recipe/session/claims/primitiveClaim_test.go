package claims

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TestPrimitiveClaim(t *testing.T) {
	primClaim := PrimitiveClaim{
		Key: "test",
		fetchValue: func(userId string, userContext supertokens.UserContext) interface{} {
			return map[string]interface{}{}
		},
	}
	payload := map[string]interface{}{}
	payload = primClaim.AddToPayload_internal(payload, "world", nil)

	validators := primClaim.GetValidators()
	assert.True(t, validators.HasValue("world", nil).Validate(payload, nil).IsValid)
	assert.False(t, validators.HasValue("hello", nil).Validate(payload, nil).IsValid)

	assert.True(t, validators.HasFreshValue("world", 5, nil).Validate(payload, nil).IsValid)
	time.Sleep(2 * time.Second)
	assert.False(t, validators.HasFreshValue("world", 1, nil).Validate(payload, nil).IsValid)
}
