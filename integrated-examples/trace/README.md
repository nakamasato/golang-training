# Trace on GCP (Cloud Run)

![](diagram.drawio.svg)

## Prerequisite

```
export PROJECT=PROJECT_ID
export LOCATION_ID=asia-northeast1
export QUEUE_ID=helloworld
```

```
gcloud config set project $PROJECT
gcloud config set compute/region $LOCATION_ID
```

## Components

1. Go app `helloworld` ([ref](https://cloud.google.com/run/docs/quickstarts/build-and-deploy/deploy-go-service))
    1. `/helloworld`: just return `Hello World!`
    2. `/cloudtask`: Create a cloud task
1. PubSub `helloworld`: triggers Cloud Run service `helloworld`
1. Cloud Tasks: Cloud Run service enqueue and invoke `helloworld`

## Deployment

### Cloud Tasks

```
gcloud tasks queues create $QUEUE_ID --location $LOCATION_ID --project $PROJECT
```

### Cloud Run - server

build image

```
export KO_DOCKER_REPO=$LOCATION_ID-docker.pkg.dev/$PROJECT/cloud-run-source-deploy/helloworld
ko build --bare ./helloworld
```

Run the following command in `helloworld` directory:

Set [deterministic URL](https://cloud.google.com/run/docs/triggering/https-request#deterministic) of Cloud Run service:

```
PROJECT_NUMBER=$(gcloud projects describe $PROJECT --format="value(projectNumber)")
CLOUD_RUN_URL=https://helloworld-$PROJECT_NUMBER.$LOCATION_ID.run.app
```

```
gcloud run deploy helloworld --image $KO_DOCKER_REPO --set-env-vars=CLOUD_TASK_TARGET_URL=${CLOUD_RUN_URL}/helloworld,PROJECT_ID=$PROJECT,LOCATION_ID=$LOCATION_ID,QUEUE_ID=$QUEUE_ID --allow-unauthenticated --region $LOCATION_ID --project $PROJECT
```

Click on the url -> You'll `Hello, World`

```
curl $CLOUD_RUN_URL/helloworld
Hello World!
```

### PubSub

1. Create topic

    ```
    gcloud pubsub topics create helloworld --project $PROJECT
    ```

    ```
    gcloud pubsub topics list --project $PROJECT
    ```

1. Create subscription

    ```
    gcloud pubsub subscriptions create helloworld --topic helloworld --push-endpoint ${CLOUD_RUN_URL}/cloudtask --project $PROJECT
    ```

    ```
    gcloud pubsub subscriptions list --filter=topic=projects/$PROJECT/topics/helloworld --project $PROJECT
    ```

1. Publish message

    ```
    gcloud pubsub topics publish helloworld --message="helloworld" --project $PROJECT
    ```

    or with go (https://pkg.go.dev/cloud.google.com/go/pubsub)

    ```
    gcloud auth application-default login
    PROJECT_ID=$PROJECT go run pubsubpublisher/main.go
    ```

### Cloud Run Job - pubsubpublisher

```
export KO_DOCKER_REPO=$LOCATION_ID-docker.pkg.dev/$PROJECT/cloud-run-source-deploy/helloworld-pubsubpublisher
```

```
ko build --bare ./pubsubpublisher
```

```
gcloud run jobs deploy helloworld-pubsubpublisher --image $KO_DOCKER_REPO --set-env-vars=PROJECT_ID=$PROJECT,OTEL_SERVICE_NAME=helloworld-pubsubpublisher --region $LOCATION_ID --project $PROJECT
```

> [!NOTE]
> `OTEL_SERVICE_NAME` is used for `service.name` for OpenTelemetry

![](publisher-trace.png)


```
gcloud run jobs execute helloworld-pubsubpublisher --region $LOCATION_ID --project $PROJECT
```

## Cleanup

1. Cloud Run

    ```
    gcloud run services delete helloworld --region $LOCATION_ID
    ```

1. Cloud Task (cannot be recreated within 7 days)

    ```
    gcloud tasks queues delete $QUEUE_ID --location $LOCATION_ID --project $PROJECT
    ```

1. PubSub

    ```
    gcloud pubsub subscriptions delete helloworld
    gcloud pubsub topics delete helloworld
    ```

## Tips

1. OpenTelemetry + Cloud Trace: https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/blob/main/exporter/trace/README.md

## References

1. [Cloud Pub/Sub経由でトレースを取得する](https://zenn.dev/google_cloud_jp/articles/20230626-pubsub-trace) (ToDo)
1. [仕様と実装から理解するOpenTelemetryの全体像](https://zenn.dev/ymtdzzz/articles/37c2856f46ea10)
1. [OpenTelemetry で始める分散トレース](https://qiita.com/atsu_kg/items/c3ee8141e4638957a947)
1. [opentelemetry-cloud-run](https://github.com/GoogleCloudPlatform/opentelemetry-cloud-run) (ToDo)
1. [pubsub: extract trace information on push subscriptions #10828](https://github.com/googleapis/google-cloud-go/issues/10828)
    1. [feat(pubsub): allow trace extraction from protobuf message #10827](https://github.com/googleapis/google-cloud-go/pull/10827)
    1. https://github.com/einride/cloudrunner-go/pull/700
1. [Generate traces and metrics with Go](https://cloud.google.com/trace/docs/setup/go-ot)
1. [propagation/trace_context.go](https://github.com/open-telemetry/opentelemetry-go/blob/main/propagation/trace_context.go)
1. https://www.w3.org/TR/trace-context/
