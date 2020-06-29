#!/usr/local/bin/python
import socket

def send_message(sock, message, address=""):
    if address == "":
        address = "<broadcast>"
    server.sendto(bytes(message, 'utf-8'), (address, port))
    server.close()

def listener(sock, size=1024):
    while True:
        data, addr = client.recvfrom(size)
        if addr:
            handler(sock, data, addr)

def handler(sock, data, addr):
    command = str(data, 'utf-8').split()
    cmd = command[0].upper()
    print("Recv "+cmd+" fr: "+addr[0]+":"+addr[1])
    if cmd == "BING":
        send_message(sock, "BONG")
    if cmd == "BONG":
        send_message(sock, "PING", addr)
    if cmd == "PING":
        send_message(sock, "PONG", addr)

if __name__ == "__main__":
	port = 37419
	size = 1024
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEPORT, 1)
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
    sock.settimeout(0.2)
    sock.bind(("", port))
    send_message(sock, "BING")
    listener(sock, size)
    sock.close()