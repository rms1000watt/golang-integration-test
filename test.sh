#!/usr/bin/env bash

set -e

help() {
  echo "Usage:

./test.sh      # Defaults to: ./test.sh go
./test.sh curl # Runs tests using: curl
./test.sh go   # Runs tests using: go
"
}


wait_for_postgres() {
  echo "Waiting 60s for Postgres"
  for i in {1..60}; do
    if pg_isready -h localhost -U postgres; then
      break
    fi

    sleep 1 || exit 1

    if [[ $i = 60 ]]; then
      echo "TIMEOUT: 60s.. Exiting.."
      exit 1
    fi
  done
}


cleanup_integration_environment() {
  echo "Removing integration environment in docker compose"

  docker-compose -f docker-compose.db.yml down -t 1 ||:
  docker-compose -f docker-compose.svc.yml down -t 1 ||:
  docker network rm integration ||:
}


setup_integration_environment() {
  # Destroy all
  cleanup_integration_environment

  echo "Setting up integration environment in docker compose"

  # Rebuild network + db
  docker network create --attachable integration
  docker-compose -f docker-compose.db.yml up -d

  wait_for_postgres

  # Hydrate db
  createdb -h localhost -U postgres person
  for file in sql/*.sql; do
    echo "Loading file: ${file}"
    psql -h localhost -U postgres -d person -v ON_ERROR_STOP=1 < "${file}"
  done

  # Rebuild svc
  docker-compose -f docker-compose.svc.yml up -d
  echo "Wait 3s for services to startup..."
  sleep 3
}


cmd_curl() {
  setup_integration_environment
  trap cleanup_integration_environment EXIT

  echo "Running curl..."

  # Validate these requests don't return failing exit code
  curl -f -XPOST "localhost:9999/person?name=ryan&age=88"
  curl -f -XGET "localhost:9999/person?name=ryan"

  echo "Tests Passed"
}


cmd_go() {
  setup_integration_environment
  trap cleanup_integration_environment EXIT

  echo "Running go..."

  # Run a bunch of tests managed by go
  export PERSON_SVC_TEST_SERVER_URL="http://localhost:9999"
  go test ./test
}


case "$1" in
  curl)
    cmd_curl ${@:2};;
  go)
    cmd_go ${@:2};;
  help)
    ;&
  -h)
    ;&
  --h)
    ;&
  --help)
    help; exit 1;;
  *)
    cmd_go; exit 1;;
esac