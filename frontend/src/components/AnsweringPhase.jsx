export default function AnsweringPhase({
  question, answer, setAnswer, hasSubmitted, submittedAnswer, setSubmittedAnswer, setHasSubmitted, answeredCount, totalPlayers, sendMessage
}) {
  return (
    <div className="phase answering-phase">
      <div className="question-container">
        <h2>Pergunta</h2>
        <div className="question-card">
          <p>{question}</p>
        </div>
      </div>
      <div className="answer-input-container">
        <textarea
          placeholder="Sua resposta..."
          value={hasSubmitted ? submittedAnswer : answer}
          onChange={e => !hasSubmitted && setAnswer(e.target.value)}
          className="answer-textarea"
          maxLength={500}
          disabled={hasSubmitted}
        />
        <div className="char-count">{(hasSubmitted ? submittedAnswer.length : answer.length)}/500</div>
        <button
          className="btn-primary btn-large"
          onClick={() => {
            if (hasSubmitted) return;
            const text = answer.trim();
            sendMessage('SUBMIT_ANSWER', text);
            setSubmittedAnswer(text);
            setHasSubmitted(true);
          }}
          disabled={!answer.trim() || hasSubmitted}
        >
          Enviar
        </button>
        {hasSubmitted && (
          <div className="answer-sent">Enviado: "{submittedAnswer}"</div>
        )}
        <div className="progress-bar">
          <p>Respondidos: {answeredCount}/{totalPlayers}</p>
          <div className="bar-fill" style={{width: `${totalPlayers ? (answeredCount / totalPlayers) * 100 : 0}%`}}></div>
        </div>
      </div>
    </div>
  );
}
