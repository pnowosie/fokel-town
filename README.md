# FokelTown: open the source that should stay closed

## Description 

Write a Go API service that allows user registration. Each user will have a unique identifier, first name, and last name. Use an in-memory Merkle Trie to store the user data, eliminating the need for an additional database.


## Requirements

1. Create an in-memory Merkle Trie data structure to store user data. Each user should be identified by a unique identifier (make sure it's your own merkle trie implementation (we're skipping the one from go-ethereum)).

2. Develop an API endpoint that enables the registration of a new user. The endpoint should receive user data (identifier, first name, last name) and store it in the in-memory Merkle Trie.

3. Implement an API endpoint that allows retrieving user data based on their identifier. The endpoint should search the in-memory Merkle Trie to find the user data and return it as the API response.


## How to run

==TODO==

Program accepts following command line arguments:
- host - host to listen on, default is `localhost`
- port - port to listen on, default is `4000`

Without arguments, service will start listening on `localhost:4000`.


## Exposed endpoints

### `GET /health`

```bash
http :4000/v0/health
```

**Response:**

```json
{
  "name": "merkle-service",
  "root": "0000000000000000000000000000000000000000000000000000000000000000",
  "uptime": 123,
  "version": "0.0.1"
}
```

**Status codes:**
- `200` - service is healthy

### `GET /v0/user/:id`

```bash
http :4000/v0/user/beef0c
```

**Response:**

```json
{
    "id": "beef0c",
    "firstname": "John",
    "lastname": "Doe"
}
```
**Status codes:**
- `200` - user found
- `400` - invalid user id
- `404` - user not found
- `500` - internal server error


### `PUT/v0/user`

```bash
http PUT :4000/v0/user id=beef0c firstname=John lastname=Doe
```

**Request Body:**

```json
{
    "id": "beef0c",
    "firstName": "John",
    "lastName": "Doe"
}
```
**Response:**

No response body


**Status codes:**
- `201` - user entry created
- `302` - user with given id already exists
- `400` - invalid user data
- `500` - internal server error
