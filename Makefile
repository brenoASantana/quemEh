.PHONY: help install build dev backend frontend clean production ngrok

help:
	@echo "Quem É? - Comandos disponíveis:"
	@echo ""
	@echo "  make dev          Inicia backend + frontend (dev)"
	@echo "  make backend      Apenas servidor Go"
	@echo "  make frontend     Apenas React (dev)"
	@echo "  make build        Compila tudo com go:embed"
	@echo "  make production   Build otimizado para produção"
	@echo "  make install      Instala dependências"
	@echo "  make ngrok        Guia para usar ngrok"
	@echo "  make clean        Remove binários"

install:
	cd backend && go mod download
	cd frontend && npm install

build:
	cd frontend && npm run build
	cd backend && go build -o ../quemEh main.go

production: build
	@echo "✅ Binário pronto! Execute: ./quemEh"

backend:
	cd frontend && npm run build
	cd backend && go build -o ../quemEh main.go && ../quemEh

frontend:
	cd frontend && npm run dev

dev:
	cd backend && go build -o ../quemEh main.go && ../quemEh &
	sleep 2
	cd frontend && npm run dev

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
