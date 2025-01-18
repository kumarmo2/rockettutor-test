package main

import (
	"fmt"
	"net"
	"time"
)

func StartUpdApp() (chan bool, error) {

	udpAddr, err := net.ResolveUDPAddr("udp", "35.159.152.17:5000")

	if err != nil {
		// fmt.Println(err)
		return nil, err
	}

	// Dial to the address with UDP
	conn, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	connInitMessage := &Package{
		Version:     1,
		MessageType: 1, // connInit
		SeqNum:      0,
		Payload:     []byte{},
		Timestamp:   uint64(time.Now().UnixNano()),
	}

	err = connInitMessage.WriteToUdp(conn)
	if err != nil {
		panic(err)
	}

	doneChan := make(chan bool, 1)

	go func() {
		heartBeatPackage := &Package{
			Version:     1,
			MessageType: 2, // heartBeat
			SeqNum:      0,
			Payload:     []byte{},
			Timestamp:   uint64(time.Now().UnixNano()),
		}
		for {
			err = heartBeatPackage.WriteToUdp(conn)
			if err != nil {
				fmt.Println("error writing heartBeatPackage", err)
			} else {
				fmt.Println("successfully sent the heartBeat")
			}
			time.Sleep(5 * time.Second)
		}

	}()
	//TODO: need to do proper error handling.

	go func() {
		buf := make([]byte, 1024)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("error while reading udp msg")
				panic(err)
			}

			pkg, err := DeserializePackageFromBytes(buf[0:n])
			if err != nil {
				panic(err)
			}
			fmt.Printf("payload: %+v\n", *pkg._deserializedMetricsData)
			time.Sleep(1 * time.Second)
		}
	}()

	return doneChan, nil

}

func main() {

	fmt.Println("hello world")
	udpClientDoneChan, err := StartUpdApp()
	if err != nil {
		panic(err)
	}

	httpServerDoneChan := make(chan bool, 1)
	server := newServer("0.0.0.0:9001")
	server.run(httpServerDoneChan)

	select {
	case <-udpClientDoneChan:
		fmt.Println("received from  udpClientDoneChan ")
	case <-httpServerDoneChan:
	}
}
