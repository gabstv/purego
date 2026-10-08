// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build (darwin || linux) && (amd64 || arm64)

package purego_test

import (
	"math"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/internal/load"
	"github.com/ebitengine/purego/internal/testlib"
)

func libmPath() string {
	switch runtime.GOOS {
	case "android":
		return "libm.so"
	case "darwin":
		return "/usr/lib/libSystem.B.dylib"
	default:
		return "libm.so.6"
	}
}

func TestSyscallMixed(t *testing.T) {
	libm, err := purego.Dlopen(libmPath(), purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		t.Fatal(err)
	}
	fmax, err := purego.Dlsym(libm, "fmax")
	if err != nil {
		t.Fatal(err)
	}
	ldexp, err := purego.Dlsym(libm, "ldexp")
	if err != nil {
		t.Fatal(err)
	}

	// double fmax(double, double): two float arguments and a float result.
	var ints [16]uintptr
	floats := [8]uintptr{uintptr(math.Float64bits(2.5)), uintptr(math.Float64bits(7.5))}
	_, f1, _, _, _ := purego.SyscallMixed(fmax, &ints, &floats)
	if got := math.Float64frombits(uint64(f1)); got != 7.5 {
		t.Errorf("fmax(2.5, 7.5) = %v, want 7.5", got)
	}

	// double ldexp(double, int): a float and an integer argument.
	ints = [16]uintptr{3}
	floats = [8]uintptr{uintptr(math.Float64bits(1.5))}
	_, f1, _, _, _ = purego.SyscallMixed(ldexp, &ints, &floats)
	if got := math.Float64frombits(uint64(f1)); got != 12 {
		t.Errorf("ldexp(1.5, 3) = %v, want 12", got)
	}

	if n := testing.AllocsPerRun(100, func() { purego.SyscallMixed(fmax, &ints, &floats) }); n != 0 {
		t.Errorf("SyscallMixed: got %v allocs per call, want 0", n)
	}
}

func TestSyscallMixedStret(t *testing.T) {
	libFileName := filepath.Join(t.TempDir(), "abitest.so")
	if err := testlib.BuildSharedLib(t, "CC", libFileName, filepath.Join("testdata", "abitest", "abi_test.c")); err != nil {
		t.Fatal(err)
	}
	lib, err := load.OpenLibrary(libFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := load.CloseLibrary(lib); err != nil {
			t.Error(err)
		}
	}()
	fn, err := load.OpenSymbol(lib, "return_three_words")
	if err != nil {
		t.Fatal(err)
	}

	var ints [16]uintptr
	var floats [8]uintptr
	if runtime.GOARCH == "amd64" {
		ints[1] = 10 // ints[0] is the hidden result pointer
	} else {
		ints[0] = 10
	}
	got := purego.SyscallMixedStret(fn, &ints, &floats)
	if want := [4]uintptr{10, 11, 12, 0}; got != want {
		t.Errorf("return_three_words(10) = %v, want %v", got, want)
	}

	if n := testing.AllocsPerRun(100, func() { purego.SyscallMixedStret(fn, &ints, &floats) }); n != 0 {
		t.Errorf("SyscallMixedStret: got %v allocs per call, want 0", n)
	}
}
