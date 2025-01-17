package main

import (
	"encoding/binary"
	// "fmt"
	"net"
)

// ConnectionInit = 0x01,
//     HeartBeat = 0x02,
//     MetricsData = 0x03,

type HttpRouteMetricsPayload struct {
	Route        string
	Method       byte
	StatusCode   uint16
	ResponseTime float64
}

func DeserializeMetricsDataFromBytes(buf []byte) (*HttpRouteMetricsPayload, error) {
	// buf must only contain the dynamic payload
	// fmt.Println("len of payload buf: ", len(buf))
	var statusCode uint16
	_, err := binary.Decode(buf[0:2], binary.BigEndian, &statusCode)
	if err != nil {
		return nil, err
	}
	// fmt.Println("statusCode: ", statusCode)
	var responseTime float64
	_, err = binary.Decode(buf[2:10], binary.BigEndian, &responseTime)
	if err != nil {
		return nil, err
	}
	// fmt.Println("responseTime: ", responseTime)
	var method byte
	_, err = binary.Decode(buf[10:11], binary.BigEndian, &method)
	if err != nil {
		return nil, err
	}
	// fmt.Println("method: ", method)
	route := string(buf[11:])
	// fmt.Println(statusCode, responseTime, method, route)
	return &HttpRouteMetricsPayload{
		Route:        route,
		Method:       method,
		StatusCode:   statusCode,
		ResponseTime: responseTime,
	}, nil
}

// func DeserializePackageFromBytes(buf []byte, payloadSize int) (*Package, error) {
func DeserializePackageFromBytes(buf []byte) (*Package, error) {
	// buf: contains complete package(along with the dynamic payload)
	// payloadSize: this is just the size of the dynamic payload

	var version byte
	_, err := binary.Decode(buf[0:1], binary.BigEndian, &version)
	if err != nil {
		return nil, err
	}

	// fmt.Println("version: ", version)
	var messageType byte
	_, err = binary.Decode(buf[1:2], binary.BigEndian, &messageType)
	if err != nil {
		return nil, err
	}
	// fmt.Println("messageType: ", messageType)
	var seqNum uint32
	_, err = binary.Decode(buf[2:6], binary.BigEndian, &seqNum)
	if err != nil {
		return nil, err
	}
	// fmt.Println("seqNum: ", seqNum)

	var timeStamp uint64
	_, err = binary.Decode(buf[6:14], binary.BigEndian, &timeStamp)
	if err != nil {
		return nil, err
	}
	// fmt.Println("timeStamp: ", timeStamp)
	payload := buf[14:]
	deserializedMetricsData, err := DeserializeMetricsDataFromBytes(payload)
	if err != nil {
		return nil, err
	}

	pkg := &Package{
		_deserializedMetricsData: deserializedMetricsData,
		Version:                  version,
		MessageType:              messageType,
		SeqNum:                   seqNum,
		Payload:                  payload,
	}

	return pkg, nil
}

type Package struct {
	Version                  byte // 1
	MessageType              byte //
	SeqNum                   uint32
	Payload                  []byte
	Timestamp                uint64
	_deserializedMetricsData *HttpRouteMetricsPayload
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
		return err
	}

	_, err = binary.Encode(buffer[6:14], binary.BigEndian, pkg.Timestamp)
	if err != nil {
		return err
	}
	if payloadLen > 0 {
		// write the payload and update totalBytesWritten
		return err
	}

	_, err = writer.Write(buffer)
	if err != nil {
		return err
	}
	return nil

}
