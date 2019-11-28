package sqlstore_test

import (
	"os"
	"testing"
)

//... вызывается 1 раз перед всеми тестами в конкретном пакете
//йен переменная ???


var (
	databaseURL string
)
func TestMain(m *testing.M)  {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost port=5431 dbname=restapi_test sslmode=disable"
	}
	os.Exit(m.Run())
}
