package cmake

import (
	"reflect"
	"testing"

	"github.com/zelviner/cgear/config"
)

func TestConfigArgSupportsAllCMakeBuildTypes(t *testing.T) {
	arg := ConfigArg{
		Toolchain: &config.Toolchain{
			Compiler: config.Compiler{C: "clang", CXX: "clang++"},
		},
		BuildType: "RelWithDebInfo",
	}

	want := []string{
		"-DCMAKE_C_COMPILER:FILEPATH=clang",
		"-DCMAKE_CXX_COMPILER:FILEPATH=clang++",
		"-DCMAKE_BUILD_TYPE=RelWithDebInfo",
	}
	if got := arg.toStringSlice(); !reflect.DeepEqual(got, want) {
		t.Fatalf("toStringSlice() = %#v, want %#v", got, want)
	}
}

func TestBuildArgIncludesTargetForMSVC(t *testing.T) {
	arg := BuildArg{
		BuildPath: "build",
		Target:    "install",
		BuildType: "Release",
		IsMSVC:    true,
	}

	want := []string{"--build", "build", "--config", "Release", "--target", "install"}
	if got := arg.toStringSlice(); !reflect.DeepEqual(got, want) {
		t.Fatalf("toStringSlice() = %#v, want %#v", got, want)
	}
}
