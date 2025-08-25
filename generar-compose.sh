#!/bin/bash

function error_output(){
    echo "Error: $1"
    exit 1
}

FILE=$1
CLIENTS=$2

function check_input_parameters(){
    if [ $# -ne 2 ]; then
        error_output "Cantidad incorrecta de parametros. Uso: $0 <archivo_salida.yaml> <cantidad_clientes>"
    fi

    if ! [[ $CLIENTS =~ ^[0-9]+$ ]]; then
        error_output "El parametro de cantidad de clientes debe ser un numero entero positivo."
    fi

    if [[ "$FILE" != *.yaml ]]; then
        error_output "El archivo de salida debe tener extension .yaml"
    fi
}

echo "Archivo de salida: $FILE"
echo "Cantidad de clientes: $CLIENTS"
python3 mi-generador.py $FILE $CLIENTS || error_output "Fallo la generacion del archivo YAML."