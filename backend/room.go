import "sync"

// Inicia o jogo
func (r *Room) StartGame() {
       r.Mutex.Lock()
       r.QuestionDeck = repo.GetShuffledQuestions()
       r.QuestionIndex = 0
       r.State = StateAnswer
       if len(r.QuestionDeck) == 0 {
	       r.Question = "Sem perguntas cadastradas"
       } else {
	       r.Question = r.QuestionDeck[0]
       }
       r.Answers = make(map[string]string)
       r.VotedGuesses = make(map[string]map[int]string)
       r.ShuffledAnswers = []AnswerWithIndex{}
       r.Mutex.Unlock()
       broadcastGameState(r)
}

// Processa resposta de um jogador
func (r *Room) HandleAnswer(p *Player, answer string) {
       r.Mutex.Lock()
       if _, alreadyAnswered := r.Answers[p.ID]; alreadyAnswered {
	       r.Mutex.Unlock()
	       return
       }
       if answer == "" {
	       r.Mutex.Unlock()
	       return
       }
       r.Answers[p.ID] = answer
       allAnswered := len(r.Answers) == len(r.Players)
       if allAnswered {
	       r.State = StateVoting
	       shuffledAnswers := shuffleAnswers(r.Answers)
	       r.ShuffledAnswers = shuffledAnswers
	       r.VotedGuesses = make(map[string]map[int]string)
       }
       r.Mutex.Unlock()
       broadcastGameState(r)
}

// Processa palpite de um jogador
func (r *Room) HandleGuess(p *Player, data map[string]interface{}) {
       answerIdFloat, ok1 := data["answerId"].(float64)
       guessedPlayerId, ok2 := data["guessedPlayerId"].(string)
       if !ok1 || !ok2 {
	       return
       }
       answerId := int(answerIdFloat)
       r.Mutex.Lock()
       if r.VotedGuesses[p.ID] == nil {
	       r.VotedGuesses[p.ID] = make(map[int]string)
       }
       r.VotedGuesses[p.ID][answerId] = guessedPlayerId
       expectedVotes := len(r.ShuffledAnswers) - 1
       allVoted := true
       for _, player := range r.Players {
	       if len(r.VotedGuesses[player.ID]) < expectedVotes {
		       allVoted = false
		       break
	       }
       }
       r.Mutex.Unlock()
       if allVoted {
	       r.CalculateScores()
	       r.Mutex.Lock()
	       if r.State != "GAME_OVER" {
		       r.State = StateResults
	       }
	       r.Mutex.Unlock()
       }
       broadcastGameState(r)
}

// Mostra resultados e verifica vitória
func (r *Room) ShowResults() {
       r.CalculateScores()
       r.Mutex.Lock()
       winner := ""
       for _, p := range r.Players {
	       if p.Score >= 50 {
		       winner = p.Name
		       break
	       }
       }
       if winner != "" {
	       r.State = "GAME_OVER"
       } else {
	       r.State = StateResults
       }
       r.Mutex.Unlock()
       broadcastGameState(r)
}

// Próxima rodada
func (r *Room) NextRound() {
       r.Mutex.Lock()
       r.QuestionIndex++
       if r.QuestionIndex >= len(r.QuestionDeck) {
	       r.QuestionDeck = repo.GetShuffledQuestions()
	       r.QuestionIndex = 0
       }
       r.State = StateAnswer
       if len(r.QuestionDeck) == 0 {
	       r.Question = "Sem perguntas cadastradas"
       } else {
	       r.Question = r.QuestionDeck[r.QuestionIndex]
       }
       r.Answers = make(map[string]string)
       r.VotedGuesses = make(map[string]map[int]string)
       r.ShuffledAnswers = []AnswerWithIndex{}
       r.Mutex.Unlock()
       broadcastGameState(r)
}

// Calcula pontuação
func (r *Room) CalculateScores() {
       r.Mutex.RLock()
       defer r.Mutex.RUnlock()
       correct := make(map[int]string)
       for _, ans := range r.ShuffledAnswers {
	       correct[ans.DisplayID] = ans.RealID
       }
       for playerId, guesses := range r.VotedGuesses {
	       for answerDisplayId, guessedPlayerId := range guesses {
		       if correctPlayerId, exists := correct[answerDisplayId]; exists && correctPlayerId == guessedPlayerId {
			       if player, exists := r.Players[playerId]; exists {
				       player.Score += 10
			       }
		       }
	       }
       }
       for _, p := range r.Players {
	       if p.Score >= 50 {
		       r.Mutex.RUnlock()
		       r.Mutex.Lock()
		       r.State = "GAME_OVER"
		       r.Mutex.Unlock()
		       r.Mutex.RLock()
		       break
	       }
       }
}
package main
