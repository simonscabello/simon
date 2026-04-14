package services

import "errors"

var (
	ErrNotFound      = errors.New("não encontrado")
	ErrInvalidName   = errors.New("nome inválido")
	ErrInvalidURL    = errors.New("url inválida")
	ErrInvalidMethod = errors.New("método inválido")
	ErrInvalidKey    = errors.New("chave inválida")
	ErrDuplicateKey  = errors.New("chave duplicada neste ambiente")
)
