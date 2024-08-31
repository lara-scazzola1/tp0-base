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

Para probar que verdaderamente funcionaba se agrego el comando docker-compose-start en el makefile. Se probo cambiando el puerto del servidor y se inicio el docker-compose con el comando make docker-compose-start y se pudo ver que no se lograron conectar, por lo que cambio la configuracion sin necesidad de hacer el build.

