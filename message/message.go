package message

const (
	LOG_CLI_IMAGE           = "Ekara installation based on the Docker image: %s\n"
	LOG_COMMAND_COMPLETED   = "Command completed"
	LOG_CONFIG_CONFIRMATION = "Launching ekara CLI for %s:%s"
	LOG_FLAG_CONFIRMATION   = "Flagged, %s %s"

	LOG_OUTPUT_FILE_IGNORED = "The output file %s will not be create because the output is turned off\n"
	LOG_CONTAINER_LOG_IS    = "The container logs will be written into %s\n"

	LOG_INIT_FLAGGED_DOCKER_CLIENT = "Flagged docker client initialization..."
	LOG_INIT_DOCKER_CLIENT         = "Docker client initialization..."

	LOG_SSH_PUBLIC_CONFIRMATION  = "Launching ekara CLI with the public SSH key %s"
	LOG_SSH_PRIVATE_CONFIRMATION = "Launching ekara CLI with the private SSH key %s"

	//Actions
	LOG_CHECKING_FROM   = "Checking from: %s \n"
	LOG_DUMPING_FROM    = "Dumping from: %s \n"
	LOG_CREATING_FROM   = "Creating from: %s \n"
	LOG_DEPLOYING_FROM  = "Deploying from: %s \n"
	LOG_INSTALLING_FROM = "Installing from: %s \n"

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

	LOG_GETTING_HTTP_PROXY  = "Getting HttpProxy from environment: %s\n"
	LOG_GETTING_HTTPS_PROXY = "Getting HttpsProxy from environment: %s\n"
	LOG_GETTING_NO_PROXY    = "Getting NoProxy from environment: %s\n"

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
