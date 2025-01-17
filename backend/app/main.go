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
	// panic("done writing inti")

	//     pub enum MessageType {
	//     ConnectionInit = 0x01,
	//     HeartBeat = 0x02,
	//     MetricsData = 0x03,
	// }
	// // Expect Data Package for any package
	// pub struct Package {
	//     pub version: u8,
	//     pub message_type: MessageType,
	//     pub seq_num: u32,
	//     pub payload: Vec<u8>,
	//     pub time_stamp: u64,
	// }

	doneChan := make(chan bool, 1)

	go func() {
		buf := make([]byte, 1024)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("error while reading udp msg")
				panic(err)
			}
			panic(fmt.Sprintf("read %v bytes from udp", n))

		}

	}()

	return doneChan, nil

}

func main() {

	fmt.Println("hello world")
	doneChan, err := StartUpdApp()
	if err != nil {
		panic(err)
	}
	<-doneChan
}
