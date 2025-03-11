package ssh

type Flags struct {
	Username string
	KeyFile  string

	ProxyJump ProxyJumpFlags

	ProxyCommand string

	Command string
}

type ProxyJumpFlags struct {
	Address  string
	Username string
	KeyFile  string
}
