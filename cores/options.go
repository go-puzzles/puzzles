package cores

import "strings"

func WithService(name string) ServiceOption {
	return func(o *Options) {
		segs := strings.SplitN(name, ":", 2)
		if len(segs) < 2 {
			o.ServiceName = name
		} else {
			o.ServiceName = segs[0]
			o.Tags = append(o.Tags, segs[1])
		}
	}
}

func WithTag(tag string) ServiceOption {
	return func(o *Options) {
		o.Tags = append(o.Tags, tag)
	}
}
