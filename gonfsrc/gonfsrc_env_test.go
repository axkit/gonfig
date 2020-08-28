package gonfsrc_test

import (
	"testing"

	"github.com/axkit/gonfig"
	"github.com/axkit/gonfig/gonfsrc"
)

func TestEnvSource_CopyTo(t *testing.T) {

	cfg := gonfig.New()

	if err := gonfsrc.NewEnvSource("GO", true).ApplyTo(cfg, false); err != nil {
		t.Error(err)
	}

	if cfg.IsExist("gopath") == false {
		t.Error("no gopath var")
	}
}
