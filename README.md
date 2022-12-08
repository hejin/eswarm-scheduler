            how to demonstrate the sample scheduler

0. prepare the environment
  - x86-64 ubuntu 22.04 LTS jammy
  - minikube 1.25.2 (k8s 1.23.3)
  - kubectl 1.23.14 (standalone instead of the minikube kubectl)
  - if you are in China
      - setup your own aliyuncs docker image registry, and create the namespace for the project 
      - if possible, make the registry public to save the work to setup imagePullSecrets : check the deploy/

1. build it (if you finished setup your own image registry)
  - replace the image name with your own
  - make push

2. start the minikube 1.25.2 with single node
  - if you are in China, launch it with '--image-mirror-country' and '--image-repository' parameters

3. after minikube up, deploy and check the eswarm scheduler
  - kubectl apply -f deploy/
  - kubectl get deployments -n kube-system |grep eswarm

4. after eswarm scheduler ready, check if it really works with example pods (busybox)
  - kubectl apply -f example/

5. run 'kubectl get pods' to check if pod busybox running

