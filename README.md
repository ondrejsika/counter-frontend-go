# counter-frontend-go

## Container Images

- `ghcr.io/ondrejsika/counter-frontend-go`


## Configuration

- `API_ORIGIN` - Origin of the API server (default: `http://127.0.0.1:8000`)
- `FONT_COLOR` - Color of the counter text (default: `#000000`, black)
- `BACKGROUND_COLOR` - Color of the counter background (default: `#ffffff`, white)
- `FAIL_ON_ERROR` - Exit with fatal error when API request fails (set to `1` to enable, useful for Kubernetes probes)
