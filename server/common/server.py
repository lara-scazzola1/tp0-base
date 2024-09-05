from common.socket import *
from common.protocol import *
from common.utils import *
import logging
import signal


class Server:
    def __init__(self, port, listen_backlog):
        self._server_socket = Socket()
        self._server_socket.bind_and_listen(port, listen_backlog)
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
            client_sock = self._server_socket.accept()
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
            command, data_size = protocol.receive_command()
            if command == BET_COMMAND:
                bet = protocol.receive_bet(data_size)
                store_bets([bet])
                logging.info(f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}")
                protocol.send_response_bet()
                return
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            protocol.close()

    
    def stop_server(self, signum, frame):
        """
        Stops the server
        """
        self._server_socket.close()
        self._stop = True

