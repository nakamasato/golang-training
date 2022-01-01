# k8sclient

## Prerequisite

## Steps

1. Prepare k8s cluster with `kind`.
    ```
    kind create cluster --name test --kubeconfig kubeconfig
    ```
1. Run `main.go` (get pod).
    ```
    go run main.go
    test-pod pod doesn't exist in default namespace
    kube-controller-manager-test-control-plane pod exists in kube-system namespace
    ```
1. Delete kind cluster.
    ```
    kind delete clusters test
    ```
