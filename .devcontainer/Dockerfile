FROM mcr.microsoft.com/devcontainers/go:0-1.19
# Install the xz-utils package
RUN apt-get update && apt-get upgrade -y
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
  && unzip awscliv2.zip \
  && sudo ./aws/install

RUN curl -fsSL https://deb.nodesource.com/setup_19.x | sudo -E bash - &&\
  sudo apt-get install -y nodejs

RUN npm install -g aws-cdk