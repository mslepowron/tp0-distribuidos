#!/bin/bash
SERVER_CONFIG_FILE="server/config.ini"

TEST_MESSAGE="Hello, Server!"

#busca la expresion regular, por ej: SERVER_PORT = nro_puerto, hace un trim para quedarse
#con la segunda parte (despues del =) y elimina los espacios extra
get_config() {
    local key="$1"
    grep -E "^$key" "$SERVER_CONFIG_FILE" | cut -d'=' -f2 | tr -d '[:space:]'
}

SERVER_PORT=$(get_config "SERVER_PORT")
SERVER_IP=$(get_config "SERVER_IP")

RESPONSE=$(docker run --rm --network tp0_testing_net alpine:latest sh -c "apk add --no-cache netcat-openbsd >/dev/null 2>&1 && echo '$TEST_MESSAGE' | nc $SERVER_IP $SERVER_PORT")

if [ "$RESPONSE" = "$TEST_MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi