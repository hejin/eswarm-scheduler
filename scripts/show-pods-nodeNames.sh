#!/bin/bash

pods="`kubectl get pods|grep busybox |awk {'print $1'}`"; 

for pod in $pods; do 
	kubectl describe pod $pod |grep -i node\:; 
done
