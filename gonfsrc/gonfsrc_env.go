package gonfsrc

import (
	"os"
	"strings"

	"github.com/axkit/gonfig"
)

// EnvSource implements logic or reading application parameters
// from the environment variables.
type EnvSource struct {
	prefix string

	tolower bool
}

// NewEnvSource returns EnvSource. if tolower is true, the envvar code
// will be lower cased before applying to config container.
func NewEnvSource(prefix string, tolower bool) *EnvSource {
	return &EnvSource{prefix: strings.ToUpper(prefix), tolower: tolower}
}

// CopyTo copies environment variables starting with prefix.
func (s *EnvSource) ApplyTo(g gonfig.Configer, ow bool) error {
	return s.applyTo(g, ow)
}

func (s *EnvSource) applyTo(g gonfig.Configer, ow bool) error {

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if !strings.HasPrefix(pair[0], s.prefix) {
			continue
		}

		code := pair[0]
		if s.tolower {
			code = strings.ToLower(code)
		}

		var err error
		if !ow && g.IsExist(code) {
			continue
		}

		err = g.MustParam(code, gonfig.AString).Parse(pair[1])
		if err != nil {
			return err
		}
	}
	return nil
}
