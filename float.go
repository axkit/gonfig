package gonfig

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync/atomic"
	"unsafe"
)

var (
	// NonBindedFloat is returned by method Val
	// if Float is not initialized yet.
	NonBindedFloat float64 = 0.0

	// DefaultFmtStyle used in String method.
	DefaultFmtStyle = "%f"
)

// A Float implements Valuer as atomic float64,
// implemented using uint64 inside.
type Float struct {
	ref *uint64
}

// NewFloat returns atomic float.
func NewFloat() *Float {
	return &Float{new(uint64)}
}

// Kind returns AFloat.
func (a *Float) Kind() AKind {
	return AFloat
}

// Set assigns value atomically.Initializes if was not before.
func (a *Float) Set(f float64) {
	fu := math.Float64bits(f)
	ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref)))
	if ptr != nil {
		atomic.StoreUint64((*uint64)(ptr), fu)
		return
	}
	n := NewFloat()
	a.Bind(n)
	atomic.StoreUint64(n.ref, fu)
}

// Val returns value atomically. Returns NonBindedFloat
// if it's not binded to params container.
func (a *Float) Val() float64 {
	ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref)))
	if ptr == nil {
		return NonBindedFloat
	}
	fu := atomic.LoadUint64((*uint64)(ptr))
	return math.Float64frombits(fu)
}

// IsBinded returns true if Float bineded to params container.
func (a *Float) IsBinded() bool {
	return atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref))) != nil
}

// Bind binds current atomic variable to variable identified by to.
// As a result two Float address to the same variable.
func (a *Float) Bind(to *Float) {
	ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&to.ref)))
	if ptr != nil {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&a.ref)), ptr)
	}
}

// String implements Stringer interface. Returns value as string in decimal
// format. DefaultFmtStyle is used for styling.
func (a *Float) String() string {
	return fmt.Sprintf(DefaultFmtStyle, a.Val())
}

// Parse converts input argument and assigns to value.
func (a *Float) Parse(s string) error {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		a.Set(0.0)
		return nil
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	a.Set(f)
	return nil
}

// MarshalJSON implement Marshaller interface.
func (a Float) MarshalJSON() ([]byte, error) {
	return []byte(a.String()), nil
}

// UnmarshalJSON implement Unmarshaller interface.
func (a *Float) UnmarshalJSON(buf []byte) error {
	return a.Parse(string(buf))
}
