# FokelTown: open the source that should stay closed

## Description 

Write a Go API service that allows user registration. Each user will have a unique identifier, first name, and last name. Use an in-memory Merkle Trie to store the user data, eliminating the need for an additional database.


## Requirements

1. Create an in-memory Merkle Trie data structure to store user data. Each user should be identified by a unique identifier (make sure it's your own merkle trie implementation (we're skipping the one from go-ethereum)).

2. Develop an API endpoint that enables the registration of a new user. The endpoint should receive user data (identifier, first name, last name) and store it in the in-memory Merkle Trie.

3. Implement an API endpoint that allows retrieving user data based on their identifier. The endpoint should search the in-memory Merkle Trie to find the user data and return it as the API response.