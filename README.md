# Tekton Tasks Operator

## Github
https://github.com/ksimon1/kubevirt-tekton-tasks-operator (WIP)

## Purpose
deploy kubevirt tekton tasks (https://github.com/kubevirt/kubevirt-tekton-tasks) and example pipelines (ideas here https://docs.google.com/document/d/1q7pDGyGLB7KHghtecQrm3ktU_QLdjt0DDXXOqbwIpmw/edit#heading=h.ky0d3evm3f99)

## Prerequisites
Tekton - https://github.com/tektoncd
KubeVirt - https://github.com/kubevirt/kubevirt
CDI - https://github.com/kubevirt/containerized-data-importer

## Operator framework
https://sdk.operatorframework.io/
This framework was chosen, because we have the most experiences with it, 
it is easily modified, written in golang and consume small amount of resources.

## Certification rotation
https://github.com/qinqon/kube-admission-webhook 

https://github.com/kubevirt/hyperconverged-cluster-operator/pull/1717 will not be currently implemented. We are planning to add support in future.

## What is missing: 
e2e tests, unit tests, single node support, prometheus implementation, csv generator

## QA involvement
operator is marked as tech preview, so currently product will be tested only by unit and e2e tests

## How will the operator be deployed
operator will be deployed as a part of HCO. Operator will not be deployed by default. Users will have to set a special bool attribute in HCO CR (e.g. deployTektonOperator). Deployment of operator will be possible even after installation of cnv. In case a user would like to use an operator, he/she has to install tekton pipelines / openshift pipelines into cluster. If tekton is not installed, the operator will be deployed and will be waiting for installation of tekton and there will be a message of required crds in the operator log. After tekton is installed, the operator will deploy tasks and pipelines.

## Installation
The `Tekton Tasks Operator` can run on both Kubernetes and OpenShift. However to be able to 
use all tasks, we recommend to use OpenShift (due to support for templates).

## Building

The Make will try to install kustomize, however if it is already installed it will not reinstall it.
In case of an error, make sure you are using at least v3 of kustomize, available here: https://kustomize.io/

To build the container image run:
```shell
make container-build
```

To upload the image to the default repository run:
```shell
make container-push
```
After the image is pushed to the repository,
manifests and the operator can be deployed using:
```shell
make deploy
```