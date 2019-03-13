# Terraless

[![Build Status](https://travis-ci.org/Odania-IT/terraless.svg?branch=master)](https://travis-ci.org/Odania-IT/terraless)

[![codecov](https://codecov.io/gh/Odania-IT/terraless/branch/master/graph/badge.svg)](https://codecov.io/gh/Odania-IT/terraless)

Helper to deploy projects with Lambda. This helps to unify the backend and provider sections.

For this there are 2 config files. One global config file and one project specific config file.

In the project specific config you can reference parts of the global config or add specific parts.

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

terraless -config examples/terraless-project.yml -global-config examples/terraless.yml

# TODO

* Separately define global and local config
* More providers ;)
* Tests ;)
* YAML config: Can it be written nicer?
* Make it nicer...
* Cleanup
