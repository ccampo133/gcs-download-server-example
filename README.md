# gcs-download-server-example

Example Go server to download files from a Google Cloud Storage bucket.

# Usage

You need a sufficiently privileged [GCP service account](https://cloud.google.com/docs/authentication/getting-started).
For example:

```
export GOOGLE_APPLICATION_CREDENTIALS="KEY_PATH"
```

Then just run the application:

```
go run main.go -project="example-gcp-project"
```

This will start an HTTP server on port 8080 with the single route `GET /download/{bucket}/{filename:.*}`.

For more information, see: https://cloud.google.com/storage/docs/streaming
