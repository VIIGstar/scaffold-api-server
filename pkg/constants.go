package pkg

type AppName string

const (
	// APIAppName is the visible name of the application.
	APIAppName = AppName("Scaffold-API")
	// MysqlConnectorAppName is the visible name of the application.
	MysqlConnectorAppName = AppName("MysqlConnector")
)

func (name AppName) ToString() string {
	return string(name)
}
