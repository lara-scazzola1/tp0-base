from common.socket import *
from common.protocol import *
from common.utils import *
import logging
import signal

TOTAL_CONNECTIONS = 5

class Server:
    def __init__(self, port, listen_backlog):
        self._socket = Socket()
        self._socket.bind_and_listen(port, listen_backlog)
        self._stop = False
        self._protocol = Protocol()
        self._connections = []
        self._amount_wait_winners = 0

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
                self._connections.append(client_sock)
                self.__handle_client_connection(client_sock)
        


    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            while not self._stop:
                command = self._protocol.receive_command(client_sock)

                if command == BATCH_COMMAND:
                    amount_bets_send, bets = self._protocol.receive_batch(client_sock)
                    ok = amount_bets_send == len(bets)
                    if ok:
                        #logging.info(f"action: apuesta_recibida | result: success | cantidad: {amount_bets_send}")
                        #print(bets)
                        store_bets(bets)
                    else:
                        logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(bets)}")  
                    self._protocol.send_response_batch(amount_bets_send, client_sock)

                if command == BET_COMMAND:
                    bet = self._protocol.receive_bet(client_sock)
                    store_bets([bet])
                    logging.info(f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}")
                    self._protocol.send_response_bet(client_sock)

                if command == DISCONNECT_COMMAND:
                    logging.info("action: desconexion | result: success")
                    break

                if command == WAIT_WINNERS_COMMAND:
                    logging.info("action: esperar_ganador | result: success")
                    self._amount_wait_winners += 1
                    if self._amount_wait_winners == TOTAL_CONNECTIONS:
                        for connection in self._connections:
                            self._protocol.send_response_winners(connection)
                    break

        except OSError as e:
            logging.error("action: receive_message | result: fail | error: ", e)
        except Exception as e:
            logging.error("action: receive_message | result: fail | error: ", e)
        finally:
            if self._amount_wait_winners == TOTAL_CONNECTIONS:
                for connection in self._connections:
                    connection.close()

    
    def stop_server(self, signum, frame):
        self._server_socket.close()
        self._stop = True

