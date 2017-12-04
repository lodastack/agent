package common

import (
	"os"
	"testing"
)

func isDocker() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

func isLXC() bool {
	return os.Getenv("container") == "lxc"
}

func skipInContainer(t *testing.T) {
	if isDocker() {
		t.Skip("skip this test in Docker container")
	}
	if isLXC() {
		t.Skip("skip this test in LXC container")
	}
}

func Test_SN(t *testing.T) {
	// skipInContainer(t)
	// sn := SN()
	// if sn == "" {
	// 	t.Fatalf("get sn fatal")
	// }
	// t.Logf("SN:%s", sn)
}
