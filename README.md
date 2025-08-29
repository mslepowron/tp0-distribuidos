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

---

### Ejercicio 5

Como planteo inicial se implentó un modulo _agency message_ para el cliente. Los datos de la apuesta se definen como variables de entorno, y se utilizan para construir un mensaje con el siguiente formato:

 ``` go
msg := fmt.Sprintf("%s|%s|%s|%s|%s|%s", bet.AgencyId, bet.Name, bet.LastName, bet.Document, bet.BirthDate, bet.Number) 
```

Por otro lado, en la funcion ```StartClientLoop()``` del modulo client, se crea el socket, se formatea el mensaje especificado arriba y se llama a una funcion ```SendClientMessage()``` que se encarga de enviarle el mensaje con los datos de la apuesta al servidor.

En cuanto al envío de los datos, en go, utilizar
``` go
conn.Write(message)
```
a secas no nos garantiza que se envíen todos los datos que queremos. Puede suceder que la cantidad de bytes enviados sea menor al largo total de los datos (short write), por ejemplo, porque el buffer del socket estaba lleno y solo acepto y envió una parte.

Para evitar que esto suceda, se tuvieron en cuenta dos cosas:

Por un lado, establecer un limite de tamaño máximo de mensaje a enviar (8KB). Esto evita saturar el buffer del socket, aporta algo de robustez y seguridad (por ejemplo evitando que in cliente malicioso mande mesajes excesivamente grandes que saturen el servidor o provoquen un DoS), y aporta simplicidad; mantener mensajes de tamaño moderado facilita el manejo de buffers y control de errores.

Por otro lado, se implementó una función ```WriteFull()``` donde se repite la operación de escritura en un loop hasta que efectivamente se hallan enviado todos los bytes del mensaje, evitando el envío de mensajes incompletos.

El mensaje con los datos de la apuesta es precedido por un mensaje de tamaño de 4 bytes que contiene la longitud del mensaje del cliente. Este paso se realiza para que el server sepa cuántos bytes tiene que leer y así evitar un short read.

Es decir, el protocolo esta conformado por 4 bytes (message length) + payload (client message)

El cliente espera un mensaje de confirmación de parte del servidor para verificar que haya recibio correctamente los datos de la apuesta. Este mismo contiene los datos que presenta el log del server cuando persiste una apuesta correctamente; es decir, el DNI y el Numero del cliente.

Si el cliente no recibe este mensaje de confirmación, retorna un error y logguea el problema.

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


Ideas refactorizacion:
  - que la serializacion del mensaje sea con un json (mas flexible, legible, no es necesario saber el orden de los campos, es una heraamienta comun entre todos los lenguajes etc)
  - que si falla el envio del client se reintente (3 veces por ej) para + tolerancia a fallos; excepto que supere la capacidad maxima (eso igual creo q se puede resolver en el 6 con los chunks)
  - uso de constantes para tamanos maximos, errores, etc.
  - chequeo de errores terminales (si se cerro la conexion abortar y termina) ?
  - chequear codigos de error correspondientes en cada caso
