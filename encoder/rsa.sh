#!/usr/bin/env zsh

set -e

# PRIVATE KEY
openssl genrsa -out id_rsa 1024

# PUBLIC KEY
openssl rsa -in id_rsa -pubout -out id_rsa.pub
