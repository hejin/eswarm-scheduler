#!/bin/bash

kubectl label node minikube-m02 node-role.kubernetes.io/worker=
kubectl label node minikube-m02 power-core=Y

kubectl label node minikube-m03 node-role.kubernetes.io/worker=
kubectl label node minikube-m03 power-core=N
