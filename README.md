# yakit [![Build Status](https://travis-ci.com/egegunes/yakit.svg?branch=master)](https://travis-ci.com/egegunes/yakit) [![Go Report Card](https://goreportcard.com/badge/github.com/egegunes/yakit)](https://goreportcard.com/report/github.com/egegunes/yakit)

## Setup

```
$ sudo docker network create yakit

# Run Postgresql container
$ sudo make db

# Import DB structure
$ sudo make initial

# Import dummy data
$ sudo make dummy

# Run yakit
$ sudo make run
```
