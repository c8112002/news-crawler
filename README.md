# News Crawler

## Getting started

### dev

```bash
docker-compose up -d
docker-compose exec app go run main.go
```

migration

```bash
./bin/dev.sh migrate! up
```


### production

```bash
docker build -t crawler -f docker/go/Dockerfile .
docker run crawler
```