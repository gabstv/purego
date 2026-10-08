// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build darwin || freebsd || linux || netbsd || windows

package purego

// The fixed-arity variants of SyscallN below do not allocate: a variadic argument slice
// always escapes to the heap because SyscallN has //go:uintptrescapes, but the slice
// literals here stay on the stack.

// Syscall0 is like [SyscallN] but takes exactly zero arguments.
//
//go:uintptrescapes
func Syscall0(fn uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, nil)
}

// Syscall1 is like [SyscallN] but takes exactly one argument.
//
//go:uintptrescapes
func Syscall1(fn, a1 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1})
}

// Syscall2 is like [SyscallN] but takes exactly two arguments.
//
//go:uintptrescapes
func Syscall2(fn, a1, a2 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1, a2})
}

// Syscall3 is like [SyscallN] but takes exactly three arguments.
//
//go:uintptrescapes
func Syscall3(fn, a1, a2, a3 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1, a2, a3})
}

// Syscall4 is like [SyscallN] but takes exactly four arguments.
//
//go:uintptrescapes
func Syscall4(fn, a1, a2, a3, a4 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1, a2, a3, a4})
}

// Syscall5 is like [SyscallN] but takes exactly five arguments.
//
//go:uintptrescapes
func Syscall5(fn, a1, a2, a3, a4, a5 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1, a2, a3, a4, a5})
}

// Syscall6 is like [SyscallN] but takes exactly six arguments.
//
//go:uintptrescapes
func Syscall6(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1, a2, a3, a4, a5, a6})
}

// Syscall7 is like [SyscallN] but takes exactly seven arguments.
//
//go:uintptrescapes
func Syscall7(fn, a1, a2, a3, a4, a5, a6, a7 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1, a2, a3, a4, a5, a6, a7})
}

// Syscall8 is like [SyscallN] but takes exactly eight arguments.
//
//go:uintptrescapes
func Syscall8(fn, a1, a2, a3, a4, a5, a6, a7, a8 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1, a2, a3, a4, a5, a6, a7, a8})
}

// Syscall9 is like [SyscallN] but takes exactly nine arguments.
//
//go:uintptrescapes
func Syscall9(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1, a2, a3, a4, a5, a6, a7, a8, a9})
}

// Syscall10 is like [SyscallN] but takes exactly ten arguments.
//
//go:uintptrescapes
func Syscall10(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1, a2, a3, a4, a5, a6, a7, a8, a9, a10})
}

// Syscall11 is like [SyscallN] but takes exactly eleven arguments.
//
//go:uintptrescapes
func Syscall11(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11})
}

// Syscall12 is like [SyscallN] but takes exactly twelve arguments.
//
//go:uintptrescapes
func Syscall12(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12})
}

// Syscall13 is like [SyscallN] but takes exactly thirteen arguments.
//
//go:uintptrescapes
func Syscall13(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13})
}

// Syscall14 is like [SyscallN] but takes exactly fourteen arguments.
//
//go:uintptrescapes
func Syscall14(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14})
}

// Syscall15 is like [SyscallN] but takes exactly fifteen arguments.
//
//go:uintptrescapes
func Syscall15(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15 uintptr) (r1, r2, err uintptr) {
	return syscallN(fn, []uintptr{a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15})
}
