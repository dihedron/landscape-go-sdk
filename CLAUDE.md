# Ubuntu Landscape Go SDK

## Description

The project contains an SDK, located under the package `pkg/landscape`; the SDK wraps the REST API calls to Canonical Ubuntu Landscape as Go objects and methods.

The SDK exposes a single API client in package `landscape`; the API client is object `API` and contains:

* a `go-resty` v3 client reference, which will be used for all API HTTP requests;
* a set of services, e.g. `ComputerService`, each exposing one entity from the REST API.

Services embed a `Service` struct which acts as a common base and simply wraps the HTTP client (as a reference to the original `go-resty` client in the `API` object).

## Authentication

The `API`can be initialised either using `email` and `password`, plus an optional `account` info (which are used fo basic authentication), or via a JWT `token` acquired in a previous authentication flow.

This can be achieved by providing authentication info through options to the `New` method.

When the `email`, `password` (and optional `account`) values are provided in the authentication options, the client will also perform an authentication cycle and store the resulting JWT token in the client itself, for the caller to retrieve ans tore away.

When the authetication options contain the JWT token directly, the `New` method will directly initialise the client's token.

The client will provide a `Token()` method that allows the caller to retrieve the authentication token.

## How a service is implemented

Each service should only expose SDK specific types and errors. Implementation details such as the error messages returned by the REST API calls must always be wrapped as go `error`s and be made available through a specific method (e.g. `Message`).

Each service basically implements the CRUD of the associated entity and some methods on the collection (e.g. `List`,`Search`, `DeleteAll` if available in the REST API).

Services do not stub methods that are not available in the API.

When the REST API provides a way to performs actions on resources, e.g. via pseudo-entities such as a restart endpoint for computers, they are exposed throuhgh the Service as methods.

## The Command Line Interface

The project will also contain a command line, called `landscape`, which allows to manipulate Ubuntu Landscape entities.

The command line interface is built around entities and verbs, similar to how Docker works, e.g.

```bash
$> landscape computer list 
```

or 

```bash
$> landscape computer create --data=computer-info.json 
```

All methods on SDK entities are exposed through command line commands. Each command is an instance of a `Command`.
Commands are structs; struct fields are annotated according to `github.com/jessevdk/go-flags` conventions
and are gouped together in a multi-level command line structure.




