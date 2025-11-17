# 🔗 Интеграция proto-common в существующие сервисы

## Пример: book-service

### 1. Установите модуль

```bash
cd /path/to/book-service
go get github.com/yourusername/proto-common@latest
```

### 2. Обновите books.proto

Измените импорт в `internal/proto/book/books.proto`:

```protobuf
// Было:
import "internal/proto/common/common.proto";

// Стало:
import "common/common.proto";
```

### 3. Обновите Makefile

Добавьте путь к proto-common модулю:

```makefile
# Получаем путь к proto-common из go.mod
PROTO_COMMON_PATH := $(shell go list -m -f '{{.Dir}}' github.com/yourusername/proto-common)

.PHONY: gen
gen:
	protoc -I . -I $(PROTO_COMMON_PATH) -I C:\Users\wakar\googleapis-master \
		--go_out=. --go_opt=module=book-service \
		--go-grpc_out=. --go-grpc_opt=module=book-service \
		--grpc-gateway_out=. --grpc-gateway_opt=logtostderr=true --grpc-gateway_opt=module=book-service \
		--openapi_out=output_mode=source_relative:. \
		internal/proto/book/books.proto
	echo servers: >> internal/proto/book/books.openapi.yaml
	echo     - url: http://localhost:8080 >> internal/proto/book/books.openapi.yaml
```

### 4. Удалите локальный common

```bash
# Удаляем локальные proto файлы
rm -rf internal/proto/common

# Обновляем импорты в Go коде
# Было:
import "book-service/internal/proto/common"

# Стало:
import "github.com/yourusername/proto-common/common"
```

### 5. Обновите Go импорты

Замените все импорты в `.go` файлах:

**PowerShell команда для автозамены:**

```powershell
$files = Get-ChildItem -Path . -Recurse -Filter *.go | Where-Object { $_.FullName -notlike "*\vendor\*" }

foreach ($file in $files) {
    $content = Get-Content $file.FullName -Raw
    $newContent = $content -replace '"book-service/internal/proto/common"', '"github.com/yourusername/proto-common/common"'
    if ($content -ne $newContent) {
        Set-Content -Path $file.FullName -Value $newContent -NoNewline
        Write-Host "Updated: $($file.FullName)"
    }
}
```

### 6. Регенерируйте proto и проверьте

```bash
make gen
go mod tidy
go build ./cmd/main.go
```

## ✅ Готово!

Теперь ваш сервис использует централизованный proto-common модуль.

## 🔄 Обновление proto-common

Когда в proto-common появятся изменения:

```bash
# Обновите зависимость
go get github.com/yourusername/proto-common@latest

# Или конкретную версию
go get github.com/yourusername/proto-common@v0.2.0

# Регенерируйте proto
make gen

# Обновите зависимости
go mod tidy
```

## 📋 Чеклист миграции

- [ ] Установлен proto-common модуль
- [ ] Обновлён импорт в books.proto
- [ ] Обновлён Makefile с путём к proto-common
- [ ] Удалена локальная директория internal/proto/common
- [ ] Обновлены все Go импорты
- [ ] Регенерированы proto файлы
- [ ] Проект компилируется без ошибок
- [ ] Удалены старые файлы (update-module.ps1, README.md в common/)

