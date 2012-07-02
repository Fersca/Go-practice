#!/bin/bash
c=1
while [ $c -le 10000 ]
do
	echo "Welcone $c times"
	curl -v -H "Content-Type: application/json" -X POST -d '{"name":"fer"}' http://localhost:8080/$c
	(( c++ ))
done
echo ""
echo "Fin cache test"
