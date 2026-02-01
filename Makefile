.PHONY: help install build dev backend frontend clean production ngrok

help:
	@echo "Quem É? - Comandos disponíveis:"
	@echo ""
	@echo "  make dev          Inicia backend + frontend (tudo em uma porta)"
	@echo "  make backend      Inicia backend + frontend (mesmo que dev)"
	@echo "  make frontend     Inicia backend + frontend (mesmo que dev)"
	@echo "  make build        Compila tudo (sem rodar)"
	@echo "  make production   Build otimizado para produção"
	@echo "  make install      Instala dependências"
	@echo "  make kill         Mata processo na porta 8080"
	@echo "  make clean        Remove binários"
	@echo "  make ngrok        Guia para usar ngrok"

install:
	cd backend && go mod tidy
	cd frontend && npm install

build:
	cd frontend && npm run build
	cd backend && go build -o ../quemEh .

production: build
	@echo "✅ Binário pronto! Execute: ./quemEh"

backend:
	cd frontend && npm run build
	cd backend && go build -o ../quemEh . && ../quemEh

frontend: backend

dev: backend

kill:
	@echo "Finalizando processo na porta 8080..."
	@lsof -t -i:8080 | xargs kill -9 2>/dev/null || echo "Nenhum processo encontrado na porta 8080."

clean:
	rm -f quemEh
	rm -rf frontend/dist

ngrok:
	@echo "Para usar ngrok:"
	@echo "1. Download: https://ngrok.com/download"
	@echo "2. Execute: ngrok http 8080"
	@echo "3. Compartilhe o link gerado!"
	@echo ""
	@echo "Veja mais em: NGROK.md"

.DEFAULT_GOAL := help
