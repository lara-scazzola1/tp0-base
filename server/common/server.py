from common.socket import *
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
            # TODO: Modify the receive to avoid short-reads
            msg = client_sock.recv(1024).rstrip().decode('utf-8')
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # TODO: Modify the send to avoid short-writes
            client_sock.send("{}\n".format(msg).encode('utf-8'))
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    
    def stop_server(self, signum, frame):
        self._server_socket.close()
        self._stop = True

