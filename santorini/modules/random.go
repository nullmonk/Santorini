package modules

import (
	"crypto/rand"
	"math/big"

	"go.starlark.net/starlark"
)

var Random = NewModule("random", map[string]Function{
	"choice": func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var list *starlark.List
		if err := starlark.UnpackPositionalArgs("", args, kwargs, 1, &list); err != nil {
			return nil, err
		}
		i, _ := rand.Int(rand.Reader, big.NewInt(int64(list.Len())))
		return list.Index(int(i.Int64())), nil
	},
	"rand_int": nil,
})
