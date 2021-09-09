import socket

addr = "test.socket"

with socket.socket(socket.AF_UNIX, socket.SOCK_STREAM) as s:
    s.connect(addr)
    s.sendall(b'Hello, world\n')
    data = s.recv(1024)

print('Received', repr(data))