import socket
import logging

class Socket:
    def __init__(self, skt=None):
        """
        Socket class constructor.
        """
        if skt:
            self._skt = skt
        else:
            self._skt = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    def bind_and_listen(self, port, listen_backlog):
        """
        Bind the socket to the given port and start listening.
        """
        self._skt.bind(('', port))
        self._skt.listen(listen_backlog)

    def connect(self, address):
        """
        Connect the socket to the given address.
        """
        self._skt.connect(address)

    def accept(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """
        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        try:
            client_skt, addr = self._skt.accept()
        except OSError as e:
            return None
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return Socket(client_skt)

    def recvall(self, bufsize):
        """
        Receive data from the socket until the buffer is
        filled or the connection is closed.
        """
        data = b''
        while len(data) < bufsize:
            packet = self._skt.recv(bufsize - len(data))
            if not packet:
                break
            data += packet
        return data

    def sendall(self, data):
        """
        Send all the data through the socket.
        """
        total_sent = 0
        while total_sent < len(data):
            sent = self._skt.send(data[total_sent:])
            if sent == 0:
                raise RuntimeError("Socket connection broken")
            total_sent += sent

    def close(self):
        """
        Close the socket.
        """
        self._skt.close()

    def recv(self, bufsize):
        """
        Receive data from the socket.
        """
        return self._skt.recv(bufsize)
    
    def send(self, data):
        """
        Send data through the socket.
        """
        self._skt.send(data)

    def getpeername(self):
        return self._skt.getpeername()
