#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# ./start.sh --profile dev --network full --action up
# ./start.sh --profile dev --network local --action up
# ./start.sh --profile dev --network netbird --action up
#
# ./start.sh --profile prod --network full --action up
# ./start.sh --profile prod --network local --action up
# ./start.sh --profile prod --network netbird --action up
#
# ./start.sh --profile dev  --action down
# ./start.sh --profile prod --action down

PROFILE="dev"
NETWORK="local"
ACTION="up"

while [[ $# -gt 0 ]]; do
  case $1 in
    --profile)
      PROFILE="$2"
      shift 2
      ;;
    --network)
      NETWORK="$2"
      shift 2
      ;;
    --action)
      ACTION="$2"
      shift 2
      ;;
    *)
      echo "invalid value $1"
      echo "help: $0 [--profile dev|prod] [--network local|netbird|full] [--action up|down]"
      exit 1
      ;;
  esac
done

if [ "$NETWORK" == "netbird" ]; then
  HOST_IP=$(netbird status --ipv4 | awk '{print $1}')
  if [ -z "$HOST_IP" ]; then
    echo "fail to get netbird ip."
    exit 1
  fi
elif [ "$NETWORK" == "full" ]; then
  HOST_IP="0.0.0.0"
else
  HOST_IP="127.0.0.1"
fi

echo "final IP: $HOST_IP"
echo "Profile: $PROFILE"
echo "Action: $ACTION"

if [ "$PROFILE" == "prod" ]; then
  COMPOSE_FILE="docker/docker-compose.prod.yml"
else
  COMPOSE_FILE="docker/docker-compose.yml"
fi

if command -v docker-compose &> /dev/null; then
  DOCKER_COMPOSE="docker-compose"
else
  DOCKER_COMPOSE="docker compose"
fi

if [ ! -f .env ]; then
  echo ".env not found — creating a minimal one"
  cat > .env <<'EOF'
PORT=8080
EOF
fi

if [ "$ACTION" == "up" ]; then
  HOST_IP=$HOST_IP $DOCKER_COMPOSE -f "$COMPOSE_FILE" --env-file .env up --build -d
elif [ "$ACTION" == "down" ]; then
  HOST_IP=$HOST_IP $DOCKER_COMPOSE -f "$COMPOSE_FILE" --env-file .env down
else
  echo "Fail. action should be up or down."
  exit 1
fi
