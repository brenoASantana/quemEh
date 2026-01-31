#!/bin/bash

# Build para produção

set -e

cd "$(dirname "$0")/.."

echo "Compilando backend..."
cd backend && go build -o quemEh-prod

echo "Build do frontend..."
cd ../frontend && npm run build

echo ""
echo "Pronto!"
echo "Backend: backend/quemEh-prod"
echo "Frontend: frontend/dist/"
