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

Tests are a good thing and should be written where practical.  To run tests:

```bash
make test
```

An example of a test can be found in pkg/database/dto/user_test.go. Use either the built-in testing framework
or [stretchr/testify](https://github.com/stretchr/testify).

## Installation

There are a few ways to use this... the preferred is using a Kubernetes deployment (see the manifests directory). However, you can also run this locally.

### Hard Coded Discord Webhook Names

These should be configured in the config.yaml file with appropriate webhook URLs.  The webhook names are:

- default - If the webhook name is not found, this is used.. if not set, will be logged to console only.
- online - Webhook for displaying online controllers
- seniorstaff - Webhook for displaying senior staff messages
- uploads - Action logs for uploaded files
- pending_feedback - Notification about new feedback submissions
- broadcast_feedback - Approved feedback reports are broadcast to this webhook, intended to be public praises
- role - Notifications about added/removed roles
- visitng_application - Notifications related to Visiting Applications (new, approved, denied, etc), should be restricted
- staffing_request - Webhook for displaying requests for future staffing

### Kubernetes

The `manifests/deployment.yaml` is an example of the deployment. We use an init-container to do environment variable substitution using notation similar to that of Helm.
The `manifests/configmap.yaml` is your Config Map that specifies the template used. Pass environment variables by either specifying them in `env` or mounting secrets.

### Local

To run locally, you'll need to setup a MySQL instance.

1. Download the desired release from github
2. Create a database, user and password in MySQL for the API to utilize
3. Rename config.yaml.example and edit filling in the values as appropriate
4. Run `./api bootstrap` to bootstrap the database, this will load the basics into the database, sync the Roster from VATUSA, and add staff roles based on the facility in VATUSA
5. Run `./api add-role --cid 123456 --role wm` if you need to manually add a role... once the environment is live, it is not recommended to use this command. Default roles are: atm, datm, ta, fe, wm, ec, mtr, and events.
6. Run `./api server` to start

### FAQ

1. How do I start the API automatically on boot?

    - Use a systemd service file, or utilize Docker/Containerd/Kubernetes.

        - For systemd, see the `systemd` directory for an example service file or visit [this link](https://www.digitalocean.com/community/tutorials/how-to-use-systemctl-to-manage-systemd-services-and-units) for more information.

2. How do I update the API?

    - Download the latest release, stop the API, replace the binary, and start the API.
    - If using Docker/Containerd/Kubernetes, follow the standard procedure for updating a container. The container is ephemeral, so data will not be lost.

3. How do I update the API without downtime?

    - Use a load balancer that supports zero downtime deployments.  This is not a requirement, but is recommended.
    - Use Kubernetes and utilize the rolling update strategy.

4. Is there Swagger/OpenAPI documentation available?

Yes, we have Swagger (OpenAPI 2.0) documentation available. When you start the API, the documentation can be found by visiting the root directory.  For example, if you are running the API on localhost on port 3000, you can visit <http://localhost:3000/> to view the documentation.

## Email Templates

The following are the coded email templates and available variables:

- activity_warning

    - FirstName
    - LastName

- visitor_accepted

    - FirstName
    - LastName
    - Rating

- visitor_denied

    - FirstName
    - LastName
    - Rating
    - Reason

### Email Template Format

We use Go's templating engine to generate emails. More information can be found at:

- [Go html/template](https://pkg.go.dev/html/template)
- [Building Web Applications](https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/07.4.html)

Emails send in HTML format.

### Email Template Functions

- findRoles - Returns an array of names for a given role in signature form ("First Last, Role")
- urlEscape - Returns a URL escaped string

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
