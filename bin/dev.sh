#! /bin/bash

cd "$(dirname "$0")/.."

case $1 in
  "migrate") docker-compose exec migrator sql-migrate ${2} --config=db/dbconf.yml;;
  "migrate!") docker-compose exec migrator sql-migrate ${2} --config=db/dbconf.yml -limit=0;;
  "seed") docker-compose exec app go run seeder.go ;;
         *) docker-compose $@;;
esac