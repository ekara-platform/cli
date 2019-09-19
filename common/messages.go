package common

const (
	LOG_CLI_IMAGE         = "Ekara installer image: %s\n"
	LOG_COMMAND_COMPLETED = "Command completed"
	LOG_CONTAINER_LOG_IS  = "The installer logs will be written into %s\n"
	LOG_SSH_CONFIRMATION  = "Using specified SSH keys"

	//Actions
	LOG_VALIDATING_ENV = "Validating environment: %s \n"
	LOG_DUMPING_ENV    = "Dumping environment: %s \n"
	LOG_APPLYING_ENV   = "Applying environment: %s \n"

	//DOCKER IMAGE AND CONTAINER
	LOG_GET_IMAGE             = "Get or refresh the latest CLI image"
	LOG_WAITING_DOWNLOAD      = "Waiting for the download to be completed"
	LOG_DOWNLOAD_COMPLETED    = "Image download completed"
	LOG_WAITING_STOP          = "Waiting for the container to stop"
	LOG_STOPPED               = "Container stopped"
	LOG_CONTAINER_LOG_WRITTEN = "The container logs have been written into %s\n"

	// Prompt messages
	PROMPT_RESTART             = "Are you sure you want to recreate the starter container now (Y/N) "
	LOG_FAIL_ON_PROMPT_RESTART = "Cannot go further if you refuse to stop the running container!"

	LOG_PASSING_CONTAINER_ENVARS = "Env passed to the container %v\n"

	LOG_QUALIFIED_NAME = "The environment qualified name is :%s"

	ERROR_REQUIRED_DESCRIPTOR_URL = "The \"descriptor url\" should be defined"
	ERROR_REQUIRED_FLAG           = "The flag \"%s\" should be defined"
	ERROR_REQUIRED_ENV            = "The environment variable \"%s\" should be defined"
	ERROR_REQUIRED_SSH_PUBLIC     = "The flag \"public_ssh\" must be defined"
	ERROR_REQUIRED_SSH_PRIVATE    = "The flag \"private_ssh\" must be defined"
	ERROR_COPYING_SSH_PUB         = "Error copying the SSH public key %s"
	ERROR_COPYING_SSH_PRIV        = "Error copying the SSH public key %s"

	ERROR_CREATING_EXCHANGE_FOLDER  = "Error creating the exchange folder %s"
	ERROR_UNREACHABLE_PARAM_FILE    = "Error, the file \"%s\" cannot be located"
	ERROR_CREATING_EKARA_ENGINE     = "Error creating the Ekara engine %s"
	ERROR_INITIALIZING_EKARA_ENGINE = "Error initializing the Ekara engine %s"
)
