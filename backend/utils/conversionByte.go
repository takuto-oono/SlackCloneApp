package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"encoding/gob"

	"backend/controllerUtils"
)

func SendMessageInputToByte(b []byte) *controllerUtils.SendMessageInput {
	reader := bytes.NewReader(b)
	res := new(controllerUtils.SendMessageInput)
	binary.Read(reader, binary.LittleEndian, res)
	return res
}

func ByteToStruct(s interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(s); err != nil {
		fmt.Println(err)
	}
	return buf.Bytes()

}
