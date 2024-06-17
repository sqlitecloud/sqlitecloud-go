FROM mcr.microsoft.com/devcontainers/go:1-1-bullseye

# Python environment
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
&& apt-get -y install python3-pip git zlib1g-dev unzip

# Pre-Commit to run pre-commit tools
RUN pip3 install pre-commit
