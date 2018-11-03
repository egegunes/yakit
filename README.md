# yakit

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
