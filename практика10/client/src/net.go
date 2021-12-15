package src

import (
	"fmt"
	"net"
	"strings"
)

var pc net.PacketConn

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

func allIP(address string) string {
	parts := strings.Split(address, ".")
	parts[3] = "255"
	return parts[0] + "." + parts[1] + "." + parts[2] + "." + parts[3]
}

func listen(ch chan string) string {
	myIP := localIp()
	fmt.Println(myIP.String() + " listen")
	pc, _ = net.ListenPacket("udp", allIP(myIP.String())+":12345")
	//if err != nil {
	//	mylog.Write([]byte(err.Error()))
	//}
	for {
		buf := make([]byte, 10000)
		_, _, err := pc.ReadFrom(buf)
		if err != nil {
			mylog.Write([]byte(err.Error()))
			continue
		}
		msg := decode(string(buf))
		ch <- msg
	}
}

func takeDB() string {
	myIP := localIp()
	fmt.Println(myIP.String() + " takeDB")
	querys := []string{"select * from clients", "select * from cassettes", "select * from films", "select * from librarians", "select * from givings"}
	var r string
	for _, q := range querys {
		send(allIP(myIP.String()), encode(q), mylog)
		buf := make([]byte, 10000)
		_, _, err := pc.(*net.UDPConn).ReadFromUDP(buf)
		if err != nil {
			mylog.Write([]byte(err.Error()))
			return ""
		}
		r += "|||" + decode(string(buf))
	}
	return r
}
