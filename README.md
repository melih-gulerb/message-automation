# Message Automation Service

## Overview

The **Message Automation Service** is a Go-based project which employed to automate 
the sending of messages per specific period of time and specific count of messages.

## Tech Stack

- **Go**: The primary programming language used for the service.
- **MSSQL**: Relational database used for data storage.
- **Redis**: In-memory data structure store used for caching.

## Features

- **Automated Message Sending**: Unsent messages are sent asynchronously using goroutines.
- **Project Deployment Ingestion**: All unsent messages are sent when the system initialized.
- **Transactional Database Operations**: Ensures data consistency with MSSQL.
- **Redis Integration**: Caches messageId's and their timestamps.

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/melih-gulerb/message-automation.git
   cd message-automation
   ```
2. **Install dependencies**
   ```bash
   go mod tidy
   ```
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   cd src
   swag init
   ```
3. **Initialize**
   ```bash
   go run ./main.go
   ```

## Running with Docker
```bash 
docker build -t message-automation . 
```
```bash 
docker run -p 3030:3030 -e message-automation 
```


## Endpoints

### 1. Automate Message Sending

- **Endpoint**: `/message/automation`
- **Method**: `POST`
- **Query Parameters**:
    - `isActive`:
- **Summary**: Starts or stops the message automation based on the `isActive` status.

#### Example:

```bash
curl -X POST "http://localhost:3030/message/automation?isActive=true"
```
### 2. Retrieve Sent Messages

- **Endpoint**: `/message/messages`
- **Method**: `GET`
- **Query Parameters**:
    - `messageId`
    - `limit`
- **Summary**: Retrieves sent messages with filtering `messageId` and `limit` parameters.

#### Example:

```bash
curl -X GET "http://localhost:3030/message/messages?limit=4"
```