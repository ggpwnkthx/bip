#!/usr/local/bin/python
import socket
import threading

def broadcast_message(message, port=37419, encoding='utf-8'):
    server = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
    server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEPORT, 1)
    server.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
    server.settimeout(0.2)
    server.sendto(bytes(message, encoding), ("<broadcast>", port))
    server.close()

def unicast_message(addr, message, port=37419, encoding='utf-8'):
    server = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
    server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEPORT, 1)
    server.settimeout(0.2)
    server.sendto(bytes(message, encoding), (addr, port))
    server.close()

def listener(port=37419, size=1024):
    client = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
    client.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEPORT, 1)
    client.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
    client.bind(("", port))
    while True:
        data, addr = client.recvfrom(size)
        if addr:
            x = threading.Thread(target=handler, args=(data,addr,), daemon=True)
            x.start()
    client.close()

def handler(data, addr, encoding='utf-8'):
    command = str(data, encoding).split()
    if command[0].upper() is "BING":
        broadcast_message("BONG")
    if command[0].upper() is "PING":
        unicast_message(addr, "PONG")

if __name__ == "__main__":
    broadcast_message("BING")
    listener()
