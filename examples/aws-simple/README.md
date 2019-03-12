# Simplest Configuration

The simplest config just has the terraless-project.yml defined.

[Project Configuration: terraless-project.yml](terraless-project.yml)

## Deploy infrastructure

You can deploy the infrastructure with:

    terraless --config examples/aws-simple/terraless-project.yml deploy

Every terraform file you would add to the folder would automatically be executed.
