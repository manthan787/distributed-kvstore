# distributed-kvstore
A distributed in-memory key-value store.


## Design

The project includes two main components:

- Proxy/Router written in golang

- Server written in Python

**Proxy** is responsible for routing client requests for adding and fetching key-value pairs to servers based on key hashes.

**Server** code is responsible for spinning up new http servers and enabling servers to store key-value pairs in-memory.


## Running the code

Install `python 2.7` (required for server)

Install `golang` (required for proxy)

Run `make install` to install dependencies and `make check` to run example application with a few servers spawned on localhost.