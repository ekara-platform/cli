package main

const (
	ERROR_REQUIRED_CONFIG string = "the flag \"config\" sould be definded"
	ERROR_REQUIRED_ENV    string = "the environment variable \"%s\" sould be definded"
	ERROR_REQUIRED_FLAG   string = "the flag \"%s\" sould be definded"

	PROMPT_RESTART string = "Are you sure you want to recreate the starter container now (Y/N) "

	LOG_FAIL_ON_PROMPT_RESTART     string = "Cannot go further is you refuse to stop the running container!"
	LOG_CONFIG_CONFIRMATION        string = "Launching lagoon started for %s:%s"
	LOG_FLAG_CONFIRMATION          string = "Flagged, %s %s"
	LOG_INIT_FLAGGED_DOCKER_CLIENT string = "Flagged docker client initialization..."
	LOG_INIT_DOCKER_CLIENT         string = "Docker client initialization..."
	LOG_OK_STARTED                 string = "Starter container stated!"
	LOG_GET_IMAGE                  string = "Get or refresh the latest starter image"

	LOG_WAITING_DOWNLOAD   string = "waiting for the download to be completed"
	LOG_DOWNLOAD_COMPLETED string = "image dowload completed"
	LOG_WAITING_START      string = "waiting for the container to start"
	LOG_STARTED            string = "container started"
	LOG_WAITING_STOP       string = "waiting for the container to stop"
	LOG_STOPPED            string = "container stopped"
)
