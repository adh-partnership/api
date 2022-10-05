# ADH-Parnership API

ADH Partnership Monolithic API.

## Introduction

This monolithic API is meant to break apart the previous Laravel web application into a backend and frontend.  This backend will provide data to the frontend as well as other web components in the future, ie: IDS.

This API is designed to be run in Kubernetes with a few helper tools. Ideally, we would be using Vault to feed Redis and Database credentials, and external secrets for other information, however, for our purpose this is not necessary so we will be using normal secrets. These secrets can be replaced with External Secrets, Sealed Secrets, etc. in the future or if necessary.

The API follows standard RESTful API designs. The API documentation is accessible through root, and the methods are standard:

- GET
- POST (create)
- PUT (replace)
- PATCH (update)
- DELETE

## Tests

Tests are a good thing and should be written were practical.  To run tests:

```bash
make test
```

An example of a test can be found in pkg/database/dto/user_test.go. Use either the built-in testing framework
or [stretchr/testify](https://github.com/stretchr/testify).
