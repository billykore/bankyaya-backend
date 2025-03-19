package ctxt

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.bankyaya.org/app/backend/pkg/data"
)

func TestContextWithUserAndUserFromContext(t *testing.T) {
	ctx := context.Background()
	uctx := ContextWithUser(ctx, data.User{
		CIF: "123456789",
		Id:  5,
	})
	assert.NotNil(t, uctx)

	user, ok := UserFromContext(uctx)
	assert.True(t, ok)
	assert.Equal(t, 5, user.Id)
	assert.Equal(t, "123456789", user.CIF)
}
