package ctxt

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextWithUserAndUserFromContext(t *testing.T) {
	ctx := context.Background()
	uctx := ContextWithUser(ctx, &User{
		ID:  5,
		CIF: "123456789",
	})
	assert.NotNil(t, uctx)

	user, ok := UserFromContext(uctx)
	assert.True(t, ok)
	assert.Equal(t, 5, user.ID)
	assert.Equal(t, "123456789", user.CIF)
}
