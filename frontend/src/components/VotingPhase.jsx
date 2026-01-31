export default function VotingPhase({
  answers, players, currentPlayer, hasSubmittedVotes, localVotes, setLocalVotes, gameState, sendMessage, totalPlayers
}) {
  return (
    <div className="phase voting-phase">
      <h2>🕵️ Quem respondeu?</h2>
      <p className="voting-instruction">Tente adivinhar quem escreveu cada resposta!</p>
      <div className="answers-container">
        {(answers || [])
          .filter(ans => ans.realId !== currentPlayer?.id)
          .map((ans) => {
            const myVote = hasSubmittedVotes
              ? (gameState.votedGuesses?.[currentPlayer?.id]?.[ans.id] || null)
              : (localVotes[ans.id] || null);
            const alreadyVotedPlayers = Object.entries(hasSubmittedVotes ? (gameState.votedGuesses?.[currentPlayer?.id] || {}) : localVotes)
              .filter(([answerId]) => parseInt(answerId) !== ans.id)
              .map(([, playerId]) => playerId);
            return (
              <div
                key={ans.id}
                className={`answer-card ${myVote ? 'selected' : ''}`}
              >
                <p className="answer-text">"{ans.text}"</p>
                {myVote && (
                  <div className="voted-indicator">
                    ✅ Você escolheu: {players.find(p => p.id === myVote)?.name}
                  </div>
                )}
                <div className="player-options">
                  {players
                    .filter(p => p.id !== currentPlayer?.id)
                    .map(p => {
                      const isAlreadyVoted = alreadyVotedPlayers.includes(p.id);
                      return (
                        <button
                          key={p.id}
                          className={`guess-btn ${myVote === p.id ? 'active' : ''} ${isAlreadyVoted && myVote !== p.id ? 'disabled' : ''}`}
                          onClick={() => {
                            if (!isAlreadyVoted || myVote === p.id) {
                              setLocalVotes(prev => ({ ...prev, [ans.id]: p.id }));
                            }
                          }}
                          disabled={hasSubmittedVotes || (isAlreadyVoted && myVote !== p.id)}
                          title={isAlreadyVoted && myVote !== p.id ? 'Você já votou nesta pessoa em outra resposta' : ''}
                        >
                          {p.name}
                        </button>
                      );
                    })}
                </div>
              </div>
            );
          })}
      </div>
      {hasSubmittedVotes ? (
        <div className="answer-sent">Votos enviados! Aguardando outros jogadores...</div>
      ) : (
        <button
          className="btn-primary btn-large"
          onClick={() => {
            Object.entries(localVotes).forEach(([answerId, playerId]) => {
              sendMessage('SUBMIT_GUESS', {
                answerId: parseInt(answerId),
                guessedPlayerId: playerId
              });
            });
            // O setHasSubmittedVotes deve ser feito no componente pai
          }}
          disabled={Object.keys(localVotes).length < (answers?.filter(a => a.realId !== currentPlayer?.id).length || 0)}
          style={{ marginTop: '20px' }}
        >
          📤 ENVIAR VOTOS
        </button>
      )}
      <div className="progress-bar">
        <p>Votaram: {Object.keys(gameState.votedGuesses || {}).length}/{totalPlayers}</p>
        <div className="bar-fill" style={{width: `${totalPlayers ? (Object.keys(gameState.votedGuesses || {}).length / totalPlayers) * 100 : 0}%`}}></div>
      </div>
    </div>
  );
}
