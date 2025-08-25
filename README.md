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