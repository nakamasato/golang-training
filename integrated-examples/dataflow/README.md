# Dataflow

## Batch

https://cloud.google.com/dataflow/docs/quickstarts/create-pipeline-go


### Run on Local

```
go run integrated-examples/dataflow/wordcount/main.go --output output.txt
```

### Run on Dataflow

```
PROJECT=
REGION=asia-northeast1
BUCKET_NAME=
```

```
gcloud auth application-default login
```

```
go run integrated-examples/dataflow/wordcount/main.go --input gs://dataflow-samples/shakespeare/kinglear.txt \
    --output gs://$BUCKET_NAME/results/outputs \
    --runner dataflow \
    --project $PROJECT \
    --region $REGION \
    --staging_location gs://$BUCKET_NAME/binaries/
```

Output: `gs://$BUCKET_NAME/results/outputs`


```
gcloud storage cp gs://$BUCKET_NAME/results/outputs - | head
feature: 1
block: 1
Cried: 1
scatter'd: 1
she: 44
sudden: 1
silly: 1
More: 6
out: 68
believe: 3
```
