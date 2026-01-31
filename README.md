# Quem É?

Jogo de festa online. Jogadores respondem perguntas de forma anônima e votam tentando adivinhar quem respondeu cada uma. Primeiro com 50 pontos vence.

## Índice

1. [Quick Start](#quick-start)
2. [Como Jogar](#como-jogar)
3. [Stack de Tecnologia](#stack-de-tecnologia)
4. [Arquitetura](#arquitetura)
5. [Fluxo do Jogo](#fluxo-do-jogo)
6. [Comunicação em Tempo Real](#comunicação-em-tempo-real)
7. [Avatares de Usuário](#avatares-de-usuário)
8. [Compartilhando com Ngrok](#compartilhando-com-ngrok)
9. [Como Funciona Cada Parte](#como-funciona-cada-parte)
10. [Segurança](#segurança)

---

## Quick Start

### Backend
```bash
cd backend
go run main.go
```
Servidor roda em `http://localhost:8080`

### Frontend
```bash
cd frontend-quemEh
npm install
npm run dev
```
Acessa em `http://localhost:5173`

---

## Como Jogar

1. Abra o navegador em `http://localhost:8080`
2. Digite seu nome e um código de sala
3. Espere outros jogadores entrarem
4. Host clica "Iniciar Jogo"
5. Responda as perguntas
6. Vote quem respondeu cada uma
7. Primeiro com 50 pontos ganha

---

## Stack de Tecnologia

### Backend

| Tecnologia | Uso |
|-----------|-----|
| **Go** | Linguagem principal do servidor |
| **Gorilla WebSocket** | Comunicação em tempo real (bidirecional) |
| **net/http** | Servidor HTTP para servir arquivos estáticos |
| **sync** | Mutexes para sincronizar acesso concorrente |
| **JSON** | Serialização de dados |

### Frontend

| Tecnologia | Uso |
|-----------|-----|
| **React 18** | Framework de UI |
| **Vite** | Build tool e dev server (muito rápido) |
| **WebSocket** | API nativa do navegador para conexão real-time |
| **CSS3** | Estilos (avatares, cards, etc) |

### Infraestrutura

| Ferramenta | Uso |
|-----------|-----|
| **Ngrok** | Tunel seguro para compartilhar servidor local pela internet |
| **GitHub** | Versionamento de código |

---

## Arquitetura

```
┌─────────────────────┐
│   NAVEGADOR (FE)    │
│  React + Vite       │
└──────────┬──────────┘
           │ HTTP + WebSocket
           │
┌──────────▼──────────┐
│  SERVIDOR (BE)      │
│  Go + Gorilla WS    │
│  Port: 8080         │
└─────────────────────┘
           │
           ▼
   Frontend/dist
   (HTML, JS, CSS)
```

### Backend (Go)

Estrutura simplificada:
- **main.go**: HTTP server, WebSocket handler, lógica de conexão
- **room.go**: Lógica do jogo (respostas, votos, pontuação)
- **types.go**: Estruturas de dados
- **questions.go**: Carregamento de perguntas do JSON

### Frontend (React)

Estrutura modular:
- **App.jsx**: Componente principal, roteamento de estados
- **useGameSocket.js**: Hook customizado para gerenciar conexão WebSocket
- **components/**: Componentes para cada fase do jogo

---

## Fluxo do Jogo

### 1. Fase de Lobby
```
[Jogador entra] → [Nome + Código da sala] → [Conecta via WebSocket]
                     ↓
                 [Servidor cria sala]
                     ↓
                 [Host é o primeiro]
                     ↓
            [Aguarda outros jogadores]
```

### 2. Fase de Resposta (ANSWERING)
```
[Pergunta é enviada] → [Todos escrevem resposta]
                          ↓
                 [Servidor coleta respostas]
                          ↓
         [Quando todos responderam, embaralha]
```

### 3. Fase de Votação (VOTING)
```
[Respostas aparecem sem autor] → [Cada jogador vota]
                                      ↓
                          [Voto: quem fez resposta X?]
                                      ↓
                     [Servidor valida: sem self-vote]
```

### 4. Resultados (RESULTS)
```
[Mostra quem respondeu cada uma] → [Calcula pontos]
                                        ↓
                                  [Atualiza ranking]
                                        ↓
                    [Host vai para próxima rodada]
```

### 5. Vitória (GAME_OVER)
```
[Alguém atinge 50 pontos] → [Jogo termina]
                               ↓
                        [Host pode resetar]
```

---

## Comunicação em Tempo Real

### WebSocket vs HTTP

**HTTP (Tradicional):**
- Cliente pede, servidor responde
- Polling = muitas requisições desnecessárias
- Lento para atualizações em tempo real

**WebSocket:**
- Conexão bidirecional aberta
- Servidor envia atualizações quando algo muda
- Perfeito para jogos multiplayer
- Menos overhead de rede

### Protocolo de Mensagens

O jogo usa mensagens JSON simples:

**Cliente → Servidor:**
```json
{
  "type": "SUBMIT_ANSWER",
  "payload": "Minha resposta para pergunta"
}
```

**Servidor → Todos os Clientes:**
```json
{
  "type": "GAME_STATE",
  "payload": {
    "state": "VOTING",
    "players": [{...}, {...}],
    "question": "Qual é seu maior medo?",
    "answers": [
      {"id": 0, "text": "Resposta 1"},
      {"id": 1, "text": "Resposta 2"}
    ],
    "votedGuesses": {"player-id": {0: "player-id-resposta1"}}
  }
}
```

### Tipos de Mensagem

| Type | Enviado por | O que faz |
|------|-------------|----------|
| START_GAME | Host | Inicia o jogo |
| SUBMIT_ANSWER | Jogador | Envia resposta |
| SUBMIT_GUESS | Jogador | Vota em quem respondeu |
| SHOW_RESULTS | Host | Mostra resultados |
| NEXT_ROUND | Host | Próxima pergunta |
| RESET_GAME | Host | Reinicia tudo |
| GAME_STATE | Servidor | Estado atual do jogo |

---

## Avatares de Usuário

### Como Funciona

O projeto gera avatares automáticos baseados no **nome do jogador**.

```javascript
// Frontend (useGameSocket.js)
const getAvatarColor = (name) => {
  // Usa o nome para gerar uma cor consistente
  // Mesmo nome = mesma cor sempre
}
```

### Por que é legal?

1. **Consistência**: Mesmo jogador sempre tem a mesma cor
2. **Sem servidor extra**: Gerado localmente no navegador
3. **Sem dependência**: Não precisa chamar API de avatar
4. **Rápido**: Zero latência

### Implementação

```javascript
// Exemplo simplificado
function getAvatarColor(name) {
  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash);
  }
  const color = '#' + (hash & 0xFFFFFF).toString(16);
  return color;
}

// Renderiza no React
<div style={{
  backgroundColor: getAvatarColor(playerName),
  width: '60px',
  height: '60px',
  borderRadius: '50%'
}}>
  {playerName[0].toUpperCase()}
</div>
```

---

## Compartilhando com Ngrok

### O que é Ngrok?

Ngrok cria um **túnel seguro** que expõe seu servidor local para a internet.

```
Internet → ngrok.io (domínio público) → localhost:8080 (seu PC)
```

### Por que usar?

Sem ngrok, você só pode acessar de:
- Localhost: `http://localhost:8080`
- Mesma rede: `http://192.168.1.100:8080`

Com ngrok, pessoas de qualquer lugar do mundo podem jogar!

### Como usar

1. **Baixar ngrok**
   ```bash
   # macOS
   brew install ngrok

   # Linux
   wget https://bin.equinox.io/c/bNyj1mQVY4c/ngrok-v3-stable-linux-amd64.zip
   ```

2. **Criar conta gratuita**
   - Visite: https://ngrok.com
   - Crie conta
   - Copie seu token de autenticação

3. **Autenticar**
   ```bash
   ngrok config add-authtoken SEU_TOKEN_AQUI
   ```

4. **Executar ngrok**
   ```bash
   # Em outro terminal, enquanto seu servidor Go roda
   ngrok http 8080
   ```

5. **Compartilhar o link**
   ```
   Forwarding: https://abcd-1234-ef56-gh78.ngrok-free.app
   ```

   Envie esse link para amigos! Eles acessam em:
   ```
   https://abcd-1234-ef56-gh78.ngrok-free.app
   ```

### Segurança com Ngrok

O ngrok é seguro por padrão:
- HTTPS (criptografado)
- Túnel único por sessão
- Pode regenerar link
- Logs de requisições disponíveis

**Dica**: Feche o ngrok quando não estiver usando!

---

## Como Funciona Cada Parte

### 1. Inicialização (main.go)

```go
func main() {
    loadQuestions()      // Carrega perguntas.json

    // Encontra frontend/dist
    staticDir := "frontend/dist"

    // Setup HTTP
    http.HandleFunc("/ws", handleConnections)    // WebSocket
    http.HandleFunc("/", serveStatic)             // Arquivos estáticos

    http.ListenAndServe(":8080", nil)
}
```

### 2. Conexão WebSocket (main.go)

```go
func handleConnections(w http.ResponseWriter, r *http.Request) {
    // Atualiza conexão HTTP para WebSocket
    conn, _ := websocket.Upgrade(w, r, nil)

    // Extrai parâmetros
    roomID := r.URL.Query().Get("room")
    playerName := r.URL.Query().Get("name")

    // Cria jogador
    player := &Player{
        ID:   "único-id",
        Name: playerName,
        Conn: conn,
    }

    // Adiciona à sala
    room := getOrCreateRoom(roomID)
    room.Players[player.ID] = player

    // Primeiro a entrar = host
    if len(room.Players) == 1 {
        player.IsHost = true
    }

    // Começa a escutar mensagens
    go handlePlayerMessages(player)
}
```

### 3. Lógica do Jogo (room.go)

```go
// Quando todos responderam
func (r *Room) HandleAnswer(p *Player, answer string) {
    r.Answers[p.ID] = answer

    if len(r.Answers) == len(r.Players) {
        // Todos responderam!
        r.State = StateVoting
        r.ShuffledAnswers = ShuffleAnswers(r.Answers)
        r.VotedGuesses = make(map[string]map[int]string)
    }

    // Envia novo estado para todos
    broadcastGameState(r)
}
```

### 4. Validação de Votos

O backend **previne trapaças**:

```go
func (r *Room) HandleGuess(p *Player, data map[string]interface{}) {
    answerId := data["answerId"].(float64)
    guessedPlayerId := data["guessedPlayerId"].(string)

    // Não pode ser a resposta do próprio jogador
    if guessedPlayerId == p.ID {
        return  // Ignora self-vote
    }

    // Não pode votar duas vezes na mesma resposta
    if r.VotedGuesses[p.ID][answerId] != "" {
        return  // Já votou aqui
    }

    // Voto válido
    r.VotedGuesses[p.ID][answerId] = guessedPlayerId
}
```

### 5. Cálculo de Pontos

```go
func (r *Room) CalculateScores() {
    // Para cada resposta, quem acertou?
    for answerAuthor, votes := range r.VotedGuesses {
        for _, guessedPlayer := range votes {
            if guessedPlayer == answerAuthor {
                // Acertou!
                r.Players[guessedPlayer].Score += 10
            }
        }
    }

    // Alguém atingiu 50?
    for _, p := range r.Players {
        if p.Score >= 50 {
            r.State = StateGameOver
        }
    }
}
```

### 6. Frontend - Hook WebSocket (useGameSocket.js)

```javascript
export const useGameSocket = () => {
    const [gameState, setGameState] = useState(null);
    const [isConnected, setIsConnected] = useState(false);

    const connect = (name, room) => {
        const wsUrl = `ws://localhost:8080/ws?room=${room}&name=${name}`;
        const ws = new WebSocket(wsUrl);

        ws.onopen = () => {
            setIsConnected(true);
        };

        ws.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            if (msg.type === 'GAME_STATE') {
                setGameState(msg.payload);  // Atualiza React
            }
        };
    };

    const sendMessage = (type, payload) => {
        ws.send(JSON.stringify({ type, payload }));
    };

    return { connect, gameState, sendMessage, isConnected };
};
```

### 7. Fluxo de Estado no React

```javascript
// App.jsx
const { gameState, sendMessage } = useGameSocket();

// Renderizar baseado no estado
if (gameState.state === 'LOBBY') {
    return <LobbyPhase />;
}

if (gameState.state === 'ANSWERING') {
    return <AnsweringPhase />;
}

if (gameState.state === 'VOTING') {
    return <VotingPhase />;
}

if (gameState.state === 'RESULTS') {
    return <ResultsPhase />;
}
```

---

## Fluxo de Dados Completo

### Exemplo: Um jogador responde uma pergunta

1. **Frontend**
   - Usuário digita resposta
   - Clica "Enviar"
   - Envia: `{ type: "SUBMIT_ANSWER", payload: "Minha resposta" }`

2. **Backend - handlePlayerMessages()**
   ```go
   case "SUBMIT_ANSWER":
       p.Room.HandleAnswer(p, msg.Payload.(string))
   ```

3. **Backend - HandleAnswer()**
   ```go
   r.Answers[p.ID] = "Minha resposta"

   if len(r.Answers) == len(r.Players) {
       // Todos responderam!
       r.State = StateVoting
       r.ShuffledAnswers = ShuffleAnswers(r.Answers)
   }

   broadcastGameState(r)  // Notifica todos
   ```

4. **Backend - broadcastGameState()**
   ```go
   msg := Message{
       Type: "GAME_STATE",
       Payload: PublicGameState{
           State: "VOTING",
           Answers: [
               {id: 0, text: "Resposta 1", author: hidden},
               {id: 1, text: "Resposta 2", author: hidden}
           ]
       }
   }
   r.Broadcast <- msg
   ```

5. **Backend - runRoom()**
   ```go
   for msg := range r.Broadcast {
       // Envia para cada jogador
       for _, p := range r.Players {
           p.Conn.WriteJSON(msg)
       }
   }
   ```

6. **Frontend - useGameSocket()**
   ```javascript
   ws.onmessage = (event) => {
       const msg = JSON.parse(event.data);
       setGameState(msg.payload);  // React re-renderiza
   }
   ```

7. **Frontend - React Re-renderiza**
   - App.jsx vê `gameState.state === 'VOTING'`
   - Mostra `<VotingPhase />` ao invés de `<AnsweringPhase />`
   - Usuário vê respostas sem saber quem respondeu

---

## Segurança

### O que está protegido

✅ **Votos**: Validado no servidor (sem self-vote, sem duplicate)
✅ **Respostas**: Anônimas durante votação
✅ **Autenticação**: Identificação por nome + ID único

### O que NÃO está protegido (considerar em produção)

❌ Usuário pode falsificar nome na URL
❌ Sem login real / token JWT
❌ Sem HTTPS em localhost
❌ Sem rate limiting
❌ Sem validação de entrada

### Melhorias Futuras

```go
// Usar JWT para login seguro
type Claims struct {
    UserID   string
    Username string
    exp      time.Time
}

// Rate limiting por IP
type RateLimiter struct {
    limiter map[string]*time.Ticker
}

// Input validation
func ValidateAnswer(s string) error {
    if len(s) > 500 {
        return errors.New("resposta muito longa")
    }
    // ... mais validações
}
```

---

## Escalabilidade

### Atual (Single Server)

```
Servidor 8080
  ├── Sala 1 (5 jogadores)
  ├── Sala 2 (3 jogadores)
  └── Sala 3 (8 jogadores)
```

**Limite**: ~100-500 jogadores simultâneos em um servidor comum

### Se crescer muito

Opções:

1. **Redis Pub/Sub**: Múltiplos servidores Go sincronizados
2. **Load Balancer**: Distribuir conexões entre servidores
3. **Message Queue**: RabbitMQ, Kafka para comunicação entre servidores
4. **Serverless**: AWS Lambda, Google Cloud Functions
