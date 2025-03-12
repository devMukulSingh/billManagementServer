package dbconnection

import (
	"context"
	"os"

	"github.com/devMukulSingh/billManagementServer.git/database"
	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Queries* database.Queries = nil
var Ctx context.Context
var Connection  *pgxpool.Pool

func ConnectDb() error {
	dbUrl := os.Getenv("DB_URL")
	ctx := context.Background()
	Ctx = ctx;
	pool,err := pgxpool.New(Ctx,dbUrl)
	if err != nil {
		return err
	}
	  Connection = pool
	
	queries := database.New(pool)
	Queries = queries;
	return nil
}
