import struct
from common.utils import Bet

# comandos que recibe
BATCH_COMMAND = 2
WAIT_WINNERS_COMMAND = 4
CLIENT_ID = 5

# comandos que manda
RESPONSE_BATCH_COMMAND_OK    = 2
RESPONSE_BATCH_COMMAND_ERROR = 3
RESPONSE_WINNERS_COMMAND = 4
RESPONSE_CLIENT_ID = 5

class Protocol:
    def receive_command(self, skt):
        command = struct.unpack('B', skt.recvall(1))[0]
        return command

    def receive_batch(self, skt):
        """
        Receive a batch of bets from a client socket.
        """

        # Receive the size of the data
        data_size = struct.unpack('>I', skt.recvall(4))[0]

        # Receive the data of the batch
        data = skt.recvall(data_size)

        # Deserialize the bets
        bets = []
        bytes_read = 0
        amount_bets_send = 0
        while bytes_read < data_size:
            amount_bets_send += 1

            # Read the size of the bet
            bet_size = struct.unpack('B', data[bytes_read:bytes_read+1])[0]
            bytes_read += 1

            # Read the bet data
            bet_data = data[bytes_read:bytes_read+bet_size]
            bytes_read += bet_size

            # Deserialize the bet
            bet = Bet.deserialize(bet_data)
            if bet is not None:
                bets.append(bet)

        return amount_bets_send, bets
    
    def receive_client_id(self, skt):
        data = skt.recvall(1)
        return struct.unpack('B', data)[0]

    def send_response_batch(self, ok, skt):
        if ok:
            skt.sendall(struct.pack('B', RESPONSE_BATCH_COMMAND_OK))
        else:
            skt.sendall(struct.pack('B', RESPONSE_BATCH_COMMAND_ERROR))

    def send_response_winners_by_agency(self, skt, winning_documents):
        """
        Send the winners documents of the corresponding lottery to each agency
        Encoding:
        - Command (1 byte)
        - Size data (4 bytes) (big endian)
        - Data (4 bytes per document) (big endian)
        """
        data_documents = b''
        for document in winning_documents:
            data_documents += struct.pack('>I', document)

        data = b''
        data += struct.pack('B', RESPONSE_WINNERS_COMMAND)
        data += struct.pack('>I', len(data_documents))
        data += data_documents
        skt.sendall(data)