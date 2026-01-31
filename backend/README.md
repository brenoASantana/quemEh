# Backend

Servidor em Go.

```bash
go build
./quemEh
```

Roda em `localhost:8080`

Endpoints:
- `GET /health` - Health check
- `WS /ws?room=CODIGO&name=NOME` - Conexão WebSocket
