#!/bin/bash

echo "creating application namespace"

kubectl create namespace ra-portal
kubectl config set-context --current --namespace=ra-portal

echo "deploying minio"
kubectl apply -f minio

echo "Deploying MySQL"
kubectl apply -f mysql

sleep 20s

echo "Creating RefArchPortal"
kubectl apply -f ra-portal

echo "Bringing up app"

sleep 20s

kubectl get service
