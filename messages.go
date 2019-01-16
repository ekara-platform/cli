package main

const (
	CLI_DESCRIPTION string = "Ekara Command Line Interface."

	//	// Error messages
	ERROR_CREATING_EXCHANGE_FOLDER string = "Error creating the exchange folder %s"

	ERROR_REQUIRED_CONFIG string = "the flag \"config\" should be defined"

	ERROR_REQUIRED_CLIENT           string = "the flag \"client\" should be defined"
	ERROR_REQUIRED_ENV              string = "the environment variable \"%s\" should be defined"
	ERROR_REQUIRED_FLAG             string = "the flag \"%s\" should be defined"
	ERROR_SESSION_NOT_CLOSED        string = "Unable to complete the logout! You can force the logout manually deleting the file \"%s\""
	ERROR_CLIENT_SESSION_NOT_CLOSED string = "Unable to complete client \"%s\" logout! You can force it manually deleting the file \"%s\""
	ERROR_NO_PROVIDED_USER          string = "No user has been provided using --%s and we cannot use the system user"
	ERROR_READING_PASSWORD          string = "Error reading the password"

	ERROR_REQUIRED_SSH_PUBLIC  string = "the flag \"public_ssh\" must be defined"
	ERROR_REQUIRED_SSH_PRIVATE string = "the flag \"private_ssh\" must be defined"

	ERROR_COPYING_SSH_PUB  string = "Error copying the SSH public key %s"
	ERROR_COPYING_SSH_PRIV string = "Error copying the SSH public key %s"

	ERROR_UNREACHABLE_PARAM_FILE string = "Error, the file \"%s\" cannot be located"

	ERROR_CREATING_EKARA_ENGINE     string = "Error creating the Ekara engine %s"
	ERROR_INITIALIZING_EKARA_ENGINE string = "Error initializing the Ekara engine %s"

	//	// Prompt messages
	PROMPT_RESTART  string = "Are you sure you want to recreate the starter container now (Y/N) "
	PROMPT_PASSWORD string = "Please enter password for %s:"

	//	// Log messages
	LOG_FAIL_ON_PROMPT_RESTART   string = "Cannot go further if you refuse to stop the running container!"
	LOG_CONFIG_CONFIRMATION      string = "Launching ekara CLI for %s:%s"
	LOG_SSH_PUBLIC_CONFIRMATION  string = "Launching ekara CLI with the public SSH key %s"
	LOG_SSH_PRIVATE_CONFIRMATION string = "Launching ekara CLI with the private SSH key %s"

	LOG_FLAG_CONFIRMATION          string = "Flagged, %s %s"
	LOG_INIT_FLAGGED_DOCKER_CLIENT string = "Flagged docker client initialization..."
	LOG_INIT_DOCKER_CLIENT         string = "Docker client initialization..."
	LOG_GET_IMAGE                  string = "Get or refresh the latest CLI image"
	LOG_WAITING_DOWNLOAD           string = "waiting for the download to be completed"
	LOG_DOWNLOAD_COMPLETED         string = "image download completed"
	LOG_WAITING_STOP               string = "waiting for the container to stop"
	LOG_STOPPED                    string = "container stopped"
	LOG_COMMAND_COMPLETED          string = "Command completed"
	LOG_LOGGED_AS                  string = "You are currently logged as %s on %s"
	LOG_ALREADY_LOGGED_AS          string = "You are already logged as %s on %s"
	LOG_ALREADY_LOGGED_OUT         string = "You are already logged out"
	LOG_LOGOUT_REQUIRED            string = "You must to logout first to be able to create an environment"
	LOG_LOGIN_REQUIRED             string = "You must be logged in be able to perform this command"
	LOG_CREATING_FROM              string = "Creating from: %s \n"
	LOG_DEPLOYING_FROM             string = "Deploying from: %s \n"
	LOG_INSTALLING_FROM            string = "Installing from: %s \n"
	LOG_UPDATING_FROM              string = "Updating from: %s \n"

	LOG_CHECKING_FROM string = "Checking from: %s \n"

	LOG_GETTING_STATUS        string = "Getting the status of %s \n"
	LOG_LOGGING_INTO          string = "Logging into: %s \n"
	LOG_CONTAINER_LOG_WRITTEN string = "The container logs have been written into %s\n"

	LOG_LOADING_EXTRA_ENVARS string = "Setting environment variables fron client: %s=%s"

	LOG_QUALIFIED_NAME string = "The environemt qualified name is :%s"
)
