# goph-keeper
Выпускной проект курса "Продвинутый Go разработчик" (Яндекс Практикум). Содержит следующие компоненты:
- сервис хранения секретов (`keeper`);
- клиент командной строки (`keepctl`);
- библиотека для работы с сервисом по `gRPC API` (`goph`).

<img src="./docs/arch/Product arch.drawio.svg">

## Структура проекта
```
├── .dockerignore
├── .gitignore
├── .github
│   └── workflows            # конфигурация Github CI
├── .golangci-lint
├── .pre-commit-config.yaml
├── LICENSE
├── Makefile
├── README.md
├── api
│   └── proto                # спецификация gRPC API
├── build
│   └── docker               # конфигурация для упаковки проекта в docker
├── cmd
├── deployments              # конфигурация для разворачивания проекта
│   ├── docker-compose.yaml
│   ├── keepctl.env          # переменные окружения для клиента keepctl, используются для локального запуска клиента
│   └── keeper.env           # переменные окружения для сервиса keeper, используются для запуска в docker compose
├── dist                     # скомпилированные исполняемые файлы клиента и сервера
├── docs
│   ├── api                  # документация gRPC API
│   └── arch                 # архитектура проекта
├── internal
│   ├── libraries            # общие внутренние библиотеки клиента и сервера
│   │   ├── creds            # общие типы безопасного использования паролей внутри приложения
│   │   └── gophtest         # набор фикстур и хэлперов для тестирования проекта, не предполагает покрытие тестами
│   ├── keepctl              # код клиента командной строки
│   │   ├── app              # реализация клиентского приложения keepctl
│   │   ├── config           # конфигурация клиента
│   │   ├── controller       # слой контроллеров по чистой архитектуре, содержит реализацию интерфейса командной строки
│   │   ├── entity           # ядро по чистой архитектуре, содержит основные структуры данных
│   │   ├── infra            # внешний инфраструктурный слой по чистой архитектуре, инкапсулирует клиентское соединение gRPC и т.п.
│   │   ├── repo             # слой репозиториев по чистой архитектуре, фасад для работы с данными из внешних источников
│   │   └── usecase          # слой бизнес логики по чистой архитектуре
│   └── keeper               # код сервиса хранения паролей
│       ├── app              # основная точка входа для запуска сервиса keeper
│       ├── config           # конфигурация сервиса
│       ├── controller       # слой контроллеров по чистой архитектуре, содержит реализацию API хэндлеров
│       ├── entity           # ядро по чистой архитектуре, содержит основные структуры данных
│       ├── infra            # внешний инфраструктурный слой по чистой архитектуре, инкапсулирует базу данных, gRPC сервер и т.п.
│       ├── repo             # слой репозиториев по чистой архитектуре, фасад для работы с данными из внешних источников
│       └── usecase          # слой бизнес логики по чистой архитектуре
├── migrations               # код миграций для формирования базы данных
├── pkg
│   └── goph                 # библиотека для работы с сервисом keeper по gRPC API
├── scripts                  # вспомогательные скрипты для сборки и работы проекта
├── ssl                      # конфигурация для выпуска SSL сертификатов
├── go.mod
└── go.sum
```

## Запуск сервиса keeper
1. Установите `docker compose` по [инструкции](https://docs.docker.com/compose/install/).
2. Сгенерируйте сертификаты для клиента и сервера:
    ```bash
    make ssl
    ```
3. Для сборки и запуска сервиса с помощью `docker compose` выполните команду:
    ```bash
    make run
    ```
4. Для остановки сервиса выполните команду:
    ```bash
    make stop
    ```

## Конфигурация сервиса keeper
Переменные окружения для сервиса `keeper` описаны в файле `deployments/keeper.env`.  
(!) Опции командной строки имеют более высокий приоритет по сравнению с переменными окружения.

## Сборка клиента
1. Сгенерируйте сертификаты для клиента и сервера:
    ```bash
    make ssl
    ```
2. Соберите клиент для всех поддерживаемых платформ, выполнив команду:
    ```bash
    make keepctl
    ```

## Разработка
### Генерация кода для gRPC
1. Установите `protoc` по [инструкции](https://grpc.io/docs/protoc-installation/).
2. Установите плагины для `gRPC`:
    ```bash
    make install-tools
    ```

### Контроль качества кода
1. Установите `gofumpt` (улучшенное форматирование кода) по [инструкции](https://github.com/mvdan/gofumpt#installation).
2. Установите линтер `golangci-lint` по [инструкции](https://golangci-lint.run/usage/install/).
3. Установите линтер `shellcheck` (проверка `bash`-скриптов) по [инструкции](https://github.com/koalaman/shellcheck#installing).
4. Установите линтер `hadolint` (проверка `Dockerfile`) по [инструкции](https://github.com/hadolint/hadolint#install).
5. Установите `pre-commit` (запуск линтеров перед коммитом) по [инструкции](https://pre-commit.com/#install), затем выполните команду:
    ```bash
    make install-tools
    ```

### Работа с базой данных
1. Для ручной работы с миграциями (вне контейнера docker) установите утилиту `golang-migrate`:
    ```bash
    go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```
2. Для применения миграций выполните команду:
    ```bash
    migrate -database ${DATABASE_DSN} -path ./migrations up
    ```
3. Для возврата базы данных в первоначальное состояние выполните команду:
    ```bash
    migrate -database ${DATABASE_DSN} -path ./migrations down -all
    ```

## Лицензия
Copyright (c) 2023 Alexander Kurbatov

Лицензировано по [GPLv3](LICENSE).
