package main

import (
	"fmt"
	"net"
	"strings"
	"strconv"
)

func getLocalIPs() ([]string) {
	output := make([]string, 0)
	ifaces,_ := net.Interfaces()
	for _, i := range ifaces {
		addrs,_ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
					ip = v.IP
			case *net.IPAddr:
					ip = v.IP
			}
			output = append(output, ip.String())
		}
	}
	return output
}
func inSlice(slice []string, val string) (bool) {
    for _, item := range slice {
        if item == val {
            return true
        }
    }
    return false
}

func send_message(local net.PacketConn, message string, remote net.Addr) {
	fmt.Printf("Sent %s to: %s\n", message, remote.String())
	local.WriteTo([]byte(message), remote)
}
func broadcast_message(local net.PacketConn, message string) {
	localUDP,_ := net.ResolveUDPAddr(local.LocalAddr().Network(), local.LocalAddr().String())
	remote,_ := net.ResolveUDPAddr("udp4", "255.255.255.255:"+strconv.Itoa(localUDP.Port))
	send_message(local, message, remote)
}
func listener(local net.PacketConn, size int) {
	filter := getLocalIPs()
	for {
		data := make([]byte, size)
		len,remote,_ := local.ReadFrom(data)
		remoteUDP,_ := net.ResolveUDPAddr(remote.Network(), remote.String())
		if ! inSlice(filter, remoteUDP.IP.String()) {
			handler(local, string(data[:len]), remote)
		}
	}
}
func handler(local net.PacketConn, data string, remote net.Addr) {
	command := strings.Fields(data)
	cmd := strings.ToUpper(command[0])
	fmt.Printf("Recv %s fr: %s\n", cmd, remote.String())
	switch cmd {
		case "BING":
			broadcast_message(local, "BONG")
		case "BONG":
			send_message(local, "PING", remote)
		case "PING":
			send_message(local, "PONG", remote)
	}
}

func main() {
	port := 37419
	size := 1024
	socket,_ := net.ListenPacket("udp4", ":"+strconv.Itoa(port))
	broadcast_message(socket, "BING")
	listener(socket, size)
}