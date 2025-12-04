package src

const freeMemoryConst = "MemAvailable:"

const (
	system             = "/usr/lib/libSystem.B.dylib"
	host_vm_info       = 2
	host_cpu_load_info = 3
	host_vm_info_count = 0xf
	kern_success       = 0
	hostStatisticsSym  = "host_statistics"
	machHostSelfSym    = "mach_host_self"
	mbtouintsize       = 131072
)

const (
	session_token_length = 20
	refresh_token_length = 32
	parse_start          = 10
	parse_stop           = 64
	binary_sys_number    = 1024
)

const (
	Memory_cutoff           = 20
	DefaultSessionTokenTime = 5
	DefaultRefreshTokenTime = 1
	Allowed_Sessions  = 3
)
