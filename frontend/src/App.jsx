import { useState } from 'react';
import './App.css';
import { useGameSocket } from './useGameSocket';

const Avatar = ({ name, size = "md" }) => (
  <img
    src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${name}`}
    alt="avatar"
    className={`avatar avatar-${size}`}
  />
);

function App() {
  const { connect, isConnected, gameState, sendMessage } = useGameSocket();
  const [name, setName] = useState('');
  const [room, setRoom] = useState('');
  const [answer, setAnswer] = useState('');
  const [guess, setGuess] = useState(null);

  if (!isConnected) {
    return (
      <div className="container login-screen">
        <div className="login-card">
          <h1>Quem É?</h1>
          <p className="subtitle">O jogo das respostas misteriosas</p>

          <div className="form">
            <input
              placeholder="Seu nome"
              value={name}
              onChange={e => setName(e.target.value)}
              onKeyPress={e => e.key === 'Enter' && name && room && connect(name, room)}
            />
            <input
              placeholder="Código da sala"
              value={room}
              onChange={e => setRoom(e.target.value)}
              onKeyPress={e => e.key === 'Enter' && name && room && connect(name, room)}
            />
            <button
              onClick={() => connect(name, room)}
              disabled={!name || !room}
              className="btn-primary"
            >
              Entrar
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (!gameState) {
    return <div className="loading">Conectando...</div>;
  }

  const me = gameState.players?.find(p => p.name === name);
  const amHost = me?.isHost;

  return (
    <div className="container game-screen">
      <header className="game-header">
        <div className="header-left">
          <div className="room-code">Sala: <strong>{room}</strong></div>
          <div className="game-state">{getStateLabel(gameState.state)}</div>
        </div>
        <div className="header-right">
          <Avatar name={name} size="sm" />
          <div className="player-info">
            <p className="player-name">{name} {amHost ? '👑' : ''}</p>
            <p className="player-score">{me?.score || 0} pts</p>
          </div>
        </div>
      </header>

      <main className="game-content">
        {gameState.state === 'LOBBY' && (
          <div className="phase lobby-phase">
            <h2>Aguardando jogadores...</h2>
            <div className="players-grid">
              {gameState.players.map(p => (
                <div key={p.id} className="player-card">
                  <Avatar name={p.name} />
                  <p className="player-name">{p.name}</p>
                  {p.isHost && <span className="host-badge">HOST</span>}
                  <p className="player-score">{p.score}</p>
                </div>
              ))}
            </div>

            {amHost ? (
              <button
                className="btn-primary btn-large"
                onClick={() => sendMessage('START_GAME', null)}
              >
                Começar
              </button>
            ) : (
              <div className="waiting-message">
                Aguardando o host iniciar...
              </div>
            )}
          </div>
        )}

        {gameState.state === 'ANSWERING' && (
          <div className="phase answering-phase">
            <div className="question-container">
              <h2>Pergunta</h2>
              <div className="question-card">
                <p>{gameState.question}</p>
              </div>
            </div>

            <div className="answer-input-container">
              <textarea
                placeholder="Sua resposta..."
                value={answer}
                onChange={e => setAnswer(e.target.value)}
                className="answer-textarea"
                maxLength={500}
              />
              <div className="char-count">{answer.length}/500</div>

              <button
                className="btn-primary btn-large"
                onClick={() => {
                  sendMessage('SUBMIT_ANSWER', answer);
                  setAnswer('');
                }}
                disabled={!answer.trim()}
              >
                Enviar
              </button>

              <div className="progress-bar">
                <p>Respondidos: {Object.keys(gameState.answers || {}).length}/{gameState.players.length}</p>
                <div className="bar-fill" style={{width: `${(Object.keys(gameState.answers || {}).length / gameState.players.length) * 100}%`}}></div>
              </div>
            </div>
          </div>
        )}

        {gameState.state === 'VOTING' && (
          <div className="phase voting-phase">
            <h2>Quem respondeu?</h2>
            <p className="voting-instruction">Tente adivinhar</p>

            <div className="answers-container">
              {gameState.answers && gameState.answers.map((ans) => (
                <div
                  key={ans.id}
                  className={`answer-card ${guess === ans.id ? 'selected' : ''}`}
                  onClick={() => {
                    setGuess(ans.id);
                    sendMessage('SUBMIT_GUESS', { answerId: ans.id.toString() });
                  }}
                >
                  <div className="answer-number">#{ans.id + 1}</div>
                  <p className="answer-text">"{ans.text}"</p>
                  <div className="player-options">
                    {gameState.players.map(p => (
                      <button
                        key={p.id}
                        className={`guess-btn ${guess === ans.id ? 'active' : ''}`}
                        onClick={(e) => {
                          e.stopPropagation();
                          setGuess(ans.id);
                          sendMessage('SUBMIT_GUESS', { answerId: ans.id.toString() });
                        }}
                      >
                        {p.name}
                      </button>
                    ))}
                  </div>
                </div>
              ))}
            </div>

            <div className="progress-bar">
              <p>Votados: {Object.keys(gameState.votedGuesses || {}).length}/{gameState.players.length}</p>
              <div className="bar-fill" style={{width: `${(Object.keys(gameState.votedGuesses || {}).length / gameState.players.length) * 100}%`}}></div>
            </div>
          </div>
        )}

        {gameState.state === 'RESULTS' && (
          <div className="phase results-phase">
            <h2>Resultados</h2>

            <div className="results-container">
              <div className="scores-summary">
                <h3>Placar</h3>
                <div className="scores-list">
                  {gameState.players
                    .sort((a, b) => b.score - a.score)
                    .map((p, i) => (
                      <div key={p.id} className="score-row">
                        <span className="rank">#{i + 1}</span>
                        <Avatar name={p.name} size="sm" />
                        <span className="name">{p.name}</span>
                        <span className="score">{p.score} pts</span>
                      </div>
                    ))}
                </div>
              </div>

              <div className="answers-review">
                <h3>Respostas</h3>
                {gameState.answers && gameState.answers.map((ans) => {
                  const writer = gameState.players.find(p => p.id === ans.id);
                  return (
                    <div key={ans.id} className="answer-reveal">
                      <p className="answer-text">"{ans.text}"</p>
                      <p className="answer-writer">Por: {writer?.name}</p>
                    </div>
                  );
                })}
              </div>
            </div>

            {amHost && (
              <button
                className="btn-primary btn-large"
                onClick={() => sendMessage('NEXT_ROUND', null)}
              >
                Próxima rodada
              </button>
            )}
          </div>
        )}
      </main>
    </div>
  );
}

function getStateLabel(state) {
  const map = {
    'LOBBY': 'Aguardando',
    'ANSWERING': 'Respondendo',
    'VOTING': 'Votando',
    'RESULTS': 'Resultados'
  };
  return map[state] || state;
}

export default App;