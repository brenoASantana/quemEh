# Usando ngrok para Acesso Remoto

## O que é ngrok?

ngrok cria um túnel seguro para expor seu servidor local na internet. Perfeito para que amigos acessem o jogo de qualquer lugar.

## Instalação

1. Baixe ngrok: https://ngrok.com/download
2. Extraia o arquivo
3. (Opcional) Adicione ao PATH: `export PATH="$PATH:~/ngrok"`

## Como usar

### 1. Inicie o servidor Go
```bash
cd /home/user/Coding/Pessoal/quemEh
./quemEh
```

O servidor rodará em `http://localhost:8080`

### 2. Abra outro terminal e execute ngrok
```bash
ngrok http 8080
```

Você verá algo como:
```
ngrok                                       (Ctrl+C to quit)

Session Status                online
Session Expires               2h 59m 47s
Version                       3.0.0
Region                        us (United States)
Forwarding                    https://1234-5678-abcd.ngrok.io -> http://localhost:8080
```

### 3. Compartilhe o link

Seu link público é: `https://1234-5678-abcd.ngrok.io`

Envie este link para seus amigos!

## Características

✅ **Automático**: A URL se ajusta dinamicamente (funciona com qualquer host)
✅ **Seguro**: Usa HTTPS automaticamente
✅ **WebSocket**: Funciona perfeitamente com WebSocket (wss://)
✅ **Sem configuração**: Nenhuma alteração necessária no código

## Dicas

- O binário compilado (`./quemEh`) já inclui todo o frontend (go:embed)
- Não precisa servir arquivos separadamente
- A URL muda cada vez que você reinicia ngrok (versão gratuita)
- Para manter a mesma URL, use plano pago do ngrok
