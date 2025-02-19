# Use latest Ubuntu as the base image
FROM ubuntu:latest

RUN apt update
RUN apt install -y curl git

RUN curl -L https://raw.githubusercontent.com/matter-labs/foundry-zksync/main/install-foundry-zksync | bash
