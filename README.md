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