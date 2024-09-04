import struct
from common.socket import Socket
from common.utils import Bet

# comandos que recibe
BET_COMMAND = 9
BATCH_COMMAND = 19
DISCONNECT_COMMAND = 29

# comandos que manda
RESPONSE_BET_COMMAND = 9
RESPONSE_BATCH_COMMAND_OK    = 19
RESPONSE_BATCH_COMMAND_ERROR = 20

class Protocol:
    def __init__(self, sock: Socket):
        self.socket = sock

    def receive_command(self):
        command = struct.unpack('B', self.socket.recvall(1))[0]
        return command
    
    def receive_data_size(self):
        data_size = struct.unpack('>I', self.socket.recvall(4))[0]
        return data_size

    def receive_bet(self, data_size: int):
        data = self.socket.recvall(data_size)
        return Bet.deserialize(data)
    
    def receive_batch(self, data_size: int):
        data = self.socket.recvall(data_size)
        bets = []
        bytes_read = 0
        amount_bets_send = 0
        while bytes_read < data_size:
            amount_bets_send += 1
            bet_size = struct.unpack('>I', data[bytes_read:bytes_read+4])[0]
            bytes_read += 4
            bet_data = data[bytes_read:bytes_read+bet_size]
            bytes_read += bet_size
            bet = Bet.deserialize(bet_data)
            if bet is not None:
                bets.append(bet)
        return amount_bets_send, bets

    def send_response_bet(self):
        self.socket.sendall(struct.pack('B', RESPONSE_BET_COMMAND))

    def send_response_batch(self, ok: bool):
        if ok:
            self.socket.sendall(struct.pack('B', RESPONSE_BATCH_COMMAND_OK))
        else:
            self.socket.sendall(struct.pack('B', RESPONSE_BATCH_COMMAND_ERROR))

    def close(self):
        self.socket.close()
