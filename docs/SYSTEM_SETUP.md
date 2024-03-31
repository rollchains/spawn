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
```

# Ubuntu

```bash
# Base
sudo apt-get install make gcc

# Github CLI - https://github.com/cli/cli
curl -sS https://webi.sh/gh | sh
gh auth login

# Golang
wget https://go.dev/dl/go1.22.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
export PATH=$PATH:/usr/local/go/bin

# Docker
sudo apt -y install docker.io
```