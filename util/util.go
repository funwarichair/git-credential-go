package util

import (
	"bytes"
	"strings"

	g "github.com/funwarichair/git-credential-go"
)

func ToString(g *g.GitCrendential) string {
	var buf bytes.Buffer // bytes.Buffer.Write always returns (length of p, nil).
	g.Write(&buf)
	return buf.String()
}

func FromString(input string) (*g.GitCrendential, error) {
	return g.Parse(strings.NewReader(input))
}
