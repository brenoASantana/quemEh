import AnsweringPhase from './components/AnsweringPhase';
import LobbyPhase from './components/LobbyPhase';
import VotingPhase from './components/VotingPhase';

import { useEffect, useState } from 'react';
import './App.css';
import AnswersReview from './components/AnswersReview';
import Header from './components/Header';
import Ranking from './components/Ranking';
import { useGameSocket } from './useGameSocket';

function App() {
  const { connect, isConnected, gameState, sendMessage } = useGameSocket();
  const [name, setName] = useState('');
  const [room, setRoom] = useState('');
  const [answer, setAnswer] = useState('');
  const [guess, setGuess] = useState(null);
  const [votedAnswer, setVotedAnswer] = useState(null);
  const [submittedAnswer, setSubmittedAnswer] = useState('');
  const [hasSubmitted, setHasSubmitted] = useState(false);
  const [localVotes, setLocalVotes] = useState({});
  const [hasSubmittedVotes, setHasSubmittedVotes] = useState(false);

  useEffect(() => {
    if (gameState?.state === 'ANSWERING') {
      setHasSubmitted(false);
      setSubmittedAnswer('');
      setAnswer('');
    }
    if (gameState?.state === 'VOTING') {
      setVotedAnswer(null);
      setLocalVotes({});
      setHasSubmittedVotes(false);
    }
  }, [gameState?.state, gameState?.question]);

  if (!isConnected) {
    return (
      <div className="container login-screen">
        <div className="login-card">
          <h1>🤔 Quem É?</h1>
          <p className="subtitle">O jogo das respostas misteriosas</p>

          <div className="form">
            <input
              placeholder="Seu nome"
              value={name}
              onChange={e => setName(e.target.value)}
              onKeyPress={e => e.key === 'Enter' && name && room && connect(name, room)}
            />
            <input
              placeholder="Código da sala (ex: GLOBO)"
              value={room}
              onChange={e => setRoom(e.target.value)}
              onKeyPress={e => e.key === 'Enter' && name && room && connect(name, room)}
            />
            <button
              onClick={() => connect(name, room)}
              disabled={!name || !room}
              className="btn-primary"
            >
              ENTRAR NA SALA
            </button>
          </div>

          <div className="tips">
            <p>💡 Dica: Use o mesmo código para acessar a mesma sala</p>
          </div>
        </div>
      </div>
    );
  }

  if (!gameState) {
    return (
      <div className="container">
        <div className="loading" style={{ padding: '40px', textAlign: 'center', color: '#666' }}>
          <p>⏳ Conectando ao jogo...</p>
        </div>
      </div>
    );
  }

  let currentPlayer;
  try {
    currentPlayer = gameState.players?.find(p => p.name === name);
  } catch (e) {
    console.error('Erro ao encontrar jogador atual:', e);
    return (
      <div className="container">
        <div className="loading" style={{ padding: '40px', textAlign: 'center', color: 'red' }}>
          <p>❌ Erro ao carregar jogo. Recarregue a página.</p>
        </div>
      </div>
    );
  }
  const isHost = currentPlayer?.isHost;
  const answeredCount = gameState.answersCount ?? Object.keys(gameState.answers || {}).length;
  const totalPlayers = gameState.players?.length || 0;
  const votedCount = Object.keys(gameState.votedGuesses || {}).length;
  const hasVotedForAllAnswers = votedCount === totalPlayers;

  return (
    <div className="container game-screen">
      <Header room={room} name={name} isHost={isHost} currentPlayer={{...currentPlayer, stateLabel: getStateLabel(gameState.state)}} />

      <main className="game-content">
        {gameState.state === 'LOBBY' && (
          <LobbyPhase players={gameState.players} isHost={isHost} sendMessage={sendMessage} />
        )}

        {gameState.state === 'ANSWERING' && (
          <AnsweringPhase
            question={gameState.question}
            answer={answer}
            setAnswer={setAnswer}
            hasSubmitted={hasSubmitted}
            submittedAnswer={submittedAnswer}
            setSubmittedAnswer={setSubmittedAnswer}
            setHasSubmitted={setHasSubmitted}
            answeredCount={answeredCount}
            totalPlayers={totalPlayers}
            sendMessage={sendMessage}
          />
        )}

        {gameState.state === 'VOTING' && (
          <VotingPhase
            answers={gameState.answers}
            players={gameState.players}
            currentPlayer={currentPlayer}
            hasSubmittedVotes={hasSubmittedVotes}
            localVotes={localVotes}
            setLocalVotes={setLocalVotes}
            gameState={gameState}
            sendMessage={sendMessage}
            totalPlayers={totalPlayers}
          />
        )}

        {(gameState.state === 'RESULTS' || gameState.state === 'GAME_OVER') && (
          <div className="phase results-phase">
            <h2>{gameState.state === 'GAME_OVER' ? '🏆 Fim de Jogo' : 'Resultados'}</h2>

            <div className="results-container">
              <Ranking players={gameState.players} gameState={gameState} />

              <AnswersReview answers={gameState.answers} players={gameState.players} votedGuesses={gameState.votedGuesses} />
            </div>

            {isHost && (
              gameState.state === 'GAME_OVER' ? (
                <button
                  className="btn-primary btn-large"
                  onClick={() => {
                    // Zera os pontos de todos e reinicia o jogo
                    sendMessage('RESET_GAME', null);
                  }}
                >
                  Nova partida
                </button>
              ) : (
                <button
                  className="btn-primary btn-large"
                  onClick={() => sendMessage('NEXT_ROUND', null)}
                >
                  Próxima rodada
                </button>
              )
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