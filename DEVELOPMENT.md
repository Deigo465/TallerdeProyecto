# Development

To setup the dependencies for repository we recommend using Ubuntu either natively or through `WSL2`.

- [Development](#development)
  - [Project structure](#project-structure)
  - [Dependencies](#dependencies)
    - [SQLite](#sqlite)
    - [Go](#go)
    - [Bun](#bun)
  - [Testing](#testing)
    - [Unit tests](#unit-tests)
    - [Coverage tests](#coverage-tests)
      - [Coverage of all package](#coverage-of-all-package)
  - [TailwindCSS](#tailwindcss)

## Project structure

The project is divided into two main folders:
- `cmd` contains the main entry points for the project
- `pkg` contains all the packages that are used in the project

This project follows Uncle Bob's _clean architecture_ principles, where the main business logic is separated from the delivery mechanism (web server) and the data access layer (repositories).

The tree structure of the project is as follows:
| Path | Description |
| --- | --- |
| `cmd/web-server` | Main entry point for the web server |
| `pkg/domain/` | Wrapper for all the main business logic of the project  |
| `pkg/domain/entities` | All the entities of the project |
| `pkg/domain/usecases` | All the use cases (business logic) of the project |
| `pkg/web` | Everything related to the deliver of web, such as handlers, routes and views |
| `pkg/web/views/components` | Vue components |
| `pkg/repositories` | The repositories layer  |
| `pkg/mocks` | Mocks used in testing, this is probably going to be moved in the future  |
| `pkg/interfaces` | The interfaces that connect UseCases with the implementations  |
| `public` | Images, CSS and JS you may need in the project  |
| `docs` | This project documentation, architecture, etc.  |



## Dependencies

### SQLite
```bash
sudo apt-get install sqlite3
```

### Go
Go v1.20 is needed to build this project
check [Go installation guide](https://golang.org/doc/install) for more information on how to install Go.


Once Go is installed you can install `Air` to improve development experience, this is the prefer method of developing this project.
```bash
go install github.com/cosmtrek/air@latest
```
Otherwise you can run the webserver with Go:
```bash
go run -v ./cmd/server -port 3000      
```

### Bun

[Bun](https://bun.sh/) is fast javascript runtime env, to install it  in Ubuntu run:
```bash
curl -fsSL https://bun.sh/install | bash
```

Then install all the project JS dependencies with:
```bash
bun install
```


## Testing
All tests are located within the packages they are testing.

### Unit tests
To run all the unit tests in the app run the following command:
```bash
go test ./...
```

### Coverage tests
Run the following script to get the coverage tests for the whole app
```bash
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out && rm coverage.out
```
#### Coverage of all package
go test -v -coverpkg=./... -coverprofile=profile.cov ./...
go tool cover -func profile.cov

> **Tip**: You can add the following alias to your `.bashrc` or `.zshrc` to run the coverage tests with a single command:
> ```bash
> alias cover='go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out && rm coverage.out'
> ```


## TailwindCSS

```bash
tailwindcss --input public/css/tailwind.css --output public/css/output.css --watch
```
