import socket
import logging
import signal
import sys
import threading
from common import utils, messages

class Server:
    def __init__(self, port, listen_backlog, clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.shutdown = False
        self.lottery_finished = False
        self.finished_clients = 0
        self.client_winners = {}
        self.total_clients = clients
        self.lottery_results_lock = threading.Lock()
        self.bet_storting_lock = threading.Lock()
        self.threads = []


    def shutdown_server(self):
        if self.shutdown:
            return
        
        self.shutdown = True
        logging.info('action: shutdown | result: in_progress')
        try:
            self._server_socket.close()
        except Exception as e:
            logging.error(f'action: shutdown | result: fail')
        self.__join_client_threads()
        logging.info(f'action: shutdown | result: success')


    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while not self.shutdown:
            try:
                client_sock = self.__accept_new_connection()
                if client_sock:
                    client_thread = threading.Thread(
                    target=self.__handle_client_connection, args=(client_sock,), daemon=True)
                    client_thread.start()
                    self.threads.append(client_thread)
            except:
                if self.shutdown:
                    break
        
        self.__join_client_threads(self)

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            message = messages.recieve_client_messasge(client_sock)
            if message == "BETS":
                agency_bets = 0
                while True:
                    message = messages.recieve_client_messasge(client_sock)
                    
                    is_eof, agency_id = messages.is_end_of_agency_file(message)
                    if is_eof:
                        logging.info(f"action: apuesta_recibida | result: success | cantidad: {agency_bets}")
                        
                        ack_str = "{};{}\n".format(agency_bets, agency_id)
                        ack_bytes = ack_str.encode("utf-8")
                        messages.send_ack_client(client_sock, ack_bytes)
                        with self.lottery_results_lock:
                            self.finished_clients += 1   
                        break
                    else:
                        try:
                            bets = messages.decode_batch_bets(message)
                            with self.bet_storting_lock:
                                utils.store_bets(bets)
                                #logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
                                agency_bets += len(bets)

                            ack_str = "BATCH_OK;{}\n".format(len(bets))
                            ack_bytes = ack_str.encode("utf-8")
                            messages.send_ack_client(client_sock, ack_bytes)
                        except Exception as e:
                            batch_size = len(bets)
                            logging.error(f"action: apuesta_recibida | result: fail | cantidad: {batch_size}")

                            ack_str = "ERROR_BATCH;{}\n".format(batch_size)
                            ack_bytes = ack_str.encode("utf-8")
                            messages.send_ack_client(client_sock, ack_bytes)
            if message.startswith("LOTERY_WINNER;"):
                agency_id = message.split(";")[1]
                with self.lottery_results_lock:
                    if int(self.finished_clients) == int(self.total_clients):
                        if not self.lottery_finished:
                            self.lottery_finished = True
                            self.__process_lottery_winners()
                                
                            logging.info(f'action: sorteo | result: success')
                    
                        if self.lottery_finished:
                            winners = self.client_winners.get(agency_id, [])
                            response_winners = "WINNERS;" + ";".join(winners) + "\n"
                            messages.send_ack_client(client_sock, response_winners.encode("utf-8"))
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()


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

    def __process_lottery_winners(self):
        """
        loads de agencys bets and uses has_won to calculate the
        agency's winners
        """
        all_bets = list(utils.load_bets())
        winners_by_agency = {}

        agencies = set(str(bet.agency) for bet in all_bets)
        for agency_id in agencies:
            winners_by_agency[agency_id] = []

        for bet in all_bets:
            if utils.has_won(bet):
                agency_id = str(bet.agency)
                winners_by_agency[agency_id].append(bet.document)

        self.client_winners = winners_by_agency

    def __join_client_threads(self):
        for thread in self.threads:
            thread.join()