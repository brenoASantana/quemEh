import Avatar from "../Avatar";

export default function Ranking({ players, gameState }) {
  return (
    <div className="scores-summary">
      <h3>Ranking</h3>
      <div className="scores-list">
        {players
          .sort((a, b) => b.score - a.score)
          .map((p, i) => (
            <div key={p.id} className="score-row">
              <span className="rank">#{i + 1}</span>
              <Avatar name={p.name} size="sm" />
              <span className="name">{p.name}</span>
              <span className="score">{p.score} pts</span>
              {gameState.state === 'GAME_OVER' && p.score >= 50 && (
                <span className="winner-badge">🏆 Vencedor!</span>
              )}
            </div>
          ))}
      </div>
    </div>
  );
}
