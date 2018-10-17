# Ekara environment descriptor URL
___

Some CLI commands require an environment descriptor. The URL of the descriptor will be passed to the CLI command through the `descriptor` argument.

### Usage: 

`ekara create [<flags>] <descriptor> `


## Using a simple url

`$ cli create http://path.to.your.project ...`

Using the simple URL we will look for a environment descriptor file named `ekara.yaml` into the `master` branch of your repository.

> You can use the flag `--file` in order to specify a descriptor file name other than  `ekara.yaml` 


## Using tag or branch

Adding `@` to the URL allows you to specify a **tag** or a **branch** where to look for the descriptor 


### example:

`$ cli create http://path.to.your.project@yourTagOrBranchName ...`


Specifying a tag or a branch we will intent to checkout an environment descriptor file named `ekara.yaml` trying within tags matching the given name first and then branches.





