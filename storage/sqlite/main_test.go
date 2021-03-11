package sqlite

import (
	"github.com/ipfs-force-community/venus-wallet/config"
	"github.com/ipfs-force-community/venus-wallet/storage"
	"os"
	"testing"
)

var mockRouterStore storage.RouterStore

func TestMain(m *testing.M) {
	file := "./mockSqlite.sqlit"
	os.Remove(file)
	defer os.Remove(file)
	conn, err := NewSQLiteConn(&config.DBConfig{
		Conn: file,
	})
	if err != nil {
		log.Fatal(err)
	}
	mockRouterStore = NewRouterStore(conn)
	m.Run()
}
