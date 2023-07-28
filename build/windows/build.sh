#!/bin/bash

# Ensure the script is run without superuser permissions
if [[ $EUID -eq 0 ]]; then
   echo "This script should not be run as root"
   exit 1
fi

# Ensure the Go environment is properly set up
if ! [ -x "$(command -v go)" ]; then
  echo 'Error: Go is not installed.' >&2
  exit 1
fi

# Set the GoOS and GoARCH for Windows
export GOOS=windows
export GOARCH=amd64

# Specify your Go file's name here
go_file="main.go"

# Ensure that ~/bin exists
mkdir -p $HOME/bin

# Build the Go file
go build -o $HOME/bin/pc.exe $go_file

# If successful, notify the user that the binary is ready
if [ $? -eq 0 ]; then
    echo "Build successful. Please run 'pc.exe help'"
else
    echo "Build failed."
fi

# Ensure that ~/bin is in the PATH
if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
    if ! grep -q "$HOME/bin" "$HOME/.zshrc"; then
        echo 'export PATH="$HOME/bin:$PATH"' >> $HOME/.zshrc
        echo "PATH has been updated in ~/.zshrc. Please run 'source ~/.zshrc' or start a new terminal session for the changes to take effect."
    fi
fi
