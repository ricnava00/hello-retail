# WILL.IAM
Workflow Integration Alleviates Identity and Access Management

## About
WILL.IAM is a workflow-aware access control model and reference monitor that satisfies the functional requirements of the serverless computing paradigm.

## Using this Repository
This repository contains the code for building OpenFaaS with WILL.IAM and deploying a sample application (Hello, Retail!). It is possible to use this repository to build the platform components and sample application, alternativley it is possible to deploying the platform and sample application using prebuilt images. If you would like to use the pre-built images skip down to the deployment section.

### Build Process
In order to build the platform and sample application simply run `./build.sh` from the root of the repository. You will need to be connected to the internet, as the build process will pull source code from the OpenFaaS platform and patch WILL.IAM into that codebase. Additionally you will need to have `docker` installed, as it is used heavily during the build process. Once you build process has finished, you can run `./retag_dev_to_prod.sh` from the root of the repository. Running this script will make the deployment scripts utilizing the locally built images instead of pulling the pre-built images at deployment time.

### Deployment Process
In order to deploy the platform and sample application you will need a `kubectl` configuration which can connect to a kubernetes cluster. Additionally you will need to install `faas-cli` from the OpenFaaS project. Once read to deploy, `cd` into the `kubernetes/` directory and run `./deploy.sh`. Once the script has finished, you interact with the running application. You can look at the last few lines of the deployment script for some sample `curl` commands to run.

### Teardown Process
Once finished with running the deployment, you can fully cleanup the platform and sample application by running `./teardown.sh`.
