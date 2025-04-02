package ctxt

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func a(ctx context.Context) error {
	err := HandleCancellation(ctx)
	if err != nil {
		return errors.New("HandleCancellation thrown error")
	}
	return nil
}

func b() error {
	time.Sleep(2 * time.Second)
	return errors.New("b error")
}

func TestHandleCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := b()
	if err != nil {
		cancel()
	}
	err = a(ctx)
	if err != nil {
		t.Log(err)
	}
	assert.Error(t, err)
}
