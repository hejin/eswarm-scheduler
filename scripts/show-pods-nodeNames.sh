#!/bin/bash
app=$1
pods="`kubectl get pods|grep $app |awk {'print $1'}`"; 
for pod in $pods; do 
	kubectl describe pod $pod |grep -i node\:; 
done
