# cz-ip-service

## Introduction

A service that provides API requests using CZ88 IP Community Database as the foundation.

## Getting Keys

Get developer key, fileKey, and secretKey from [https://www.cz88.com/geo-public](https://www.cz88.com/geo-public)

1. Developer key
   ![dev_key.png](img/dev_key.png)
2. fileKey
   > When you click copy, you'll get a download link. We only need the key part from the link:  
   > `https://www.cz88.net/api/communityIpAuthorization/communityIpDbFile?fn=czdb&key=1234567890`  
   > 1234567890 is the fileKey we need

   ![file_key.png](img/file_key.png)

3. secretKey
   ![secret_key.png](img/secret_key.png)

## Software Architecture

This project is an IP address lookup service developed in Go, providing geographic location information using the CZ88 IP database. It adopts a layered architecture design, mainly consisting of the following parts:

1. **API Layer**
   - Provides interfaces in both JSON and MessagePack data formats
   - Includes single IP lookup and batch IP lookup functions
   - Web service implemented based on the Gin framework

2. **Service Layer**
   - Core business logic processing
   - IP database updating, loading, and querying
   - Supports both IPv4 and IPv6 address lookups

3. **Data Access Layer**
   - Encapsulates access to the CZ88 IP database
   - Supports file decryption and random access
   - Efficient IP address lookup algorithms

4. **Utility Layer**
   - Provides common functions such as logging, encryption/decryption, and file operations
   - Configuration management and environment variable parsing

5. **Web Interface**
   - Simple HTML frontend interface
   - Provides intuitive IP lookup functionality

System Dependencies:

- Gin framework: Provides web service and routing functionality
- MessagePack: Efficient binary serialization format for API data transfer
- Zap: High-performance logging library
- ENV: Environment variable parsing library

Data Flow:

1. Client sends an IP lookup request
2. API layer receives the request and forwards it to the service layer
3. Service layer calls the data access layer to query data from the IP database
4. Query results are returned to the client through the API layer

Project Features:

- Supports automatic IP database updates
- Provides API interfaces in multiple data formats
- Efficient memory-mapped database access
- Supports containerized deployment

## Installation Guide

1. Download dependencies

    ```shell
    go mod tidy
    ```

2. Build

    ```shell
    go build main.go
    ```

3. Start

    ```shell
    ./main -developerKey=developerKey -fileKey=fileKey -secretKey=secretKey
    ```

4. Access

   ```shell
   curl http://127.0.0.1/json?ip=1.1.1.1
   curl http://127.0.0.1
   ```

## Contribution

1. Fork this repository
2. Create a new Feat_xxx branch
3. Submit your code
4. Create a new Pull Request
