package gonfig_test

import (
	"encoding/json"
	"testing"

	"github.com/axkit/gonfig"
)

func TestConfig(t *testing.T) {

	type API struct {
		IsAdmin gonfig.Bool `cfg:"is_admin"`
	}

	type Backend struct {
		Listen gonfig.String `cfg:"listen"`
		Port   gonfig.Int    `cfg:"port"`
		API
		G struct {
			ReportURL gonfig.String `cfg:"report_url"`
		}
		Noinit  gonfig.Int   `cfg:"noinit"`
		Percent gonfig.Float `cfg:"float-test"`

		Balance gonfig.Float `cfg:"balance"`

		NonValuer int `cfg:"nonvaluer"`
	}

	type param struct {
		ID       string
		Type     gonfig.AKind
		RawValue string
	}

	input := []param{
		{
			ID:       "listen",
			Type:     gonfig.AString,
			RawValue: "127.0.0.1",
		},
		{
			ID:       "port",
			Type:     gonfig.AInt,
			RawValue: "8080",
		},
		{
			ID:       "report_url",
			Type:     gonfig.AString,
			RawValue: "https://report.mdevteam.com",
		},
		{
			ID:       "is_admin",
			Type:     gonfig.ABool,
			RawValue: "true",
		},
		{
			ID:       "border",
			Type:     gonfig.AFloat,
			RawValue: "10.11",
		},
	}

	var b Backend
	var c Backend

	cfg := gonfig.New()
	for i := range input {
		if err := cfg.MustParam(input[i].ID, input[i].Type).Parse(input[i].RawValue); err != nil {
			t.Error(err)
		}
	}

	if err := cfg.MustParam("float-test", gonfig.AFloat).Parse("1.1"); err != nil {
		t.Error(err)
		return
	}

	cfg.BindStruct(&b)
	cfg.BindStruct(&c)

	p, ok := cfg.Get("float-test")
	if !ok {
		t.Error("Get failed")
	}
	if p.Kind() != gonfig.AFloat {
		t.Error("Get failed. Wrong kind")
	}

	l, ok := cfg.Get("listen")
	if !ok {
		t.Error("get failed. ")
	}
	if err := l.Parse("192.160.0.1"); err != nil {
		t.Error("parse string failed")
	}

	var listen gonfig.String
	if err := cfg.BindVar("listen", &listen); err != nil {
		t.Error(err)
	}

	listen.Set("127.0.0.2")

	if listen.Val() != b.Listen.Val() || listen.Val() != c.Listen.Val() {
		t.Error("sync failed")
	}
	t.Log("b.Listen", b.Listen.String())

	//t.Logf("%#v\n", b)
	t.Logf("b.listen:%s port:%d, reporturl=%s, is_admin=%t", b.Listen.Val(), b.Port.Val(), b.G.ReportURL.Val(), b.IsAdmin.Val())

	b.Port.Set(8081)
	t.Logf("b.listen:%s port:%d, reporturl=%s, is_admin=%t", b.Listen.Val(), b.Port.Val(), b.G.ReportURL.Val(), b.IsAdmin.Val())

	cfg.MustParam("port", gonfig.AInt).Parse("8082")
	b.IsAdmin.Set(false)
	t.Logf("b.listen:%s port:%d, reporturl=%s, is_admin=%t", b.Listen.Val(), b.Port.Val(), b.G.ReportURL.Val(), b.IsAdmin.Val())
	t.Logf("c.listen:%s port:%d, reporturl=%s, is_admin=%t", c.Listen.Val(), c.Port.Val(), c.G.ReportURL.Val(), c.IsAdmin.Val())

	var port gonfig.Int
	if err := cfg.BindVar("port", &port); err != nil {
		t.Error(err)
	}

	t.Logf("noinit: %d", b.Noinit.Val())
	t.Logf("port-var: %d", port.Val())
	port.Set(8083)
	t.Logf("port-var: %d", port.Val())
	t.Logf("c.listen:%s port:%d, reporturl=%s, is_admin=%t", c.Listen.Val(), c.Port.Val(), c.G.ReportURL.Val(), c.IsAdmin.Val())
	t.Logf("test-float: %f", c.Percent.Val())
}

func TestInt_UnmarshalJSON(t *testing.T) {
	var a struct {
		A gonfig.Int
	}

	err := json.Unmarshal([]byte(`{"A": 42}`), &a)
	if err != nil {
		t.Error(err)
	}
	t.Log(a.A.Val())
}

func TestConfig_BindStructInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
		} else {
			t.Error("panic expected")
		}
	}()

	cfg := gonfig.New()
	var h int

	cfg.BindStruct(&h)
}

func TestConfig_IsExist(t *testing.T) {

	cfg := gonfig.New()

	if cfg.IsExist("today_limit") {
		t.Error("expected false, got true")
	}

	cfg.MustParam("today_limit", gonfig.AFloat).Parse("987.65")

	if !cfg.IsExist("today_limit") {
		t.Error("expected true, got false")
	}
}

func TestConfig_Get(t *testing.T) {

	cfg := gonfig.New()

	_, ok := cfg.Get("today_limit")
	if ok {
		t.Error("expected nil, got Valuer")
	}

	cfg.MustParam("today_limit", gonfig.AFloat).Parse("987.65")

	tl, ok := cfg.Get("today_limit")
	if !ok {
		t.Error("expected Valuer, got nil")
	}

	if tl.(*gonfig.Float).Val() != 987.65 {
		t.Error("wrong value stored")
	}
}

func TestConfig_Walk(t *testing.T) {
	cfg := gonfig.New()
	cfg.MustParam("today_limit", gonfig.AFloat).Parse("987.65")
	cfg.MustParam("is_virtual", gonfig.ABool).Parse("true")
	cfg.MustParam("stop_loss", gonfig.AInt).Parse("70")

	keys := map[string]bool{"today_limit": false, "is_virtual": false, "stop_loss": false}

	var tl gonfig.Float
	if err := cfg.BindVar("today_limit", &tl); err != nil {
		t.Error(err)
	}
	if tl.Val() != 987.65 {
		t.Error("BindVar failed")
	}

	cfg.Walk(func(code string, v gonfig.Valuer, inited, asked int) {
		keys[code] = true
		println("code=", code, inited, asked)
	})

	for k, v := range keys {
		if v == false {
			t.Errorf("lost param '%s'", k)
		}
	}
}

func TestConfig_BindStructNotInitied(t *testing.T) {

	type atype struct {
		A gonfig.Int `cfg:"a"`
	}

	var a, b atype

	cfg := gonfig.New()

	if a.A.IsBinded() {
		t.Error("it's not binded yet")
	}

	cfg.BindStruct(&a)

	if !a.A.IsBinded() {
		t.Error("it's binded already")
	}

	cfg.BindStruct(&b)

	a.A.Set(1)
	if b.A.Val() != 1 {
		t.Error("not shared a single memory")
	}

	_, err := cfg.Param("a", gonfig.AFloat)
	if err == nil {
		t.Error("rewriting akind failed")
	}

	v, ok := cfg.Get("a")
	if !ok {
		t.Error("existing param not found")
	}

	if err := v.Parse("3"); err != nil {
		t.Error(err)
	}

	if b.A.Val() != 3 {
		t.Error("post parse failed", a.A.Val(), b.A.Val())
	}
	t.Log(a.A.Val())
}

/*
func TestConfig_Race(t *testing.T) {

	type Backend struct {
		Listen gonfig.String `cfg:"listen"`
		Port   gonfig.Int    `cfg:"port"`
	}

	var b Backend

	mockstore, _ := cfgstore.NewMockStorage()
	cfg := cfgstore.New(mockstore)

	if err := cfg.Cache(); err != nil {
		t.Error(err)
	}

	cfg.BindStruct(&b)

	go func() {
		var i int64
		for {
			i++
			b.Listen.Set("aaa")
		}
	}()

	go func() {
		var i int64
		for {
			i++
			b.Listen.Set("bbb")
		}
	}()
	var aa, bb int64

	go func() {
		defer func() {
			t.Log("a, b = ", aa, bb)
		}()
		for {
			k := b.Listen.Val()
			if k == "aaa" {
				aa++
			}
			if k == "bbb" {
				bb++
			}
		}
	}()

	//time.Sleep(10 * time.Second)
}
*/
