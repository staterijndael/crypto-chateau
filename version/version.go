package version

const (
	ProtocolVersion byte = 0b00000001 // 4bits: min 0 max 7
	CodegenVersion       = "1.0.0"    // TODO: get from git on build
)

func NewProtocolByte() byte {
	return ProtocolVersion & (0b11110000) // first 4 bits reserved
}
