package server

import (
	"github.com/valyala/fasthttp"
)

type headerWraper struct {
	raw *fasthttp.RequestHeader
}

func (hw *headerWraper) Get(key string) string {
	return string(hw.raw.Peek(key))
}

func (hw *headerWraper) Set(key string, value string) {
	hw.raw.Set(key, value)
}

func (hw *headerWraper) Keys() []string {
	raw := hw.raw.PeekKeys()

	res := make([]string, len(raw))
	for i := range raw {
		res[i] = string(raw[i])
	}

	return res
}
