package gonfig

import (
	"errors"
	"strings"
	"sync/atomic"
	"unsafe"
)

// NonBindedBool is returned by method Val
// if Int is not initialized yet.
var NonBindedBool bool = false

// A Bool implements Valuer as atomic bool,
// implemented using int32 inside.
type Bool struct {
	ref *int32
}

// ErrInvalidBool indicated failed parsing.
var ErrInvalidBool = errors.New("invalid value. Expected true or false")

// NewBool returns atomic bool.
func NewBool() *Bool {
	return &Bool{new(int32)}
}

// Kind return ABool.
func (a *Bool) Kind() AKind {
	return ABool
}

// Set assigns value atomically.Initializes if was not before.
func (a *Bool) Set(b bool) {

	var i int32
	if b {
		i = 1
	}

	ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref)))
	if ptr != nil {
		atomic.StoreInt32((*int32)(ptr), i)
		return
	}

	n := NewBool()
	a.Bind(n)
	atomic.StoreInt32(n.ref, int32(i))
}

// Val returns value atomically. Returns NonBindedBool
// if it's not binded to params container.
func (a *Bool) Val() bool {
	ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref)))
	if ptr == nil {
		return NonBindedBool
	}
	return int(atomic.LoadInt32((*int32)(ptr))) == 1
}

// IsBinded returns true if Bool bineded to params container.
func (a *Bool) IsBinded() bool {
	return atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref))) != nil
}

// Bind binds current atomic variable to variable identified by to.
// As a result two Bool address to the same variable.
func (a *Bool) Bind(b *Bool) {
	ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&b.ref)))
	if ptr != nil {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref)), ptr)
	}
}

// String implements Stringer interface.
func (a *Bool) String() string {
	if a.Val() {
		return "true"
	}
	return "false"
}

// MarshalJSON implement Marshaller interface.
func (a Bool) MarshalJSON() ([]byte, error) {
	return []byte(a.String()), nil
}

// UnmarshalJSON implement Unmarshaller interface.
func (a *Bool) UnmarshalJSON(buf []byte) error {
	s := strings.ToLower(string(buf))
	m := map[string]bool{"true": true, "false": false}
	v, ok := m[s]
	if !ok {
		return ErrInvalidBool
	}
	a.Set(v)
	return nil
}

// Parse converts input argument and assigns to value.
// Accepts Y,N, T,F, TRUE,FALSE, YES,NO, 1,0 in any register.
func (a *Bool) Parse(s string) error {
	s = strings.ToLower(s)
	switch s {
	case "y", "t", "true", "yes", "1":
		a.Set(true)
	case "n", "f", "false", "no", "0":
		a.Set(false)
	default:
		return ErrInvalidBool
	}
	return nil
}
