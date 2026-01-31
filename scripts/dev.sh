#!/bin/bash

# Apenas compila, não inicia

set -e

cd "$(dirname "$0")/.."

echo "Compilando backend..."
cd backend && go build -o quemEh

echo "Preparando frontend..."
cd ../frontend
[ ! -d "node_modules" ] && npm install

echo ""
echo "Pronto!"
echo "Para rodar, em dois terminais:"
echo "  Terminal 1: cd backend && ./quemEh"
echo "  Terminal 2: cd frontend && npm run dev"
echo ""
echo "Ou execute: ./scripts/start.sh"
