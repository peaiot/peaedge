package e212

import (
	"fmt"
	"strconv"
	"strings"
)

func Crc(dataString string) string {
	data := []byte(dataString)
	crc := 0xFFFF
	length := len(data)
	for i := 0; i < length; i++ {
		// hj212取寄存器的高8位参与异或运算
		crc = (crc >> 8) ^ int(data[i])
		for j := 0; j < 8; j++ {
			flag := crc & 0x0001
			crc >>= 1
			if flag == 1 {
				crc ^= 0xA001
			}
		}
	}
	// 因为是基于右移位运算的结果，得到的本身就是低字节在前高字节在后的结果
	// 不足4位的需要在开头补0
	hex := strconv.FormatInt(int64(crc), 16)
	if len(hex) < 4 {
		hex = fmt.Sprintf("%04s", hex)
	}
	return strings.ToUpper(hex)
}

// unsigned int CRC16_Checkout ( unsigned char *puchMsg, unsigned int usDataLen )
// {
// 	unsigned int i, j, crc_reg, check;
// 	crc_reg = 0xFFFF;
// 	for (i = 0;i<usDataLen; i++) {
// 		crc_reg = (crc_reg>>8) ^ puchMsg[i];
// 		for (j = 0; j<8; j++)
// 		{
// 			check = crc_reg & 0x0001;
// 			crc_reg >>= 1;
// 			if (check==0x0001)
// 			{
// 				crc_reg ^= 0xA001;
// 			}
// 		}
// 	}
// 	return crc_reg;
// }
