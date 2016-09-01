#Crea el docker para memcached
#docker run --detach --name memcached1 memcached
#docker run --detach --name memcached1 --publish=11211:11211 memcached 
docker start memcached1

#Corre el rabbit, usa una imagen que tiene la consola de operaciones
#docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management
docker start some-rabbit

#Compila el código del webserver y lo mete en una imagen
docker build -t my-golang-app .

#Crea un docker con la imagen recientemente creada y lo linkea al memcached1
docker run --publish=8080:8080 --link memcached1 --link some-rabbit -it --rm --name my-running-app my-golang-app
