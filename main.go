package main

import (
	"fmt"
	"net"
	"strings"
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
	fmt.Printf("Sent %s to: %s\n", message, remote.String())
	local.WriteTo([]byte(message), remote)
}
func broadcast_message(local net.PacketConn, data []byte) {
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
			handler(local, data[:len], remote)
		}
	}
}
func handler(local net.PacketConn, data []byte, remote net.Addr) {
	input_cmd := binary.BigEndian.Uint64(data[1:8])
	input_payload = data[8:]
	fmt.Printf("Recv %s fr: %s\n", input_payload.String, remote.String())
	switch input_cmd {
		case 0:
			output_cmd := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(1))
			output_payload := b"BONG"
			output := append(output_cmd, output_payload)
			broadcast_message(local, output)
		case 1:
			output_cmd := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(2))
			output_payload := b"PING"
			output := append(output_cmd, output_payload)
			send_message(local, output, remote)
		case 2:
			output_cmd := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(3))
			output_payload := b"PONG"
			output := append(output_cmd, output_payload)
			send_message(local, output, remote)
	}
}

func main() {
	port := 37419
	size := 1024
	socket,_ := net.ListenPacket("udp4", ":"+strconv.Itoa(port))
	broadcast_message(socket, "BING")
	listener(socket, size)
}