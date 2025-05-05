package gsync_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/drshriveer/gtools/gsync"
)

func TestExecutor(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	executor, done := gsync.NewSliceExecutor[int](ctx)
	defer done()
	assert.NoError(t, executor.AddTask(func(_ context.Context) (int, error) {
		return 1, nil
	}))
	assert.NoError(t, executor.AddTask(func(_ context.Context) (int, error) {
		return 2, nil
	}))
	assert.NoError(t, executor.WaitForCompletion())
	result := executor.Result()
	assert.ElementsMatch(t, []int{1, 2}, result)
}
