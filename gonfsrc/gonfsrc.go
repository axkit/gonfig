package gonfsrc

import (
	"github.com/axkit/gonfig"
)

// ConfigSource is an interface wrapping a single method ApplyTo.
//
// ApplyTo reads parameters from the source: database, file, env, etc.
// and adds them into config param container. Value overwrites if overwrite is true.
type ConfigSource interface {
	ApplyTo(g gonfig.Configer, overwrite bool) error
}
