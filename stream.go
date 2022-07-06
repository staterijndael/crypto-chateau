package crypto_chateau

type Stream struct {
	peer *Peer
}

func (s *Stream) Read() (int, error) {
	buf := make([]byte, 1024)
	n, err := s.peer.Read(buf)

	return n, err
}

func (s *Stream) Write(data []byte) (int, error) {
	return s.peer.Write(data)
}
