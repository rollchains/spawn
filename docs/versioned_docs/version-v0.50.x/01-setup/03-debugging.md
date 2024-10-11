---
title: Debugging
sidebar_label: Debugging
sidebar_position: 3
slug: /install/debugging
---

This section will contain common setup problems and how to resolve them.

## Golang

### /bin/sh: 1: go: not found

Just add the following lines to `~/.bashrc` (or `~/.zshrc` if MacOs) and this will persist. [Source](https://stackoverflow.com/a/21012349)
If you run the above in your terminal, it will apply to the current session but not on new terminal sessions.

```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

Then apply it with `source ~/.bashrc` or `source ~/.zshrc`

### build constraints excluded all Go files in /usr/local/go/ ...

Your Go install is not properly setup. Follow the install instructions above or install directly from source with [go.dev](https://go.dev/doc/install).

### make: heighliner: Permission denied

```bash
make get-heighliner
chmod +x $(which heighliner)
```

If the above does not work, your user or directory permissions may not be setup. Or your `ls -la $(go env GOPATH)/bin` path is to a bad.

If using WSL, try https://superuser.com/questions/1352207/windows-wsl-ubuntu-sees-wrong-permissions-on-files-in-mounted-disk.

---

## Windows / WSL

### make: /mnt/c/Program: No such file or directory

Delete your GOMODCACHE directory: `go clean -modcache` or run the direct command `rm -rf $(go env GOMODCACHE)`.

---

## Docker

### Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?.

Start the docker daemon. Run [docker engine](https://docs.docker.com/engine/) or `systemctl start docker && systemctl enable docker` for Linux.

### docker: Got permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock

You don't have permissions to interact with the Docker daemon.

1) Install properly with https://docs.docker.com/get-started/get-docker/

2)
```bash
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker

reboot # if you still get the error
```

Technically you can also `sudo chmod 666 /var/run/docker.sock` but this is NOT advised. -->

## Generation

### remote: Repository not found. fatal: reposity not found

This error is due to not having properly `make proto-gen`ed the project. View the [Application](#running-the-binary-gives-me-panic-reflect-newnil) section for the solution.

## Application

### Running the binary gives me `panic: reflect: New(nil)`

The `make proto-gen` command was either not run, or is causing issues. This could be due to your users permissions or the filesystem. By default, the protoc docker image uses your current users id and group. Try switching as a super user (i.e. `su -`) or fixing your permissions. A very ugly hack is to run `chmod a+rwx -R ./rollchain` where `./rollchain` is the project you generated. This will cause git to change all files, but it does fix it. Unsure of the long term side effects that may come up from this.
