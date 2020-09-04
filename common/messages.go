package common

const (
	LOG_CLI_IMAGE         = "Ekara installer image: %s\n"
	LOG_COMMAND_COMPLETED = "Command completed"
	LOG_CONTAINER_LOG_IS  = "The installer logs will be written into %s\n"
	LOG_SSH_CONFIRMATION  = "Using specified SSH keys"

	//Actions
	LOG_VALIDATING_ENV = "Validating environment..."
	LOG_DUMPING_ENV    = "Dumping environment..."
	LOG_APPLYING_ENV   = "Applying environment..."
	LOG_DESTROYING_ENV   = "Destroying environment..."

	//DOCKER IMAGE AND CONTAINER
	LOG_WAITING_DOWNLOAD      = "Waiting for the download to be completed"
	LOG_DOWNLOAD_COMPLETED    = "Image download completed"
	LOG_WAITING_STOP          = "Waiting for the container to stop"
	LOG_STOPPED               = "Container stopped"
	LOG_CONTAINER_LOG_WRITTEN = "The container logs have been written into %s\n"

	// Prompt messages
	PROMPT_RESTART             = "Are you sure you want to recreate the starter container now (y/n) ? "
	LOG_FAIL_ON_PROMPT_RESTART = "cannot go further if you refuse to stop the running container"

	LOG_PASSING_CONTAINER_ENVARS = "Env passed to the container %v\n"

	INVALID_DESCRIPTOR_URL     = "Invalid descriptor URL: %s\n"
	ERROR_REQUIRED_FLAG        = "The flag \"%s\" should be defined"
	ERROR_REQUIRED_ENV         = "The environment variable \"%s\" should be defined"
	ERROR_REQUIRED_SSH_PUBLIC  = "The flag \"public_ssh\" must be defined"
	ERROR_REQUIRED_SSH_PRIVATE = "The flag \"private_ssh\" must be defined"

	ERROR_CREATING_EXCHANGE_FOLDER  = "error creating the exchange folder %s"
	ERROR_UNREACHABLE_PARAM_FILE    = "error, the file \"%s\" cannot be located"
	ERROR_CREATING_EKARA_ENGINE     = "error creating the Ekara engine %s"
	ERROR_INITIALIZING_EKARA_ENGINE = "error initializing the Ekara engine %s"
)
