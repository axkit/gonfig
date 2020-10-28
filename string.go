package gonfig

import (
	"sync/atomic"
	"unsafe"
)

// NonBindedString is returned by method Val
// if String is not initialized yet.
var NonBindedString string = ""

// String implements atomic string.
type String struct {
	ref *string
}

// NewString returns atomic string implemented as atomic ptr.
func NewString() *String {
	return &String{ref: new(string)}
}

// Kind returns AString.
func (a *String) Kind() AKind {
	return AString
}

// Set assigns value atomically. Initializes if was not before.
func (a *String) Set(s string) {
	sp := new(string)
	*sp = s
	ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref)))
	if ptr != nil {
		atomic.StorePointer((*unsafe.Pointer)(ptr), (unsafe.Pointer)(sp))
		return
	}

	n := NewString()
	a.Bind(n)
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(n.ref)), (unsafe.Pointer)(sp))
}

// Val returns value atomically. Returns NonBindedString
// if it's not binded to params container.
func (a *String) Val() string {

	ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref)))
	if ptr == nil {
		return NonBindedString
	}
	s := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(ptr)))
	if s == nil {
		return NonBindedString
	}
	return *(*string)(s)
}

// IsBinded returns true if String bineded to params container.
func (a *String) IsBinded() bool {
	return atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref))) != nil
}

// Bind binds current atomic variable to variable identified by to.
// As a result two String address to the same variable.
func (a *String) Bind(i *String) {
	ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&i.ref)))
	if ptr != nil {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref)), ptr)
	}
}

// String implements Stringer interface.
func (a *String) String() string {
	return a.Val()
}

// MarshalJSON implement Marshaller interface.
func (a String) MarshalJSON() ([]byte, error) {
	val := a.Val()
	return []byte("\"" + val + "\""), nil
}

// UnmarshalJSON implement Unmarshaller interface.
func (a *String) UnmarshalJSON(buf []byte) error {
	if len(buf) < 2 {
		return nil
	}
	v := string(buf[1 : len(buf)-1])
	a.Set(v)
	return nil
}

// Parse implements Valuer interface. Calls Set.
func (a *String) Parse(s string) error {
	a.Set(s)
	return nil
}
