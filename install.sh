#!/bin/bash

if [[ $EUID -ne 0 ]]; then
	echo "Must be root."
	exit
fi
cp bin/jwpf.o /usr/local/bin/jwpf
