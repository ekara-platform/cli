# Lagoon CLI
___

**Lagoon CLI** is the command line interface tool used to build and interact with Lagoon environments.

It will be compiled in order to provide a file which will be executed on the machine where the environment creation will be launched

## Available commands

* **create** Create a new environment.
* **update** Update an existing environment.
* **check**  Check the validy of an environment descriptor.
* **login**  Login into an environment manager API.
* **logout** Logout from an environment manager API.
* **status** Status of the environment manager API.

We can distinguish two types of commands

* **Docker** commands: In order to be executed these commands require a Docker configuration. Please refer to this [section](#docker-commands-required-configuration) to get the details of the configuration.
* **API** commands: These commands can require specific parameters of flags

## Command "create"
This command allows to create a new environment based on the provided environment descriptor.

Command type: **Docker**

Argument(s):

* descriptor : The location of the descriptor of the environment to create. This location can be an ULR or a path to the file system.


Flags: 

* `--client` 
* `--cert` (can be substituted by an environment variable) 
* `--api`  (can be substituted by an environment variable)
* `--host` (can be substituted by an environment variable)
* `--env`
* `--http_proxy` (can be substituted by an environment variable)
* `--https_proxy` (can be substituted by an environment variable)
* `--no_proxy` (can be substituted by an environment variable)
* `--output`
* `--file`

Or environment variables :

* `DOCKER_CERT_PATH`
* `DOCKER_HOST`
* `DOCKER_API_VERSION`
* `HTTP_PROXY`
* `HTTPS_PROXY`
* `NO_PROXY`

Example :

`$ cli create http://path.to.my.project/ --cert ./cert_location --host tcp://192.168.99.100:2376 --api 1.30 --client myClientName` 

Example writing the container logs into `container.log`:

`$ cli create http://path.to.my.project/ --cert ./cert_location --host tcp://192.168.99.100:2376 --api 1.30 --client myClientName --output`

Example writing the container logs into a specific file:

`$ cli create http://path.to.my.project/ --cert ./cert_location --host tcp://192.168.99.100:2376 --api 1.30 --client myClientName --output --file myLogFile.log `


## Command "update"
This command allows to update an existing environment based on the provided descriptor. 

In order to perform an update the user must be logged into the environment manager API corresponding to the environment to update.

Command type: **API**

Argument(s):

* descriptor : The location of the descriptor of the environment to create. This location can be an ULR or a path to the file system.


Flags: 

* `--output`
* `--file`



Example :

`$ lagoon update http://path.to.my.project/ --output --file myLogFile.log `


## Command "check"
This command allows to check the validity of an environment descriptor. 

Command type: **Docker**

Argument(s):

* descriptor : The location of the descriptor of the environment to validate. This location can be an ULR or a path to the file system.


Flags: 

* `--cert` (can be substituted by an environment variable) 
* `--api`  (can be substituted by an environment variable)
* `--host` (can be substituted by an environment variable)
* `--env`
* `--http_proxy` (can be substituted by an environment variable)
* `--https_proxy` (can be substituted by an environment variable)
* `--no_proxy` (can be substituted by an environment variable)
* `--output`
* `--file`

Or environment variables :

* `DOCKER_CERT_PATH`
* `DOCKER_HOST`
* `DOCKER_API_VERSION`
* `HTTP_PROXY`
* `HTTPS_PROXY`
* `NO_PROXY`

Example :

`$ cli check http://path.to.my.project/ --cert ./cert_location --host tcp://192.168.99.100:2376 --api 1.30 ` 

Example writing the container logs into `container.log`:

`$ cli check http://path.to.my.project/ --cert ./cert_location --host tcp://192.168.99.100:2376 --api 1.30 --output`

Example writing the container logs into a specific file:

`$ cli check http://path.to.my.project/ --cert ./cert_location --host tcp://192.168.99.100:2376 --api 1.30 --output --file myLogFile.log `



## Command "login"
This command performs a login into an environment manager API.

Command type: **API**

Argument(s):

* url : The url of the environment manager API where to login

Flag(s):

* `--user` login ID. If missing then the CLI it will use the shell's user ID.

Example :

`$ lagoon login http://path.to.the.api --user usrXX`

## Command "logout"
This command performs a logout from an environment manager API.

Command type: **API**

Example :

`$ lagoon logout`

## Command "status"
This command returns the status of the environment manager API where the user is logged in.

The user must be logged into an environment manager API to get its status.


Command type: **API**

Example :

`$ lagoon status`

## Docker commands required configuration
___
In order to interact with docker **Lagoon CLI** requires the following configuration:

* The address of the docker host wherein you want to create or update an environment
* The docker host's certificates location
* The version of the docker host API we will deal with


___
### Docker Configuration using `Flags`

The Docker flags exposed by **Lagoon CLI** are:

* `--cert` : the Docker certificates location
* `--api` : the Docker API version
* `--host` : the Docker host address 

 
> If you decide to configure **Lagoon CLI** using flags then remember that all these 3 flags must be setted.

___
### Docker Configuration using `environment variables`

If you want to create an environment on your own docker host you can take advantage of these predefined environment variables.

* `DOCKER_CERT_PATH` : the certificates location
* `DOCKER_HOST` : the docker host address 

You will need to create an enviroment variable

* `DOCKER_API_VERSION` : the API version


---
### How Lagoon CLI will decide to use Docker flags or environment variables

If any of `--cert` , `--api` or `--host` is setted then **Lagoon CLI** will use flags to establish the connection with the docker daemon, if not then the environment variables will be used.
