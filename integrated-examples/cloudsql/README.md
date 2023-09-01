# Cloud Run with Cloud SQL

![](diagram.drawio.svg)

## 1. Prerequisite

```
PROJECT=xxxxx
```

```bash
gcloud auth login --update-adc
INSTANCE_NAME=naka-test
ZONE=asia-northeast1-b
REGION=asia-northeast1
SA_NAME=helloworld # for IAM_SERVICE_ACCOUNT user
CLOUDSQLUSER=helloworld # BUILT_IN user
DB_NAME_FOR_IAM_AUTH_USER=helloworld_auth
DB_NAME_FOR_BUILTIN_USER=helloworld
SECRET_NAME=cloudsqluser_pass_helloworld
CLOUDRUN_SERVICE_NAME_FOR_IAM_AUTH_USER=helloworld-auth
CLOUDRUN_SERVICE_NAME_FOR_BUILTIN_USER=helloworld
CLOUDRUN_SERVICE_NAME_FOR_IAM_AUTH_WITH_PROXY=helloworld-multi-container
IMAGE_NAME=helloworld # build image for Cloud Run with gcloud builds
```

## 2. Local run

1. Run
    1. Local run with local db:

        ```
        cd helloworld
        DB_HOST=localhost DB_USER=postgres DB_NAME=test_db DB_PASS=postgres go run main.go
        ```

    1. Local run with cloud sql (after completing the setup below)

        ```
        cd helloworld
        DB_HOST=localhost DB_USER=your@gmail.com DB_NAME=test_db DB_PASS=$(gcloud sql generate-login-token) go run main.go
        ```

1. Check
    ```
    curl localhost:8080/get
    ```

## 3. Set Up Cloud SQL (Instance, Database, User)

1. Create Cloud SQL instance

    ```
    gcloud sql instances create ${INSTANCE_NAME} \
    --database-version=POSTGRES_15 \
    --cpu=1 \
    --memory=3840MiB \
    --zone=$ZONE \
    --root-password=test \
    --database-flags=cloudsql.iam_authentication=on \
    --project ${PROJECT}
    ```

1. Create databases for iam user and built-in user

    ```
    gcloud sql databases create $DB_NAME_FOR_IAM_AUTH_USER --instance=${INSTANCE_NAME} --project ${PROJECT}
    gcloud sql databases create $DB_NAME_FOR_BUILTIN_USER --instance=${INSTANCE_NAME} --project ${PROJECT}
    ```

1. Create a Service Account for Cloud Run, which is also used for IAM database authentication for Cloud SQL

    ```
    gcloud iam service-accounts create ${SA_NAME} \
        --description="hello world cloud run service" \
        --display-name="helloworld" \
        --project=${PROJECT}
    ```

1. Create Cloud SQL user for the service account

    ```
    gcloud sql users create ${SA_NAME}@${PROJECT}.iam \
    --instance=${INSTANCE_NAME} \
    --type=cloud_iam_service_account --project ${PROJECT}
    ```
1. Grant roles to the service account

    ```
    gcloud projects add-iam-policy-binding ${PROJECT} \
        --member=serviceAccount:${SA_NAME}@${PROJECT}.iam.gserviceaccount.com \
        --role=roles/cloudsql.instanceUser
    ```
    ```
    gcloud projects add-iam-policy-binding ${PROJECT} \
        --member=serviceAccount:${SA_NAME}@${PROJECT}.iam.gserviceaccount.com \
        --role=roles/cloudsql.client
    ```

    1. [IAM add-iam-policy-binding](https://cloud.google.com/sdk/gcloud/reference/projects/add-iam-policy-binding)
    1. [Cloud SQL IAM Roles](https://cloud.google.com/sql/docs/postgres/iam-roles)

1. Create table `$DB_NAME_FOR_IAM_AUTH_USER.accounts` for `<sa_name>@<project>.iam`

    ```
    gcloud auth application-default login
    cloud-sql-proxy ${PROJECT}:${REGION}:${INSTANCE_NAME}
    ```

    Log in with root (password is set above)

    ```
    psql --host=localhost --username=postgres --dbname=$DB_NAME_FOR_IAM_AUTH_USER
    ```

    ```sql
    CREATE TABLE accounts (
        user_id serial PRIMARY KEY,
        username VARCHAR ( 50 ) UNIQUE NOT NULL,
        password VARCHAR ( 50 ) NOT NULL,
        email VARCHAR ( 255 ) UNIQUE NOT NULL,
        created_on TIMESTAMP NOT NULL,
        last_login TIMESTAMP
    );
    ```

    ```sql
    INSERT INTO accounts VALUES (1, 'uid', 'password', 'email@gmail.com', '2023-07-12 00:00:00', '2023-07-12 00:01:00');
    ```

    ```
    psql --host=localhost --username=postgres --dbname=$DB_NAME_FOR_IAM_AUTH_USER -c "alter table accounts owner to \"${SA_NAME}@${PROJECT}.iam\";"
    ```

1. Create Cloud SQL User `helloworld` (built-in user)
    ```
    CLOUDSQLUSER_PASS=$(openssl rand -base64 32)
    ```

    ```bash
    gcloud sql users create $CLOUDSQLUSER \
        --instance=${INSTANCE_NAME} \
        --password=$CLOUDSQLUSER_PASS \
        --type=BUILT_IN \
        --project ${PROJECT}
    ```

    `--password` is necessary for built-in user

    ```
    gcloud sql users list --instance ${INSTANCE_NAME} --project ${PROJECT}
    ```

1. Create table `$DB_NAME_FOR_BUILTIN_USER.accounts` for `helloworld`

    ```
    cloud-sql-proxy ${PROJECT}:${REGION}:${INSTANCE_NAME}
    ```

    ```
    PGPASSWORD=$CLOUDSQLUSER_PASS psql --host=localhost --username=$CLOUDSQLUSER --dbname=$DB_NAME_FOR_BUILTIN_USER
    ```

    ```sql
    CREATE TABLE accounts (
        user_id serial PRIMARY KEY,
        username VARCHAR ( 50 ) UNIQUE NOT NULL,
        password VARCHAR ( 50 ) NOT NULL,
        email VARCHAR ( 255 ) UNIQUE NOT NULL,
        created_on TIMESTAMP NOT NULL,
        last_login TIMESTAMP
    );
    ```

    ```sql
    INSERT INTO accounts VALUES (1, 'uid', 'password', 'email@gmail.com', '2023-07-12 00:00:00', '2023-07-12 00:01:00');
    ```

1. Grant permission (ToDo)
    ```
    cloud-sql-proxy ${PROJECT}:${REGION}:${INSTANCE_NAME}
    ```

    ```
    psql --host=localhost --username=postgres --dbname=postgres
    ```

    ```sql
    GRANT cloudsqlsuperuser TO "helloworld@<project_id>.iam";
    ```

## 4. Deploy Cloud Run

### 4.1. Access with IAM database authentication (Service Account)

1. Deploy

    ```
    gcloud run deploy ${CLOUDRUN_SERVICE_NAME_FOR_IAM_AUTH_USER} --service-account ${SA_NAME}@${PROJECT}.iam.gserviceaccount.com \
        --source ./helloworld \
        --set-env-vars CLOUD_SQL_WITH_IAM_AUTH=true \
        --set-env-vars INSTANCE_CONNECTION_NAME=${PROJECT}:${REGION}:${INSTANCE_NAME} \
        --set-env-vars DB_NAME=$DB_NAME_FOR_IAM_AUTH_USER \
        --set-env-vars DB_IAM_USER=${SA_NAME}@${PROJECT}.iam \
        --project ${PROJECT} \
        --region ${REGION} \
        --async
    ```

1. Get URL

    ```
    URL=$(gcloud run services describe ${CLOUDRUN_SERVICE_NAME_FOR_IAM_AUTH_USER} --project $PROJECT --region $REGION --format json | jq -r .status.url); echo $URL
    ```

1. Check

    ```
    curl -H "authorization: bearer $(gcloud auth print-identity-token --project $PROJECT)" $URL
    Hello World!
    ```

    ```
    curl -H "authorization: bearer $(gcloud auth print-identity-token --project $PROJECT)" $URL/get
    Getting!
    Got accounts: [uid:1, username: uid]
    ```
### 4.2. Access with built-in auth (Postgres user)

1. Create secret `$SECRET_NAME`

    ```
    gcloud secrets create $SECRET_NAME \
        --replication-policy="automatic" --project ${PROJECT}
    echo -n "$CLOUDSQLUSER_PASS" | \
        gcloud secrets versions add $SECRET_NAME --data-file=- --project ${PROJECT}
    ```

    - [create secret](https://cloud.google.com/secret-manager/docs/creating-and-accessing-secrets#secretmanager-create-secret-gcloud)
    - [add secret version](https://cloud.google.com/secret-manager/docs/add-secret-version)

1. Grant roles to the service account

    ```
    gcloud secrets add-iam-policy-binding $SECRET_NAME \
        --member="serviceAccount:${SA_NAME}@${PROJECT}.iam.gserviceaccount.com" \
        --role="roles/secretmanager.secretAccessor" --project ${PROJECT}
    ```

1. Deploy
    ```
    gcloud run deploy ${CLOUDRUN_SERVICE_NAME_FOR_BUILTIN_USER} --service-account ${SA_NAME}@${PROJECT}.iam.gserviceaccount.com \
        --source ./helloworld \
        --set-env-vars CLOUD_SQL_WITH_BUILT_IN_USER=true \
        --set-env-vars INSTANCE_CONNECTION_NAME=${PROJECT}:${REGION}:${INSTANCE_NAME} \
        --set-env-vars DB_NAME=$DB_NAME_FOR_BUILTIN_USER \
        --set-env-vars DB_USER=$CLOUDSQLUSER \
        --set-secrets DB_PASS=${SECRET_NAME}:latest \
        --project ${PROJECT} \
        --region ${REGION} \
        --async
    ```

    - [connect from cloud run](https://cloud.google.com/sql/docs/postgres/connect-run)

1. Get URL
    ```
    URL=$(gcloud run services describe ${CLOUDRUN_SERVICE_NAME_FOR_BUILTIN_USER} --project $PROJECT --region $REGION --format json | jq -r .status.url); echo $URL
    ```

1. Check

    ```
    curl -H "authorization: bearer $(gcloud auth print-identity-token --project $PROJECT)" $URL
    Hello World!
    ```

    ```
    curl -H "authorization: bearer $(gcloud auth print-identity-token --project $PROJECT)" $URL/get
    Getting!
    Got accounts: [uid:1, username: uid]
    ```

### 4.3. With SQL Auth Proxy [WIP]

[About the Cloud SQL Auth Proxy](https://cloud.google.com/sql/docs/postgres/sql-proxy)

1. [Build an image with Google Cloud's buildpacks](https://cloud.google.com/run/docs/building/containers#buildpacks)

    ```
    cd helloworld
    gcloud builds submit --pack image=gcr.io/${PROJECT}/${IMAGE_NAME} --project $PROJ
    ECT
    ```

1. Deploy with yaml

    ```
    cat << EOF > $CLOUDRUN_SERVICE_NAME_FOR_IAM_AUTH_WITH_PROXY.yaml
    apiVersion: serving.knative.dev/v1
    kind: Service
    metadata:
      annotations:
         run.googleapis.com/launch-stage: ALPHA
      name: $CLOUDRUN_SERVICE_NAME_FOR_IAM_AUTH_WITH_PROXY
    spec:
      template:
        metadata:
          annotations:
            run.googleapis.com/execution-environment: gen2

        spec:
          serviceAccountName: ${SA_NAME}@${PROJECT}.iam.gserviceaccount.com
          containers:
          - name: helloworld
            image: gcr.io/$PROJECT/$IMAGE_NAME
            ports:
              - containerPort: 8080
            env:
              - name: CLOUD_SQL_WITH_IAM_AUTH
                value: "true"
              - name: DB_IAM_USER
                value: ${SA_NAME}@${PROJECT}.iam
              - name: DB_NAME
                value: $DB_NAME_FOR_IAM_AUTH_USER
              - name: INSTANCE_CONNECTION_NAME
                value: 127.0.0.1
              - name: DB_PORT
                value: "5432"
          - name: cloud-sql-proxy
            image: gcr.io/cloud-sql-connectors/cloud-sql-proxy:latest
            args:
                 # Ensure the port number on the --port argument matches the value of
                 # the DB_PORT env var on the my-app container.
                 - "--port=5432"
                 - "--auto-iam-authn"
                 - "${PROJECT}:${REGION}:${INSTANCE_NAME}"
    EOF
    ```

    ```
    gcloud run services replace $CLOUDRUN_SERVICE_NAME_FOR_IAM_AUTH_WITH_PROXY.yaml --project $PROJECT --region $REGION
    ```

1. Get URL
    ```
    URL=$(gcloud run services describe ${CLOUDRUN_SERVICE_NAME_FOR_IAM_AUTH_WITH_PROXY} --project $PROJECT --region $REGION --format json | jq -r .status.url); echo $URL
    ```

1. Check

    ```
    curl -H "authorization: bearer $(gcloud auth print-identity-token --project $PROJECT)" $URL
    Hello World!
    ```

    ```
    curl -H "authorization: bearer $(gcloud auth print-identity-token --project $PROJECT)" $URL/get
    Getting!
    ```

    **Error:** -> `main.go` で IAM_AUTHで接続する時にINSTANCE_NAMEで接続する選択肢しかないから、sql auth proxyで接続する方法を探す
    ```
    "getRows failed failed to connect to `host=/tmp user=xxx@xxxx.iam database=xxxx`: dial error (Config error: invalid instance connection name, expected PROJECT:REGION:INSTANCE (connection name = "127.0.0.1"))"
    ```

## 5. Connect to Cloud Run with `psql` from Local

### 5.1. With `cloud-sql-proxy`

```
gcloud auth application-default login
```

#### 5.1.1. Using IAM authentication (Service Account)

1. With service account impersonation, you need to have `roles/iam.serviceAccountTokenCreator` ([required roles](https://cloud.google.com/docs/authentication/use-service-account-impersonation#required-roles))

    ```
    gcloud projects add-iam-policy-binding $PROJECT \
        --member=user:<your>@gmail.com --role=roles/iam.serviceAccountTokenCreator \
        --condition=None
    ```
1. Run cloud-sql-proxy
    ```
    cloud-sql-proxy ${PROJECT}:${REGION}:${INSTANCE_NAME} --token $(gcloud auth print-access-token) --auto-iam-authn --login-token $(gcloud sql generate-login-token --impersonate-service-account=${SA_NAME}@${PROJECT}.iam.gserviceaccount.com)
    ```

    - `--token`: Use bearer token as a source of IAM credentials.
    - `--login-token`: Use bearer token as a database password (used with token and auto-iam-authn only)
    - `--auto-iam-authn`: (*) Enables Automatic IAM Authentication for all instances

1. Connect via cloud-sql-proxy without password
    ```
    psql "host=localhost user=${SA_NAME}@${PROJECT}.iam dbname=${DB_NAME_FOR_IAM_AUTH_USER}"
    ```

#### 5.1.2. Using built-in user (Username & Password)

1. Run cloud-sql-proxy
    ```
    cloud-sql-proxy ${PROJECT}:${REGION}:${INSTANCE_NAME}
    ```

1. Connect via cloud-sql-proxy with password

    ```
    PGPASSWORD=$CLOUDSQLUSER_PASS psql --host=localhost --username=$CLOUDSQLUSER --dbname=$DB_NAME_FOR_BUILTIN_USE
    ```
### 5.2. With public IP

```
DB_HOST=$(gcloud sql instances describe ${INSTANCE_NAME} --project ${PROJECT} --format json | jq -r '.ipAddresses[] | select(.type == "PRIMARY").ipAddress')
```

```
PGPASSWORD=$CLOUDSQLUSER_PASS psql --host=$DB_HOST --username=$CLOUDSQLUSER --dbname=$DB_NAME_FOR_BUILTIN_USER
```

If [Cloud SQL organization policies](https://cloud.google.com/sql/docs/postgres/org-policy/org-policy) is set, you cannot use this way.

### 5.3. With `gcloud sql connect` command (Public IP)

```
gcloud sql connect ${INSTANCE_NAME} --user=postgres --quiet --project ${PROJECT}
```

> The gcloud sql connect command does not support connecting to a Cloud SQL instance using private IP, or using SSL/TLS.

> This command isn't supported for Cloud SQL instances with only private IP addresses.

> NOTE: If you're connecting from an IPv6 address, or are constrained by certain organization policies (restrictPublicIP, restrictAuthorizedNetworks), consider running the beta version of this command to avoid error by connecting through the Cloud SQL proxy: gcloud beta sql connect

For more about `gcloud sql connect`, please read [gcloud sql connect](https://cloud.google.com/sdk/gcloud/reference/sql/connect)

You can also use

## Tips

```
gcloud sql users set-password postgres \
    --instance=${INSTANCE_NAME} \
    --project=${PROJECT} \
    --prompt-for-password
```

```
gcloud sql connect ${INSTANCE_NAME} --user=postgres --quiet --project ${PROJECT}
```

## Clean up

```
gcloud run services delete ${CLOUDRUN_SERVICE_NAME_FOR_BUILTIN_USER} --project ${PROJECT} --region ${REGION}
gcloud run services delete ${CLOUDRUN_SERVICE_NAME_FOR_IAM_AUTH_USER} --project ${PROJECT} --region ${REGION}
gcloud run services delete ${CLOUDRUN_SERVICE_NAME_FOR_IAM_AUTH_WITH_PROXY} --project ${PROJECT} --region ${REGION}
gcloud sql instances delete ${INSTANCE_NAME} --project ${PROJECT}
gcloud iam service-accounts delete ${SA_NAME}@${PROJECT}.iam.gserviceaccount.com --project ${PROJECT}
gcloud secrets delete $SECRET_NAME --project ${PROJECT}
```

## Ref

1. Connect Cloud SQL from Cloud Run
    1. [Connect from Cloud Run](https://cloud.google.com/sql/docs/postgres/connect-run)
    1. [Connect to Cloud SQL for PostgreSQL from Cloud Run](https://cloud.google.com/sql/docs/postgres/connect-instance-cloud-run)
    1. [Cloud Run now supports sidecar deployments — monitoring agents, proxies and more](https://cloud.google.com/blog/products/serverless/cloud-run-now-supports-multi-container-deployments)
1. [create secret](https://cloud.google.com/secret-manager/docs/creating-and-accessing-secrets#secretmanager-create-secret-gcloud)
1. [add secret version](https://cloud.google.com/secret-manager/docs/add-secret-version)
1. [IAM add-iam-policy-binding](https://cloud.google.com/sdk/gcloud/reference/projects/add-iam-policy-binding)
1. [Cloud SQL IAM Roles](https://cloud.google.com/sql/docs/postgres/iam-roles)
1. [Cloud SQL organization policies](https://cloud.google.com/sql/docs/postgres/org-policy/org-policy): `constraints/sql.restrictPublicIp`
1. [About the Cloud SQL Auth Proxy](https://cloud.google.com/sql/docs/postgres/sql-proxy)
