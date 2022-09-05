package utility

import (
	"bytes"
	"encoding/binary"
	"log"
)

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Fatalln("something wrong in convert number to hex")
	}
	return buff.Bytes()
}