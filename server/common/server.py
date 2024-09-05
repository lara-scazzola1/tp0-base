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
        self._protocol = Protocol()
        self._client_sock = None

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        signal.signal(signal.SIGTERM, self.stop_server)

        while not self._stop:
            self._client_sock = self._server_socket.accept()
            if self._client_sock:
                self.__handle_client_connection()
        
    def __handle_receive_batch(self):
        """
        Receive a batch of bets from a client socket and store them.
        Return a tuple with the result of the operation and the amount 
        of bets received
        """
        amount_bets_send, bets_received = self._protocol.receive_batch(self._client_sock)
        store_bets(bets_received)
        ok = amount_bets_send == len(bets_received)
        self._protocol.send_response_batch(ok, self._client_sock)
        return ok, len(bets_received)

    def __handle_client_connection(self):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            while not self._stop:
                command = self._protocol.receive_command(self._client_sock)

                if command == BATCH_COMMAND:
                    ok, amount = self.__handle_receive_batch()
                    if ok:
                        logging.info(f"action: apuesta_recibida | result: success | cantidad: {amount}")
                    else:
                        logging.error(f"action: apuesta_recibida | result: fail | cantidad: {amount}")  

                if command == DISCONNECT_COMMAND:
                    logging.info("action: desconexion | result: success")
                    break

        except OSError as e:
            logging.error("action: receive_message | result: fail | error: ", e)
        finally:
            self._client_sock.close()

    
    def stop_server(self, signum, frame):
        self._server_socket.close()
        self._stop = True
        if self._client_sock:
            self._client_sock.close()

