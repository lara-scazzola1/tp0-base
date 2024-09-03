import struct
from common.socket import Socket
from common.utils import Bet

# comandos que recibe
BET_COMMAND = 0

# comandos que manda
RESPONSE_BET_COMMAND = 0

class Protocol:
    def __init__(self, sock: Socket):
        self.socket = sock

    def receive_command(self):
        command = struct.unpack('B', self.socket.recvall(1))[0]
        data_size = struct.unpack('>I', self.socket.recvall(4))[0]
        return command, data_size

    def receive_bet(self, data_size: int):
        data = self.socket.recvall(data_size)
        return Bet.deserialize(data)

    def send_response_bet(self):
        self.socket.sendall(struct.pack('B', RESPONSE_BET_COMMAND))

    def close(self):
        self.socket.close()
