package main

import (
	"fmt"
	"os"
)

func main() {
	appConfig := NewAppConfig()
	err := appConfig.Server.Start(os.Getenv("PORT"))
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}
