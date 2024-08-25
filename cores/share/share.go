package share

var (
	UseConsul  func() bool
	ConsulAddr func() string
)

func GetConsulAddr() string {
	if ConsulAddr == nil {
		return ""
	}

	return ConsulAddr()
}

func GetConsulEnable() bool {
	if UseConsul == nil {
		return false
	}

	return UseConsul()
}
