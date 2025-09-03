# TP0: Docker + Comunicaciones + Concurrencia

En el presente repositorio se provee un esqueleto básico de cliente/servidor, en donde todas las dependencias del mismo se encuentran encapsuladas en containers. Los alumnos deberán resolver una guía de ejercicios incrementales, teniendo en cuenta las condiciones de entrega descritas al final de este enunciado.

 El cliente (Golang) y el servidor (Python) fueron desarrollados en diferentes lenguajes simplemente para mostrar cómo dos lenguajes de programación pueden convivir en el mismo proyecto con la ayuda de containers, en este caso utilizando [Docker Compose](https://docs.docker.com/compose/).

## Instrucciones de uso
El repositorio cuenta con un **Makefile** que incluye distintos comandos en forma de targets. Los targets se ejecutan mediante la invocación de:  **make \<target\>**. Los target imprescindibles para iniciar y detener el sistema son **docker-compose-up** y **docker-compose-down**, siendo los restantes targets de utilidad para el proceso de depuración.

Los targets disponibles son:

| target  | accion  |
|---|---|
|  `docker-compose-up`  | Inicializa el ambiente de desarrollo. Construye las imágenes del cliente y el servidor, inicializa los recursos a utilizar (volúmenes, redes, etc) e inicia los propios containers. |
| `docker-compose-down`  | Ejecuta `docker-compose stop` para detener los containers asociados al compose y luego  `docker-compose down` para destruir todos los recursos asociados al proyecto que fueron inicializados. Se recomienda ejecutar este comando al finalizar cada ejecución para evitar que el disco de la máquina host se llene de versiones de desarrollo y recursos sin liberar. |
|  `docker-compose-logs` | Permite ver los logs actuales del proyecto. Acompañar con `grep` para lograr ver mensajes de una aplicación específica dentro del compose. |
| `docker-image`  | Construye las imágenes a ser utilizadas tanto en el servidor como en el cliente. Este target es utilizado por **docker-compose-up**, por lo cual se lo puede utilizar para probar nuevos cambios en las imágenes antes de arrancar el proyecto. |
| `build` | Compila la aplicación cliente para ejecución en el _host_ en lugar de en Docker. De este modo la compilación es mucho más veloz, pero requiere contar con todo el entorno de Golang y Python instalados en la máquina _host_. |

### Servidor

Se trata de un "echo server", en donde los mensajes recibidos por el cliente se responden inmediatamente y sin alterar. 

Se ejecutan en bucle las siguientes etapas:

1. Servidor acepta una nueva conexión.
2. Servidor recibe mensaje del cliente y procede a responder el mismo.
3. Servidor desconecta al cliente.
4. Servidor retorna al paso 1.


### Cliente
 se conecta reiteradas veces al servidor y envía mensajes de la siguiente forma:
 
1. Cliente se conecta al servidor.
2. Cliente genera mensaje incremental.
3. Cliente envía mensaje al servidor y espera mensaje de respuesta.
4. Servidor responde al mensaje.
5. Servidor desconecta al cliente.
6. Cliente verifica si aún debe enviar un mensaje y si es así, vuelve al paso 2.

### Ejemplo

Al ejecutar el comando `make docker-compose-up`  y luego  `make docker-compose-logs`, se observan los siguientes logs:

```
client1  | 2024-08-21 22:11:15 INFO     action: config | result: success | client_id: 1 | server_address: server:12345 | loop_amount: 5 | loop_period: 5s | log_level: DEBUG
client1  | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:14 DEBUG    action: config | result: success | port: 12345 | listen_backlog: 5 | logging_level: DEBUG
server   | 2024-08-21 22:11:14 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°3
client1  | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°3
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°5
client1  | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°5
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:40 INFO     action: loop_finished | result: success | client_id: 1
client1 exited with code 0
```


## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1:
Definir un script de bash `generar-compose.sh` que permita crear una definición de Docker Compose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```

En el archivo de Docker Compose de salida se pueden definir volúmenes, variables de entorno y redes con libertad, pero recordar actualizar este script cuando se modifiquen tales definiciones en los sucesivos ejercicios.

### Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera reconstruír las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida por fuera de la imagen (hint: `docker volumes`).


### Ejercicio N°3:
Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un echo server, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se pueden exponer puertos del servidor para realizar la comunicación (hint: `docker network`). `


### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.



#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).


### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). 
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los archivos deberán ser inyectados en los containers correspondientes y persistido por fuera de la imagen (hint: `docker volumes`), manteniendo la convencion de que el cliente N utilizara el archivo de apuestas `.data/agency-{N}.csv` .

En el servidor, si todas las apuestas del *batch* fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB. 

Por su parte, el servidor deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.

### Ejercicio N°7:

Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo no se podrán responder consultas por la lista de ganadores con información parcial.

Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

No es correcto realizar un broadcast de todos los ganadores hacia todas las agencias, se espera que se informen los DNIs ganadores que correspondan a cada una de ellas.

## Parte 3: Repaso de Concurrencia
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

### Ejercicio N°8:

Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo. En caso de que el alumno implemente el servidor en Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).

## Condiciones de Entrega
Se espera que los alumnos realicen un _fork_ del presente repositorio para el desarrollo de los ejercicios y que aprovechen el esqueleto provisto tanto (o tan poco) como consideren necesario.

Cada ejercicio deberá resolverse en una rama independiente con nombres siguiendo el formato `ej${Nro de ejercicio}`. Se permite agregar commits en cualquier órden, así como crear una rama a partir de otra, pero al momento de la entrega deberán existir 8 ramas llamadas: ej1, ej2, ..., ej7, ej8.
 (hint: verificar listado de ramas y últimos commits con `git ls-remote`)

Se espera que se redacte una sección del README en donde se indique cómo ejecutar cada ejercicio y se detallen los aspectos más importantes de la solución provista, como ser el protocolo de comunicación implementado (Parte 2) y los mecanismos de sincronización utilizados (Parte 3).

Se proveen [pruebas automáticas](https://github.com/7574-sistemas-distribuidos/tp0-tests) de caja negra. Se exige que la resolución de los ejercicios pase tales pruebas, o en su defecto que las discrepancias sean justificadas y discutidas con los docentes antes del día de la entrega. El incumplimiento de las pruebas es condición de desaprobación, pero su cumplimiento no es suficiente para la aprobación. Respetar las entradas de log planteadas en los ejercicios, pues son las que se chequean en cada uno de los tests.

La corrección personal tendrá en cuenta la calidad del código entregado y casos de error posibles, se manifiesten o no durante la ejecución del trabajo práctico. Se pide a los alumnos leer atentamente y **tener en cuenta** los criterios de corrección informados  [en el campus](https://campusgrado.fi.uba.ar/mod/page/view.php?id=73393).

# Resolución

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
  Los archivos de cada agencia ```agency-${id}.csv``` deben estar en el directorio .data
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
  Los archivos de cada agencia ```agency-${id}.csv``` deben estar en el directorio .data

## Parte 3

### Ejercicio 8

En este ejercicio se modificó el servidor para poder procesar las consultas de los clientes concurrentemente, utilizando multithreading, implementado con la librearía de ```threading``` de Python.

Se crea un nuevo hilo para cada conexión que se establece con un cliente, y se ejecuta la función 

```python
__handle_client_connection()
```

Es así que los clientes se conectan y se permite que se envíen sus apuestas en paralelo.

Decidí usar multithreading más allá de las limitaciones de Python mencionadas del GIL (Global Interpreter Lock) porque considero que para esta implementacion alcanza con resolver la concurrencia con multithreading, ya que las operaciones con los clientes son principalmente de I/O (esperar apuestas y consultas de sorteo y esccribir resultados en archivos) y no de cómputo intensivo.

Considero que, teniendo en cuenta que para el TP son como máximo 5 clientes, el overhead del context switching no es tal como para jusitifcar usar multiprocessing, aunque también sería una solución posible.

Para evitar problemas de concurrencia como race conditions, implemente dos locks sobre las secciones críticas, donde se acceden a recursos compartidos:
  - _lottery_results_lock_: para asegurar que el conteo de los clientes que ya terminaron de mandar todas sus apuestas, y el cálculo de la lotería se haga de manera consistente, sin que se intente correr el sorteo por más de un cliente a la vez-
  - _bet_storig_lock_ : protege la operación del guardado de apuestas (```utils.store_bets```) para que no se corrompan los datos guardados.

Para cerrar las conexiones y liberar los recursos una vez finalizado el programa, el server se guarda una lista con los threads activos de los clientes, y cuando se realiza el shutdown, se hace un ```join``` de los client threads.

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

Los archivos de cada agencia ```agency-${id}.csv``` deben estar en el directorio .data