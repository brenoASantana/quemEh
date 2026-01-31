# Quem É?

Jogo de festa onde você responde perguntas e tenta adivinhar quem respondeu cada uma. Simples e divertido.

## Começar

### Requisitos
- Go 1.16+
- Node.js 16+

### Instalação

```bash
git clone https://github.com/brenoASantana/quemEh
cd quemEh

cd frontend
npm install
npm run build
cd ..

cd backend
go build -o ../quemEh
cd ..

./quemEh
```

Abre em http://localhost:8080

## Como funciona

1. Entra numa sala com um código
2. O host começa o jogo
3. Todos respondem a pergunta
4. Todos votam em quem respondeu cada uma
5. Ganha quem fizer 50 pontos primeiro

## Estrutura

```
backend/   - Go + WebSocket
frontend/  - React + Vite
```

## Tecnologia

- Go
- React
- WebSocket

## Regras

- Não pode votar em si mesmo
- Não pode votar duas vezes na mesma resposta
- 50 pontos pra ganhar

---

Feito por Breno Santana
✅ **ngrok-ready** - Funciona perfeitamente com URLs dinâmicas
✅ **WebSocket** - Comunicação em tempo real
✅ **Responsivo** - Funciona em desktop e mobile

## Tech Stack

- Go 1.22 + Gorilla WebSocket + go:embed
- React 18 + Vite
- CSS responsivo

## Desenvolvido para impressionar

Demonstra React, Go, WebSocket e como integrar frontend e backend de verdade.
