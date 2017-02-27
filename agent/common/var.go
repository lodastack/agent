package common

import (
	"sync"
)

const (
	EXCEPTION_NS = "collect.exception.monitor.loda"
	HOST_SUFFIX  = ""
)

const (
	TYPE_PROC     = "PROC"
	TYPE_PORT     = "PORT"
	TYPE_PLUGIN   = "PLUGIN"
	TYPE_LOG      = "LOG"
	TYPE_GOPLUGIN = "GOPLUGIN"
	TYPE_RUN      = "RUN"

	TYPE_CPU  = "CPU"
	TYPE_DISK = "DISK"
	TYPE_MEM  = "MEM"
	TYPE_NET  = "NET"
	TYPE_FS   = "FS"
	TYPE_TIME = "TIME"
	TYPE_SYS  = "SYS"

	TYPE_COREDUMP = "COREDUMP"
)

var (
	DEFAULT_INTERVAL = map[string]int{
		TYPE_CPU:      10,
		TYPE_DISK:     10,
		TYPE_MEM:      10,
		TYPE_NET:      10,
		TYPE_PROC:     60,
		TYPE_PORT:     60,
		TYPE_FS:       300,
		TYPE_TIME:     300,
		TYPE_COREDUMP: 60,
	}

	SYS_TYPES = []string{TYPE_CPU, TYPE_DISK, TYPE_MEM, TYPE_NET, TYPE_COREDUMP, TYPE_FS, TYPE_TIME}
)

var (
	idLock = new(sync.RWMutex)
	UUID   []string

	nsLock     = new(sync.RWMutex)
	Namespaces = []string{}

	portsLock = new(sync.RWMutex)
	ports     = []PortMonitor{}

	procsLock = new(sync.RWMutex)
	procs     = []ProcMonitor{}

	pluginLock = new(sync.RWMutex)
	pluginInfo = map[string]bool{}
)

func SetUUID(ids []string) {
	idLock.Lock()
	UUID = ids
	idLock.Unlock()
}

func GetUUID() []string {
	idLock.Lock()
	defer idLock.Unlock()
	return UUID
}

func SetNamespaces(ns []string) {
	nsLock.Lock()
	defer nsLock.Unlock()
	Namespaces = ns
}

func GetNamespaces() []string {
	nsLock.Lock()
	defer nsLock.Unlock()
	return Namespaces
}

func SetPorts(newPorts []PortMonitor) {
	portsLock.Lock()
	defer portsLock.Unlock()
	ports = newPorts
}

func ReportPorts() []PortMonitor {
	portsLock.Lock()
	defer portsLock.Unlock()
	return ports
}

func SetProcs(newProcs []ProcMonitor) {
	procsLock.Lock()
	defer procsLock.Unlock()
	procs = newProcs
}

func ReportProcs() []ProcMonitor {
	procsLock.Lock()
	defer procsLock.Unlock()
	return procs
}

func SetPluginInfo(p map[string]bool) {
	pluginLock.Lock()
	defer pluginLock.Unlock()
	pluginInfo = p
}

func GetPluginInfo() map[string]bool {
	pluginLock.Lock()
	defer pluginLock.Unlock()
	return pluginInfo
}
