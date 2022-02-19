# scrumpoker_api

Add your app description

## Makefile and environment variables

Run **make help** to see all available tasks.

### environment variables

Check the Makefile to see the default values.  
You can override env-vars when calling make (except APP_NAME and VERSION)
Example:
```
PORT=1234 make run
```

__env-vars__
| env-var           | description                                                       |
| :---              | :---                                                              |
| APP_NAME          | app name, will be used as db-name and docker image/container-name |
| VERSION           | app version, will be added added as a tag to the docker image     |
| HOST              | host ip/name to listen for requests                               |
| PORT              | port to listen for requests                                       |
| DOCKER_NETWORK    | docker network to join in order to reach other containers         |
| DATABASE_URL      | database connection url                                           |
| TEST_DATABASE_URL | by default same as DATABASE_URL with the table postfix **_test**  |

### DATABASE_URL

value will be passed to gorm, check https://gorm.io/docs/connecting_to_the_database.html to see databases that work out of the box.

## run

run api locally

```
HOST=localhost \
PORT=8080 \
DATABASE_URL=postgresql://@localhost/db_override \
make run
```

## run docker container locally

Run docker container in an environment similar to production, locally.  
No ports are exposed, use a revers proxy.
For example Caddy: https://caddyserver.com/docs/quick-starts/reverse-proxy

```
HOST=localhost \
PORT=8080 \
DATABASE_URL=postgresql://@localhost/app_name \
make run
```

## deployment

### build & upload image

Builds a docker image for linux/amd64 & uploads it to the given GitHub Container Registry.  
Image name is APP_NAME:VERSION, check Makefile.

```
GHCR_PAT="GitHub Container Registry Private Access Token" \
GHCR_USER="GitHub Container Registry Username" \
make release
```

### deploy to server

1. Connects to REMOTE_SERVER via ssh.
2. Pulls the image APP_NAME:VERSION & :lastest from GitHub Container Registry.
3. Stop/Remove old container & start new one with :lastest image.

No ports are exposed, use a revers proxy.
For example Caddy: https://caddyserver.com/docs/quick-starts/reverse-proxy

```
REMOTE_SERVER="username@example.com" \
GHCR_PAT="GitHub Container Registry Private Access Token" \
GHCR_USER="GitHub Container Registry Username" \
DATABASE_URL=postgresql://@localhost/app_name \
make deploy
```

## code structure

### cmd

Contains entry points of your app. By default there is only one, api/main.go

### docs

Add all the document describing your app in this directory.

**docs/swagger.yaml** contains basic info about your app. Placeholder APP_NAME & VERSION will be replaced by the values contained in the Makefile.  

**docs/swagger-gen.sh** will be called by make docs. Generates handler & models for your domains based on the domains swagger .yaml definition.  
All swagger .yaml files will be assembled into one, internal/service/server/swagger/swagger_gen.yaml  
By default the swagger-ui will be served on /api/scrumpoker_api/doc

### internal/service

Add services that will be used across the codebase.

### internal/domain

Contains different logical domains of your app. The domain will then be served on the subroute named after the domain. Example /api/scrumpoker_api/session

To create a new **DOMAIN** follow these steps:
1. Add directory/package and the domains swagger specification whithin. The **DOMAIN**.yaml file defines the http endpoints and models.
2. Executing **make docs** to generate **DOMAIN**_gen.go files.
3. Add a **DOMAIN**_impl.go file and implement the ServerInterface interface to handle the generated requests.
4. Add a test case in test/**DOMAIN**_test.go, call the functions within **DOMAIN**_impl.go with the generated request body structs.
5. (optionally) add a repository.go if that domain interacts with a persistent storage.
6. Modify **internal/domain/domain.go**
  1. add the domain model to DbMigrate to generate the db-tables.
  2. add the domains ServerInterfaceImpl to RegisterHandlers

### internal/domain/common

Contains model.go and it's common.yaml swagger definition. The BaseModel got fields commonly used accross db tables.

### test

Contains your test cases. main_test.go ensures that the database connection is established before the tests start.  
The test database will be droped in the makefile before the tests start.
