package gonfig

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// AKind represents possible kinds of config param data type.
// The zero kind is not a valid kind.
type AKind uint8

const (
	// Unknown represents not specified kind (invalid).
	Unknown AKind = 0

	// AInt represents atomic int. It's int64 internally.
	AInt AKind = 1

	// ABool represents atomic bool. It's int32 internally.
	ABool AKind = 2

	// AString represents atomic string. It's pointers to string internally.
	AString AKind = 3

	// AFloat represents atomic float. It's float64 internally.
	AFloat AKind = 4
)

type action uint8

const (
	added action = 1
	asked action = 2
)

// String implements Stringer interface.
func (ak AKind) String() string {
	switch ak {
	case AInt:
		return "AInt"
	case ABool:
		return "ABool"
	case AString:
		return "AString"
	case AFloat:
		return "AFloat"
	}
	return "Unknown"
}

// A Configer is an interface what wraps following methods:
//
// Param returns Valuer identified by code. Creates if not exist.
// Returns error if exists but Kind is different.
//
// MustParam returns Valuer identified by code. Creates if not exist.
// Panics if exists but Kind is different.
//
// IsExist returns true if params identified by code is exists.
//
// Get returns Valuer by code.
//
// BindStruct walks through all struct fields and bind then
// if they implement Valuer interface.
//
// BindVar binds a single var implementing Valuer interface.
type Configer interface {
	Param(code string, ak AKind) (Valuer, error)
	MustParam(code string, ak AKind) Valuer

	IsExist(code string) bool
	Get(code string) (Valuer, bool)

	BindStruct(structAddr interface{}) []error
	BindVar(code string, v Valuer) error
	Walk(func(code string, v Valuer, inited, asked int))
}

// Valuer is an interface what wraps following methods.
//
// Kind returns data type of Valuer.
//
// Parse converts string value to internal representation.
// in accordance with Kind.
//
// IsBinded returns false is Valuer is not binded to
// Configer.
type Valuer interface {
	Kind() AKind
	Parse(string) error
	IsBinded() bool
}

type param struct {
	code   string
	av     Valuer
	inited int
	asked  int
}

// Config is in-memory config params container.
// On init step Config accepts all param and values.
// Later all Valuers bind themselves to Config.
type Config struct {
	mux  sync.RWMutex
	list []param
	idx  map[string]int
}

// New returns new container of config parameters.
// It's ok to have a single instance for the whole application.
func New() Configer {
	return &Config{idx: make(map[string]int)}
}

// IsExist returns true if parameter identified by code is in container.
func (c *Config) IsExist(code string) bool {
	c.mux.RLock()
	defer c.mux.RUnlock()
	_, ok := c.idx[code]
	return ok
}

// BindVar binds an Valuer to container. Creates param if not found.
func (c *Config) BindVar(code string, addr Valuer) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	p, err := c.param(code, addr.Kind(), false)
	if err != nil {
		return err
	}

	switch addr.Kind() {
	case AInt:
		addr.(*Int).Bind(p.(*Int))
	case ABool:
		addr.(*Bool).Bind(p.(*Bool))
	case AString:
		addr.(*String).Bind(p.(*String))
	case AFloat:
		addr.(*Float).Bind(p.(*Float))
	default:
		return ErrDifferentKind
	}

	c.setStat(code, asked)

	return nil
}

// BindStruct binds struct's attributes implementing interface Valuer
// to params from container. If param is not containerized yet
// it adds param into container.
//
// Mapping between fields and params in container via tag "cfg".
// Example
// type Listener struct {
//		Port gonfig.Int `cfg:"port"`
// }
//
// BindStruct works properly with fields as structs and
// embedded anonymous structs.
func (c *Config) BindStruct(structAddr interface{}) []error {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.bindStruct(structAddr)
}

func (c *Config) bindStruct(structAddr interface{}) []error {

	var res []error
	s := reflect.ValueOf(structAddr).Elem()
	tof := s.Type()

	if s.Kind() != reflect.Struct {
		panic("expected argument as reference to struct")
	}

	// create empty object complaint with interface Valuer.
	atype := reflect.TypeOf((*Valuer)(nil)).Elem()

	for i := 0; i < s.NumField(); i++ {
		f := tof.Field(i)

		if s.Field(i).Addr().Type().Implements(atype) {
			code := tof.Field(i).Tag.Get("cfg")
			if code == "" {
				continue
			}
			fai := s.Field(i).Addr().Interface()

			p, err := c.param(code, fai.(Valuer).Kind(), false)
			if err != nil {
				res = append(res, err)
				continue
			}

			defv := tof.Field(i).Tag.Get("default")
			_, ok := c.idx[code]
			if !ok && defv != "" {
				if err := p.Parse(defv); err != nil {
					res = append(res, err)
					continue
				}
			}

			switch p.Kind() {
			case ABool:
				fai.(*Bool).Bind(p.(*Bool))
			case AString:
				fai.(*String).Bind(p.(*String))
			case AInt:
				fai.(*Int).Bind(p.(*Int))
			case AFloat:
				fai.(*Float).Bind(p.(*Float))
			}
			c.setStat(code, asked)
			continue
		}

		if f.Anonymous || s.Field(i).Kind() == reflect.Struct {
			if f.Type.Kind() != reflect.Ptr {
				c.bindStruct(s.Field(i).Addr().Interface())
			} else {
				c.bindStruct(s.Field(i).Interface())
			}
		}
	}
	return res
}

// ErrDifferentKind indicates raises when Set trying
// overwrite value with different AKind.
var ErrDifferentKind = errors.New("different value kind")

// Param returns Valuer identified by code. Creates if not exist.
// Returns error if exists but Kind is different.
func (c *Config) Param(code string, ak AKind) (Valuer, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	v, err := c.param(code, ak, false)
	if err == nil {
		c.setStat(code, added)
	}
	return v, err
}

// MustParam returns Valuer identified by code. Creates if not exist.
// Panics if exists but Kind is different.
func (c *Config) MustParam(code string, ak AKind) Valuer {
	c.mux.Lock()
	defer c.mux.Unlock()
	res, _ := c.param(code, ak, true)
	c.setStat(code, added)
	return res
}

func (c *Config) setStat(code string, a action) {
	idx, ok := c.idx[code]
	if ok {
		if a&asked == asked {
			c.list[idx].asked++
		}
		if a&added == added {
			c.list[idx].inited++
		}
	}
}

func (c *Config) param(code string, ak AKind, dopanic bool) (Valuer, error) {
	idx, ok := c.idx[code]
	if ok {
		oak := c.list[idx].av.Kind()
		if oak == ak {
			return c.list[idx].av, nil
		}

		msg := fmt.Sprintf("different value kind of param '%s', wanted %s, got %s ", code, oak, ak)
		if dopanic {
			panic(msg)
		}
		return nil, errors.New(msg)
	}

	p := param{code: code, av: makeValuer(ak)}
	c.list = append(c.list, p)
	c.idx[code] = len(c.list) - 1
	return p.av, nil
}

// Get returns Valuer instance by code.
func (c *Config) Get(code string) (Valuer, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	idx, ok := c.idx[code]
	if ok {
		return c.list[idx].av, true
	}
	return nil, false
}

func makeValuer(ak AKind) Valuer {

	var res Valuer
	switch ak {
	case AInt:
		res = NewInt()
	case ABool:
		res = NewBool()
	case AString:
		res = NewString()
	case AFloat:
		res = NewFloat()
	default:
		panic("invalid AKind")
	}

	return res
}

// Walk calls function f() for every parameter in container.
func (c *Config) Walk(f func(code string, v Valuer, init, asked int)) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	for i := range c.list {
		f(c.list[i].code, c.list[i].av, c.list[i].inited, c.list[i].asked)
	}
}
