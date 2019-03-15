# Terraless

[![Build Status](https://travis-ci.org/Odania-IT/terraless.svg?branch=master)](https://travis-ci.org/Odania-IT/terraless)

[![codecov](https://codecov.io/gh/Odania-IT/terraless/branch/master/graph/badge.svg)](https://codecov.io/gh/Odania-IT/terraless)

Helper to deploy projects with Lambda. This helps to unify the backend and provider sections.

For this there are 2 config files. One global config file and one project specific config file.

In the project specific config you can reference parts of the global config or add specific parts.

This helps you to build serverless applications for different cloud providers.

## Configs

Example Global Config: [examples/terraless.yml](examples/terraless.yml)
Example Project Config: [examples/terraless-project.yml](examples/terraless-project.yml)

### Global Config

The global config will be automatically looked up relative to the home folder in the following locations:

		.terraless
		.config/.terraless
		.aws",
		.config/gcloud

## Execution

If you simply call "terraless" the help will be displayed.

terraless -config examples/terraless-project.yml -global-config examples/terraless.yml

### Command: deploy

This command does:

* Create Templates
* Prepares Sessions
* Packages Archive (Lambda)
* Deploys Terraform
* Processes all uploads

### Command: init

This created the Terraless Terraform Templates for the given config. Two files will be generated:

    terraless-main.tf        <-- The generated Terraform Provider and Backend will be in here
    terraless-project.yml    <-- All other resources for certificates, endpoints, functions, ... will be in this file

### Command: session

This will verify access to all providers. If AutoLogin is enabled in Settings it will also try to authenticate if there
is currently no access.

Currently for the AWS-Terraless-Provider it can also ask for your MFA-Token and assume the corresponding roles for the
providers.

### Command: upload

This executes the session command and afterwords processes all defined uploads.

### Command: version

Simply displays the version.

# TODO

* Separately define global and local config
* More providers ;)
* Tests ;)
* YAML config: Can it be written nicer?
* Make it nicer...
* Cleanup
