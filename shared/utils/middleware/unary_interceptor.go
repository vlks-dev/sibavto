package middleware

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	jwtService *auth.JWTService
}

func NewAuthInterceptor(jwtService *auth.JWTService) *AuthInterceptor {
	return &AuthInterceptor{jwtService: jwtService}
}

func (a *AuthInterceptor) AuthorizationInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no metadata provided")
	}

	// Получаем токен из метаданных
	tokens := md["authorization"]
	if len(tokens) == 0 {
		return nil, errors.New("authorization token not provided")
	}
	token := tokens[0]

	// Проверяем токен и получаем claims
	claims, err := a.jwtService.ValidateJWT(token) // Используем метод `ValidateJWT` у `JWTService`
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Определяем разрешения для ролей
	permissions := map[string][]string{
		"admin":     {"RegisterUser", "LoginUser", "GrantAdminRights", "GrantRole", "GetUserData"},
		"logistics": {"GetUserData", "CreateRequest", "GetLogisticsData"},
		"user":      {"RegisterUser", "LoginUser", "GetUserData"},
	}

	// Проверяем, имеет ли пользователь хотя бы одну роль, которая разрешает доступ к методу
	method := info.FullMethod
	hasPermission := false

	for role := range claims.Roles {
		if allowedMethods, ok := permissions[role]; ok && contains(allowedMethods, method) {
			hasPermission = true
			break
		}
	}

	// Если ни одна роль пользователя не даёт доступ, возвращаем ошибку
	if !hasPermission {
		return nil, errors.New("permission denied: insufficient role privileges")
	}

	// Продолжаем обработку запроса
	return handler(ctx, req)
}

// Вспомогательная функция для проверки наличия метода в списке разрешенных
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
