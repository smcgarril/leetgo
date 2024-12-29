# [LeetGo](https://leetgo-server.fly.dev/)
Solve LeetCode style problems in Go!\
(A learning exercise in Go programming and web development.)\
[https://leetgo-server.fly.dev](https://leetgo-server.fly.dev)

## Usage

This application is intended to check user-provided solutions to Leetcode-style problems. The server module initializes a SQLite database and seeds it with the necessary tables and a handful of problems (future versions will support user-generated problems).

A simple front-end in vanilla JS provides a prompt and sends user solutions for processing. Because validating solutions requires running arbitrary code, the actual compute is done on a separate worker service with some guardrails in place to prevent abuse (infite loops, improper recursion, fork bombs, resource exhaustion, etc.). 

Future plans include implementing a messaging queue (RabbitMQ or similar), user log in, additional problems, hints and solutions, metrics displays, improved UI (mobile-friendly), and whatever other things might sound fun to add. This has been a great learning experience and I welcome any and all feedback. Happy coding!

## Quick Start
To run locally you will need to start both the server and worker services. This can be done via docker-compose, or by building the project directly. The server application will send requests to the worker at a URL pulled from the env vars `WORKER_HOST`, `WORKER_PORT`, and `WORKER_PATH`. If empty the application will [default](https://github.com/smcgarril/leetgo/blob/main/server/api/utils.go#L9-L25) to [localhost:8081](http://localhost:8080). For docker-compose deployments this can be updated [here](https://github.com/smcgarril/leetgo/blob/main/docker-compose.yml#L9-L11), and for Dockerfile builds [here](https://github.com/smcgarril/leetgo/blob/main/server/Dockerfile#L16-L18).

### Run in Container

1. Pull latest [server](https://hub.docker.com/r/smcgarril/leetgo-server) and [worker](https://hub.docker.com/r/smcgarril/leetgo-worker) images from hub.docker.com
  ```
  $ docker run -p 8080:8080 docker.io/smcgarril/leetgo-server:1.0.0
  $ docker run -p 8081:8081 docker.io/smcgarril/leetgo-worker:1.0.0
  ```
  
  or simply use the provided docker-compose.yml
  ```
  docker-compose up -d
  ```

2. View in browser
  http://localhost:8080

### Local Build
(Assumes an [installation of Go](https://go.dev/doc/install))

1. Clone the repo
  ```
  $ git clone https://github.com/smcgarril/leetgo.git
  $ cd leetgo/server
  ```

2. Run the server program:
  ```
  $ go run .
  ```

3. In new shell cd to worker directory
  ```
  $ cd leetgo/worker
  ```

4. Run the worker service:
  ```
  $ go run .
  ```

5. View in browser
  http://localhost:8080

