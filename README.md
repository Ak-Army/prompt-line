# Prompt-line
A cross-shell themeable prompt heavily inspired by
* [oh-my-posh](https://github.com/JanDeDobbeleer/oh-my-posh)
* [powerline-go](https://github.com/justjanne/powerline-go)
* [bronze](https://github.com/reujab/bronze)

# Installation

## Recommanded
Prompt-line was designed to use [Nerd Fonts](https://www.nerdfonts.com/). 
Nerd Fonts are popular fonts that are patched to include icons. 
We recommend [Meslo LGM NF](https://github.com/ryanoasis/nerd-fonts/releases/download/v2.1.0/Meslo.zip), 
but any Nerd Font should be compatible with the standard themes.

##From source
* install and setup [Go](https://go.dev/)
* run `go install github.com/Ak-Army/prompt-line@latest`

## From pre-compiled binary
* download a binary on the [releases page](https://github.com/Ak-Army/prompt-line/releases)
* add binary to `PATH` environment variable

# Configuration

## Themes
* download themes from [releases page](https://github.com/Ak-Army/prompt-line/releases)
* unzip it 

## BASH
Add the following to `.bashrc` or `.profile`.
```bash
eval "$(~/prompt-line init -shell=bash -config=${HOME}/theme/default.toml)"
```
Once added, reload your profile  or bashrc for the changes to take effect.
```bash
source ~/.profile
source ~/.bashrc
```

## ZSH
Add the following to `.zshrc`.
```bash
eval "$(~/prompt-line init -shell=zsh -config=${HOME}/theme/default.toml)"
```
Once added, reload your config for the changes to take effect.
```bash
source ~/.zshrc
``````

## ZSH
Add the following to `~/.config/fish/config.fish`.
```bash
~/prompt-line init -shell fish -config ${HOME}/theme/default.toml) |source
```
Once added, reload your config for the changes to take effect.
```bash
. ~/.config/fish/config.fish
```
