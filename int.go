package gonfig

import (
	"strconv"
	"strings"
	"sync/atomic"
	"unsafe"
)

// NonBindedInt is returned by method Val
// if Int is not initialized yet.
var NonBindedInt int = 0

// A Int implements atomic int.
type Int struct {
	ref *int64
}

// NewInt returns atomic int implemented using int64.
func NewInt() *Int {
	return &Int{new(int64)}
}

// Kind returns AInt.
func (a *Int) Kind() AKind {
	return AInt
}

// Set assigns value atomically. Initializes if was not before.
func (a *Int) Set(i int) {
	ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref)))
	if ptr != nil {
		atomic.StoreInt64((*int64)(ptr), int64(i))
		return
	}

	n := NewInt()
	a.Bind(n)
	atomic.StoreInt64(n.ref, int64(i))
}

// Val returns value atomically. Returns NonBindedInt
// if it's not binded to params container.
func (a *Int) Val() int {
	ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref)))
	if ptr == nil {
		return NonBindedInt
	}
	return int(atomic.LoadInt64((*int64)(ptr)))
}

// IsBinded returns true if Int bineded to params container.
func (a *Int) IsBinded() bool {
	return atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref))) != nil
}

// Bind binds current atomic variable to variable identified by to.
// As a result two Int address to the same variable.
func (a *Int) Bind(to *Int) {
	ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&to.ref)))
	if ptr != nil {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref)), ptr)
	}
}

// String implements Stringer interface.
func (a *Int) String() string {
	return strconv.Itoa(a.Val())
}

// Parse converts input argument and assigns to value.
func (a *Int) Parse(s string) error {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		a.Set(0)
		return nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	a.Set(i)
	return nil
}

// MarshalJSON implement Marshaller interface.
func (a Int) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(a.Val())), nil
}

// UnmarshalJSON implement Unmarshaller interface.
func (a *Int) UnmarshalJSON(buf []byte) error {
	i, err := strconv.Atoi(string(buf))
	if err != nil {
		return err
	}
	a.Set(i)
	return nil
}
