package aes_256

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Test_KeyExpansion(t *testing.T) {
	hasher := sha256.New()
	randSymbol := string([]rune{rune(rand.Intn(122-97+1) + 97)})
	hasher.Write([]byte(randSymbol))
	shaHash := hex.EncodeToString(hasher.Sum(nil))

	shaHashRunes := []rune(shaHash)[:len(shaHash)/2]

	kExp, err := keyExpansion(shaHashRunes)
	assert.NoError(t, err)
	t.Log(shaHash)
	t.Log(kExp)
}
