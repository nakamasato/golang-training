# Cloud Run

## Connect Cloud SQL

Please check [Cloud SQL](../cloudsql/README.md), which covers:

1. Access with Cloud SQL BUILT_IN user (with Go Connector)
1. Access with Cloud SQL IAM_SERVICE_ACCOUNT user (with Go Connector)
1. Access with Cloud SQL BUILT_IN user via SQL Auth Proxy (Multi Container)
1. Access with Cloud SQL IAM_SERVICE_ACCOUNT user via SQL Auth Proxy (Multi Container)


## Ref
1. [Cloud Run now supports sidecar deployments — monitoring agents, proxies and more](https://cloud.google.com/blog/products/serverless/cloud-run-now-supports-multi-container-deployments)
1. [Cloud SQL Auth Proxy Sidecar](https://github.com/GoogleCloudPlatform/cloud-sql-proxy/tree/main/examples/multi-container/ruby)
1. https://cloud.google.com/run/docs/reference/yaml/v1
1. https://zenn.dev/google_cloud_jp/articles/cloud-run-multi-container-features
