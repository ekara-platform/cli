package message

const (
	LOG_CLI_IMAGE           string = "Ekara installation based on the Docker image: %s\n"
	LOG_COMMAND_COMPLETED   string = "Command completed"
	LOG_CONFIG_CONFIRMATION string = "Launching ekara CLI for %s:%s"
	LOG_FLAG_CONFIRMATION   string = "Flagged, %s %s"

	LOG_OUTPUT_FILE_IGNORED string = "The output file %s will not be create because the output is turned off\n"
	LOG_CONTAINER_LOG_IS    string = "The container logs will be written into %s\n"

	LOG_INIT_FLAGGED_DOCKER_CLIENT string = "Flagged docker client initialization..."
	LOG_INIT_DOCKER_CLIENT         string = "Docker client initialization..."

	LOG_SSH_PUBLIC_CONFIRMATION  string = "Launching ekara CLI with the public SSH key %s"
	LOG_SSH_PRIVATE_CONFIRMATION string = "Launching ekara CLI with the private SSH key %s"

	//Actions
	LOG_CHECKING_FROM   string = "Checking from: %s \n"
	LOG_DUMPING_FROM    string = "Dumping from: %s \n"
	LOG_CREATING_FROM   string = "Creating from: %s \n"
	LOG_DEPLOYING_FROM  string = "Deploying from: %s \n"
	LOG_INSTALLING_FROM string = "Installing from: %s \n"

	//DOCKER IMAGE AND CONTAINER
	LOG_GET_IMAGE             string = "Get or refresh the latest CLI image"
	LOG_WAITING_DOWNLOAD      string = "waiting for the download to be completed"
	LOG_DOWNLOAD_COMPLETED    string = "image download completed"
	LOG_WAITING_STOP          string = "waiting for the container to stop"
	LOG_STOPPED               string = "container stopped"
	LOG_CONTAINER_LOG_WRITTEN string = "The container logs have been written into %s\n"

	// Prompt messages
	PROMPT_RESTART             string = "Are you sure you want to recreate the starter container now (Y/N) "
	LOG_FAIL_ON_PROMPT_RESTART string = "Cannot go further if you refuse to stop the running container!"

	LOG_PASSING_CONTAINER_ENVARS string = "Env passed to the container %v\n"

	LOG_QUALIFIED_NAME string = "The environemt qualified name is :%s"

	LOG_GETTING_HTTP_PROXY  string = "Getting HttpProxy from environment: %s\n"
	LOG_GETTING_HTTPS_PROXY string = "Getting HttpsProxy from environment: %s\n"
	LOG_GETTING_NO_PROXY    string = "Getting HttpsProxy from environment: %s\n"

	ERROR_REQUIRED_DESCRIPTOR_URL        = "the \"descriptor url\" should be defined"
	ERROR_REQUIRED_FLAG           string = "the flag \"%s\" should be defined"
	ERROR_REQUIRED_ENV            string = "the environment variable \"%s\" should be defined"
	ERROR_REQUIRED_SSH_PUBLIC     string = "the flag \"public_ssh\" must be defined"
	ERROR_REQUIRED_SSH_PRIVATE    string = "the flag \"private_ssh\" must be defined"
	ERROR_COPYING_SSH_PUB         string = "Error copying the SSH public key %s"
	ERROR_COPYING_SSH_PRIV        string = "Error copying the SSH public key %s"

	ERROR_CREATING_EXCHANGE_FOLDER  string = "Error creating the exchange folder %s"
	ERROR_UNREACHABLE_PARAM_FILE    string = "Error, the file \"%s\" cannot be located"
	ERROR_CREATING_EKARA_ENGINE     string = "Error creating the Ekara engine %s"
	ERROR_INITIALIZING_EKARA_ENGINE string = "Error initializing the Ekara engine %s"

)
