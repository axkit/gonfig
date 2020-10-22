package gonfigenv_test

import (
	"testing"

	"github.com/axkit/gonfig"
	"github.com/axkit/gonfig/gonfigenv"
)

func TestEnvSource_CopyTo(t *testing.T) {

	cfg := gonfig.New()

	if err := gonfigenv.NewEnvSource("GO", true).ApplyTo(cfg, false); err != nil {
		t.Error(err)
	}

	if cfg.IsExist("gopath") == false {
		t.Error("no gopath var")
	}
}
