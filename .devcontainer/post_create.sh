#!/bin/bash

# Define the path to your Zsh profile
zshrc_path="$HOME/.zshrc"
bashrc_path="$HOME/.bashrc"

echo "export PATH=\"$HOME/.local/bin:$PATH\"" >> "$zshrc_path"
echo "export PATH=\"$HOME/.local/bin:$PATH\"" >> "$bashrc_path"

cat "$HOME"/.zshrc
export PATH="$HOME/.local/bin:$PATH"

# install the requirements
# pip install --upgrade pip
pip install --user -r requirements-dev.txt

# install the pre-commit hooks
pre-commit autoupdate
pre-commit install

# install the npm packages
npm install -g cspell
npm install -g markdownlint-cli2
