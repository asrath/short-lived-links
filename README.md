# short-lived-links
Short lived links to share passwords, files and other stuff

## Usage

The provided docker image can be customized by bind mounting custom config file to `/opt/sll/app.yaml` or using environment variables.

The supported environment variables are:
* SLL_TITLE (default: Short Lived Links)
* SLL_LOGO_TEXT (default: SLL)
* SLL_PASTES_PATH (default: /var/sll/pastes)

## Development

This repository is prepared to work out of the box in VSCode with Go extension and Dev Containers.

## TODO
* Enable more expirations and implement expired pastes cleanup
* Refactor UI to use vue or react
* Implement client side encryption and add option in config