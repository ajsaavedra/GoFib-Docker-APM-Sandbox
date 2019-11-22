# GoFib-Docker-APM-Sandbox
Dockerized Fibonacci Calculator

This sandbox spins up a Go api that interacts within a containerized environment and communicates MySQL, Redis, and can be utilized via the command line or using the React client application.

The code is made up of 3 separate subprojects which run on separate containers and is meant to help illustrate tracing, span generation, tagging, and correlated logs in the Datadog APM.

What's containerized:
- Datadog Agent enables APM and Logs
- Go API server
  - Debug mode optional
  - Automatic instrumentation traces API calls
  - Includes Redis client that publishes requests to Redis 
- Go DB worker
  - Debug mode optional
  - Manual instrumentation traces work done by the server
  - Subscribed to Redis and calculates Fibonacci values
- React client app
  - Provides a UI to easily calculate, delete, and view all values
- Redis server
- MySQL server

## Sections

- Step 1: Create a custom .env file
- Step 2: Build images
- Step 3: Spin up containers
- Step 4: Curling endpoints
- Step 5: Spin down containers

### Step 1: Create a custom .env file

Make sure that in your ~ directory, you have a file called sandbox.docker.env that contains:

```
DD_API_KEY=<Your API Key>
```

### Step 2: Build images

1. Add your API Key to docker-compose.yml
2. Have Docker running, then 

```
docker-compose build --no-cache
```

### Step 3: Spin up containers

```
docker-compose up
```

### Step 4: Curling endpoints

You can use the browser to access these endpoints. Here are the endpoints you can use if you want to use the terminal:

1. GET a list of all values

```
curl localhost:80/api/all
```

2. POST a new value

```
curl -X POST localhost:80/api/fib -d '{"value": <fib-value>}'
```

3. DELETE a value

```
curl -X DELETE localhost:80/api/fib/<fib-value>
```

4. GET a value

```
curl localhost:80/api/fib/<fib-value>
```

At this point, you should have traces in the Datadog UI outlining how some recursive functions can have overlapping subproblems and how to optimize them.

### Step 5: Spin down containers

```
docker-compose down
```
