# Quem É?

Um jogo online pra jogar com amigos. Vocês entram em uma sala, respondem perguntas e depois tentam adivinhar quem escreveu cada resposta. É divertido mesmo.

## Começar

### Desenvolvimento
```bash
make dev
```

### Produção (compilado com go:embed)
```bash
make build
./quemEh
```

Acesse [http://localhost:8080](http://localhost:8080).

## Compartilhar com amigos remotos

Use ngrok para dar acesso remoto:

```bash
ngrok http 8080
```

Compartilhe o link que aparecer. Veja mais em [NGROK.md](NGROK.md).

## Como funciona

1. Alguém cria uma sala (vira o host/anfitrião)
2. Amigos entram com o código da sala
3. Host inicia o jogo
4. Cada um responde uma pergunta (anônimo)
5. Todos tentam adivinhar quem respondeu
6. Ganha quem acertar mais

## Estrutura

```
backend/      - Servidor Go (WebSocket + go:embed)
frontend/     - Interface React
docs/         - Documentação
NGROK.md      - Guia de acesso remoto
```

## Características

✅ **go:embed** - Frontend embutido no binário Go (compilação simplificada)
✅ **ngrok-ready** - Funciona perfeitamente com URLs dinâmicas
✅ **WebSocket** - Comunicação em tempo real
✅ **Responsivo** - Funciona em desktop e mobile

## Tech Stack

- Go 1.22 + Gorilla WebSocket + go:embed
- React 18 + Vite
- CSS responsivo

## Desenvolvido para impressionar

Demonstra React, Go, WebSocket e como integrar frontend e backend de verdade.
