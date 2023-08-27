# Google Drive

1. Set GCP project

    ```
    PROJECT_ID=xxxx
    gcloud auth login
    gcloud config set project $PROJECT_ID
    gcloud config list
    ```

1. Enable Google Drive API

    https://console.cloud.google.com/apis/library/drive.googleapis.com?project=$PROJECT or command

    ```
    gcloud services enable drive.googleapis.com --project $PROJECT
    ```


1. Set `Application Default Credentials` with the necessary scope:

    ```
    gcloud auth application-default login --scopes "openid,https://www.googleapis.com/auth/userinfo.email,https://www.googleapis.com/auth/cloud-platform,https://www.googleapis.com/auth/sqlservice.login,https://www.googleapis.com/auth/drive"
    ```
1. Go run
    ```
    go run integrated-examples/gdrive/main.go --id-to-copy-from "1A3Yuew33YREFU5uLv05p9-PmwsewrC8OFTem8pWQcP8" --new-file-name "Copied by Golang Script"
    created 1RdcOgtrZRX70VwDzmUM6wS5BpmC6WR5dgfcSSgMAS_Q
    ```

1. Unset GCP project
    ```
    gcloud auth revoke
    gcloud config unset project
    ```

## Ref


1. Golang
    1. https://pkg.go.dev/google.golang.org/api/drive/v3
    1. https://pkg.go.dev/google.golang.org/api/drive/v3#NewFilesService
    1. https://github.com/googleapis/google-api-go-client/blob/main/examples/main.go
    1. https://github.com/sinmetal/other_than_gcp/tree/dcca7fbe509476b84ccd27aa7e09a617344cb42d
    1. https://www.loginradius.com/blog/engineering/google-authentication-with-golang-and-goth/
    1. [Google Drive上のファイルをOCRする](https://qiita.com/shin1ogawa/items/559ec58f6d7840721e5a)
1. GitHub Issues
    1. https://github.com/hashicorp/terraform-provider-google/issues/12774
    1. https://github.com/gcpug/nouhau/issues/124
    1. https://github.com/hashicorp/terraform-provider-google/issues/12774
1. Google Auth
    1. https://cloud.google.com/sdk/gcloud/reference/auth/print-access-token
1. Google Drive API
    1. https://developers.google.com/drive/api/reference/rest/v3/files/copy
