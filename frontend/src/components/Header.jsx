import Avatar from "../Avatar";

export default function Header({ room, name, isHost, currentPlayer }) {
  return (
    <header className="game-header">
      <div className="header-left">
        <div className="room-code">Sala: <strong>{room}</strong></div>
        <div className="game-state">{currentPlayer?.stateLabel}</div>
      </div>
      <div className="header-right">
        <Avatar name={name} size="sm" />
        <div className="player-info">
          <p className="player-name">{name} {isHost ? '👑' : ''}</p>
          <p className="player-score">{currentPlayer?.score || 0} pts</p>
        </div>
      </div>
    </header>
  );
}
