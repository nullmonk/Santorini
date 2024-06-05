package modules

import (
	"fmt"
	"reflect"

	"go.starlark.net/starlark"
)

type Function func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)

type Module starlark.StringDict

func NewModule(name string, funcs map[string]Function) Module {
	m := Module{}
	for k, f := range funcs {
		if f == nil {
			f = notImplemented(name, k)
		}
		m[k] = starlark.NewBuiltin(name+"."+k, f)
	}
	return m
}

func (m Module) Hash() (uint32, error) {
	return 0, fmt.Errorf("library is unhashable")
}

func (m Module) Freeze() {}

func (m Module) String() string {
	return ""
}

func (m Module) Type() string {
	return "module"
}

func (m Module) Truth() starlark.Bool {
	return starlark.True
}

func (m Module) Attr(name string) (starlark.Value, error) {
	fn, ok := m[name]
	if ok && fn != nil {
		return fn, nil
	}

	return nil, nil
}

// AttrNames returns the set of methods provided by the library.
func (m Module) AttrNames() []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// ToStarlarkValue converts a Go value to the corresponding Starlark value
func ToStarlarkValue(value interface{}) (starlark.Value, error) {
	switch v := value.(type) {
	case nil:
		return starlark.None, nil
	case bool:
		return starlark.Bool(v), nil
	case int:
		return starlark.MakeInt(v), nil
	case int64:
		return starlark.MakeInt64(v), nil
	case float64:
		if v == float64(int64(v)) {
			return starlark.MakeInt(int(v)), nil
		}
		return starlark.Float(v), nil
	case string:
		return starlark.String(v), nil
	case []string:
		list := starlark.NewList(nil)
		for _, elem := range v {
			list.Append(starlark.String(elem))
		}
		return list, nil
	case []int:
		list := starlark.NewList(nil)
		for _, elem := range v {
			list.Append(starlark.MakeInt(elem))
		}
		return list, nil
	case []any:
		list := starlark.NewList(nil)
		for _, elem := range v {
			stValue, err := ToStarlarkValue(elem)
			if err != nil {
				return nil, err
			}
			list.Append(stValue)
		}
		return list, nil
	case []map[string]interface{}:
		list := starlark.NewList(nil)
		for _, elem := range v {
			stValue, err := ToStarlarkValue(elem)
			if err != nil {
				return nil, err
			}
			list.Append(stValue)
		}
		return list, nil
	case map[string]interface{}:
		dict := starlark.NewDict(len(v))
		for key, elem := range v {
			stKey := starlark.String(key)
			stValue, err := ToStarlarkValue(elem)
			if err != nil {
				return nil, err
			}
			dict.SetKey(stKey, stValue)
		}
		return dict, nil
	default:
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value))
	}
}

// ToGolangValue converts a Starlark value to a Go interface
func ToGolangValue(val starlark.Value) (interface{}, error) {
	switch v := val.(type) {
	case starlark.NoneType:
		return nil, nil
	case starlark.Bool:
		return bool(v), nil
	case starlark.Int:
		i, ok := v.Int64()
		if !ok {
			return nil, fmt.Errorf("int value %s is 'not exactly representable'", v.String())
		}
		return i, nil
	case starlark.Float:
		return float64(v), nil
	case starlark.String:
		return string(v), nil
	case starlark.Tuple:
		var list []interface{}
		iter := v.Iterate()
		defer iter.Done()
		var x starlark.Value
		for iter.Next(&x) {
			i, err := ToGolangValue(x)
			if err != nil {
				return nil, err
			}
			list = append(list, i)
		}
		return list, nil
	case *starlark.List:
		var list []interface{}
		iter := v.Iterate()
		defer iter.Done()
		var x starlark.Value
		for iter.Next(&x) {
			i, err := ToGolangValue(x)
			if err != nil {
				return nil, err
			}
			list = append(list, i)
		}
		return list, nil
	case *starlark.Dict:
		dict := make(map[string]interface{})
		iter := v.Iterate()
		defer iter.Done()
		var x starlark.Value
		for iter.Next(&x) {
			k, ok := x.(starlark.String)
			if !ok {
				return nil, fmt.Errorf("cannot use %s as key in json dictionary", reflect.TypeOf(x))
			}
			val, ok, _ := v.Get(x)
			if !ok {
				continue
			}
			i, err := ToGolangValue(val)
			if err != nil {
				return nil, err
			}
			dict[k.GoString()] = i
		}
		return dict, nil
	default:
		// Handle other Starlark types or return an error as needed
		return nil, fmt.Errorf("unsupported Starlark type: %v", val.Type())
	}
}

func notImplemented(module, name string) Function {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		return nil, fmt.Errorf("%s.%s not implemented", module, name)
	}
}
