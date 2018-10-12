package sysinfo

import (
	"fmt"
	"testing"
)

func TestWtmpMetrics(t *testing.T) {
	wtmpFile = "./testdata/wtmp"
	res := WtmpMetrics()
	fmt.Println(res)
	if len(res) == 0 {
		t.Fatal("parse wtmp log failed", len(res))
	}
}
