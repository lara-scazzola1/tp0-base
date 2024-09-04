import struct
from common.utils import Bet

# comandos que recibe
BET_COMMAND = 9
BATCH_COMMAND = 19
DISCONNECT_COMMAND = 29
WAIT_WINNERS_COMMAND = 39

# comandos que manda
RESPONSE_BET_COMMAND = 9
RESPONSE_BATCH_COMMAND_OK    = 19
RESPONSE_BATCH_COMMAND_ERROR = 20
RESPONSE_WINNERS_COMMAND = 39

class Protocol:
    def receive_command(self, skt):
        command = struct.unpack('B', skt.recvall(1))[0]
        return command

    def receive_bet(self, skt):
        data_size = struct.unpack('>I', skt.recvall(4))[0]
        
        data = skt.recvall(data_size)
        return Bet.deserialize(data)
    
    def receive_batch(self, skt):
        data_size = struct.unpack('>I', skt.recvall(4))[0]

        data = skt.recvall(data_size)
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

    def send_response_bet(self, skt):
        skt.sendall(struct.pack('B', RESPONSE_BET_COMMAND))

    def send_response_batch(self, ok, skt):
        if ok:
            skt.sendall(struct.pack('B', RESPONSE_BATCH_COMMAND_OK))
        else:
            skt.sendall(struct.pack('B', RESPONSE_BATCH_COMMAND_ERROR))

    def send_response_winners(self, skt):
        skt.sendall(struct.pack('B', RESPONSE_WINNERS_COMMAND))
