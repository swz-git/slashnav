# Slashnav

Terminal folder navigator inspired by https://github.com/souvikinator/lsx.

## Install

Currently there isn't any easy way to install. You need to compile main.go into a binary called slashnav and put that into your path (could be done using `go install` in a cloned repo), then add `source path/to/your/script.sh` to your .bashrc or similar. You should now be able to use the `slash` command!

## Todo

- [x] Fix scrolling in large directories
- [ ] Better command line with options etc