# PoW Shield

## Overview

PoW Shield is a proof-of-work based rate limiting service that helps protect APIs from abuse and DDoS attacks. It uses the Scrypt algorithm for proof-of-work verification.

The service requires clients to solve a computational puzzle before accessing protected endpoints. This creates a computational cost for each request, effectively deterring automated attacks while allowing legitimate traffic through.

### Architecture

![Architecture](./arch.png)

### Why Scrypt?

[Scrypt Wikipedia](https://ru.wikipedia.org/wiki/Scrypt)

Scrypt was chosen as the proof-of-work algorithm because:

- It is designed to be computationally expensive and memory-intensive
- It is configurable and allows to tune the difficulty flexibly
- In future mechanism can be improved by using module for decision trusted and untrusted clients by checking IP address and request rate and we can use less CPU and memory for proof-of-work verification for trusted clients, to decrease the load on client side.
- Makes it much harder to solve quickly on specialized hardware like GPUs or ASICs because it requires more memory to compute the hash
- Provides strong protection against parallel brute-force attacks
- Well-tested and cryptographically secure

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Building and Running with Docker Compose

```bash
docker-compose up --build
```
