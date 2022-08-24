package claims

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TestBooleanClaim(t *testing.T) {
	boolClaim := BooleanClaim(
		"test",
		func(userId string, userContext supertokens.UserContext) (interface{}, error) {
			return map[string]interface{}{}, nil
		},
		nil,
	)

	payload := map[string]interface{}{}
	payload = boolClaim.AddToPayload_internal(payload, true, nil)

	validators := boolClaim.Validators
	assert.True(t, validators.IsTrue(nil).Validate(payload, nil).IsValid)
	assert.False(t, validators.IsFalse(nil).Validate(payload, nil).IsValid)

	maxAge := int64(1)
	assert.True(t, validators.IsTrue(&maxAge).Validate(payload, nil).IsValid)
	time.Sleep(2 * time.Second)
	assert.False(t, validators.IsTrue(&maxAge).Validate(payload, nil).IsValid)
}
