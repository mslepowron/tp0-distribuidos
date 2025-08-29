import socket
import logging
import signal
import sys
from common import utils, messages



class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.shutdown = False
        self.client_sockets = []


    def shutdown_server(self):
        if self.shutdown:
            return
        
        self.shutdown = True
        logging.info('action: shutdown | result: in_progress')
        try:
            self._server_socket.close()
            for client_sock in self.client_sockets:
                try:
                    client_sock.close()
                except Exception as e:
                    logging.info(f'action: close_client_socket | result: fail')
                else:
                    logging.info(f'action: close_client_socket | result: success')
            logging.info(f'action: shutdown | result: success')
        except Exception as e:
            logging.error(f'action: shutdown | result: fail')


    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while not self.shutdown:
            try:
                client_sock = self.__accept_new_connection()
                self.client_sockets.append(client_sock)
                self.__handle_client_connection(client_sock)
            except:
                if self.shutdown:
                    break

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # # TODO: Modify the receive to avoid short-reads
            # msg = client_sock.recv(1024).rstrip().decode('utf-8')
            # addr = client_sock.getpeername()
            # logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # # TODO: Modify the send to avoid short-writes
            # client_sock.send("{}\n".format(msg).encode('utf-8'))
            
            # decode cliet msg
            message = messages.recieve_client_messasge(client_sock)
            addr = client_sock.getpeername()
            logging.info(f'action: recieved_client_message | result: success | ip: {addr[0]} | msg: {message}')
            #save bet
            bet = messages.decode_message(message)
            utils.store_bets([bet])
            #log save bet
            logging.info(f'action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}')
            #send ack
            ack_str = "{};{}\n".format(bet.document, bet.number)
            ack_bytes = ack_str.encode("utf-8")
            messages.send_ack_client(client_sock, ack_bytes)
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()
            if client_sock in self.client_sockets:
                self.client_sockets.remove(client_sock)

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
