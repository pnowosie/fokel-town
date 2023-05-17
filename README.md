# FokelTown: open the source that should stay closed

## Description 

Write a Go API service that allows user registration. Each user will have a unique identifier, first name, and last name. Use an in-memory Merkle Trie to store the user data, eliminating the need for an additional database.


## Requirements

1. Create an in-memory Merkle Trie data structure to store user data. Each user should be identified by a unique identifier (make sure it's your own merkle trie implementation (we're skipping the one from go-ethereum)).

2. Develop an API endpoint that enables the registration of a new user. The endpoint should receive user data (identifier, first name, last name) and store it in the in-memory Merkle Trie.

3. Implement an API endpoint that allows retrieving user data based on their identifier. The endpoint should search the in-memory Merkle Trie to find the user data and return it as the API response.


## How to run

### Run in docker

Here is how to build an image and run the container:
```bash
 GIT_SHA=$(git rev-parse --short HEAD)
 docker build . --progress=plain -t merkle-srv:${GIT_SHA} -t merkle-srv:latest
```
Tagging with latest isn't the best practice, but it's convenient for local tests.

Once the image is ready, run the container:
```bash
docker run --rm -p 4000:4000 merkle-srv
```

### Build from sources

If you have `go` version `1.19` installed, you can run service with following command executed from the project root directory:

```bash
go run ./cmd/api -host localhost -port 4000
```

You might want to run all test first:

```bash
go test -v ./...
```
or for the better output:
```bash
gotestsum -f testname
```

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

## Trie implementation notes

The trie implemented in this service is Merkle Patricia Trie, where a variant of Radix trie also delivers cryptographic data authentication.

In my simplified implementation there are 2 types of nodes:
- _branch node_ - contains array of 16 child nodes, branch or leaf. Array elements correspond to 4-bit nibbles of the key.
- _leaf node_ - contains a key and a pointer to the user data.

The root node of a trie is Trie-structure itself and has a pointer to the "root" branch node. To safe space and limit trie height, branch nodes below the "root" one contains a path _prefix_ which holds a common key substring for all its children. 

**Trie diagram overview**
![Ptricia trie example](patricia_trie.png "Source: https://www.youtube.com/watch?v=QlawpoK4g5A")


I artificially limited key length to 3-bytes hex string, which still allows to store more than 16.7M entries.  This limit is made mainly for demonstration and service usability purposes, and can be easily repealed by changing `UserData.IsValid()` method.
