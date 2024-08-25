package pflags

import "github.com/spf13/pflag"

type Float64Getter func() float64

func (fg Float64Getter) Value() float64 {
	if fg == nil {
		return 0
	}

	return fg()
}

func Float64(key string, defaultVal float64, usage string) Float64Getter {
	pflag.Float64(key, defaultVal, usage)
	v.SetDefault(key, defaultVal)
	BindPFlag(key, pflag.Lookup(key))

	return func() float64 {
		return v.GetFloat64(key)
	}
}

func Float64Required(key, usage string) Float64Getter {
	requiredFlags = append(requiredFlags, key)

	return Float64(key, 0, usage)
}
