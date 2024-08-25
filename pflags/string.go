package pflags

import "github.com/spf13/pflag"

type StringGetter func() string

func (sg StringGetter) Value() string {
	if sg == nil {
		return ""
	}

	return sg()
}

func String(key, defaultVal, usage string) StringGetter {
	pflag.String(key, defaultVal, usage)
	v.SetDefault(key, defaultVal)
	BindPFlag(key, pflag.Lookup(key))

	return func() string {
		return v.GetString(key)
	}
}

func StringP(key, shorthand, defaultVal, usage string) StringGetter {
	pflag.StringP(key, shorthand, defaultVal, usage)
	v.SetDefault(key, defaultVal)
	BindPFlag(key, pflag.Lookup(key))

	return func() string {
		return v.GetString(key)
	}
}

func StringRequired(key, usage string) StringGetter {
	requiredFlags = append(requiredFlags, key)

	return String(key, "", usage)
}
