package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"kraken/infrastructure"
	"kraken/settings_loader"
)

func init() {
	// Carrega as variáveis do .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Erro ao carregar .env")
	}
}

func main() {
	router := mux.NewRouter()

	// Carregar as configurações
	settings := settings_loader.NewSettingsLoader()

	// Configura o projeto chamando o Setup da infraestrutura
	setupConfig, err := setup.Setup(router, settings)
	if err != nil {
		log.Fatalf("Erro ao configurar a infraestrutura: %v", err)
	}

	defer setupConfig.CloseDB()

	// Inicia o servidor
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
