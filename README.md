# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1

### Ejercicio 1

Se genero un script de bash generar-compose.sh, que permite configurar un nombre de archivo de configuracion ```.yaml``` y una cantidad de clientes determinada.

  **Uso:**  
  Se debe correr, desde la raiz del proyecto:
  ```
  ./generar-compose.sh <archivo_de_salida.yaml> <cantidad de clientes>
  ```

Se creo a su vez un sub-script de python, para el armado del archivo. Una vez verificados los parametros ingresados por el usuario al correr el comando anterior, se invoca al sub-script para el armado del archivo. Una vez finalizado el proceso, se genera un archivo valido, listo para usar con docker compose.

---

### Ejercicio 2

Para lograr que realizar cambios en el archivo de configuración no requiera reconstruir las imágenes de Docker, se agrega en el script que crea el archivo de configuración .yaml la seccion de ```volumnes```, tanto en el service del servidor como el del cliente.

Al utilizar Docker Volumes logramos persistir datos fuera del contenedor, tal que los archivos de configuracion se mantengan fuera de la imagen de Docker, y puedan ser modificados en el host. Cualquier cambio en estos archivos se refleja automáticamente en el contenedor sin necesidad de reconstruir la imagen.

Para correrlo, se puede continuar utilizando el generador de script del ejercicio 1, de la siguiente manera:

  **Uso:**  
  Se debe correr, desde la raiz del proyecto:
  ```
  ./generar-compose.sh <archivo_de_salida.yaml> <cantidad de clientes>
  ```

  y luego se levantan los containers de docker con el comando:

  ```
  make docker-compose-up
  ```

---

### Ejercicio 3

Se desarrolló un script ```validar-echo-server.sh ``` que testea el correcto funcionamiento del echo server con netcat.

Para ello, se toman los datos del Puerto y la IP del servidor de su respectivo archivo de configuracion, y luego se corre un comando de docker que levanta un nuvo contenedor en la red interna del tp (testing_net). Se levanta una imagen liviana de Linux (alpine) y se envía un mensaje al echo server para testear el funcionamiento. 

Se captura la respuesta del server en uan variable Response, y se verifica que esta sea exactamente la misma que el mensaje enviado. Esto determinaría que el server está funcionando correctamente.

**Uso:**  
  El usuario debe contar con los permisos necesarios para correr ambos scripts de bash. Si no los tiene:

  ```
  chmod +x generar-compose.sh
  ```

  ```
  chmod +x validar-echo-server.sh
  ```

  Se puede generar desde la raíz del proyecto el archivo de configuración inicial como se hizo en los ejercicios previos:
  ```
  ./generar-compose.sh <archivo_de_salida.yaml> <cantidad de clientes>
  ```

  Luego levantamos los contenedores correspondientes:
  ```
  make docker-compose-up
  ```

  Y corremos el script de verificacion de funcionamiento del servidor:
  ```
  ./validar-echo-server.sh
  ```

---

### Ejercicio 4

Se agregó una funcion shutdown_server() para cerrar ordenadamente los recursos del socket del servidor, y los sockets de los clientes conectados y se registra esto mismo en logs que van detallando el proceso.

El proceso luego termina de forma controlada con ``` sys.exit(0)```

En cuanto al Client, se creo un canal sigChannel que puede recibir señales ```SIGTERM```. Cuando se inicia la iteracion en loop del cliente, se chequea si en el canal se recibio alguna señal (caso en el cual se hace el shutdown).

De esta manera, el cliente tambien cierra de manera segura si se apaga el contenedor de Docker.

Podemos probar que los sistemas terminan de manera graceful utilizando el comando ```docker compose down``` acompañado del flag ```-t```

Docker va a enviar la señal ```SIGTERM```, y si los procesos saben manejarla, tienen hasta ´t´ segundos para cerrar todo ordenadamente. Sino, Docker envia ```SIGKILL``` para matar al proceso.

  **Uso:**  
  Se corre para generar el archivo de configuracion, desde la raiz del proyecto:
  ```
  ./generar-compose.sh <archivo_de_salida.yaml> <cantidad de clientes>
  ```

  Se levantan los contenedores de docker:
  ```
  make docker-compose-up
  ```
  Corremos el siguiente comando de docker para bajar los container y enviar la señal ```SIGTERM```

  ```
  docker compose -f docker-compose-dev.yaml down -t 15
  ```


## Parte 2

### Ejercicio 5

Se implentó un modulo _agency message_ para el cliente. Los datos de la apuesta se definen como variables de entorno, y se utilizan para construir un mensaje con el siguiente formato:

 ``` go
msg := fmt.Sprintf("%s;%s;%s;%s;%s;%s", bet.AgencyId, bet.Name, bet.LastName, bet.Document, bet.BirthDate, bet.Number) 
```

Por otro lado, en la funcion ```StartClientLoop()``` del modulo client, se crea el socket, se formatea el mensaje especificado arriba y se llama a una funcion ```SendClientMessage()``` que se encarga de enviarle el mensaje con los datos de la apuesta al servidor.

En cuanto al envío de los datos, en go, utilizar
``` go
conn.Write(message)
```
a secas no nos garantiza que se envíen todos los datos que queremos. Puede suceder que la cantidad de bytes enviados sea menor al largo total de los datos (short write), por ejemplo, porque el buffer del socket estaba lleno y solo acepto y envió una parte.

Para evitar que esto suceda, se tuvieron en cuenta dos cosas:

Por un lado, establecer un limite de tamaño máximo de mensaje a enviar (8KB). Esto evita saturar el buffer del socket, aporta algo de robustez y seguridad (por ejemplo evitando que un cliente malicioso mande mesajes excesivamente grandes que saturen el servidor o provoquen un DoS), y aporta simplicidad; mantener mensajes de tamaño moderado facilita el manejo de buffers y control de errores.

Por otro lado, se implementó una función ```WriteFull()``` donde se repite la operación de escritura en un loop hasta que efectivamente se hallan enviado todos los bytes del mensaje, evitando el envío de mensajes incompletos.

El mensaje con los datos de la apuesta es precedido por un mensaje de tamaño de 4 bytes que contiene la longitud del mensaje del cliente. Este paso se realiza para que el server sepa cuántos bytes tiene que leer y así evitar un short read.

Es decir, el protocolo esta conformado por 4 bytes (message length) + payload (client message)

El cliente espera un mensaje de confirmación de parte del servidor para verificar que haya recibio correctamente los datos de la apuesta. Este mismo contiene los datos que presenta el log del server cuando persiste una apuesta correctamente; es decir, el DNI y el Numero del cliente.

Si el cliente no recibe este mensaje de confirmación, se logguea el problema y tiene un maximo de 3 reintentos para conectarse al servidor y mandar los datos de la apuesta. Una vez agotados los 3 reintentos, se logguea un error final de fallo de envio del mensaje, y se finaliza el programa.

  **Uso:**  
  Se debe correr, desde la raiz del proyecto:
  ```
  ./generar-compose.sh <archivo_de_salida.yaml> <cantidad de clientes>
  ```
  por ejemplo:
  ```
  ./generar-compose.sh docker-compose-dev.yaml 1
  ```
  Luego se levanta el sistema con:
  ```
  make docker-compose-up
  ```

---

### Ejercicio 6

Para el envío de apuestas en batch al servidor se modificó la función   ```StartClientLoop()``` de tal forma que en lugar de formatear y enviar la apuesta de un cliente, se realiza una lectura del archivo csv de la agencia correspondiente al Id, se almacenan todas las apuestas de esa agencia y se las va enviando al server en chunks del tamaño MaxAmount que se encuentra especificado en el archivo de configuración del cliente.

De las apuestas correspondientes a la agencia, se van tomando chunks del tamaño batch configurado y se llama a ```FormatBatchMessage()``` para configurar el formato del mensaje que se enviará al servidor.

Se conserva el protocolo desarrollado en el ejercicio 5, y se definió que la separación de las diferentes apuestas que se envian en el batch se identifiquen con un ```\n```

Por cada batch enviado, el server logguea el proceso y le envia un ack al cliente, que puede ser de tipo ```BATCH_OK``` si todas las apuestas del batch se procesaron y almacenaron correctamente, o ```ERROR_BATCH``` si falló alguna. Para el primer caso, el cliente continua enviando el siguiente batch, y para el segundo caso, se corta la ejecución y se cierra la conexión.

Cuando finaliza el proceso de envío por batch, el cliente envía un último mensaje al server con el siguiente formato:

```END_OF_FILE;<agencyID>```

Indicando que ha finalizado el envío de datos de sus apuestas, y que corresponden a ese Id de agencia.

El server procesa los datos recibidos y almacena las apuestas. En caso de que haya recibido exitosamente toda la información del archivo de la agencia, contesta con un ack final al cliente, con la siguiente información:

```<countOfBets>;<agencyID>```

Le adjunta la cantidad de apuestas que almacenó correctamente, y el Id de esa agencia. 

Si la cantidad de apuestas recibia en el ack del server coincide con la cantidad de apuestas totales de esa agencia, la agencia logguea:

```go
log.Infof("action: apuesta_enviada | result: success | amount: %v", <agencyID>, agencyBets)
```

  **Uso:**  
  Se debe correr, desde la raiz del proyecto:
  ```
  ./generar-compose.sh <archivo_de_salida.yaml> <cantidad de clientes>
  ```
  por ejemplo:
  ```
  ./generar-compose.sh docker-compose-dev.yaml 1
  ```
  Luego se levanta el sistema con:
  ```
  make docker-compose-up
  ```

---

### Ejercicio 7

_Flujo de Mensajes_

Para este ejercicio se realizaron algunas modificaciones a el flujo de mensajes que se manejaba anteriormente. En cuanto a los _batchs_ de apuestas, se conserva el protocolo desarrollado en el ejercicio 5, y la separación de las diferentes apuestas que se envian en el batch se identifican con un ```\n```

Teniendo en cuenta que el cliente también se tiene que comunicar con el server para solicitar a los ganadores del sorteo, se tomó como convención diferenciar estos dos casos, enviando un mensaje más antes de enviar los _batch_ de apuestas.

Si el cliente se está comunicando con el server para enviarle las apuestas, antes de mandar los batches envía un mensaje ```BETS```  

Una vez que el cliente finaliza el envío de las apuestas, envía un mensaje ```END_OF_FILE;<id_cliente>``` para notificar al servidor que ha terminado. Una vez que el server recibe este mensaje, guarda internamente el dato de que ha finalizado una agencia más en el envío de datos.

Inmediatamente después, el cliente le manda al server un mensaje ```LOTERY_WINNER;<id_cliente>``` consultando si estan listos los resultados de los sorteos.

Los clientes loopean sobre una función 
```go
func (c *Client) WaitForLoteryResults(sigChannel chan os.Signal) error {
} 
```

Donde esperan el resultado del sorteo. En cada loop de la funcion, envia un mensaje ```LOTERY_WINNER;<id_cliente>```. Si el servidor le contesta con un mensaje con el header ```WINNERS```, significa que el sorteo ha finalizado (todos los clientes cargaron sus apuestas) y los resultados están listos. En caso de que no estan disponibles los resultados todavía, un cliente que ya envío sus apuestas, y esta consultando por sus resultados, continúa loopeando sobre la función ```WaitForLoteryResults``

_Conexiones_

En cuanto a la primera parte del proceso (envío de apuestas), la comunicación entre cada cliente y servidor se realiza sobre una única conexión.

En cuanto a la segunda parte, para consultar si los resultados están listos cada cliente levanta una conexión con el servidor por cada iteración de la función ``WaitForLoteryResults``:
  - Si el server recibió los datos de todas las agencias de  lotería, le contesta con un mensaje de exito al cliente con el header ```WINNER```, seguido de todos los documentos de los ganadores de esa agencia y el cliente finaliza su trabajo
  - En caso de que no esten listos los datos de todas las agencias, el cliente cierra la conexión y espera un determinado tiempo para volver a reintentar la consulta.

_Cantidad de clientes esperados configurable_

El servidor debe esperar a recibir los datos de todas las agencias conectadas antes de procesar el sorteo. Para que la cantidad de clientes se pueda manejar de manera flexible, se utiliza una variable de entorno para que el servidor sepa a cuantos clientes debe esperar antes de enviar a los ganadores. Esto se configura según el valor almacenado en el docker-compose-dev.yaml, que se configura según el script desarrollado en los primeros puntos.

_Otras Aclaraciones_

El server le manda un mensaje BATCH OK cada vez que procesa correctamente un batch del client. En el ejercicio anterior, cada vez que
el cliente recibia ese mensaje imprimia un log. Para este ejercicio se eliminan esos logs y quedo unicamente un log final cuando el
server le responde el mensaje EOF al client (es decirque proceso todos los batchs correctamente y puede seguir.)

Esto se hizo para mayor claridad a la hora de leer los logs y que se vea bien que se enviaron todos los batch, que consulta por el sorteo y hasta que no esten todos los clientes listos no recibe respuesta, y que luego de recibirla sale con exit code 0.

  **Uso:**  
  Se debe correr, desde la raiz del proyecto:
  ```
  ./generar-compose.sh <archivo_de_salida.yaml> <cantidad de clientes>
  ```
  por ejemplo:
  ```
  ./generar-compose.sh docker-compose-dev.yaml 5
  ```
  Luego se levanta el sistema con:
  ```
  make docker-compose-up
  ```