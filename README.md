## Running Redis Container

docker run -p 6379:6379 redis redis-server --appendonly yes

## Running Mongo Container

docker run -p 27017:27017 -v ~/data:/data/db mongo