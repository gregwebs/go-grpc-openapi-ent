// dal = Data Access Layer
package dal

import (
	"fmt"
	"os"

	ent "github.com/gregwebs/go-grpc-openapi-ent/ent"
)

type ConnectConf struct {
	DB       string
	User     string
	Password string
}

func ConnectDSN(connType string, conf ConnectConf) string {
	dbUser := conf.User
	if dbUser == "" {
		dbUser = os.Getenv("DB_USER")
	}
	dbName := conf.DB
	if dbName == "" {
		dbName = os.Getenv("DB_NAME")
	}
	return fmt.Sprintf("%s://%s:%s@localhost:5432/%s?sslmode=disable", connType, dbUser, conf.Password, dbName)
}

func ConnectKV(conf ConnectConf) string {
	// fmt.Println(sql.Drivers())
	dbUser := conf.User
	if dbUser == "" {
		dbUser = os.Getenv("DB_USER")
	}
	dbName := conf.DB
	if dbName == "" {
		dbName = os.Getenv("DB_NAME")
	}

	return "port=5432 sslmode=disable TimeZone=America/Chicago host=localhost" +
		" user=" + dbUser + " password=" + conf.Password + " dbname=" + dbName
}

func Connect(conf ConnectConf) (*ent.Client, error) {
	// fmt.Println(ConnectDSN("postgres", conf))
	return ent.Open("postgres", ConnectDSN("postgres", conf))
}
