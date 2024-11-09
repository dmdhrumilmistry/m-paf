# M-PAF

![M-PAF-logo](./.static/images/m-paf-logo.png)

Malicious-PAckageFinder (m-paf), The CLI tool for discovering malicious and risky packages using SBOM files.

## Demo

<https://github.com/dmdhrumilmistry/m-paf/raw/refs/heads/main/.static/videos/m-paf.mp4>

## Installation

#### Github Hosted Method

- Install latest release using below command

  ```bash
  go install -v github.com/dmdhrumilmistry/m-paf@latest
  ```

- Install main/dev branch

  ```bash
  go install -v github.com/dmdhrumilmistry/m-paf@main # install main branch
  go install -v github.com/dmdhrumilmistry/m-paf@dev  # install dev branch
  ```

#### Clone Method

- Clone repository

    ```bash
    git clone https://github.com/dmdhrumilmistry/m-paf
    ```

- Run Go install command

    ```bash
    go install ./...
    ```

## Using M-PAF

* Print Help

    ```bash
    m-paf -h
    ```

* Basic Usage

    ```bash
    m-paf -f sbom.jsom
    ```

### Open In Google Cloud Shell

- Temporary Session

  [![Open in Cloud Shell](https://gstatic.com/cloudssh/images/open-btn.svg)](https://shell.cloud.google.com/cloudshell/editor?cloudshell_git_repo=https://github.com/dmdhrumilmistry/m-paf.git&ephemeral=true&show=terminal&cloudshell_print=./DISCLAIMER.md)

- Perisitent Session

  [![Open in Cloud Shell](https://gstatic.com/cloudssh/images/open-btn.svg)](https://shell.cloud.google.com/cloudshell/editor?cloudshell_git_repo=https://github.com/dmdhrumilmistry/m-paf&ephemeral=false&show=terminal&cloudshell_print=./DISCLAIMER.md)

## Have any Ideas 💡 or issue

Create an issue *OR* fork the repo, update script and create a Pull Request

## Contributing

Refer [CONTRIBUTIONS.md](/CONTRIBUTING.md) for contributing to the project.

## LICENSE

Tool is distributed under `MIT` License. Refer [License](/LICENSE.md) for more information.
