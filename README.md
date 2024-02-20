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
2. Open [`Dev Tools Console`](http://localhost:5601/app/dev_tools#/console) in Kibana and run the following commands to create movies index and seed with data.

```
PUT movies
{
  "mappings": {
    "properties": {
      "title": { "type": "text" },
      "director": { "type": "text" },
      "releaseYear": { "type": "integer" },
      "genre": { "type": "text" },
      "description": {"type": "text"}
    }
  }
}
```

3. Seed movie data into ElasticSearch database by running the following commands:

```text
// Insert data into index
POST movies/_bulk
// paste text from seed-data.json file here

// Retrieve all movies to verify
GET movies/_search
{ "query": { "match_all": {}}}
```

4. Start local Go server by running `air`

## Using the API

```text
// GET movies
curl --location 'http://localhost:8080/movies'

// GET movies, search on title
curl --location 'http://localhost:8080/movies?title=spirited%20away'

// GET movies, search on genre
curl --location 'http://localhost:8080/movies?genre=anime'

// GET movie by id
curl --location 'http://localhost:8080/movies/REPLACE_WITH_YOUR_ID'

// POST movie
curl --location 'http://localhost:8080/movies' \
--header 'Content-Type: application/json' \
--data '{
  "title": "The Matrix",
  "director": "Lana Wachowski",
  "releaseYear": 1999,
  "genre": "Science Fiction",
  "description": "Neo (Keanu Reeves) believes that Morpheus (Laurence Fishburne), an elusive figure considered to be the most dangerous man alive, can answer his question -- What is the Matrix? Neo is contacted by Trinity (Carrie-Anne Moss), a beautiful stranger who leads him into an underworld where he meets Morpheus. They fight a brutal battle for their lives against a cadre of viciously intelligent secret agents. It is a truth that could cost Neo something more precious than his life."
}'

// PUT movie
curl --location 'http://localhost:8080/movies/REPLACE_WITH_YOUR_ID' \
--header 'Content-Type: application/json' \
--data '{
  "title": "The Matrixxxxxxx",
  "director": "Lana Wachowski",
  "releaseYear": 1999,
  "genre": "Science Fiction",
  "description": "Neo (Keanu Reeves) believes that Morpheus (Laurence Fishburne), an elusive figure considered to be the most dangerous man alive, can answer his question -- What is the Matrix? Neo is contacted by Trinity (Carrie-Anne Moss), a beautiful stranger who leads him into an underworld where he meets Morpheus. They fight a brutal battle for their lives against a cadre of viciously intelligent secret agents. It is a truth that could cost Neo something more precious than his life."
}'

Verify with call to GET movies

** NOTE: The update operation is not immediate 100% of the time... This is one of the downsides of Elasticsearch as a data store

// DELETE movie
curl --location --request DELETE 'http://localhost:8080/movies/rw9Zx40B82ztFqczEZwu' \
--data ''

Verify with call to GET movies
```

## TODOS

1. Add testing
2. Dockerize go server and add to docker-compose setup.
3. Consider a Kubernetes deployment to scale to test, int and prod
4. Improve error handling

NOTE: I discovered during development that ElasticSearch does not sync updates immediately. This is one downside that I missed during my initial design. In retrospect, going with a more mature technology (SQL) would have been a better choice to ensure atomic transactions.
