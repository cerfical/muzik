# muzik

A simple RESTful API for managing collections of music tracks.

For now, only a limited set of operations is available, and the data model is quite primitive.

The formal description is available as an [OpenAPI specification](api/openapi.yaml).

## Usage

The project conists of the two executables:
- `api` represents the API server itself and can be run with something like:
  ```shell
  go run ./cmd/api/
  ```
  At the very least, the database password must be specified via the `MUZIK_DB_PASSWORD` variable
  or in a config file.
  Other options are avaiable and can be easily inferred from the [config structure](internal/config/config.go) used to store the configs.
  If a config file is used, it must be specified as the only argument to the executable.

- `web` is a trivial (and probably broken) HTTP server that serves a single HTML index page.
  Currently, its only use is to try out the API through a friendly user interface.

A [docker-compose](docker-compose.yaml) config file is provided to readily start the API and WEB servers along with all required runtime dependencies (`postgres`, `nginx`).
