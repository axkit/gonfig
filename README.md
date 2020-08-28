# gonfig [![GoDoc](https://godoc.org/github.com/gonfig/axkit?status.svg)](https://godoc.org/github.com/axkit/gonfig) [![Build Status](https://travis-ci.org/axkit/gonfig.svg?branch=master)](https://travis-ci.org/axkit/gonfig) [![Coverage Status](https://coveralls.io/repos/github/axkit/gonfig/badge.svg)](https://coveralls.io/github/axkit/gonfig) [![Go Report Card](https://goreportcard.com/badge/github.com/axkit/gonfig)](https://goreportcard.com/report/github.com/axkit/gonfig)

Configuration for Go Application  

## Motivation
Usually, there are static and dynamic program execution parameters.
Static parameters requires application restart to apply new values.
New values of dynamic parameters can be applied without application 
restart, but getting values requires access sync.


A parameter can be declared in several places and simultaneously:
* global variables (1)
* config file (2) 
* environment variable (3)
* command line (4)

Values of parameters re-declared with higher steps overwrites previous values.  

