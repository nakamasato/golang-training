# control skaffold

Execute `skaffold dev` programmatically.

## Prerequisite

- [kind](https://kind.sigs.k8s.io/)
- [skaffold](https://skaffold.dev/)
    [install](https://skaffold.dev/docs/install/)
    Mac:
    ```
    curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-darwin-amd64 && \
    sudo install skaffold /usr/local/bin/
    ```

## Steps

1. Start kind cluster.

    ```
    kind create cluster
    ```

    check:

    ```
    kubectl cluster-info --context kind-kind
    ```

1. Run skaffold with the script.

    ```
    go run main.go
    ```

    1. Deploy a Pod with name `test` with `skaffold run`.
    1. Get Pod list in `default` namespace with `client-go`.
    1. Delete the deployed resources with `skaffold delete`.


## Example

https://github.com/GoogleContainerTools/skaffold/blob/a00ef25db88e310ae5a67409f54a6290688b2726/integration/run_test.go#L192-L197

```go
			skaffold.Run(args...).InDir(test.dir).InNs(ns.Name).WithEnv(test.env).RunOrFail(t)

			client.WaitForPodsReady(test.pods...)
			client.WaitForDeploymentsToStabilize(test.deployments...)

			skaffold.Delete().InDir(test.dir).InNs(ns.Name).WithEnv(test.env).RunOrFail(t)
```
