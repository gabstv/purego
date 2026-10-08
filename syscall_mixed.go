// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build (darwin || freebsd || linux || netbsd) && (amd64 || arm64)

package purego

import (
	"runtime"
	"unsafe"
)

// SyscallMixed calls fn with separate integer and float arguments, without allocating.
// It is for C functions that take or return floats or small structs, which SyscallN cannot pass.
//
// intArgs fills the integer argument registers and then the stack, in order: on amd64
// RDI, RSI, RDX, RCX, R8 and R9, so intArgs[6:] are the stack slots; on arm64 x0-x7,
// so intArgs[8:] are the stack slots. A struct passed in memory, such as a struct larger
// than 16 bytes on amd64, is given as its words in the stack slots. floatArgs fills the
// float registers: XMM0-XMM7 on amd64, d0-d7 on arm64.
//
// It returns the integer result register and the float result registers. On amd64 only
// f1 and f2 (XMM0 and XMM1) are returned, and f3 and f4 are zero.
//
//go:uintptrescapes
func SyscallMixed(fn uintptr, intArgs *[16]uintptr, floatArgs *[8]uintptr) (r1, f1, f2, f3, f4 uintptr) {
	s := newMixedArgs(fn, intArgs, floatArgs)
	runtime_cgocall(syscallXABI0, unsafe.Pointer(s))
	r1, f1, f2, f3, f4 = s.a1, s.f1, s.f2, s.f3, s.f4
	thePool.Put(s)
	return
}

// SyscallMixedStret is like SyscallMixed, but for a C function that returns a struct of
// at most 32 bytes in memory, through a hidden pointer, and it returns the struct's words.
//
// On amd64 the hidden pointer is the first integer argument, so intArgs[0] is ignored and
// the real arguments start at intArgs[1]. On arm64 the hidden pointer goes in x8 and
// intArgs is used as is. A struct returned in registers, such as a struct of at most
// 16 bytes, or up to four floats or doubles on arm64, must use SyscallMixed instead.
//
//go:uintptrescapes
func SyscallMixedStret(fn uintptr, intArgs *[16]uintptr, floatArgs *[8]uintptr) (ret [4]uintptr) {
	s := newMixedArgs(fn, intArgs, floatArgs)
	// s is a heap object from the pool, so &s.sret does not move during the call.
	sret := uintptr(unsafe.Pointer(&s.sret))
	if runtime.GOARCH == "arm64" {
		s.arm64_r8 = sret
	} else {
		s.a1 = sret
	}
	runtime_cgocall(syscallXABI0, unsafe.Pointer(s))
	ret = s.sret
	thePool.Put(s)
	return
}

func newMixedArgs(fn uintptr, intArgs *[16]uintptr, floatArgs *[8]uintptr) *syscallArgs {
	if fn == 0 {
		panic("purego: fn is nil")
	}
	s := thePool.Get().(*syscallArgs)
	*s = syscallArgs{
		fn: fn,
		a1: intArgs[0], a2: intArgs[1], a3: intArgs[2], a4: intArgs[3],
		a5: intArgs[4], a6: intArgs[5], a7: intArgs[6], a8: intArgs[7],
		a9: intArgs[8], a10: intArgs[9], a11: intArgs[10], a12: intArgs[11],
		a13: intArgs[12], a14: intArgs[13], a15: intArgs[14], a16: intArgs[15],
		f1: floatArgs[0], f2: floatArgs[1], f3: floatArgs[2], f4: floatArgs[3],
		f5: floatArgs[4], f6: floatArgs[5], f7: floatArgs[6], f8: floatArgs[7],
	}
	return s
}
