#!/bin/bash

# Set the URL of the GitHub repository containing the binary.
URL="https://github.com/zcubbs/zrun/releases/latest/download/zrun_Linux_x86_64.tar.gz"

# Get the file name from the URL
FILE=$(basename $URL)

echo "Installing $FILE"

# Download the binary
curl -L -O $URL

# Unpack the binary
if [[ $FILE == *.tar.gz ]]; then
  tar -xzf $FILE
elif [[ $FILE == *.zip ]]; then
  unzip $FILE
fi

# The file that has been unpacked will usually be the binary itself.
BINARY="zrun"

# Check if the file is executable
if [[ ! -x $BINARY ]]; then
  chmod +x $BINARY
fi

# Move the binary into the PATH, so it can be executed anywhere
sudo mv $BINARY /usr/local/bin/

# Check if the binary is now in the PATH and executable
which $BINARY

if [[ $? -ne 0 ]]; then
  echo "Installation failed."
  exit 1
else
  echo "Installation succeeded."
fi

# Remove the downloaded file
rm $FILE

zrun about
