package auth

import (
	"database/sql"
	"encoding/json"
	"io"
	entities "kraken/domain/entities"
	"kraken/domain/usecases"
	"kraken/infrastructure/modules/impl/http_error"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

const (
	CtxUserKey = "auth-ctx-user-data"
	JWTSecret  = "my_secret_key" // Essa chave deve ser carregada a partir do `settings_loader` em um projeto real.
)

type AuthModule struct {
	Db          *sql.DB
	Cookie      *securecookie.SecureCookie
	AuthUseCase usecases.AuthUseCase
}

// Setup configura as rotas de autenticação
func (a *AuthModule) Setup(router *mux.Router) {
	// Rotas públicas
	router.HandleFunc("/login", a.loginHandler).Methods(http.MethodPost)
	router.HandleFunc("/register", a.registerUser).Methods(http.MethodPost)

}

func (a *AuthModule) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Ler o corpo da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[loginHandler] Erro ao ler a requisição: ", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	// Desserializar o JSON da requisição para as credenciais de login
	var request entities.LoginCredentials
	if err := json.Unmarshal(body, &request); err != nil {
		log.Println("[loginHandler] Erro ao decodificar a requisição: ", err)
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Chamar o caso de uso de autenticação para efetuar o login
	user, token, err := a.AuthUseCase.Login(r.Context(), request)
	if err != nil {
		log.Println("[loginHandler] Erro ao realizar login: ", err)
		http_error.HandleError(w, err)
		return
	}

	log.Println("token", token)
	log.Println("user.ID", user.ID)
	// Configurar o valor do token no cookie
	value := map[string]string{
		"token":   token,
		"user_id": strconv.FormatInt(user.ID, 10), // Utilize o ID do usuário retornado
	}

	encoded, err := a.Cookie.Encode("auth_token", value)
	if err != nil {
		log.Println("[loginHandler] Erro ao codificar o cookie: ", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	// Definir o cookie com as configurações de segurança
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(360 * time.Hour),
	}
	http.SetCookie(w, cookie)

	// Responder com sucesso
	response := entities.NewSuccessfulRequest()
	if b, err := json.Marshal(response); err != nil {
		log.Println("[loginHandler] Erro ao codificar a resposta: ", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

type userCreatePostRequest struct {
	User entities.User `json:"user"`
}

func (a *AuthModule) registerUser(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[registerUser] Error ReadAll", err)
		http_error.HandleError(w, http_error.NewBadRequestError(http_error.InvalidParameter))
		return
	}

	var request userCreatePostRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		log.Println("[registerUser] Error Unmarshal", err)
		http_error.HandleError(w, http_error.NewBadRequestError(http_error.InvalidParameter))
		return
	}

	user := request.User
	err = a.AuthUseCase.RegisterUser(r.Context(), user)
	if err != nil {
		log.Println("[registerUser] Error RegisterUser", err)
		http_error.HandleError(w, err)
		return
	}

	b, err = json.Marshal(entities.NewSuccessfulRequest())
	if err != nil {
		log.Println("[registerUser] Error Marshal", err)
		http_error.HandleError(w, http_error.NewUnexpectedError(http_error.Unexpected))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println("[registerUser] Error Write", err)
		http_error.HandleError(w, http_error.NewUnexpectedError(http_error.Unexpected))
		return
	}
}
