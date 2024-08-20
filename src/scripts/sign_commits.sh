#!/bin/bash

# read two params username and email
username=$1
email=$2

# set the git config
git config user.name "$username"
git config user.email "$email"

#  validte key exists
if [ ! -f ~/.ssh/id_ed25519 ]; then
    echo "No SSH key found, please generate one and try again"
    exit 1
fi

eval "$(ssh-agent -s)"
ssh-add ~/.ssh/id_ed25519
git config gpg.format ssh
git config user.signingkey ~/.ssh/id_ed25519
git config commit.gpgsign true
