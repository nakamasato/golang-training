# control skaffold

Execute `skaffold dev` programmatically.

## Prerequisite

- [kind](https://kind.sigs.k8s.io/)
- [skaffold](https://skaffold.dev/)

## Steps

1. Start kind cluster.

    ```
    kind create cluster --name test --image kindest/node:v1.20.2
    ```

    check:

    ```
    kubectl cluster-info --context kind-test
    ```

1. Run skaffold with the script.

    ```
    go run main.go
    ```

    1. Deploy a Pod with name `test` with `skaffold run`.
    1. Get Pod list in `default` namespace with `client-go`.
    1. Delete the deployed resources with `skaffold delete`.
