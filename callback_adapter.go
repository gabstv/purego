// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

package purego

// CallbackAdapter calls a callback's Go function directly, without reflection.
//
// ints holds the callback's integer argument registers in order, and the returned
// value is placed in the integer result register. An adapter must convert each
// register to the corresponding parameter type, call the function, and convert
// its result (if any) to a uintptr. Booleans are passed and returned as 0 or 1.
type CallbackAdapter func(ints []uintptr) uintptr

// uintptrAdapter returns a CallbackAdapter for common callback signatures whose
// parameters and result are all uintptr, or nil if fn has another signature.
func uintptrAdapter(fn any) CallbackAdapter {
	switch f := fn.(type) {
	case func():
		return func([]uintptr) uintptr { f(); return 0 }
	case func() uintptr:
		return func([]uintptr) uintptr { return f() }
	case func(uintptr):
		return func(a []uintptr) uintptr { f(a[0]); return 0 }
	case func(uintptr) uintptr:
		return func(a []uintptr) uintptr { return f(a[0]) }
	case func(uintptr, uintptr):
		return func(a []uintptr) uintptr { f(a[0], a[1]); return 0 }
	case func(uintptr, uintptr) uintptr:
		return func(a []uintptr) uintptr { return f(a[0], a[1]) }
	case func(uintptr, uintptr, uintptr):
		return func(a []uintptr) uintptr { f(a[0], a[1], a[2]); return 0 }
	case func(uintptr, uintptr, uintptr) uintptr:
		return func(a []uintptr) uintptr { return f(a[0], a[1], a[2]) }
	case func(uintptr, uintptr, uintptr, uintptr):
		return func(a []uintptr) uintptr { f(a[0], a[1], a[2], a[3]); return 0 }
	case func(uintptr, uintptr, uintptr, uintptr) uintptr:
		return func(a []uintptr) uintptr { return f(a[0], a[1], a[2], a[3]) }
	}
	return nil
}
