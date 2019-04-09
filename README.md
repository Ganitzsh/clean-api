# Payment API

## Getting started

You can start the api with the default configuration like this

    go run .

It will start listening on port 8080, or fail if this port is already in use.

## Testing

Tests can be run like this

    go test ./...

## Configuration

This program accepts YAML or JSON configuration file, it should be placed
at the root of the app directory with the following name: `config.json` or
`config.yml`.

NOTE: Support for flag will be added

Here is an example:

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
        -  "\*"
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

### Environment variable

It also supports the following environment variable:

  - API_DEV_MODE: `bool`
  - API_TLS_ENABLED: `bool`
  - API_TLS_KEY: path to key file
  - API_TLS_CERT: path to cert file
  - API_CORS_ORIGINS: `CSV`
  - API_CORS_HEADERS: `CSV`
  - API_CORS_METHODS: `CSV`
  - API_DATABASE_TYPE: `mongo` | `inmem`
  - API_DATABASE_MONGO_DATABASE: `string`
  - API_DATABASE_MONGO_COLLECTION: `string`
  - API_DATABASE_MONGO_URI: `string`


## Storage

Two different storage types are available:

-   `inmem`: in memory storage, no persistency
-   `mongo`: a MongoDB storage

The API will use `inmem` by default if no mongo configuration is given.
This storage is not persistent and will disappear when the program will exit.

## API

### Response format

All the responses sent by the API will be of the form:

    {
      "data": {},
      "code": 200,
      "status": "success"
    }

Where `status` can have the following values:

-   `fail`: The error was due to a user input
-   `error`: System wide error
-   `success`: Everything was fine

### Error

Errors contain the following information

    {
      "data": {
        "error": "Invalid input",
        "code": "invalid_input"
      },
      "code": 400,
      "status": "fail"
    }

#### Codes

-   `invalid_input`: Get returned when the user input is invalid
-   `internal_error`: Happens when something broke internally while processing the request
-   `not_found`: Means either that a resource was not found or the route does not exists
-   `undergoing_maintenance`: Means the whole service is not available
-   `not_implemented`: The feature is not implemented yet

### Entities

#### Payment

The main resource of the API

    {
      "id": "97122344-dc12-41e0-a81a-39c234ae7449",
      "amount": "string",
      "beneficiary": {
        "accountName": "string",
        "accountNumber": "string",
        "accountNumberCode": "string",
        "address": "string",
        "bankId": "string",
        "bankIdCode": "string",
        "name": "string"
      },
      "chargesInformation": {
        "bearerCode": "string",
        "receiverChargesAmount": "string",
        "receiverChargesCurrency": "string",
        "senderCharges": [
          {
            "amount": "string",
            "currency": "string"
          }
        ]
      },
      "createdAt": "2019-03-14T09:33:18.982Z",
      "currency": "string",
      "debitorParty": {
        "accountName": "string",
        "accountNumber": "string",
        "accountNumberCode": "string",
        "address": "string",
        "bankId": "string",
        "bankIdCode": "string",
        "name": "string"
      },
      "endToEndReference": "string",
      "fx": {
        "contractReference": "string",
        "exchangeRate": "string",
        "originalAmount": "string",
        "originalCurrency": "string"
      },
      "numericReference": "string",
      "processingDate": "string",
      "purpose": "string",
      "reference": "string",
      "scheme": "string",
      "schemePaymentSubType": "string",
      "schemePaymentType": "string",
      "type": "string",
      "updatedAt": "2019-03-14T09:33:18.982Z"
    }

### Endpoints

All Endpoints require the api prefix: `/v1`

#### Pagination

The API supports pagination for some endpoints. The following parameters are
available and must be placed un the query (e.g: `/payments?lim=10&page=2`):

-   `lim`: The maximum amount of resource to return
-   `off`: The offset after which you'd like to start
-   `page`: The page number wanted, works with `lim` only if `lim` is > 0
otherwise it's ignored. If `off` is also specified, `page` will be ignored and
`off` will be used.

##### Payload

The paginated routes have the following structure:

    {
      "total": 3,
      "subTotal": 1,
      "data": [{}]
    }

#### List

|     Method    | URI              |   Body  |      Response     | Paginated | Description                |
| :-----------: | ---------------- | :-----: | :---------------: | :-------: | -------------------------- |
|     `GET`     | `/payments`      |   None  | `200` `[]Payment` |    `X`    | Gets a list of payments    |
|     `GET`     | `/payments/{id}` |   None  |  `200` `Payment`  |    `-`    | Retrieves a single payment |
|     `POST`    | `/payments`      | Payment |  `201` `Payment`  |    `-`    | Creates a new payment      |
| `POST`, `PUT` | `/payments/{id}` | Payment |  `200` `Payment`  |    `-`    | Edit a payment             |
|    `DELETE`   | `/payments/{id}` |   None  |    `204` Empty    |    `-`    | Delete a payment           |
