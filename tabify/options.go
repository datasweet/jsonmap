package tabify

import "strings"

// Options are our tabify options
type Options struct {
	KeyFormatter KeyFormatterFunc
	KeyExcluder  KeyExcluderFunc
}

// Option is an option setter
type Option func(o *Options)

// KeyFormatterFunc is a function to format a key from the input json
type KeyFormatterFunc func([]string) string

// KeyExcluderFunc is a function to exclude a key from the input json
type KeyExcluderFunc func([]string) bool

func newOptions(opt ...Option) Options {
	opts := Options{}

	for _, o := range opt {
		o(&opts)
	}

	// <!> Need at least one formatter
	if opts.KeyFormatter == nil {
		KeyFormatter(defaultFormatter)(&opts)
	}

	return opts
}

func defaultFormatter(keys []string) string {
	return strings.Join(keys, "#")
}

// KeyFormatter sets the key formatter
// default : func (keys []string) => strings.Join(keys, "#")
func KeyFormatter(v KeyFormatterFunc) Option {
	return func(opts *Options) {
		if v != nil {
			opts.KeyFormatter = v
		}
	}
}

// KeyExcluder sets the key excluder
// default : nil
func KeyExcluder(v KeyExcluderFunc) Option {
	return func(opts *Options) {
		opts.KeyExcluder = v
	}
}
