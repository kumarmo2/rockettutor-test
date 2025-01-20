package main

import (
	"fmt"
	"net"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const dbName string = "rockettutor"

func StartUpdApp(liveEventManager *LiveEventManager[Metrics], metricsDao IMetricsDao) (chan bool, error) {

	udpAddr, err := net.ResolveUDPAddr("udp", "35.159.152.17:5000")

	if err != nil {
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
				// fmt.Println("successfully sent the heartBeat")
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
			if pkg._deserializedMetricsData != nil {
				metrics := NewMetricsFromRawPayload(pkg._deserializedMetricsData)
				err = metricsDao.Create(metrics)
				if err != nil {
					fmt.Println("error while persisting metrics: ", err)
				} else {
					liveEventManager.broadcast(*metrics)
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return doneChan, nil

}

func MustCreateDb(db *sqlx.DB) {

	rows, err := db.Queryx("select 1 from pg_database where datname = $1", dbName)
	if err != nil {
		panic(err)
	} else {
		foundDb := rows.Next()
		fmt.Println("foundDb: ", foundDb)
		if !foundDb {
			query := fmt.Sprintf("create database %v", dbName)
			_, err = db.Queryx(query)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("created db")
			}
		} else {
			fmt.Println("db already present")
		}
	}
}

func SetupDb(db *sqlx.DB) {
	query := `
    create schema if not exists metrics;
    create table if not exists metrics.metrics(
        id uuid not null,
        route text,
        method text,
        responsetime real,
        statuscode smallint,
        createdon timestamp not null,
        primary key(id)
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println(">>> db setup done")
}

func main() {
	db, err := sqlx.Connect("postgres", "user=postgres password=admin sslmode=disable")
	if err != nil {
		panic(err)
	}
	MustCreateDb(db)
	db, err = sqlx.Connect("postgres", fmt.Sprintf("user=postgres password=admin sslmode=disable database=%v", dbName))
	if err != nil {
		panic(err)
	}
	SetupDb(db)

	fmt.Println("hello world")
	liveEventManager := newLiveEventManager[Metrics]()

	go liveEventManager.startBackgroundProcessing()

	metricsDao := newMetricsDao(db)

	udpClientDoneChan, err := StartUpdApp(liveEventManager, metricsDao)
	if err != nil {
		panic(err)
	}

	httpServerDoneChan := make(chan bool, 1)
	server := newServer("0.0.0.0:9001", liveEventManager, metricsDao)

	server.run(httpServerDoneChan)

	select {
	case <-udpClientDoneChan:
		fmt.Println("received from  udpClientDoneChan ")
	case <-httpServerDoneChan:
	}
}
