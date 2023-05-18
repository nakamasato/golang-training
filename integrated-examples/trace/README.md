# Trace on GCP (Cloud Run)

![](diagram.drawio.svg)

## Prerequisite

```
export PROJECT_ID=PROJECT_ID
export LOCATION_ID=asia-northeast1
export QUEUE_ID=helloworld
```

```
gcloud config set project $PROJECT_ID
gcloud config set compute/region $LOCATION_ID
```

## Components

1. Go app `helloworld` ([ref](https://cloud.google.com/run/docs/quickstarts/build-and-deploy/deploy-go-service))
    ```
    CLOUDTASK_TARGET_URL="$CLOUD_RUN_URL" go run main.go
    ```

    ```
    curl localhost:8080/helloworld
    Hello World!
    ```

1. PubSub `helloworld`: triggers Cloud Run service `helloworld`
1. Cloud Tasks: Cloud Run service enqueue and invoke `helloworld`

## Deployment

### Cloud Run

Run the following command in `helloworld` directory:

first time:

```
gcloud run deploy --source ./helloworld --set-env-vars CLOUD_TASK_TARGET_URL=example.com
```

```
CLOUD_RUN_URL=$(gcloud run services describe helloworld --region $LOCATION_ID --format json | jq -r '.status.url')
echo $CLOUD_RUN_URL
```

```
gcloud run deploy --source ./helloworld --set-env-vars=CLOUD_TASK_TARGET_URL=${CLOUD_RUN_URL}/helloworld,PROJECT_ID=$PROJECT_ID,LOCATION_ID=$LOCATION_ID,QUEUE_ID=$QUEUE_ID --region $LOCATION_ID
```

Click on the url: https://helloworld-xxxxx-an.a.run.app/helloworld -> You'll `Hello, World`

### PubSub

1. Create topic

    ```
    gcloud pubsub topics create helloworld
    ```

1. Create subscription

    ```
    gcloud pubsub subscriptions create helloworld --topic helloworld --push-endpoint ${CLOUD_RUN_URL}/cloudtask
    ```

1. Publish message

    ```
    gcloud pubsub topics publish helloworld --message="helloworld"
    ```

    or with go (https://pkg.go.dev/cloud.google.com/go/pubsub)

    ```
    gcloud auth application-default login
    PROJECT_ID=$PROJECT_ID go run pubsubpublisher/main.go
    ```

### Cloud Task

```
gcloud tasks queues create $QUEUE_ID --location $LOCATION_ID
```

## Cleanup

1. Cloud Run

    ```
    gcloud run services delete helloworld --region $LOCATION_ID
    ```

1. Cloud Task (cannot be recreated within 7 days)

    ```
    gcloud tasks queues delete $QUEUE_ID --location $LOCATION_ID
    ```

1. PubSub

    ```
    gcloud pubsub subscriptions delete helloworld
    gcloud pubsub topics delete helloworld
    ```
