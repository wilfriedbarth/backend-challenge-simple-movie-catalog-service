# backend-challenge-simple-movie-catalog-service

Backend Challenge: Simple Movie Catalog Service

## Design Decisions

### Programming Language

For this project, I decided to go with GoLang.

1. Fast compilation / small binaries
2. Lean feature set / relatively easy to learn
3. Excellent tooling (testing, formatting)
4. Handling of web requests / scaling to high traffic loads via goroutines is more efficient than other solutions (e.g. Java threadpool)

### Database

Assuming that this is a simple microservice only responsible for movies, I decided to opt
for using a NoSQL solution (ElasticSearch) over a SQL solution.

1. ElasticSearch has the ability to perform extremely fast searches using distributed inverted indices. This allows for retrieving searched data in ms vs seconds for an SQL solution.
2. ElasticSearch is easily scalable horizontally. It runs a clusters of nodes, which can be easily scaled up to handle more traffic.
3. ElasticSearch divides indices into shards, and shards can create replicas. When documents are added routing and rebalancing operations are conducted automatically.
4. Movies are stored as simple JSON documents in an index and are schema free. As documents are added with new fields, the mappings are automatically updated and older documents will have the associated fields automatically added.
5. ElasticSearch has a convenient API to search across multiple fields (multi_match query).

For comparison, if I did select a SQL solution, here are the advantages for that solution.

1. Data integrity - atomicity, consistency, isolation, durability (ACID)
2. SQL databases scale vertically via CPUs or memory
3. Efficient at processing queries and joining data across tables.
4. More mature technology (i.e. most developers know how to use it)

My decision to choose ElasticSearch over SQL solution is primarily motivated by search performance and horizontal scalability. Given the scope of the service (movie catalog) I did not
think it was necessary to opt for an SQL solution. If this had been a service with a wider
scope (i.e. supporting multiple APIs for movies, actors, reviews, etc), I would choose a SQL solution that allowed for queries across multiple tables.

### Docker / Docker-Compose

Using Docker / Docker-Compose to package the web server and database simplifies local development.

- For deploying to testing / production environments, Kubernetes or manages solutions like ECS and Fargate would be my choice to efficient handle scaling of web servers and ElasticSearch. In the interest of time, I will ignore this for now.

## Startup Instruction

1. Run `docker compose up` to start all services.

## Using the API
