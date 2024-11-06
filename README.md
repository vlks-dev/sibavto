# sibavto
Go_gGRPC Application 
project-root/
├── auth-service/                   # Сервис аутентификации
│   ├── proto/
│   │   ├── auth.proto               # Протофайл для AuthService
│   │   └── auth.pb.go               # Сгенерированный Go-код из auth.proto
│   ├── internal/
│   │   ├── handlers/                # Обработчики запросов
│   │   │   └── auth_server.go       # Реализация AuthService (регистрация, вход, генерация JWT)
│   │   ├── services/
│   │   │   └── jwt_service.go       # Логика генерации и проверки JWT
│   ├── main.go                      # Точка входа для запуска Auth-сервера
│   └── config.yaml                  # Конфигурационный файл (например, секретный ключ JWT, настройки БД)
├── logistics-service/               # Бизнес-сервис логистики
│   ├── proto/
│   │   ├── logistics.proto          # Протофайл для LogisticsService
│   │   └── logistics.pb.go          # Сгенерированный Go-код из logistics.proto
│   ├── internal/
│   │   ├── handlers/
│   │   │   └── logistics_server.go  # Реализация бизнес-логики LogisticsService
│   │   ├── services/
│   │   │   └── auth_client.go       # gRPC клиент для запроса JWT в Auth-сервис
│   ├── main.go                      # Точка входа для запуска Logistics-сервера
│   └── config.yaml                  # Конфигурация (например, настройки БД и адрес Auth-сервиса)
├── shared/
│   ├── proto/
│   │   ├── auth/                    # Общие протофайлы для взаимодействия сервисов
│   │   ├── auth.proto               # Импортируется в других сервисах
│   │   └── logistics.proto          # Импортируется в других сервисах
│   └── utils/
│       └── middleware/              # Общие middlewares для аутентификации и логирования
├── client/                          # Пример клиента для тестирования микросервисов
│   └── main.go                      # CLI клиент для тестирования gRPC запросов
├── docker-compose.yml               # Docker Compose для запуска сервисов и Redis
├── Makefile                         # Сценарии сборки и тестирования
├── go.mod                           # Файл для управления зависимостями Go
└── README.md                        # Документация проекта
