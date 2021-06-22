# Introduction
This is sample operator which will list the type of nodes based on role 
and take some actions

## Install

### Prerequisites
Ensure KUBECONFIG is pointing to a working Kubernetes cluster

### Deploy the Operator
```
$ git clone https://github.com/bpradipt/check-node-role.git
$ cd check-node-role
$ make install && make deploy IMG=quay.io/bpradipt/check-node-role
```
For deploying on Power (`ppc64le`) use the following command
```
$ make install && make deploy IMG=quay.io/bpradipt/check-node-role:ppc64le
```

This will deploy the controller POD in the `check-node-role-system`
namespace

```
$ kubectl get pods -n check-node-role-system

NAME                                                        READY   STATUS    RESTARTS   AGE
check-node-role-controller-manager-5c79c4b48-cwq97          2/2     Running   0          19m
```

### Create a Custom Resource
```
$ kubectl create -f config/samples/nodeattr_v1alpha1_noderole.yaml
```
This will update the CR status with list of nodetypes 

```
$ kubectl describe noderoles.nodeattr.power.io noderole-sample

[snip]
Spec:
  Controller:  true
  Infra:       false
  Worker:      true
Status:
  Controller Nodes:
    kube1-00
  Infra Nodes:
    kube1-00
  Worker Nodes:
    kube1-00
Events:  <none>
```


## Uninstall
```
$ make uninstall
```

## Hacking

### Build and Push the container image

Ensure you have access to a container registry like quay.io or hub.docker.com
```
export REGISTRY=<registry>
export REGISTRY_USER=<user>
```

```
$ make docker build IMG=${REGISTRY}/${REGISTRY_USER}/check-node-role
$ make docker-push IMG=${REGISTRY}/${REGISTRY_USER}/check-node-role
```

### Deploy the new image
```
$ make install && make deploy IMG=${REGISTRY}/${REGISTRY_USER}/check-node-role
```

