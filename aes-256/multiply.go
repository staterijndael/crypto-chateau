package aes_256

func mulBy02(num uint16) uint16 {
	var res uint16

	if num < 0x80 {
		res = num << 1
	} else {
		res = (num << 1) ^ 0x1b
	}

	return res % 0x100
}

func mulBy03(num uint16) uint16 {
	return mulBy02(num) ^ num
}

func mulBy09(num uint16) uint16 {
	return mulBy02(mulBy02(mulBy02(num))) ^ num
}

func mulBy0b(num uint16) uint16 {
	return mulBy02(mulBy02(mulBy02(num))) ^ mulBy02(num) ^ num
}

func mulBy0e(num uint16) uint16 {
	return mulBy02(mulBy02(mulBy02(num))) ^ mulBy02(mulBy02(num)) ^ num
}
