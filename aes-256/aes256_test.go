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

func Test_Encrypt(t *testing.T) {
	data := "Дарова"

	hasher := sha256.New()
	randSymbol := string([]rune{rune(rand.Intn(122-97+1) + 97)})
	hasher.Write([]byte(randSymbol))
	shaHash := hex.EncodeToString(hasher.Sum(nil))

	shaHashRunes := []byte(shaHash)[:len(shaHash)/2]

	encryptedData, err := Encrypt([]byte(data), shaHashRunes)
	assert.NoError(t, err)

	decryptedData, err := Decrypt(encryptedData, shaHashRunes)
	assert.NoError(t, err)

	assert.Equal(t, data, string(decryptedData))

	t.Log(string(encryptedData))
	t.Log(string(decryptedData))
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

	invStateAfter := InvSubBytes(stateAfter)

	assert.Equal(t, stateBefore, invStateAfter)
}

func Test_KeyExpansion(t *testing.T) {
	hasher := sha256.New()
	randSymbol := string([]rune{rune(rand.Intn(122-97+1) + 97)})
	hasher.Write([]byte(randSymbol))
	shaHash := hex.EncodeToString(hasher.Sum(nil))

	shaHashRunes := []byte(shaHash)[:len(shaHash)/2]

	kExp, err := keyExpansion(shaHashRunes)
	assert.NoError(t, err)
	t.Log(len(kExp[0]))

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

	invResultState := InvShiftRows(resultState)

	log.Println(invResultState)

	assert.Equal(t, stateBefore, invResultState)
}

func Test_MixColumns(t *testing.T) {
	stateBefore := [][]uint16{
		{0x1, 0x55, 0x3, 0x14, 0x5, 0x22, 0x7, 0x8},
		{0x7, 0x22, 0x3, 0x4, 0x14, 0x54, 0x7, 0x13},
		{0x9, 0x54, 0x3, 0x4, 0x54, 0x77, 0x54, 0x99},
		{0x13, 0x37, 0x3, 0x54, 0x51, 0x62, 0x7, 0x8},
	}

	cpStateBefore := make([][]uint16, len(stateBefore))
	copy(cpStateBefore, stateBefore)

	log.Println(stateBefore)
	resultState := mixColumns(cpStateBefore)

	log.Println(stateBefore)

	invResultState := InvMixColumns(resultState)

	log.Println(invResultState)

	assert.Equal(t, stateBefore, invResultState)
}
