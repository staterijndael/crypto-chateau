package message

type Message interface {
	Marshal() []byte
	Unmarshal(map[string][]byte) error
}
