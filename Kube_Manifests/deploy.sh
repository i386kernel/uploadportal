#!/bin/bash

echo "creating minio namespace"
kubectl create namespace minio

echo "changing context"
kubectl config set-context --current --namespace=minio

echo "deploying minio"
kubectl apply -f minio

sleep 10s

echo "creating application namespace"
kubectl create namespace referarchapp
kubectl config set-context --current --namespace=referarchapp

echo "Deploying MySQL"
kubectl apply -f mysql

sleep 10s

echo "Creating RefArchPortal"
kubectl apply -f referarchapp

