package rpc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

func RpcSend(conn net.Conn, data *[]byte) {
	//data := []byte("[这里才是一个完整的数据包]")
	l := len(*data)
	//fmt.Println(l)
	magicNum := make([]byte, 4)
	binary.BigEndian.PutUint32(magicNum, 0x123456)
	lenNum := make([]byte, 2)
	binary.BigEndian.PutUint16(lenNum, uint16(l))
	packetBuf := bytes.NewBuffer(magicNum)
	packetBuf.Write(lenNum)
	packetBuf.Write(*data)
	_, err := conn.Write(packetBuf.Bytes())
	if err != nil {
		fmt.Printf("write failed , err : %v\n", err)
	}
}
