package valkeyCache

import (
	"context"
	"log"
	"os"

	"github.com/valkey-io/valkey-go"
)

var Client valkey.Client
var Ctx = context.Background()

func Connect() {

	url := os.Getenv("VALKEY_URL")

	valkeyClient, err := valkey.NewClient(valkey.MustParseURL(url))
	if err != nil {
		log.Print(err)
		panic(err)
	}

	log.Print("Valkey connected")

	Client = valkeyClient
}

func SetValue(key string, value []byte) {
	if err := Client.Do(Ctx, Client.B().Set().Key(key).Value(string(value)).Nx().Build()); err.Error() != nil {
		log.Printf("Error in setting value in valkey %s ", err.Error())
	}
}

func GetValue(key string) string{
	value, err := Client.Do(Ctx, Client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		if(err.Error()=="valkey nil message"){
			return "";
		}
		log.Printf("Error in getting value in valkey %s ", err)
	}
	return value
}

func Revalidate(key string){
	if err := Client.Do(Ctx,Client.B().Del().Key(key).Build()).Error(); err!=nil{
		log.Printf("Error in revalidating cache of: %s", key)
	}
}