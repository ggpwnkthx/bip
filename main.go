package main

import (
	"fmt"
	"net"
)

func getLocalIPs() ([]net.IP) {
	output := make([]net.IP, 0)
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
			output = append(output, ip)
		}
	}
	return output
}
func inSlice(slice []net.IP, val net.IP) (bool) {
    for _, item := range slice {
        if item.Equal(val) {
            return true
        }
    }
    return false
}

func loggit(data []byte, remote net.Addr) {
	fmt.Printf("%08b ", data[0])
	fmt.Printf("%s: %s\n", remote.String(), string(data[1:]))
}

func send_message(local net.PacketConn, data []byte, remote net.Addr) {
	loggit(data, remote)
	local.WriteTo(data, remote)
}
func broadcast_message(local net.PacketConn, data []byte) {
	_,port,_ := net.SplitHostPort(local.LocalAddr().String())
	remote,_ := net.ResolveUDPAddr("udp4", "255.255.255.255:"+port)
	send_message(local, data, remote)
}
func listener(local net.PacketConn, size int) {
	filter := getLocalIPs()
	for {
		data := make([]byte, size)
		len,remote,_ := local.ReadFrom(data)
		remoteUDP,_ := net.ResolveUDPAddr(remote.Network(), remote.String())
		if ! inSlice(filter, remoteUDP.IP) {
			loggit(data, remote)
			handler(local, data[:len], remote)
		}
	}
}
func handler(local net.PacketConn, data []byte, remote net.Addr) {
	switch data[0] {
		case byte(0):
			broadcast_message(local, build_packet(1, "BONG"))
		case byte(1):
			send_message(local, build_packet(2, "PING"), remote)
		case byte(2):
			send_message(local, build_packet(3, "PONG"), remote)
	}
}

func build_packet (cmd int, payload string) []byte {
	output := []byte{byte(cmd)}
	output = append(output, []byte(payload)...)
	return output
}

func main() {
	port := "37419"
	size := 1024
	socket,_ := net.ListenPacket("udp4", ":"+port)
	broadcast_message(socket, build_packet(0, "BING"))
	listener(socket, size)
}