#!/bin/bash

# Inicia backend + frontend

set -e

cd "$(dirname "$0")/.."

echo "Compilando backend..."
cd backend && go build -o quemEh

echo "Verificando dependências do frontend..."
cd ../frontend
[ ! -d "node_modules" ] && npm install

echo ""
echo "Iniciando servidor (Go)..."
cd ../backend && ./quemEh &
SERVER_PID=$!

sleep 2

echo "Iniciando React..."
cd ../frontend && npm run dev

trap "kill $SERVER_PID 2>/dev/null" EXIT
