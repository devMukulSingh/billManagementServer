package dbconnection

import (
	"context"
	"os"

	"github.com/devMukulSingh/billManagementServer.git/database"
	"github.com/jackc/pgx/v5"
)

var Queries* database.Queries = nil
var Ctx context.Context

func ConnectDb() error {
	dbUrl := os.Getenv("DB_URL")
	ctx := context.Background()
	Ctx = ctx;
	conn, err := pgx.Connect(Ctx, dbUrl)
	if err != nil {
		return err
	}
	// defer conn.Close(Ctx)
	queries := database.New(conn)
	Queries = queries;
	return nil
}
