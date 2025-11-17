# Proto Common

Общие proto-файлы для микросервисов.

## 📦 Установка

```bash
go get github.com/yourusername/proto-common
```

## 🚀 Использование

### В вашем proto файле:

```protobuf
syntax = "proto3";

package yourservice;

import "common/common.proto";

message YourRequest {
  common.AuthContext auth = 1;
  // ... ваши поля
}

message YourListRequest {
  common.AuthContext auth = 1;
  common.PageRequest page = 2;
}

message YourListResponse {
  repeated YourItem items = 1;
  common.PageResponse page = 2;
}
```

### В Go коде:

```go
import (
    "github.com/yourusername/proto-common/common"
)

func MyHandler(auth *common.AuthContext) {
    if auth.Valid {
        // ...
    }
}
```

## 🔧 Генерация

После изменения proto файлов:

```bash
make gen
```

## 📝 Содержимое

### AuthContext
Контекст авторизации, передаваемый от auth-service:
- `valid` - валиден ли токен
- `message` - сообщение об ошибке (если не валиден)
- `id` - UUID пользователя
- `email` - email пользователя
- `role` - роль пользователя

### PageRequest / PageResponse
Курсорная пагинация:
- `cursor` - курсор для следующей страницы
- `limit` - количество элементов на странице

## 🎯 Настройка для своего проекта

1. Форкните репозиторий или создайте свой
2. Замените `github.com/yourusername/proto-common` на ваш путь:
   - В `go.mod`
   - В `common/common.proto` (в `option go_package`)
   - В `Makefile`
3. Запустите `make gen`
4. Закоммитьте и запушьте в GitHub
5. Используйте в своих сервисах через `go get`

## 📋 Первая настройка

```bash
# 1. Замените username на ваш GitHub username
sed -i 's/yourusername/Kirimatt/g' go.mod common/common.proto Makefile README.md

# Windows PowerShell:
$files = @('go.mod', 'common/common.proto', 'Makefile', 'README.md')
foreach ($file in $files) {
    (Get-Content $file) -replace 'yourusername', 'Kirimatt' | Set-Content $file
}

# 2. Генерируем proto
make gen

# 3. Инициализируем Git
git init
git add .
git commit -m "Initial commit: common proto files"

# 4. Добавляем remote и пушим
git remote add origin git@github.com:Kirimatt/proto-common.git
git branch -M main
git push -u origin main

# 5. Создаём тег версии
git tag v0.1.0
git push origin v0.1.0
```

## 🔄 Обновление в сервисах

```bash
# В вашем сервисе
go get -u github.com/yourusername/proto-common@latest

# Или конкретную версию
go get github.com/yourusername/proto-common@v0.1.0
```

