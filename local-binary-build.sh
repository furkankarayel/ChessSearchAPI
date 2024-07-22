#!/bin/bash

# Create the bin directory if it doesn't exist
mkdir -p bin

# Pull pgn-extract and build binary
git clone https://github.com/MichaelB7/pgn-extract.git bin/pgn-extract
cd bin/pgn-extract/src
make
cd ../../..

# Pull scoutfish and build binary
git clone https://github.com/mcostalba/scoutfish.git bin/scoutfish
cd bin/scoutfish/src
make build ARCH=x86-64
cd ../../..
