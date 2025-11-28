package src

const freeMemoryConst = "MemAvailable:"

const (
	System             = "/usr/lib/libSystem.B.dylib"
	HOST_VM_INFO       = 2
	HOST_CPU_LOAD_INFO = 3
	HOST_VM_INFO_COUNT = 0xf
	KERN_SUCCESS       = 0
	HostStatisticsSym  = "host_statistics"
	MachHostSelfSym    = "mach_host_self"
	MBTOUINTSIZE       = 131072
)

const (
	ONE_PROCESS = iota + 1
	TWO_PROCESS
)

const (
	SESSION_TOKEN_LENGTH = 20
	REFRESH_TOKEN_LENGTH = 32
	ZERO = 0
	ONE = 1
	TWO = 2
	DefaultSessionTokenTime = 5
	DefaultRefreshTokenTime = 1
	PARSE_START = 10
	PARSE_STOP = 64
	BINARY_SYS_NUMBER = 1024
)

const (
	MEMORY_CUTOFF =  20
)