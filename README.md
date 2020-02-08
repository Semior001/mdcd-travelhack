# gotemplate
template for future projects or projects that will be developed during hackathons

# package structure
    api_middleware                     - golang-written project, api middleware between image processing and client
    └── app                            - main application directory 
        ├── cmd                        - everything related to cli commands
        │   ├── common.go              - general parameters, general commander
        │   ├── server.go              - serve command, its description and execution
        │   ├── db.go                  - all commands to operate the database
        │   └── register_admin.go      - command to register admin
        ├── rest
        │   ├── public                 - public controllers
        │   ├── private                - private controllers
        │   ├── admin                  - admin controllers
        │   ├── http_errors.go
        │   ├── http_errors_test.go    
        │   └── server.go              - builder for web-server
        ├── store                      - describes everything, which is related to storage functions
        │   └── user                   
        │       └── user.go            - user and his structs, interface for working directly with 
        │                                 database and service methodss
        └── main.go                    - application entrypoint, processes commands and cli-arguments

# services:
* api_middleware - golang-written web services that provides api to frontend
* imgproc 

# running application
```bash
docker-compose up --build [-d]
```

# restarting application
```bash
docker-compose restart [service name]
```
