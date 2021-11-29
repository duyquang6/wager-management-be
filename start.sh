#!/bin/sh
set -e

usage()
{
    echo "usage: start.sh [[[-i/--integration] (to run integration tests too) | [-h]]"
}

# Run local migration for developing
migrate() {
  make build.migrate
  export $(cat .env_template | xargs)
  ./bin/migrate
}

# Run local app for developing
start() {
  make build
  export $(cat .env_template | xargs)
  ./bin/app
}

# Run local app for developing
bootstrap() {
  echo "----------- SETUP DOCKER -------------"
  make docker.local.stop
  make build.docker.image
  make docker.local.start

  echo "Wait 30 secs for db start complete & run migration"
  sleep 30
  echo "----------- DONE -------------"
  echo "Try hit endpoint http://0.0.0.0:8080/wagers !"
  echo "Run make docker.local.stop to stop docker instances"
}

# Bootup docker & run integration test
integration() {
  echo "----------- SETUP DOCKER -------------"
  make docker.integration.stop
  make docker.integration.start

  echo "Wait 30 secs for db start complete"
  sleep 30
  echo "----------- INTEGRATION TEST -------------"
  make test.integration

  echo "----------- CLEAN UP -------------"
  make docker.integration.stop
}

while [ "$1" != "" ]; do
    case $1 in
        -m | --migrate )        migrate
                                ;;
        -s | --start )          start
                                ;;
        -i | --integration )    integration
                                ;;
        -b | --bootstrap )      bootstrap
                                ;;
        -h | --help )           usage
                                exit
                                ;;
        * )                     usage
                                exit 1
    esac
    shift
done
