# GitHub Webhook Proxy for Jenkins

This proxy server forwards GitHub webhook requests to a Jenkins instance while validating the webhook secret. If the secret is incorrect, the server returns a 404 status code.

## Requirements

- Golang (1.16 or higher)
- A GitHub repository
- A Jenkins instance with the GitHub plugin installed

## Configuring the Proxy Server

To configure the proxy server, you can set the following environment variables:

- `GITHUB_WEBHOOK_SECRET`: The webhook secret used for validating incoming requests.
- `WEBHOOK_ENDPOINT`: The webhook endpoint path, default is `/github-webhook/`.
- `PROXY_PORT`: The port on which the proxy server listens, default is `:8080`.
- `JENKINS_URL`: The URL of your Jenkins instance, default is `http://localhost:8081`.

For example:

```bash
export GITHUB_WEBHOOK_SECRET="your_secret_here"
export WEBHOOK_ENDPOINT="/github-webhook/"
export PROXY_PORT=":8080"
export JENKINS_URL="http://localhost:8081"
```

## Running the Proxy Server

Run the following command in your terminal to start the proxy server:

```shell
go run main.go
```

The proxy server will listen on the configured port for incoming GitHub webhook requests and forward them to the configured Jenkins instance.

## Configuring GitHub Webhooks

1. Go to your GitHub repository's settings page.
2. Click on "Webhooks" in the left sidebar.
3. Click the "Add webhook" button.
4. Set the "Payload URL" to the URL of the proxy server (e.g., `https://your-proxy-server.com:8080/github-webhook/`).
5. Set the "Content type" to `application/json`.
6. Enter your webhook secret in the "Secret" field.
7. Choose the events you want to trigger the webhook.
8. Click the "Add webhook" button to save the configuration.

## Configuring the Jenkins GitHub Plugin

1. In your Jenkins instance, go to "Manage Jenkins" > "Manage Plugins".
2. Install the "GitHub Plugin" if it's not already installed.
3. Go to "Manage Jenkins" > "Configure System".
4. Scroll down to the "GitHub" section and click the "Add" button next to "GitHub Servers".
5. Enter your GitHub credentials and configure the connection settings.
6. Click "Test Connection" to ensure the connection is working.
7. Save the configuration.
