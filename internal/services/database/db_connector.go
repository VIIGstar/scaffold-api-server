package database

type ClientInterface interface {
	Ping() error
	Connect(connectionString string) error
}
