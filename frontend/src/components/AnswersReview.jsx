export default function AnswersReview({ answers, players, votedGuesses }) {
  return (
    <div className="answers-review">
      <h3>📋 Respostas e Votos</h3>
      {answers && answers.map((ans) => {
        const author = players.find(p => p.id === ans.realId);
        // Pegar todos os votos para esta resposta
        const votesForThisAnswer = [];
        if (votedGuesses) {
          Object.entries(votedGuesses).forEach(([voterId, votes]) => {
            const guessedPlayerId = votes[ans.id];
            if (guessedPlayerId) {
              const voter = players.find(p => p.id === voterId);
              const guessed = players.find(p => p.id === guessedPlayerId);
              const isCorrect = guessedPlayerId === ans.realId;
              votesForThisAnswer.push({
                voter: voter?.name,
                guessed: guessed?.name,
                isCorrect
              });
            }
          });
        }
        return (
          <div key={ans.id} className="answer-reveal">
            <div className="answer-header">
              <span className="answer-author">✍️ Escrito por: <strong>{author?.name}</strong></span>
            </div>
            <p className="answer-text">"{ans.text}"</p>
            <div className="votes-detail">
              <h4>Votos recebidos:</h4>
              {votesForThisAnswer.length > 0 ? (
                <ul className="votes-list">
                  {votesForThisAnswer.map((vote, idx) => (
                    <li key={idx} className={vote.isCorrect ? 'vote-correct' : 'vote-wrong'}>
                      <span className="voter">{vote.voter}</span>
                      <span className="arrow">→</span>
                      <span className="guess">{vote.guessed}</span>
                      {vote.isCorrect ? (
                        <span className="result correct">✅ +10 pts</span>
                      ) : (
                        <span className="result wrong">❌ 0 pts</span>
                      )}
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="no-votes">Ninguém votou nesta resposta</p>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}
