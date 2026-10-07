// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

package purego_test

import (
	"testing"

	"github.com/ebitengine/purego"
)

// TestNewCallbackUintptrNoAllocs checks that C calling a callback with a func(uintptr...)
// signature does not allocate, because NewCallback dispatches it without reflection.
//
// SyscallN allocates its argument slice by itself, so the test compares against a call
// of the same shape into a C function.
func TestNewCallbackUintptrNoAllocs(t *testing.T) {
	lib, err := purego.Dlopen("/usr/lib/libSystem.B.dylib", purego.RTLD_GLOBAL)
	if err != nil {
		t.Fatal(err)
	}
	getGlobalQueue, err := purego.Dlsym(lib, "dispatch_get_global_queue")
	if err != nil {
		t.Fatal(err)
	}
	dispatchSyncF, err := purego.Dlsym(lib, "dispatch_sync_f")
	if err != nil {
		t.Fatal(err)
	}
	free, err := purego.Dlsym(lib, "free")
	if err != nil {
		t.Fatal(err)
	}
	queue, _, _ := purego.SyscallN(getGlobalQueue, 0, 0)

	// dispatch_sync_f ignores the result. The callback returns one anyway: without
	// an adapter, a result makes the reflection path allocate.
	var got uintptr
	cb := purego.NewCallback(func(ctx uintptr) uintptr {
		got = ctx
		return ctx
	})

	// dispatch_sync_f calls cb(ctx) synchronously.
	purego.SyscallN(dispatchSyncF, queue, 42, cb)
	if got != 42 {
		t.Fatalf("callback got %d, want 42", got)
	}

	callback := testing.AllocsPerRun(1000, func() {
		purego.SyscallN(dispatchSyncF, queue, 7, cb)
	})
	// free(NULL) does nothing, so this is the cost of the call without a Go callback.
	baseline := testing.AllocsPerRun(1000, func() {
		purego.SyscallN(dispatchSyncF, queue, 0, free)
	})
	if callback != baseline {
		t.Errorf("calling the callback allocated %v times per call, want %v (the same as calling a C function)", callback, baseline)
	}
}
