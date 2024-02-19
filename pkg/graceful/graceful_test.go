package graceful

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Context(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		ctx, cancel := Context(context.Background())
		assert.NotNil(t, ctx)
		assert.NotNil(t, cancel)
	})
}

func Test_Add(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		defer Reset()
		Add()
		assert.Equal(t, "0/1", Count())
		Add()
		assert.Equal(t, "0/2", Count())
	})
}

func Test_Done(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		defer Reset()
		Add()
		Add()
		assert.Equal(t, "0/2", Count())
		Done()
		assert.Equal(t, "1/2", Count())
		Done()
		assert.Equal(t, "2/2", Count())
	})
}

func Test_Wait(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		defer Reset()
		Add()
		Add()
		go func() {
			Done()
			Done()
		}()
		Wait()
		assert.Equal(t, "2/2", Count())
	})
}

func Test_Count(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		defer Reset()
		Add()
		assert.Equal(t, "0/1", Count())
		Done()
		assert.Equal(t, "1/1", Count())
	})
}

func Test_Reset(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		defer Reset()
		Add()
		Reset()
		Wait()
		assert.Equal(t, "0/0", Count())
	})
}

func Test_CleanFn(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		defer Reset()
		Add()
		CleanFn(0, func() {})
		Wait()
		assert.Equal(t, "1/1", Count())
	})
}
