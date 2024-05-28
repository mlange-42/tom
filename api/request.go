package api

type Options struct{}

func (r *Options) ToURL(base string, loc Location) string {
	return base
}
