
# api_middleware service for the image processing application

This service provides basic methods for authenticating user and storing, processing images, using rerouting the request to the internal image processing web-service (some kind of proxy)

# structure

everything should be covered with tests

    app
    ├── cmd
    │   ├── common.go       - some common fields for every command and some common arguments
    │   ├── migrate.go      - migrates the database 
    │   └── server.go       - runs the server
    ├── rest                - everything related to REST API 
    │   ├── image.go        - controller for images
    │   ├── user.go         - controller for users
    │   ├── http_errors.go  - default methods to return the error page or JSON if error is present
    │   └── server.go       - describes the router, middleware etc. 
    ├── store               - everything related to storing
    │   ├── image           - image store and its implementations
    │   │   ├── image.go    - describes store interface and image struct
    │   │   └── postgres.go - postgres implementation of store
    │   └── user            - user store and its implementations
    │       ├── user.go     - describes store interface and user struct
    │       └── postgres.go - postgres implementation of store
    └── main.go             - entrypoint for application, delivering deps to rest and store

## Authentication

Next methods will be supported:
* internal, via JWT
* LDAP
