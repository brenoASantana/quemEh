package main

func (r *Room) Run() {
	for {
		select {
		case msg := <-r.Broadcast:
			// Envia a mensagem para todos os jogadores na sala
			for _, player := range r.Players {
				player.Conn.WriteJSON(msg)
			}
		}
	}
}

// Exemplo de func para avançar o jogo
func (r *Room) NextRound() {
	// Lógica para mudar de 'QUESTION' para 'VOTING'
	// Aqui você atualizaria o r.State e enviaria um Broadcast
}
