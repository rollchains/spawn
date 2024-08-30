# MacOS Setup

```bash
# Base
brew install make
brew install gcc

# Github CLI - https://github.com/cli/cli
brew install gh
gh auth login

# Golang
brew install go

# Docker
brew install docker

# Continue to main README.md to install spawn & local-ic
```

## Windows Setup

```bash
# Install WSL in powershell
wsl --install
reboot

# Setup WSL Ubuntu Image
wsl.exe --install Ubuntu-24.04

# Open wsl instance
wsl

# update and add snap if not already installed
sudo apt update && sudo apt install snapd

# Install Go (Snap)
sudo snap info go
sudo snap install go --channel=1.23/stable --classic

# Install Basics
sudo apt install make gcc git jq

# Install github-cli
sudo snap install gh

# Install docker
sudo snap install docker
sudo chmod 666 /var/run/docker.sock

# Setup base github config
git config --global user.email "yourEmail@gmail.com"
git config --global user.name "Your Name"

# Fix the path
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
export PATH=$PATH:$(go env GOPATH)/bin

# Continue to main README.md to install spawn & local-ic
```

# Ubuntu Setup

```bash
# Base
sudo apt install make gcc git

# Github CLI - https://github.com/cli/cli
curl -sS https://webi.sh/gh | sh
gh auth login

# Golang
GO_VERSION=1.23.0
wget https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz

# fix paths
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
export PATH=$PATH:$(go env GOPATH)/bin

# Docker
sudo apt -y install docker.io

# Setup base github config
git config --global user.email "yourEmail@gmail.com"
git config --global user.name "Your Name"

# Continue to main README.md to install spawn & local-ic
```
