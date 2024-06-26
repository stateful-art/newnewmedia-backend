#!/bin/bash

# Pull the Redis Alpine image if not already present
docker pull redis:alpine

# Create the Redis cluster nodes in a loop
for i in {1..6}; do
    docker run -d --name "redis-node$i" -p "700$i:700$i" redis:alpine \
    redis-server --port "700$i" --cluster-enabled yes \
    --cluster-config-file nodes.conf --cluster-node-timeout 5000
done

# Wait for a few seconds for the containers to start
sleep 5

# Get the IP addresses of the Redis nodes and create the cluster
NODES=""
for i in {1..6}; do
    NODE_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "redis-node$i")
    NODES+="$NODE_IP:700$i,"
done

# Remove the trailing comma
NODES=${NODES%,}

# Add replicas for each node
REPLICAS=""
for i in {1..3}; do
    REPLICAS+="$NODES,"
done

# Remove the trailing comma
REPLICAS=${REPLICAS%,}

# Create the Redis cluster
docker exec "redis-node1" \
redis-cli --cluster create $NODES --cluster-replicas 1 --cluster-yes

# Wait for a few seconds for the nodes to be added to the cluster
sleep 12

# Check if the nginx-proxyy network exists, and create it if not
if ! docker network inspect nginx-proxyy &>/dev/null; then
    docker network create nginx-proxyy
fi

# Check if .env file exists, and create it if not
ENV_FILE=".env"
if [ ! -f "$ENV_FILE" ]; then
    touch "$ENV_FILE"
fi

# Check if Redis cluster nodes are already present in .env file
if ! grep -q "REDIS_CLUSTER_NODES=" "$ENV_FILE"; then
    echo "REDIS_CLUSTER_NODES=$NODES" >> "$ENV_FILE"
else
    echo "REDIS_CLUSTER_NODES already exists in $ENV_FILE, skipping..."
fi

# Start the Docker Compose stack
docker-compose up -d
