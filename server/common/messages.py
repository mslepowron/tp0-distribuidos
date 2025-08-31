import socket
import logging
from typing import List
from common.utils import Bet

MAX_MESSAGE_SIZE = 8192
BET_FIELDS_COUNT = 6

def recieve_full(client_sock: socket.socket, size: int) -> bytes:
    """""
    Recives `size` bytes from a client socket and reads exactly that amount,
    to avoid short reads
    """""
    data = b""
    while len(data) < size:
        chunk = client_sock.recv(size - len(data))
        if not chunk:
            raise ConnectionError("Client closed connection during recv")
        data += chunk
    return data

def recieve_client_messasge(client_sock: socket.socket) -> str:
    """""
    Recives a client message containing information about
    a client's bet. A client's message cannot exceed a size of 8KB
    """""
    length_bytes = recieve_full(client_sock, 4)
    message_length = int.from_bytes(length_bytes, "big")

    if message_length > MAX_MESSAGE_SIZE:
        logging.error("action: recieve_client_message | result: fail | message is bigger than 8KB")
        raise ValueError("message is bigger than 8KB")

    payload_bytes = recieve_full(client_sock, message_length)
    message = payload_bytes.decode("utf-8")

    return message

def decode_message(message: str) -> Bet:
    """""
    Decodes a clients message according to the protocol established
    between the server and the client
    """""
    parts = message.split(";")
    if len(parts) != BET_FIELDS_COUNT:
        raise ValueError(f"Invalid bet format: {message}")

    agency, first_name, last_name, document, birthdate, number = parts
    return Bet(agency, first_name, last_name, document, birthdate, number)

#send ack modificado para evitar short-write
def send_ack_client(client_sock: socket.socket, ack: bytes) -> None:
    """""
    Sends an ACK message to the client containing their document number
    and the bet number to confirm the be was correctly recieved and stored.
    """""
    total_sent = 0
    while total_sent < len(ack):
        sent = client_sock.send(ack[total_sent:])
        if sent == 0:
            raise RuntimeError("Conexión cerrada mientras se enviaba ACK")
        total_sent += sent


def is_end_of_agency_file(message: str):
    if message.startswith("END_OF_FILE;"):
        parts = message.split(";")
        if len(parts) != 2:
            raise ValueError(f"Invalid END_OF_FILE message: {message}")
        is_eof = True
        agency_id = parts[1]
    else:
        is_eof = False
        agency_id = ""
    
    return is_eof, agency_id

def decode_batch_bets(batch_message: str) -> List[Bet]:
    """
    Decodes a batch message containing multiple bets separated by \n.
    Procceses each bet using the decode_message in order to store it.
    """
    bets = []
    lines = batch_message.strip().split("\n")
    for line in lines:
        if line.strip():
            bet = decode_message(line)
            bets.append(bet)
    return bets
