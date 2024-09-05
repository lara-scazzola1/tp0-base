from common.socket import *
from common.protocol import *
from common.utils import *
import multiprocessing
import logging
import signal

TOTAL_CONNECTIONS = 5

class Server:
    def __init__(self, port, listen_backlog):
        self._server_socket = Socket()
        self._server_socket.bind_and_listen(port, listen_backlog)
        self._stop = False
        self._protocol = Protocol()
        self._connections = {}
        self._processes = {}
        self._file_lock = multiprocessing.Lock()
        
    def __handle_receive_client_id(self, client_sock):
        """
        Receive a client id from a client socket
        """
        client_id = self._protocol.receive_client_id(client_sock)
        self._connections[client_id] = client_sock

    def __handle_receive_batch(self, client_sock):
        """
        Receive a batch of bets from a client socket and store them.
        Return the amount of bets received and if the operation was 
        successful.
        """
        amount_bets_send, bets_received = self._protocol.receive_batch(client_sock)
        with self._file_lock:
            store_bets(bets_received)
        ok = amount_bets_send == len(bets_received)
        self._protocol.send_response_batch(ok, client_sock)
        return ok, len(bets_received)

    def __get_winners_by_agency(self):
        """
        Get the winners of the lottery
        """
        bets = load_bets()
        winners_by_agency = [[] for i in range(TOTAL_CONNECTIONS)]
        for bet in bets:
            if has_won(bet):
                winners_by_agency[bet.agency - 1].append(int(bet.document))
        return winners_by_agency

    def __send_winners(self, winners_by_agency):
        """
        Send the winners of the lottery to the clients
        """
        for id, conn in self._connections.items():
            self._protocol.send_response_winners_by_agency(conn, winners_by_agency[id - 1])
            
    def __handle_close_connections(self):
        for conn in self._connections.values():
            conn.close()   

    def __handle_client_connection(self, client_sock, client_id):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            print(f"Client {client_id} connected")
            while not self._stop:
                command = self._protocol.receive_command(client_sock)

                if command == CLIENT_ID:
                    self.__handle_receive_client_id(client_sock)

                if command == BATCH_COMMAND:
                    ok, amount = self.__handle_receive_batch(client_sock)
                    if ok:
                        logging.info(f"action: apuesta_recibida | result: success | cantidad: {amount}")
                    else:
                        logging.error(f"action: apuesta_recibida | result: fail | cantidad: {amount}")  

                if command == WAIT_WINNERS_COMMAND:
                    break

        except OSError as e:
            logging.error("action: receive_message | result: fail | error: ", e)

    
    def stop_server(self, signum, frame):
        self._server_socket.close()
        self._stop = True
        self.__handle_close_connections()

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        signal.signal(signal.SIGTERM, self.stop_server)

        try:
            while not self._stop:
                client_sock = self._server_socket.accept()
                if client_sock:
                    _ = self._protocol.receive_command(client_sock)
                    client_id = self._protocol.receive_client_id(client_sock)
                    self._connections[client_id] = client_sock

                    process = multiprocessing.Process(target=self.__handle_client_connection, args=(client_sock, client_id))
                    self._processes[client_id] = process
                    process.start()

                if len(self._connections) == TOTAL_CONNECTIONS:
                    for process in self._processes.values():
                        process.join()

                    winners_by_agency = self.__get_winners_by_agency()
                    self.__send_winners(winners_by_agency)
                    logging.info("action: sorteo | result: success")

                    self.__handle_close_connections()
                    self._connections = {}

        except Exception as e:
            logging.error("action: run | result: fail | error: ", e)

