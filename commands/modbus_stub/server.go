package main

import (
	"log"
	"time"

	"github.com/goburrow/serial"
	"github.com/tbrandon/mbserver"
)

func main() {
	serv := mbserver.NewServer()
	err := serv.ListenTCP("127.0.0.1:1502")
	if err != nil {
		log.Printf("%v\n", err)
	}

	err = serv.ListenRTU(&serial.Config{
		Address:  "/dev/ttys008",
		BaudRate: 9600,
		DataBits: 8,
		StopBits: 1,
		Parity:   "N",
		Timeout:  10 * time.Second,
		// RS485: serial.RS485Config{
		// 	Enabled:            true,
		// 	DelayRtsBeforeSend: 2 * time.Millisecond,
		// 	DelayRtsAfterSend:  3 * time.Millisecond,
		// 	RtsHighDuringSend:  false,
		// 	RtsHighAfterSend:   false,
		// 	RxDuringTx:         false,
		// },
	})

	if err != nil {
		log.Printf("%v\n", err)
	}

	defer serv.Close()

	// Wait forever
	for {
		time.Sleep(1 * time.Second)
	}
}
