# gonfig [![GoDoc](https://godoc.org/github.com/gonfig/axkit?status.svg)](https://godoc.org/github.com/axkit/gonfig) [![Build Status](https://travis-ci.org/axkit/gonfig.svg?branch=master)](https://travis-ci.org/axkit/gonfig) [![Coverage Status](https://coveralls.io/repos/github/axkit/gonfig/badge.svg)](https://coveralls.io/github/axkit/gonfig) [![Go Report Card](https://goreportcard.com/badge/github.com/axkit/gonfig)](https://goreportcard.com/report/github.com/axkit/gonfig)

Configuration for Go Application Changing in Runtime

## Motivation
There are parameters types:
* **Static program execution parameters** requires application restart to apply new values. 
* **Dynamic program execution parameters** does not require application restart, but getting values requires syncronization.
* **User settings** are the same as Dynamic


A parameter can be declared in several places and simultaneously:
* hardcoded default values (1)
* config file (2) 
* environment variable (3)
* command line (4)

Values of parameters captured on higher step overwrites previous values.  

