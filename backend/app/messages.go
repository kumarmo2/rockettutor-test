package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

// ConnectionInit = 0x01,
//     HeartBeat = 0x02,
//     MetricsData = 0x03,

type Package struct {
	Version     byte // 1
	MessageType byte //
	SeqNum      uint32
	Payload     []byte
	Timestamp   uint64
}

func (pkg *Package) WriteToUdp(writer *net.UDPConn) error {
	// 14 + payload
	payloadLen := len(pkg.Payload)
	n := 14 + payloadLen
	buffer := make([]byte, n)

	buffer[0] = pkg.Version
	buffer[1] = pkg.MessageType

	_, err := binary.Encode(buffer[2:6], binary.BigEndian, pkg.SeqNum)
	if err != nil {
		panic(err)
	}

	_, err = binary.Encode(buffer[6:14], binary.BigEndian, pkg.Timestamp)
	if err != nil {
		fmt.Println("error while writing Timestamp")
		panic(err)
	}
	// totalBytesWritten := 6
	// // will start from six.
	//
	if payloadLen > 0 {
		// write the payload and update totalBytesWritten
		panic("not implmeneted payload encoding")
	}

	_, err = writer.Write(buffer)
	if err != nil {
		panic(err)

	}
	// panic("done writing package")
	// if err != nil {
	// 	panic(err)
	// }

	// buffer[2:7]

	return nil

}
