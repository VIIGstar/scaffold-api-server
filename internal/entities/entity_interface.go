package entities

type EntityInterface interface {
	GetUniqueIndexes() map[string][]string
}
