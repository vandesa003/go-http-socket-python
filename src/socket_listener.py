import socket

addr = "test.socket"

s = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
s.connect(addr)

while True:
    try:
        data = s.recv(1024)
        print(data)
    except KeyboardInterrupt as e:
        break
s.close()