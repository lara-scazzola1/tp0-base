#!/bin/bash

NETWORK_NAME="tp0_testing_net"
SERVER_CONTAINER_NAME="server"
SERVER_PORT=12345
MESSAGE="hello"
CONTAINER_NAME="container_for_test_with_netcat"

if ! docker network ls | grep -q "${NETWORK_NAME}"; then
  echo "Network doesn't exist. Please start docker compose"
  exit 1
fi

RESPONSE=$(docker run --network="${NETWORK_NAME}" --rm busybox sh -c \
  "echo '${MESSAGE}' | nc ${SERVER_CONTAINER_NAME} ${SERVER_PORT}" | tail -n 1)

if [ "${RESPONSE}" == "${MESSAGE}" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
fi
