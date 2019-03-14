# Simplest Configuration

This is a simple example with global and project config.

[Project Configuration: terraless-project.yml](terraless-project.yml)
[Global Configuration: terraless.yml](terraless.yml)

## Deploy infrastructure

You can deploy the infrastructure with:

    terraless --config examples/aws-simple/terraless-project.yml --global-config examples/aws-simple-with-global-cfg/terraless.yml deploy

Every terraform file you would add to the folder would automatically be executed.

Depending on the environment it choose the correct provider from the global config.

