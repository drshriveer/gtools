package gencommon

import (
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindIFaceDef(t *testing.T) {
	validErrInterface := types.NewInterfaceType([]*types.Func{
		types.NewFunc(
			0,
			nil,
			"Error",
			types.NewSignatureType(nil, nil, nil, nil,
				types.NewTuple(types.NewVar(0, nil, "", types.Typ[types.String])),
				false)),
	}, nil)

	iFace, err := FindIFaceDef("builtin", "error")
	assert.NoError(t, err)
	assert.True(t, types.Implements(iFace, validErrInterface))
	// check the cache
	res, ok := iFaceCache.Load("builtin.error")
	assert.True(t, ok)
	assert.Equal(t, res, iFace)
}
