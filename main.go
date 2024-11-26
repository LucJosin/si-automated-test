// ChatGPT:
//
// Using the golang language and a single file create a method that represents both controller and service.
// Testable code that has user authentication with mocked data. Mid to complex user authentication like
// using unique nickname, email and minimum password size and more. Use go 1.23 with new net/http methods.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

type User struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthService struct {
	mockUsers []User
}

// Inicializa usuários simulados para teste
func NewAuthService() *AuthService {
	return &AuthService{
		mockUsers: []User{
			{"user", "test@example.com", "F9DXIK6hvuFINjmC"},
		},
	}
}

func (s *AuthService) Authenticate(user User) error {
	// Valida se a senha tem mais de 8 caracteres
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Valida se o email contem o formato correto
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(emailRegex, user.Email); !matched {
		return errors.New("invalid email format")
	}

	// Valida se o nickname tem entre 3-16 caracteres e tem o letras, números, linhas sublinhadas e hifens
	nicknameRegex := `^[a-zA-Z0-9_-]{3,16}$`
	if matched, _ := regexp.MatchString(nicknameRegex, user.Nickname); !matched {
		return errors.New("nickname must be 3-16 characters and can include letters, numbers, underscores, or hyphens")
	}

	// Verifica se o email e nickname já não existem
	for _, mockUser := range s.mockUsers {
		if mockUser.Email == user.Email {
			return errors.New("email already exists")
		}
		if mockUser.Nickname == user.Nickname {
			return errors.New("nickname already exists")
		}
	}

	return nil
}

func (s *AuthService) Register(user User) error {
	// Chama o método para autenticar o usuário
	if err := s.Authenticate(user); err != nil {
		return err
	}

	// Se o usuário é válido, salva no "banco de dados"
	s.mockUsers = append(s.mockUsers, user)
	return nil
}

func registerHandler(authService *AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Valida se o método de HTTP é POST (Criação de algo)
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Converte o JSON para golang struct
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Chama o serviço para registrar usuário
		if err := authService.Register(user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Retorna o código 201 informando que o usuário foi criado
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "User registered successfully")
	})
}

func main() {
	// Cria um serviço de autenticação
	authService := NewAuthService()

	// Cria rota "/register"
	http.Handle("/register", registerHandler(authService))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}
