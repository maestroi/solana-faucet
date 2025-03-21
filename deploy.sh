#!/bin/bash

# Configuration
REMOTE_HOST="192.168.50.220"
REMOTE_PATH="/opt/faucet"
REMOTE_USER="root"  # Change this if using a different user

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting deployment to ${REMOTE_HOST}...${NC}"

# Create the remote directory if it doesn't exist
ssh ${REMOTE_USER}@${REMOTE_HOST} "mkdir -p ${REMOTE_PATH}"

# Sync files using rsync
# Excluding:
# - node_modules (frontend dependencies)
# - database files
# - git files
# - environment files
# - build artifacts
rsync -avz --progress \
    --exclude 'node_modules' \
    --exclude 'frontend/node_modules' \
    --exclude '.git' \
    --exclude '.env' \
    --exclude '*.db' \
    --exclude '*.db-journal' \
    --exclude 'backend/wallet/faucet-keypair.json' \
    --exclude 'docker-compose.local.yml' \
    ./ ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/

if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to sync files to remote server${NC}"
    exit 1
fi

echo -e "${GREEN}Files synced successfully${NC}"

# Copy the wallet file separately (if it exists) to preserve permissions
if [ -f "backend/wallet/faucet-keypair.json" ]; then
    echo -e "${GREEN}Copying wallet keypair file...${NC}"
    scp backend/wallet/faucet-keypair.json ${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/backend/wallet/
fi

# SSH into the remote server and start the application
echo -e "${GREEN}Starting application on remote server...${NC}"
ssh ${REMOTE_USER}@${REMOTE_HOST} "cd ${REMOTE_PATH} && docker compose -f docker-compose.prod.yml up -d --build"

if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to start application on remote server${NC}"
    exit 1
fi

echo -e "${GREEN}Deployment completed successfully!${NC}" 