package src

import (
	"fmt"
	"net"
)

func send(address, msg string, log *myLogger) {
	conn, err := net.Dial("udp", address+":12345")
	if err != nil {
		log.Write([]byte("Не удалось отправить сообщение: " + err.Error()))
		return
	}
	fmt.Fprintf(conn, msg)
}

func localIp() net.IP {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func listen() {
	myIP := localIp()
	fmt.Println(myIP.String())
	pc, err := net.ListenPacket("udp", myIP.String()+":12345")
	if err != nil {
		mylog.Write([]byte(err.Error()))
	}
	defer pc.Close()
	for {
		buf := make([]byte, 10000)
		_, address, err := pc.ReadFrom(buf)
		if err != nil {
			mylog.Write([]byte(err.Error()))
			continue
		}
		msg := decode(string(buf))
	}
}
