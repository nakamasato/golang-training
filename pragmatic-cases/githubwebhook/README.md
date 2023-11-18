# GitHub Webhook handler on Cloud Run

## Setup

1. initialize module

    ```
    go mod init github.com/nakamasato/golang-training/pragmatic-cases/githubwebhook
    ```

1. use workspace

    ```
    go work use .
    ```

1. add main.go


## Deploy

1. Prepare GCP
    ```bash
    PROJECT=<project>
    REGION=<region>
    ```
    ```
    gcloud auth login
    ```
1. Generate a secret key
    ```bash
    SECRET_KEY=$(openssl rand -base64 32)
    echo $SECRET_KEY
    ```

1. Create SecretManager secret
    ```bash
    gcloud secrets create GITHUB_WEBHOOK_SECRET \
    --replication-policy automatic \
    --data-file <(echo -n $SECRET_KEY) \
    --project $PROJECT
    ```
1. Create ServiceAccount
    ```bash
    gcloud iam service-accounts create githubwebhook \
    --project $PROJECT
    ```
1. Grant roles
    ```bash
    gcloud secrets add-iam-policy-binding GITHUB_WEBHOOK_SECRET \
    --member="serviceAccount:githubwebhook@${PROJECT}.iam.gserviceaccount.com" \
    --role="roles/secretmanager.secretAccessor" --project ${PROJECT}
    ```
1. Deploy to Cloud Run

    ```bash
    gcloud run deploy githubwebhook \
    --allow-unauthenticated \
    --source . \
    --service-account githubwebhook@${PROJECT}.iam.gserviceaccount.com \
    --set-secrets GITHUB_WEBHOOK_SECRET=GITHUB_WEBHOOK_SECRET:latest \
    --project $PROJECT \
    --region $REGION
    ```
1. [Create Webhook](https://docs.github.com/en/webhooks/using-webhooks/creating-webhooks)
    1. `Payload URL`: Cloud Run URL e.g. https://githubwebhook-xxxx-an.a.run.app/webhooks
    1. `Content type`: `application/json`
    1. `Which events would you like to trigger this webhook?`:
        1. `Check runs`
        1. `Pull requests`
        1. `Statuses`
    1. `Secret`: `$SECRET_KEY` generated above

## Clean up

```bash
gcloud run services delete githubwebhook --project $PROJECT --region $REGION
gcloud secrets delete GITHUB_WEBHOOK_SECRET --project $PROJECT
gcloud iam service-accounts delete githubwebhook@${PROJECT}.iam.gserviceaccount.com --project $PROJECT
```
