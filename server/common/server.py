from common.socket import *
from common.protocol import *
from common.utils import *
import logging
import signal


class Server:
    def __init__(self, port, listen_backlog):
        self._socket = Socket()
        self._socket.bind_and_listen(port, listen_backlog)
        self._stop = False

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        signal.signal(signal.SIGTERM, self.stop_server)

        while not self._stop:
            client_sock = self._socket.accept()
            if client_sock:
                self.__handle_client_connection(client_sock)
        


    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            protocol = Protocol(client_sock)
            while not self._stop:
                command = protocol.receive_command()

                if command == BATCH_COMMAND:
                    data_size = protocol.receive_data_size()
                    amount_bets_send, bets = protocol.receive_batch(data_size)
                    if amount_bets_send != len(bets):
                        logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(bets)}")
                        protocol.send_response_batch(amount_bets_send)
                        continue    
                    logging.info(f"action: apuesta_recibida | result: success | cantidad: {amount_bets_send}")
                    #print(bets)
                    store_bets(bets)
                    protocol.send_response_batch(len(bets))

                if command == BET_COMMAND:
                    data_size = protocol.receive_data_size()
                    bet = protocol.receive_bet(data_size)
                    store_bets([bet])
                    logging.info(f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}")
                    protocol.send_response_bet()

                if command == DISCONNECT_COMMAND:
                    print("SE RECIBE DESCONECTAR")
                    logging.info("action: desconexion | result: success")
                    break
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: ", e)
        except Exception as e:
            logging.error("action: receive_message | result: fail | error: ", e)
        finally:
            protocol.close()

    
    def stop_server(self, signum, frame):
        self._server_socket.close()
        self._stop = True

