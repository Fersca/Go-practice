#Crea el docker para memcached
#docker run --detach --name memcached1 memcached 
docker start memcached1

#Esto es lo que tenía antes el memcached, publicaba puertos, ahora no hace falta porque está linkeado
#docker run --detach --name memcached1 --publish=11211:11211 memcached

#Compila el código del webserver y lo mete en una imagen
#docker build -t my-golang-app .

#Crea un docker con la imagen recientemente creada y lo linkea al memcached1
docker run --publish=8080:8080 --link memcached1 -it --rm --name my-running-app my-golang-app
