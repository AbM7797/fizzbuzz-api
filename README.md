# Go Worker FizzBuzz API

## Overview

This project implements a simple FizzBuzz API using Go. The API allows users to compute the FizzBuzz sequence based on provided parameters and tracks the most frequently requested parameters through a Cloudflare Workers KV store.

## Features

- **FizzBuzz Calculation**: Accepts parameters to generate the FizzBuzz sequence.
- **Statistics Endpoint**: Returns the most frequently requested FizzBuzz parameters and their hit count.
- **Logging Middleware**: Logs incoming requests for monitoring.

## Getting Started

### Prerequisites

- Go 1.23.3 or later
- Make
- Wrangler (for deploying to Cloudflare Workers)

## API Endpoints

### FizzBuzz

**POST /api/fizzbuzz**

- **Request Body**: JSON object with the following structure:
  ```json
  {
    "int1": 3,
    "int2": 5,
    "limit": 100,
    "str1": "fizz",
    "str2": "buzz"
  }
- **Example Request**:
    ```
    curl -X POST http://localhost:8787/api/fizzbuzz \
        -H "Content-Type: application/json" \
        -d '{
            "int1": 3,
            "int2": 5,
            "limit": 100,
            "str1": "fizz",
            "str2": "buzz"
        }'
    ```

- **Response**:

    ***status 200***: Returns the FizzBuzz result
    ```json
    {
        "result": ["1","2","fizz","4","buzz","fizz", ...]
    }
    ```
    ***status 400***: Invalide input parameters.

### Statistics

**GET /api/stats**
- **Example Request**:
    ```
    curl -X GET http://localhost:8787/api/stats
    ```
- **Response**:

    ***status 200***: Returns the most frequently requested parameters and hit count.
    ```json
    {
        "most_frequent_request": {
            "int1": 3,
            "int2": 5,
            "limit": 100,
            "str1": "fizz",
            "str2": "buzz"
        },
        "hits": 5
    }
    ```
    ***status 404***: No requests made yet.

### Makefile Commands

The Makefile contains several commands for local development, building the application, and deploying it to staging and production environments.

**Commands**
- ***dev***: Starts the local development server.
    ```
    make dev
    ```
- ***deploy-staging***: Deploys the application to the staging environment.
    ```
    make deploy-staging
    ```
- ***deploy-prod***: Deploys the application to the production environment.
    ```
    make deploy-prod
    ```
