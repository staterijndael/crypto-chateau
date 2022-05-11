package crypto_chateau

type Handler struct {
	callFunc func(...interface{}) (Message, error)
}
