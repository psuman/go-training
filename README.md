## Running Redis Container

docker run --name some-redis -d -p 6379:6379 redis redis-server --appendonly yes

## Running Mongo Container

docker run -d -p 27017:27017 -v ~/data:/data/db mongo