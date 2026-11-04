package executor

import (
	"bytes"
	"testing"

	ebsUtils "github.com/rexray/rexray/libstorage/drivers/storage/ebs/utils"
)

func TestNVMeDevicePath(t *testing.T) {
	tests := []struct {
		name       string
		mapping    string
		padding    byte
		largeRange bool
		want       string
	}{
		{name: "boot attachment", mapping: "xvdn", padding: ' ', want: "/dev/xvdn"},
		{name: "hot attachment", mapping: "/dev/xvdn", padding: ' ', want: "/dev/xvdn"},
		{name: "NUL padding", mapping: "xvdl", want: "/dev/xvdl"},
		{name: "large namespace", mapping: "xvdba", largeRange: true, want: "/dev/xvdba"},
		{name: "large hot attachment", mapping: "/dev/xvdcz", largeRange: true, want: "/dev/xvdcz"},
		{name: "root disk", mapping: "xvda"},
		{name: "partition", mapping: "xvdn1"},
		{name: "different namespace", mapping: "xvdba"},
		{name: "default in large namespace", mapping: "xvdn", largeRange: true},
		{name: "path traversal", mapping: "/dev/../xvdn"},
		{name: "embedded device prefix", mapping: "other/dev/xvdn"},
		{name: "not an EBS mapping", mapping: "nvme1n1"},
		{name: "empty mapping", padding: ' '},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := make([]byte, 4096)
			copy(controller[3072:3104], bytes.Repeat([]byte{tt.padding}, 32))
			copy(controller[3072:3104], tt.mapping)
			deviceRE := ebsUtils.GetDeviceRange(tt.largeRange).DeviceRE
			if got := nvmeDevicePath(controller, deviceRE); got != tt.want {
				t.Fatalf("nvmeDevicePath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNVMeDevicePathTruncatedController(t *testing.T) {
	for _, size := range []int{0, 3072, 3103} {
		if got := nvmeDevicePath(make([]byte, size), ebsUtils.GetDeviceRange(false).DeviceRE); got != "" {
			t.Fatalf("controller size %d: got %q, want no mapping", size, got)
		}
	}
}
