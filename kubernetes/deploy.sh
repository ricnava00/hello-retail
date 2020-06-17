#! /bin/bash

set -e -o pipefail

echo 'creating namespaces'
kubectl apply -f ./namespaces.yaml

PASSWORD=$(head -c 12 /dev/urandom | shasum | cut -d' ' -f1)
kubectl -n openfaas create secret generic basic-auth \
	--from-literal=basic-auth-user=admin \
	--from-literal=basic-auth-password="$PASSWORD"

echo 'deploy platform'
kubectl apply -f ./openfaas.yaml

echo 'waiting for platform to come up...'
sleep 45

export OPENFAAS_URL=http://127.0.0.1:31112

echo -n $PASSWORD | faas-cli login --password-stdin

echo 'deploying sql database'
kubectl apply -f ./mysql.yaml

echo 'waiting for mysql to come up...'
sleep 30

echo 'initializing up database'
kubectl run --namespace=openfaas-fn --restart=Never --image=mysql:5.6 mysql-client -- mysql -h mysql -ppass -e "CREATE DATABASE IF NOT EXISTS helloRetail"
sleep 5
kubectl delete pod mysql-client --namespace=openfaas-fn

echo 'deploying hello-retail'
faas-cli deploy -f ./hello-retail.yaml

echo 'waiting for hello-retail to come up...'
sleep 30

echo 'testing hello-retail workflows'
curl -XPOST http://127.0.0.1:31112/function/product-catalog-builder/product -H 'Content-Type: application/json' -H 'Authorization: inventory-manager' --data '{"schema": "com.nordstrom/retail-stream/1-0-0", "origin": "hello-retail/product-producer-automation", "timeOrigin": "2017-01-12T18:29:25.171Z", "price: ": "123", "data": {"schema": "com.nordstrom/product/create/1-0-0", "id": "4579874", "brand": "POLO RALPH LAUREN", "name": "Polo Ralph Lauren 3-Pack Socks", "description": "PAGE:/s/polo-ralph-lauren-3-pack-socks/4579874", "category": "Socks for Men", "price": "123"}}'
curl -XPOST http://127.0.0.1:31112/function/product-catalog-api/products -H 'Content-Type: application/json' --data '{"queryStringParameters":{"category": "Socks for Men"}}'
curl -XPOST http://127.0.0.1:31112/function/product-purchase -H 'Content-Type: application/json' -H 'Authorization: customer' --data '{"id":"4579874", "creditCard": "1234-5678-9", "user": "prabkumar"}'
