# TP0 - Sistemas distribuidos

## Ejercicio 1

Para dar permisos para ejecutar el script de generación de docker-compose.yml, se debe ejecutar el siguiente comando:

```bash
chmod +x generar-compose.sh
```

Luego, para ejecutar el script, se debe ejecutar el siguiente comando:

```bash
./generar-compose.sh <nombre del archivo de salida> <cantidad de clientes>
```

Si nombre del archivo de salida es "docker-compose-dev.yaml" al correr los siguentes comandos se va a levantar el docker compose generado con el script y se van a ver los logs del servidor y los clientes que se hayan levantado:

```bash
make docker-compose-up
make docker-compose-logs
```

## Ejercicio 2

Se montaron los archivos config.ini para el servidor y config.yaml para el cliente en volumenes para que puedan ser modificados sin tener que volver a construir la imagen.

Para probar que verdaderamente funcionaba se agrego el comando `docker-compose-start` en el makefile. Se probo cambiando el puerto del servidor y se inicio el docker-compose con el comando make docker-compose-start y se pudo ver que no se lograron conectar, por lo que cambio la configuracion sin necesidad de hacer el build.

## Ejercicio 3

Antes de ejecutar el script de validacion se debe correr el docker-compose con el comando 

```bash
make docker-compose-up
```

Para dar permisos para ejecutar el script de generación de docker-compose.yml, se debe ejecutar el siguiente comando:

```bash
chmod +x validar-echo-server.sh
```

Luego, para ejecutar el script, se debe ejecutar el siguiente comando:

```bash
./validar-echo-server.sh
```

El script levanta un container de doker que se conecta a la red que se configuro en el docker compose (en la que se lanzo el servidor y el cliente). Luego se una netcat para enviarle un mensaje al servidor, si el servidor responde con el mismo mensaje que se envio, el script imprime por pantalla `action: test_echo_server | result: success`, de lo contrario `imprimir:action: test_echo_server | result: fail`.

## Ejercicio 4

### Cliente

Se configuro un channel para recibir señales del sistema (os.Signal). En el metodo StartClientLoop, se asegura que la conexión con el servidor se cierre correctamente cuando se recibe una señal, simplemente si se recibe la señal de exit se sale del bucle y se cierra todo antes de irse de scope. Esto se realiza en el bloque select, que escucha tanto la señal de salida como el flujo normal del programa.

### Servidor

Se configuro el manejo de la señal SIGTERM utilizando el módulo signal de Python. El metodo stop_server se ha registrado como el manejador para la señal SIGTERM. Cuando el servidor recibe esta señal, se cierra el socket del servidor y se marca la variable _stop como True, lo que detiene el bucle principal del servidor.


Para enviar la señal de terminacion se uso el comando `docker stop` sobre el container que se queria detener (client1 o server). 
```bash
docker stop <container_name>
```
Se pudo ver que el servidor se detenia correctamente y el cliente tambien.

## Ejercicio 5

Se crea la clase socket tanto en el cliente como en el servidor. Ambas clases tienen los metodos sendall y recvall para evitar short reads y short writes.

Se crea la clase protocolo en ambas partes para manejar el envio y recepcion de las apuestas.

Se crea la clase Bet en el cliente, que implementa una funcion para serializarse ella misma. En el caso del servidor se implemento la funcion para deserializar la apuesta en el archivo de utils.

