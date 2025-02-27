package valkeyCache

import (
	"context"
	"log"
	"os"
	"github.com/valkey-io/valkey-go"
)

var Client valkey.Client
var Ctx = context.Background()

func Connect() (error) {

	url := os.Getenv("VALKEY_URL")
	valkeyClient, err := valkey.NewClient(valkey.MustParseURL(url))
	if err != nil {
		return err
	}
	log.Print("Valkey connected")
	Client = valkeyClient
	return nil
}

func SetValue(key string, value []byte) error{
	if result := Client.Do(Ctx, Client.B().Set().Key(key).Value(string(value)).Nx().Build()); result.Error() != nil {
		return result.Error()
	}
	return nil
}

func GetValue(key string) (string,error){
	value, err := Client.Do(Ctx, Client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		return "",err
	}
	return value,nil
}

func Revalidate(key string)error{
	if err := Client.Do(Ctx,Client.B().Del().Key(key).Build()).Error(); err!=nil{
		return err
	}
	return nil
}