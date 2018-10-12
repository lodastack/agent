package sysinfo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

var wtmpFile = "/var/log/wtmp"
var btmpFile = "/var/log/btmp"

const fileLimitSize = 500 * 1024 * 1024

const (
	empty        = 0x0
	runLevel     = 0x1
	bootTime     = 0x2
	newTime      = 0x3
	oldTime      = 0x4
	initProcess  = 0x5
	loginProcess = 0x6
	userProcess  = 0x7
	deadProcess  = 0x8
	accounting   = 0x9
)

const (
	lineSize = 32
	nameSize = 32
	hostSize = 256
)

// utmp structures
// see man utmp
type exitStatus struct {
	Termination int16
	Exit        int16
}

type timeVal struct {
	Sec  int32
	Usec int32
}

type utmp struct {
	Type int16
	// alignment
	_       [2]byte
	Pid     int32
	Device  [lineSize]byte
	ID      [4]byte
	User    [nameSize]byte
	Host    [hostSize]byte
	Exit    exitStatus
	Session int32
	Time    timeVal
	Addr    [4]int32
	// Reserved member
	Reserved [20]byte
}

// Read utmps
func read(file io.Reader) ([]*utmp, error) {
	var us []*utmp

	for {
		u, readErr := readLine(file)
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return nil, readErr
		}
		us = append(us, u)
	}

	return us, nil
}

// read utmp
func readLine(file io.Reader) (*utmp, error) {
	u := new(utmp)

	err := binary.Read(file, binary.LittleEndian, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

type goExitStatus struct {
	Termination int
	Exit        int
}

type goUtmp struct {
	Type    int
	Pid     int
	Device  string
	ID      string
	User    string
	Host    string
	Exit    goExitStatus
	Session int
	Time    time.Time
	Addr    string
}

// Convert Utmp to GoUtmp
func newGoUtmp(u *utmp) *goUtmp {
	return &goUtmp{
		Type:   int(u.Type),
		Pid:    int(u.Pid),
		Device: string(u.Device[:getByteLen(u.Device[:])]),
		ID:     string(u.ID[:getByteLen(u.ID[:])]),
		User:   string(u.User[:getByteLen(u.User[:])]),
		Host:   string(u.Host[:getByteLen(u.Host[:])]),
		Exit: goExitStatus{
			Termination: int(u.Exit.Termination),
			Exit:        int(u.Exit.Exit),
		},
		Session: int(u.Session),
		Time:    time.Unix(int64(u.Time.Sec), 0),
		Addr:    addrToString(u.Addr),
	}
}

// Integer ip address to string
func addrToString(addr [4]int32) string {
	if addr[1] == 0 && addr[2] == 0 && addr[3] == 0 {
		return fmt.Sprintf(
			"%d.%d.%d.%d",
			addr[0]&0xFF,
			(addr[0]>>8)&0xFF,
			(addr[0]>>16)&0xFF,
			(addr[0]>>24)&0xFF,
		)
	} else {
		return fmt.Sprintf(
			"%x:%x:%x:%x:%x:%x:%x:%x",
			addr[0]&0xffff,
			(addr[0]>>16)&0xffff,
			addr[1]&0xffff,
			(addr[1]>>16)&0xffff,
			addr[2]&0xffff,
			(addr[2]>>16)&0xffff,
			addr[3]&0xffff,
			(addr[3]>>16)&0xffff,
		)
	}
}

// get byte \0 index
func getByteLen(byteArray []byte) int {
	n := bytes.IndexByte(byteArray[:], 0)
	if n == -1 {
		return 0
	}

	return n
}
