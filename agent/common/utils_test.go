package common

import (
	"os/exec"
	"testing"
	"time"
)

func Test_Hostname(t *testing.T) {
	_, err := Hostname()
	if err != nil {
		t.Fatalf("get hostname fatal: %s", err.Error())
	}
}

func Test_normalizedHostname(t *testing.T) {
	hostname := "www_cdn256_BJ.test.com"
	finnal := normalizedHostname(hostname)
	if finnal != hostname {
		t.Fatalf("normalizedHostname fatal: %s - %s", finnal, hostname)
	}
}

func Test_GetIpList(t *testing.T) {
	ips := GetIpList()
	if len(ips) < 1 {
		t.Fatalf("get IP fatal: get ip num = %d", len(ips))
	}
	for _, ip := range ips {
		if ip == "127.0.0.1" || ip == "0.0.0.0" {
			t.Fatalf("get IP fatal: get a loop ip %s", ips)
		}
	}
}

func Test_CmdRunWithTimeout(t *testing.T) {
	cmd := exec.Command("pwd")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("exec cmd fatal: %s", err.Error())
	}
	err, isTimeout := CmdRunWithTimeout(cmd, time.Duration(100)*time.Millisecond)
	if err != nil {
		t.Fatalf("CmdRunWithTimeout cmd fatal: %s", err.Error())
	}

	if isTimeout {
		t.Fatalf("exec cmd timeout")
	}
}

func Test_StrTagsToMap(t *testing.T) {
	tags := "k1=v1,k2=v2,k3=v3"
	m := StrTagsToMap(tags)
	if m["k1"] != "v1" || m["k2"] != "v2" || m["k3"] != "v3" {
		t.Fatalf("string tags convert to map fatal")
	}
}

func Test_GitPath(t *testing.T) {
	MustConfig()
	correct := "git@git.test.com:user/project.git"
	finnal := GitPath("user/project")
	if finnal != correct {
		t.Fatalf("git path fatal: %s - %s", finnal, correct)
	}
}
