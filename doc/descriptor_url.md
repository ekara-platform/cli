# Lagoon environment descriptor URL
___

Some CLI commands require a environment descriptor. The URL of the descriptor will be passed to the CLI command through the `descriptor` argument.

### Usage: 

`lagoon create [<flags>] <descriptor> `


## Using a simple url

`$ cli create http://path.to.your.project ...`

Using the simple URL we will look for a environment descriptor file named `lagoon.yaml` into the `master` branch of your repository.

> You can use the flag `--file` in order to specify a descriptor file name other than  `lagoon.yaml` 


## Using tag or branch

Adding `@` to the URL allows you to specify a **tag** or a **branch** where to look for the descriptor 


### example:

`$ cli create http://path.to.your.project@yourTagOrBranchName ...`


Specifying a tag or a branch we will intent to checkout an environment descriptor file named `lagoon.yaml` trying within tags matching the given name first and then branches.





