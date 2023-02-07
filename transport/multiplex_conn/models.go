package multiplex_conn

type ToWriteMsg struct {
	RequestID uint32
	Data      []byte
}
