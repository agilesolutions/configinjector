# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Robert Rong <robert.rong@agile-solutions.ch>"

# Set the Current Working Directory inside the container
WORKDIR /app

# first GO build and then copy this into the workdir
COPY bin/configinjector . 

# extend PATH
ENV PATH="${PATH}:/app"

# first GO build and then copy this into the workdir
RUN chmod 777 *

# dont forget to install cat, used by jenkins agents to tail a process on that container
RUN apt-get install coreutils

# Expose port 8080 to the outside world
EXPOSE 8080
