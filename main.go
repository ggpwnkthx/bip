package main

import (
	"fmt"
	"net"
	"strconv"
	"encoding/binary"
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

func send_message(local net.PacketConn, data []byte, remote net.Addr) {
	fmt.Printf("Sent %s to: %s\n", string(data[8:]), remote.String())
	local.WriteTo(data, remote)
}
func broadcast_message(local net.PacketConn, data []byte) {
	localUDP,_ := net.ResolveUDPAddr(local.LocalAddr().Network(), local.LocalAddr().String())
	remote,_ := net.ResolveUDPAddr("udp4", "255.255.255.255:"+strconv.Itoa(localUDP.Port))
	send_message(local, data, remote)
}
func listener(local net.PacketConn, size int) {
	filter := getLocalIPs()
	for {
		data := make([]byte, size)
		len,remote,_ := local.ReadFrom(data)
		remoteUDP,_ := net.ResolveUDPAddr(remote.Network(), remote.String())
		if ! inSlice(filter, remoteUDP.IP.String()) {
			handler(local, data[:len], remote)
		}
	}
}
func handler(local net.PacketConn, data []byte, remote net.Addr) {
	cmd := binary.BigEndian.Uint64(data[1:8])
	fmt.Printf("Recv %s fr: %s\n", string(data[8:]), remote.String())
	switch cmd {
		case 0:
			broadcast_message(local, build_packet(uint64(1), "BONG"))
		case 1:
			send_message(local, build_packet(uint64(2), "PING"), remote)
		case 2:
			send_message(local, build_packet(uint64(3), "PONG"), remote)
	}
}

func build_packet (cmd uint64, payload string) []byte {
	output := make([]byte, 8)
	binary.LittleEndian.PutUint64(output, uint64(cmd))
	output = append(output, []byte(payload)...)
	return output
}

func main() {
	port := 37419
	size := 1024
	socket,_ := net.ListenPacket("udp4", ":"+strconv.Itoa(port))
	broadcast_message(socket, build_packet(0, "BING"))
	listener(socket, size)
}