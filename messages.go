package main

const (
	CLI_DESCRIPTION string = "Lagoon Command Line Interface."

	ERROR_REQUIRED_CONFIG string = "the flag \"config\" should be defined"
	ERROR_REQUIRED_ENV    string = "the environment variable \"%s\" should be defined"
	ERROR_REQUIRED_FLAG   string = "the flag \"%s\" should be defined"

	ERROR_SESSION_NOT_CLOSED string = "Enable to complete the logout! You can force the logout manually deleting the file \"%s\""
	ERROR_NO_PROVIDED_USER   string = "No user has been provided using --%s and we cannot use the system user"
	ERROR_READING_PASSWORD   string = "Error reading the password"

	ERROR_PARSING_ENVIRONMENT string = "Error parsing the environment: %s"

	ERROR_PARSING_DESCRIPTOR string = "Error parsing the descriptor %s"

	PROMPT_RESTART  string = "Are you sure you want to recreate the starter container now (Y/N) "
	PROMPT_PASSWORD string = "Please enter password for %s:"

	LOG_FAIL_ON_PROMPT_RESTART     string = "Cannot go further if you refuse to stop the running container!"
	LOG_CONFIG_CONFIRMATION        string = "Launching lagoon CLI for %s:%s"
	LOG_FLAG_CONFIRMATION          string = "Flagged, %s %s"
	LOG_INIT_FLAGGED_DOCKER_CLIENT string = "Flagged docker client initialization..."
	LOG_INIT_DOCKER_CLIENT         string = "Docker client initialization..."
	LOG_OK_STARTED                 string = "CLI container stated!"
	LOG_GET_IMAGE                  string = "Get or refresh the latest CLI image"
	LOG_START_CREATION             string = "Starting the environment creation"
	LOG_START_UPDATE               string = "Starting the environment update"

	LOG_WAITING_DOWNLOAD   string = "waiting for the download to be completed"
	LOG_DOWNLOAD_COMPLETED string = "image download completed"
	LOG_WAITING_START      string = "waiting for the container to start"
	LOG_STARTED            string = "container started"
	LOG_WAITING_STOP       string = "waiting for the container to stop"
	LOG_STOPPED            string = "container stopped"
	LOG_COMMAND_COMPLETED  string = "Command completed"

	LOG_LOGGED_AS          string = "You are currently logged as %s on %s"
	LOG_ALREADY_LOGGED_AS  string = "You are already logged as %s on %s"
	LOG_ALREADY_LOGGED_OUT string = "You are already logged out"
	LOG_LOGOUT_REQUIRED    string = "You must to logout first to be able to create an environment"
	LOG_LOGIN_REQUIRED     string = "You must be logged in be able to perform this command"

	LOG_DEPLOYING_FROM string = "Deploying from: %s \n"
	LOG_UPDATING_FROM  string = "Updating from: %s \n"
	LOG_CHECKING_FROM  string = "Checking from: %s \n"
	LOG_GETTING_STATUS string = "Getting the status of %s \n"
	LOG_LOGGING_INTO   string = "Logging into: %s \n"

	LOG_PARSING string = "Parsing the descriptor \n"
)
