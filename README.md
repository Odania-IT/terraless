# Terraless

[![Build Status](https://travis-ci.org/Odania-IT/terraless.svg?branch=master)](https://travis-ci.org/Odania-IT/terraless)

[![codecov](https://codecov.io/gh/Odania-IT/terraless/branch/master/graph/badge.svg)](https://codecov.io/gh/Odania-IT/terraless)

Helper to deploy projects with Lambda. This helps to unify the backend and provider sections.

For this there are 2 config files. One global config file and one project specific config file.

In the project specific config you can reference parts of the global config or add specific parts.

This helps you to build serverless applications for different cloud providers.

## Why?

There are different goals i wanted to achieve:

### Cloud independent resources

I wanted to be able to provision resources for different cloud providers. Terraform has the ability to manage these
resources. That is why Terraless uses Terraform as a base.

With Terraform you are able to provision multiple clouds in one run with the same language. So you are also able to
create a DNS-Zone in Google Cloud and set the entries to a server in AWS. This makes it quite flexible.

### Manage Providers in a central place

In Terraform you need to define the providers you use in every project. This leads to duplications in provider and
backend configuration. I wanted to be able to manage them in a central place. If i need to make changes, i just have to
update the config and terraless will automatically do the rest during deployment.

### Automatically sign in to accounts

If you use multiple providers at once in the aws you need to assume every role for every provider. To keep this simple
Terraless can automatically log you in and assume these roles. It also can ask for the MFA-Token if required.

### Easier Lambda Function Deployment

Deploying a Lambda Function requires several resources in Terraform. In order to make it simple, you can define them
in the Terraless Config. This makes it easy to define the functions and events without manually specifying every Terraform
ressource.

### Uploads

You can define Uploads in the Terraless Config. Currently for AWS it can also automatically add a Cloudfront Distribution
to it.

Also the Lambda@Edge to help with the resources can be automatically deployed. For this you have two options:

* index.html as Default for "subfolders": This is helpful for static pages there you do not want the /sub/index.html in the url
* /index.html as Default for "subfolders": This is helpful for Frameworks like Angular. It will redirect every "subfolder" to the index.html in the root.

## Install

You can download the executable for your system from the [Releases page](https://github.com/Odania-IT/terraless/releases)

Currently it is available for Mac, Linux and Windows (AMD64). A Docker Image will also be created it is available on 
[https://hub.docker.com/r/odaniait/terraless](https://hub.docker.com/r/odaniait/terraless)

You can use the downloader script from [scripts/godownloader-terraless.sh](scripts/godownloader-terraless.sh)

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

# Status

Currently there is only the AWS-Provider. It currently is able to:

* Create Terraform Providers and Backend
* Create Terraform Lambda Function Resources (Create Archive, Deploy Api-Gateway)
* Auto sign in to AWS (Assume role to different accounts and ask for MFA-Token)
* Upload to S3 (Serve S3 via Cloudfront)
* Create certificates that can automatically be used with Cloudfront and Api-Gateway

# TODO

* Separately define global and local config
* More providers ;)
* Tests ;)
* YAML config: Can it be written nicer?
* Make it nicer...
* Cleanup
