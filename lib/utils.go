package utils
import "os"

func GetBaseUrlClient() string{
	env := os.Getenv("ENV");
	if(env=="production"){
		return "https://invoice-management-c.vercel.app"
	}
	return "http://localhost:5173"
}