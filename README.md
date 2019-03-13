# Payment API

## Getting started

You can start the api with the default configuration like this

    go run .

It will start listening on port 8080, or fail if this port is already in use.

## Configuration

This program accepts YAML, JSON or XML configuration file.

Here is an example:

```
# Set to true if running in dev mode, it will output additional logs
dev_mode: false

# Simple tls configuration, set enabed to true with a proper key-cert pair and
# the api will start running and will be accessible over HTTPS
tls:
  enabled: true
  key:     'key.pem'
  cert:    'cert.pem'

# This part is to tweak the CORS headers to suit your needs
cors:
  origins:
    - http://localhost
  headers:
    - '*'
  methods:
    - GET
    - POST
    - PUT
    - OPTIONS

# Here is where all the database related configuration is
database:

  # You have two different types of available storage: 'inmem' and 'mongo'
  type: mongo

  mongo:
    database: api
    collection: payments
    uri: user:password@localhost
```

## Storage

Two different storage types are available:

  - `inmem`: in memory storage, no persistency
  - `mongo`: a MongoDB storage

## Endpoints
