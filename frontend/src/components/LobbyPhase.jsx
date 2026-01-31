import Avatar from '../Avatar';

export default function LobbyPhase({ players, isHost, sendMessage }) {
  return (
    <div className="phase lobby-phase">
      <h2>Aguardando jogadores...</h2>
      <div className="players-grid">
        {players.map(p => (
          <div key={p.id} className="player-card">
            <Avatar name={p.name} />
            <p className="player-name">{p.name}</p>
            {p.isHost && <span className="host-badge">HOST</span>}
            <p className="player-score">{p.score}</p>
          </div>
        ))}
      </div>
      {isHost ? (
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
  );
}
