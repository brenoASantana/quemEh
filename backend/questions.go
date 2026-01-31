package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

// --- MODELO ---
// Aqui definimos o contrato. O jogo só precisa saber que existe um método "GetNextQuestion".
type QuestionRepository interface {
	LoadQuestions(path string) error
	GetShuffledQuestions() []string
}

// --- IMPLEMENTAÇÃO (File/Memory) ---
type FileQuestionRepository struct {
	AllQuestions []string
}

// Variável global para armazenar as perguntas carregadas
var repo = &FileQuestionRepository{}

// Função para carregar o JSON (usando o sistema de arquivos embutido ou local)
func (r *FileQuestionRepository) LoadQuestions(path string) error {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileData, &r.AllQuestions)
	if err != nil {
		return err
	}

	// Semente para garantir que o aleatório seja sempre diferente
	rand.Seed(time.Now().UnixNano())
	return nil
}

func (r *FileQuestionRepository) GetShuffledQuestions() []string {
	// Cria uma cópia para não alterar a lista original
	shuffled := make([]string, len(r.AllQuestions))
	copy(shuffled, r.AllQuestions)

	// Algoritmo Fisher-Yates para embaralhar
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}
