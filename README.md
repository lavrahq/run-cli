# Lavra CLI

A CLI for deploying, configuring, and extending Lavra products.

# Commands

The following commands are prefixed with `lavra`.

## Basics

`start`                         Provides a quick start wizard which will guide the user through setting up their first local context.
`version`                       Provides version information for the CLI. Displays the latest version as well as the version installed.
`update`                        Updates the CLI to the latest version, or, exits stating the current version is the latest.
`issue`                         Allows creating a Github issue from the command line using a wizard.

## Authentication (Not Yet Implemented)

Authentication commands allow the user to authenticate with Lavra SSO. This allows the user to interact with Cloud options when
they become available.

`auth login`                    Login to Lavra Login.
`auth logout`                   Logout of the CLI and terminate the Lavra Login session.
`auth whoami`                   Informs the user who the currently logged in user is.

## Managing Contexts

Contexts provide a means for configuring where the Lavra CLI tool connects in order to deploy, configure, and administer
existing Lavra products, or new products that have not been deployed. By default, when `start` is ran, a `local` context is
created that points to the local Docker engine or Kubernetes cluster.

`contexts`                      Lists the available contexts, with an indicator on which is the current context.
`contexts ls`                   Alias to the above.
`contexts add <name>`           Adds a new context
`contexts update <name>`        Updates the named context with new settings.
`contexts test <name>`          Tests whether connectivity to the context is working.
`contexts switch <name>`        Switches to the names context.
`context <name>`                Alias to the above.

## Managing Projects

Projects are the deployed services that the user can manage. Since Projects will be managed via source control so that they are
stateless (aside from the database which will hold state), projects will be a directory full of files and a config.yml file.

`projects`                      Lists the projects currently managed by the CLI tool.
`projects ls`                   Alias to the above.
`projects new <dir=.>`          Creates a project at the directory specified. Defaults to the current directory.
`projects import <dir=.>`       Imports an existing project that was created manually or from a different user, such as a Git-versioned project in Github that the user cloned.
`projects open <dir=.>`         Opens the current project or the project at the path specified in the default editor.
`projects remove <dir=.>`       Removes the project from the CLI, but keeps the files.
`projects destroy <dir=.>`      Removes the project from the CLI and removes the files.

## Deployments

Deployments are projects that are deployed to Docker Engine or a Kubernetes Cluster. These commands must be ran in the project they are intended for.

`deploy`                        Deploys the project, or, updates the existing project on the Docker Engine or Kubernetes Cluster.
`deploy undo`                   Undeploy an existing project, will retain data and secrets unless -e or --everything is provided.
`deploy open`                   Opens the deployment in a web browser, for example, the Response Console, or the Identity app.
`deploy version`                Gets the current version of the deployment. Will tell the user if the version is out of date.
`deploy logs`                   Will subscribe to the logs of the deployment, outputting them to the console as they happen. (provided by OpenFaaS)

`up`                            Alias to `deploy`
`down`                          Alias to `deploy undo`                          
