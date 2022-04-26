package aes_256

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Test_SubBytes(t *testing.T) {
	stateBefore := [][]uint16{
		{0x10, 0x32, 0x16, 0x22},
		{0xff, 0x1b, 0x2b, 0x32},
		{0x10, 0x32, 0x16, 0x22},
		{0xff, 0x1c, 0x2d, 0x12},
	}

	tmpArr := make([][]uint16, len(stateBefore))
	for i := range stateBefore {
		for j := range stateBefore[i] {
			tmpArr[i] = append(tmpArr[i], stateBefore[i][j])
		}
	}

	stateAfter := subBytes(tmpArr)

	assert.NotEqualValues(t, stateBefore, stateAfter)

	for i, row := range stateAfter {
		for j, elem := range row {
			sboxElem := InvSbox[elem]
			stateAfter[i][j] = sboxElem
		}
	}

	assert.Equal(t, stateBefore, stateAfter)
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

func Test_ShiftRows(t *testing.T) {
	stateBefore := [][]uint16{
		{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8},
		{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8},
		{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8},
		{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8},
	}

	resultState := shiftRows(stateBefore)

	log.Println(resultState)
}
