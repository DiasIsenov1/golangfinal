package service_test

import (
	"authh/pkg/utils"
	"context"
	"net/http"
	"testing"

	"authh/pkg/db"
	_ "authh/pkg/models"
	"authh/pkg/pb"
	"authh/pkg/services"
)

func TestRegisterUser(t *testing.T) {
	// Создаем временную базу данных для тестов.
	tmpDB := db.Init("postgres://postgres:samsungtab7@localhost:5432/auth-service")

	// Инициализируем обработчик базы данных.
	handler := db.Handler{DB: tmpDB.DB}

	// Создаем новый сервер службы с использованием временной базы данных.
	server := services.Server{H: handler}

	// Задаем данные для нового пользователя.
	req := &pb.RegisterRequest{
		Email:    "test@example.com",
		Password: "testPassword",
	}

	// Регистрируем пользователя с использованием Register службы.
	resp, err := server.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Проверяем успешность регистрации пользователя.
	if resp.Status != http.StatusCreated {
		t.Errorf("Expected status %d, but got %d", http.StatusCreated, resp.Status)
	}

	// Дополнительные проверки, если необходимо.
	// Например, можно проверить, что пользователь успешно добавлен в базу данных.

	// Закрываем временную базу данных после завершения теста.
	tmpDB.Close()
}
func TestLoginUser(t *testing.T) {
	// Создаем временную базу данных для тестов.
	tmpDB := db.Init("postgres://postgres:samsungtab7@localhost:5432/auth-service")

	// Инициализируем обработчик базы данных.
	handler := db.Handler{DB: tmpDB.DB}

	// Создаем новый сервер службы с использованием временной базы данных.
	server := services.Server{H: handler}

	// Задаем данные для существующего пользователя.
	req := &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "testPassword",
	}

	// Пытаемся войти в систему с использованием Login службы.
	resp, err := server.Login(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}

	// Проверяем успешность входа пользователя.
	if resp.Status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, resp.Status)
	}

	// Дополнительные проверки, если необходимо.
	// Например, можно проверить, что токен доступа был успешно создан и возвращен.

	// Закрываем временную базу данных после завершения теста.
	tmpDB.Close()
}

func TestValidateToken(t *testing.T) {
	// Инициализируем временную базу данных для тестов.
	tmpDB := db.Init("postgres://postgres:samsungtab7@localhost:5432/auth-service")

	// Закрываем временную базу данных после завершения теста.

	// Инициализируем обработчик базы данных.
	handler := db.Handler{DB: tmpDB.DB}

	// Создаем новый сервер службы с использованием временной базы данных.
	server := services.Server{H: handler}

	// Генерируем токен для существующего пользователя.
	token, _ := utils.GenerateTokenMock("test@example.com")

	// Проверяем валидность токена с использованием Validate службы.
	resp, err := server.Validate(context.Background(), &pb.ValidateRequest{Token: token})
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Проверяем успешность валидации токена.
	if resp.Status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, resp.Status)
	}
}
