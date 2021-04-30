# short-lived-links
Short lived links to share passwords, files and other stuff

## Usage

Image available at https://hub.docker.com/r/ashrath/short-lived-links

The provided docker image can be customized by bind mounting custom config file to `/opt/sll/app.yaml` or using environment variables.

Also the image defines a volume pointing to `/var/sll/pastes`

The supported environment variables are:
* SLL_TITLE (default: Short Lived Links)
* SLL_LOGO_TEXT (default: SLL)
* SLL_PASTES_PATH (default: /var/sll/pastes)

## Helm

To preview the templates generated
```shell
helm template short-lived-links helm/short-lived-links
```

Basic command to deploy into kubernetes cluster
```shell
helm upgrade short-lived-links helm/short-lived-links \
--install \
--create-namespace \
--namespace default \
--reset-values \
--set image.tag=latest
```



## Development

This repository is prepared to work out of the box in VSCode with Go extension and Dev Containers.

## CI
Github Actions is used for CI (source: https://dev.to/techschoolguru/how-to-setup-github-actions-for-go-postgres-to-run-automated-tests-81o)

## TODO
* Enable more expirations and implement expired pastes cleanup (https://github.com/go-co-op/gocron)
* Refactor UI to use vue or react
* Implement client side encryption and add option in config