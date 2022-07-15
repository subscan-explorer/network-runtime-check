package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceExclude(t *testing.T) {
	t1o, t1s := []string{"Account", "Value"}, []string{"Account"}
	assert.Equal(t, SliceEqual(t1o, t1s), false)

	t2o, t2s := []string{"Account"}, []string{"Account"}
	assert.Equal(t, SliceEqual(t2o, t2s), true)

	t3o, t3s := []string{"Account", "Value", "Account"}, []string{"Account", "Value", "Account", "Account"}
	assert.Equal(t, SliceEqual(t3o, t3s), false)

	t4o, t4s := []string{"Account", "Value", "Account"}, []string{"Account", "Account"}
	assert.Equal(t, SliceEqual(t4o, t4s), false)

	t5o, t5s := []string{"Account", "Account"}, []string{"Value", "Value"}
	assert.Equal(t, SliceEqual(t5o, t5s), false)

	t6o, t6s := []string{"Account", "Account", "Value"}, []string{"Value", "Account", "Account"}
	assert.Equal(t, SliceEqual(t6o, t6s), true)
}
