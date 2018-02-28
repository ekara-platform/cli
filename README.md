# Lagoon Starter
___

**Lagoon Starter** is the command line interface tool used to build a Lagoon environment.

It will be compiled in order to provide a file which will be executed on the machine where the environment creation will be launched


## Required configuration
___
In order to create an environment **Lagoon Starter** requires the following configuration:

* The location of your environment deployment descriptor
* The address of the docker host wherein you want to create the environment
* The docker host's certificates location
* The version of the docker host API we will deal with


___
### Configuration using `Flags`

The flags exposed by **Lagoon Starter** are:

* `--cert` : the certificates location
* `--api` : the API version
* `--host` : the docker host address 
* `--config` : the deployment descriptor location

Example :

`lagoonstarter --config http://blablabla.com/mydescriptor.yaml --host tcp://192.168.99.100:2376 --api 1.30 --cert C:\Users\xxx\.docker\machine\machines\default
`
 
> If you decide to configure **Lagoon Starter** using flags then remember that all these 4 flags must be setted.

___
### Configuration using a mix of `Flags` and environment variables

If you want to create an environment on your own docker host you can take advantage of these predefined environment variables.

* `DOCKER_CERT_PATH` : the certificates location
* `DOCKER_HOST` : the docker host address 

You will need to create an enviroment variable

* `DOCKER_API_VERSION` : the API version

> Even if you decide to configure **Lagoon Starter** using environment variables the flag `--config` is still mandatory to reference the deployment descriptor location.

---
### How Lagoon Starter will decide to use flags or environment variables

If any of `--cert` , `--api` or `--host` is setted then **Lagoon Starter** will use flags to establish the connection with the docker daemon.