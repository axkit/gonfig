package gonfig_test

import (
	"encoding/json"
	"testing"

	"github.com/axkit/gonfig"
)

func TestBool(t *testing.T) {
	type store struct {
		IsActive gonfig.Bool `cfg:"is_active" json:"is_active"`
	}

	var a store

	cfg := gonfig.New()

	if err := cfg.MustParam("is_active", gonfig.ABool).Parse("xxx"); err == nil {
		t.Error(err)
	}

	if err := cfg.MustParam("is_active", gonfig.ABool).Parse(""); err == nil {
		t.Error(err)
	}

	if err := cfg.MustParam("is_active", gonfig.ABool).Parse("yes"); err != nil {
		t.Error(err)
	}

	cfg.BindStruct(&a)

	if a.IsActive.IsBinded() == false {
		t.Error("BindStruct() failed")
	}

	a.IsActive.Set(false)

	if s := a.IsActive.String(); s != "false" {
		t.Errorf("String() returned wrong data. expected false, got:%s", s)
	}

	if err := json.Unmarshal([]byte(`{"is_active" : true}`), &a); err != nil {
		t.Error(err)
	}

	if !a.IsActive.Val() {
		t.Error("Unmarshal() failed. expected true, got false")
	}

	if err := json.Unmarshal([]byte(`{"is_active" : "abc"}`), &a); err == nil {
		t.Error(err)
	}

	var b gonfig.Bool

	if b.IsBinded() == true {
		t.Error("is not binded yet")
	}

	if b.Val() != gonfig.NonBindedBool {
		t.Error("not binded var does not return NonBindedBool")
	}

	b.Set(true)

	if b.Val() != true {
		t.Error("Set failed")
	}

	if b.String() != "true" {
		t.Error("String failed")
	}

}

func TestInt(t *testing.T) {
	type store struct {
		A gonfig.Int `cfg:"a" json:"a"`
	}

	var a store

	cfg := gonfig.New()

	if a.A.Val() != gonfig.NonBindedInt {
		t.Error("wrong non binded Int value")
	}

	if a.A.IsBinded() == true {
		t.Error("is not binded yet")
	}

	if err := cfg.MustParam("a", gonfig.AInt).Parse("xxx"); err == nil {
		t.Error(err)
	}

	if err := cfg.MustParam("a", gonfig.AInt).Parse(""); err != nil {
		t.Error(err)
	}

	if err := cfg.MustParam("a", gonfig.AInt).Parse("1"); err != nil {
		t.Error(err)
	}

	cfg.BindStruct(&a)

	if a.A.IsBinded() == false {
		t.Error("BindStruct() failed")
	}

	a.A.Set(2)

	if s := a.A.String(); s != "2" {
		t.Errorf("String() returned wrong data. expected 2, got:%s", s)
	}

	if err := json.Unmarshal([]byte(`{"a" : 3}`), &a); err != nil {
		t.Error(err)
	}

	if a.A.Val() != 3 {
		t.Errorf("Unmarshal() failed. expected 3, got:%d", a.A.Val())
	}

	if err := json.Unmarshal([]byte(`{"a" : "abc"}`), &a); err == nil {
		t.Error(err)
	}

}

func Benchmark_intassign(b *testing.B) {
	ref := new(int)
	var i int
	for i = 0; i < b.N; i++ {
		*ref = i
	}
}

func Benchmark_IntSet(b *testing.B) {
	a := gonfig.NewInt()
	for i := 0; i < b.N; i++ {
		a.Set(i)
	}
}

func Benchmark_intread(b *testing.B) {
	ref := new(int)
	*ref = 90
	for i := 0; i < b.N; i++ {
		k := *ref
		_ = k
	}
}

func Benchmark_IntVal(b *testing.B) {
	a := gonfig.NewInt()
	a.Set(90)
	for i := 0; i < b.N; i++ {
		k := a.Val()
		_ = k
	}
}
