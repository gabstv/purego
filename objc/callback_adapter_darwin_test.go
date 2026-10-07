// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

package objc_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

var callbackAdapterTestClassID atomic.Uint64

// TestIMPNoAllocs checks that Objective-C calling a Go method with a common signature
// does not allocate, because NewIMP dispatches it without reflection.
//
// SyscallN allocates its argument slice by itself, so the test compares against a call
// of the same shape into a method implemented in C.
func TestIMPNoAllocs(t *testing.T) {
	addSel := objc.RegisterName("callbackAdapterAdd:")
	positiveSel := objc.RegisterName("callbackAdapterIsPositive:")
	// Tests run several times in one process, and a class name can be registered only once.
	name := fmt.Sprintf("PuregoCallbackAdapterTest%d", callbackAdapterTestClassID.Add(1))
	class, err := objc.RegisterClass(
		name,
		objc.GetClass("NSObject"),
		nil,
		nil,
		[]objc.MethodDef{
			{
				Cmd: addSel,
				Fn: func(self objc.ID, _ objc.SEL, arg objc.ID) objc.ID {
					return arg + 1
				},
			},
			{
				Cmd: positiveSel,
				Fn: func(self objc.ID, _ objc.SEL, arg objc.ID) bool {
					return arg > 0
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	obj := objc.ID(class).Send(objc.RegisterName("new"))

	lib, err := purego.Dlopen("/usr/lib/libobjc.A.dylib", purego.RTLD_GLOBAL)
	if err != nil {
		t.Fatal(err)
	}
	msgSend, err := purego.Dlsym(lib, "objc_msgSend")
	if err != nil {
		t.Fatal(err)
	}

	if r, _, _ := purego.SyscallN(msgSend, uintptr(obj), uintptr(addSel), 41); r != 42 {
		t.Fatalf("callbackAdapterAdd:41 = %d, want 42", r)
	}
	// Only the low byte of a BOOL result is defined.
	if r, _, _ := purego.SyscallN(msgSend, uintptr(obj), uintptr(positiveSel), 1); r&0xff != 1 {
		t.Fatalf("callbackAdapterIsPositive:1 = %d, want 1", r&0xff)
	}
	if r, _, _ := purego.SyscallN(msgSend, uintptr(obj), uintptr(positiveSel), 0); r&0xff != 0 {
		t.Fatalf("callbackAdapterIsPositive:0 = %d, want 0", r&0xff)
	}

	method := testing.AllocsPerRun(1000, func() {
		purego.SyscallN(msgSend, uintptr(obj), uintptr(addSel), 1)
	})
	// isEqual: is implemented by NSObject in C.
	isEqualSel := objc.RegisterName("isEqual:")
	baseline := testing.AllocsPerRun(1000, func() {
		purego.SyscallN(msgSend, uintptr(obj), uintptr(isEqualSel), uintptr(obj))
	})
	if method != baseline {
		t.Errorf("calling the Go method allocated %v times per call, want %v (the same as calling a C method)", method, baseline)
	}
}

// TestBlockNoAllocs checks that invoking a Go block with a common function type
// does not allocate, because NewBlock dispatches it without reflection.
func TestBlockNoAllocs(t *testing.T) {
	lib, err := purego.Dlopen("/usr/lib/libSystem.B.dylib", purego.RTLD_GLOBAL)
	if err != nil {
		t.Fatal(err)
	}
	getGlobalQueue, err := purego.Dlsym(lib, "dispatch_get_global_queue")
	if err != nil {
		t.Fatal(err)
	}
	dispatchSync, err := purego.Dlsym(lib, "dispatch_sync")
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

	var calls int
	block := objc.NewBlock(func(objc.Block) {
		calls++
	})
	defer block.Release()

	// dispatch_sync invokes the block synchronously.
	purego.SyscallN(dispatchSync, queue, uintptr(block))
	if calls != 1 {
		t.Fatalf("block was called %d times, want 1", calls)
	}

	invoke := testing.AllocsPerRun(1000, func() {
		purego.SyscallN(dispatchSync, queue, uintptr(block))
	})
	// dispatch_sync_f with free(NULL) is the same kind of call without a Go block.
	baseline := testing.AllocsPerRun(1000, func() {
		purego.SyscallN(dispatchSyncF, queue, 0, free)
	})
	if invoke != baseline {
		t.Errorf("invoking the Go block allocated %v times per call, want %v (the same as calling a C function)", invoke, baseline)
	}
}
