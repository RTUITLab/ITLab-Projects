package updater_test

import (
	"context"
	"testing"

	"github.com/ITLab-Projects/pkg/updater"
)

func TestFunc_CheckContext(t *testing.T) {
	ctx := updater.WithUpdateContext(context.Background())

	if !updater.IsUpdateContext(ctx) {
		t.Log("Assert error")
		t.FailNow()
	}
}

func TestFunc_CheckContext_False(t *testing.T) {
	ctx := context.Background()

	if updater.IsUpdateContext(ctx) {
		t.Log("Assert error")
		t.FailNow()
	}
}