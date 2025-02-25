# Todo API

This is a simple Todo API built with Go (Golang) and ScyllaDB, using the Gin web framework. This API allows to create, read, update, delete, and list Todo items with pagination.

## Table of Contents

- [TODO API](#todo-api)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Setup](#setup)
  - [Build and Run the API](#build-and-run-the-api)
  - [Testing with Postman](#testing-with-postman)

## Prerequisites

- [Go](https://golang.org/doc/install) (1.22.4)
- [ScyllaDB](https://www.scylladb.com/download/)
- [Docker](https://www.docker.com/products/docker-desktop/)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

## Setup

### 1. Clone the Repository
    ```sh
    git clone https://github.com/your-username/todo-api.git
    cd todo-api
    ```

### 2. Setup ScyllaDB

- Start ScyllaDB on your machine or in a Docker container.
   ```sh
   docker run --name scylla -d scylladb/scylla
   ```
   
- Access the ScyllaDB container:
   ```sh
   docker exec -it scylla cqlsh
   ```

- Create a keyspace and table in ScyllaDB:
   ```sh
   CREATE KEYSPACE todo WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};
   USE todo;
      
   CREATE TABLE items (
          id UUID PRIMARY KEY,
          user_id UUID,
          title TEXT,
          description TEXT,
          status TEXT,
          created TIMESTAMP,
          updated TIMESTAMP
   );
    ```

## Build and Run the API

- Build the Docker image:
   ```sh
   docker build -t todo-api .
   ```
   
- Run the API container:
   ```sh
   docker run --name todo-api --link scylla:scylla -p 8080:8080 -d todo-api
   ```

## Testing with Postman

### Create a TODO Item

1. Method: POST
2. URL: `http://localhost:8080/todos`
3. Body (JSON):
   ```sh
   {
     "user_id": "b4ff8577-0b4f-4033-8ab0-3d1b2e4a7f25",
     "title": "Make TODO API",
     "description": "Using Golang and scyllaDB",
     "status": "Pending"
   }
   ```
4. Send the Request :
   You should receive a response with the TODO item, including its `id`.

### Get a TODO Item

1. Method: POST
2. URL: `http://localhost:8080/todos/{id}`
   (Replace `{id}` with the `id` from the previous response.)
3. Send the Request :
   You should receive the details of the TODO item.

### Update a TODO Item

1. Method: PUT
2. URL: `http://localhost:8080/todos/{id}`
   (Replace `{id}` with the `id` from the previous response.)
3. Body (JSON):
   ```sh
   {
     "title": "Make TODO API",
     "description": "Using Golang and scyllaDB",
     "status": "Done"
   }
   ```
4. Send the Request :
   You should receive a response indicating the TODO item was updated.

### Delete a TODO Item

1. Method: DELETE
2. URL: `http://localhost:8080/todos/{id}`
   (Replace `{id}` with the `id` from the previous response.)
3. Send the Request :
   You should receive a response indicating the TODO item was deleted.

### List TODO Items with Pagination

1. Method: PUT
2. URL: `http://localhost:8080/todos?user_id={user_id}&page=1&limit=10`
   (Replace `{user_id}` with the `user_id` you used when creating the TODO item, e.g., `b4ff8577-0b4f-4033-8ab0-3d1b2e4a7f25`. For filtering based on status, add a parameter `&status={status}` in the URL. The status parameter is optional and can be Pending, Done, etc.)
3. Send the Request :
   You should receive a response with a list of TODO items and pagination details.

## Notes

- Replace "b4ff8577-0b4f-4033-8ab0-3d1b2e4a7f25" with actual UUIDs for testing.
- Ensure ScyllaDB and the API are correctly linked and running.
