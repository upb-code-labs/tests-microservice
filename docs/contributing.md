# General organization contribution guidelines

Please, read the [General Contribution Guidelines](https://github.com/upb-code-labs/docs/blob/main/CONTRIBUTING.md) before contributing.

# Tests micro-service contribution guidelines

## Project structure / architecture

The tests micro-service's architecture is based on the Hexagonal Architecture in order to make it easier to maintain and extend.

The project structure is as follows:

- `domain`: Contains the business logic of the tests micro-service.
  - `definitions`: Contains the interfaces (contracts) of the tests micro-service. 
  - `entities`: Contains the entities of the tests micro-service.
  - `dtos`: Contains the data transfer objects of the tests micro-service.

- `application`: Contains the application logic (use cases) of the tests micro-service.

- `infrastructure`: Contains the implementation of the tests micro-service's interfaces (contracts) and the external dependencies of the tests micro-service.
  - `implementations`: Contains the implementations of the `domain` interfaces (contracts).
  - `rabbitmq`: Contains the RabbitMQ related code.
  - `static_files`: Contains the static files micro-service related code.

- `utils`: Contains the utility functions of the tests micro-service.

Note that, as the `application` layer cannot depend on the `infrastructure` layer, the `application` layer only uses the interfaces (contracts) defined in the `domain` layer, so, any implementation in the `infrastructure` layer can injected and used by the `application` layer without any problem.

The above allows us, for instance, to use multiple `LanguageTestsRunner` implementations (`JavaTestsRunner`, `TypescriptTestsRunner`, etc.) without having to change the `application` layer.

## Local development

### Dependencies

The following dependencies are required to run the tests micro-service locally:

- [Go 1.21.5](https://golang.org/doc/install)
- [Podman](https://podman.io/getting-started/installation) (To build and test the container image)

Please, note that `Podman` is a drop-in replacement for `Docker`, so you can use `Docker` instead if you prefer.

Additionally, you may want to install the following dependencies to make your life easier:

- [Air](https://github.com/cosmtrek/air) (for live reloading)

### Running the tests micro-service locally

As the role of the tests micro-service is to listen for messages in the `submissions` queue and run the tests for the submissions, you will need to run the [gateway](https://github.com/UPB-Code-Labs/main-api) project first in order to initialize the queue, database and the other micro-services and send submissions by using the REST API.

After you have the gateway running, you can start the tests micro-service by running the following command:

```bash
air 
```

This will start the tests micro-service and will watch for changes in the source code and restart the service automatically.

Additionally, you may want to generate a `.air.toml` file and add the `tests_exec_dir/` directory to the `exclude_dir` list in order to avoid restarting the service when the tests are executed, to do this, run the following command or refer to the [Air documentation](https://github.com/cosmtrek/air)

```bash
air init
```