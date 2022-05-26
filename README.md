### Requirements:
- Go version 1.16

### RUN:

```
$ make install // load dependencies
$ cp .env.example .env
## => change HOST variable in .env to current location
$ mkdir config && cp conf.toml config/conf.toml
$ make docker-compose-up
$ make run-api
```

### Struct folder layout
```
/docs

- OpenAPI/Swagger specs, JSON schema files, protocol definition files.


/cmd

- Main applications for this project.
- The directory name for each application should match the name of the executable you want to have (e.g., `/cmd/myapp`).

/internal

- Private application and library code.


/pkg

Library code that's ok to use by external applications (e.g., `/pkg/mypubliclib`). Other projects will import these libraries expecting them to work, so think twice before you put something here :-) Note that the `internal` directory is a better way to ensure your private packages are not importable because it's enforced by Go.
```