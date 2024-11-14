package share

var (
	ConsulAddr func() string
)

func GetConsulAddr() string {
	if ConsulAddr == nil {
		return ""
	}
	
	return ConsulAddr()
}
