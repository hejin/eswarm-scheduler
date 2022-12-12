#!/bin/bash
pod=$1
echo "pods in x86-64 node ..."
kubectl describe node worker-x86 |grep -i $pod | awk {'print $2'}
echo ""
echo "pods in arm64  node ..."
kubectl describe node worker-arm |grep -i $pod | awk {'print $2'}
