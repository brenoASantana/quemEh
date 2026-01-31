# Documentação - Quem É?

## Como funciona

1. Alguém cria uma sala (vira o HOST)
2. Amigos entram com o código
3. HOST clica em "COMEÇAR JOGO"
4. Cada um responde uma pergunta (anônimo)
5. Tentam adivinhar quem respondeu
6. Ganha quem acertar mais

## Setup

```bash
make install   # Instala dependências
make dev       # Inicia tudo
```

## Portas

- Backend: `localhost:8080`
- Frontend: `localhost:5173`

## Mudando a porta

Edite `backend/main.go`:
```go
log.Fatal(http.ListenAndServe(":OUTRA_PORTA", nil))
```

## Deploy

Build para produção:
```bash
make build
```

Outputs:
- Backend: `backend/quemEh-prod`
- Frontend: `frontend/dist/`

Você pode servir os arquivos estáticos do `frontend/dist/` junto com o servidor Go.
