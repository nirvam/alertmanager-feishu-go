# Alertmanager Feishu Go

A Go implementation of the Alertmanager Feishu webhook service. Converts Alertmanager alerts into Feishu messages (Text or Interactive Cards).

## Features

- **Standard Markdown Support**: Uses Feishu Card Schema 2.0.
- **Color Coded**: Alerts are color-coded based on status (Firing/Resolved) and severity.
- **Signature Verification**: Supports Feishu webhook secret for security.
- **Lightweight**: Built with Go and packaged in a distroless Docker image.

## Installation

### From Source

```bash
go install github.com/nirvam/alertmanager-feishu-go@latest
```

### Docker

```bash
docker pull ghcr.io/nirvam/alertmanager-feishu-go:latest
```

## Configuration

The service is configured via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `FEISHU_WEBHOOK_URL` | **Required**. Your Feishu bot webhook URL. | - |
| `FEISHU_SECRET` | Optional. Feishu bot signature secret. | - |
| `MESSAGE_TYPE` | Message format: `interactive` or `text`. | `interactive` |
| `APP_HOST` | Host to bind the service to. | `0.0.0.0` |
| `APP_PORT` | Port to bind the service to. | `8080` |

## Usage

### Start the Service

```bash
alertmanager-feishu-go serve
```

### Send a Test Message

```bash
alertmanager-feishu-go test
```

## License

MIT License
