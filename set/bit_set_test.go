package set_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/drshriveer/gtools/set"
)

type testBitSet uint64

const (
	testBitSetA testBitSet = 1 << iota
	testBitSetB
	testBitSetC
	testBitSetD
)

func TestBitSet_Add(t *testing.T) {
	s := set.MakeBitSet(testBitSetA)
	assert.True(t, s.Has(testBitSetA))
	assert.False(t, s.Has(testBitSetB))
	assert.False(t, s.Has(testBitSetC))
	assert.False(t, s.Has(testBitSetD))
	assert.True(t, s.Add(testBitSetB))
	assert.False(t, s.Add(testBitSetB))
	assert.True(t, s.Has(testBitSetA))
	assert.True(t, s.Has(testBitSetB))
	assert.False(t, s.Has(testBitSetC))
	assert.False(t, s.Has(testBitSetD))
	s.Add(testBitSetC, testBitSetD)
	assert.True(t, s.Has(testBitSetA))
	assert.True(t, s.Has(testBitSetB))
	assert.True(t, s.Has(testBitSetC))
	assert.True(t, s.Has(testBitSetD))
}

func TestBitSet_Remove(t *testing.T) {
	s := set.MakeBitSet(testBitSetA, testBitSetB, testBitSetC, testBitSetD)
	assert.True(t, s.Has(testBitSetA))
	assert.True(t, s.Has(testBitSetB))
	assert.True(t, s.Has(testBitSetC))
	assert.True(t, s.Has(testBitSetD))
	s.Remove(testBitSetB)
	assert.True(t, s.Has(testBitSetA))
	assert.False(t, s.Has(testBitSetB))
	assert.True(t, s.Has(testBitSetC))
	assert.True(t, s.Has(testBitSetD))
	s.Remove(testBitSetC, testBitSetD)
	assert.True(t, s.Has(testBitSetA))
	assert.False(t, s.Has(testBitSetB))
	assert.False(t, s.Has(testBitSetC))
	assert.False(t, s.Has(testBitSetD))
}

func TestBitSet_MaskOf(t *testing.T) {
	s := set.MakeBitSet(testBitSetA, testBitSetB, testBitSetC, testBitSetD)
	mask := s.MaskOf(testBitSetA | testBitSetC)
	assert.True(t, mask.Has(testBitSetA))
	assert.False(t, mask.Has(testBitSetB))
	assert.True(t, mask.Has(testBitSetC))
	assert.False(t, mask.Has(testBitSetD))

	s = set.MakeBitSet(testBitSetA, testBitSetB)
	mask = s.MaskOf(testBitSetA | testBitSetC)
	assert.True(t, mask.Has(testBitSetA))
	assert.False(t, mask.Has(testBitSetB))
	assert.False(t, mask.Has(testBitSetC))
	assert.False(t, mask.Has(testBitSetD))
}
