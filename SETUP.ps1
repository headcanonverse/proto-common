# Скрипт первоначальной настройки proto-common репозитория

Write-Host "=== Настройка proto-common репозитория ===" -ForegroundColor Green

# Запрашиваем GitHub username
$username = Read-Host "Введите ваш GitHub username"

if ([string]::IsNullOrWhiteSpace($username)) {
    Write-Host "Ошибка: username не может быть пустым" -ForegroundColor Red
    exit 1
}

Write-Host "`nЗаменяем 'yourusername' на '$username'..." -ForegroundColor Yellow

# Файлы для замены
$files = @(
    'go.mod',
    'common/common.proto',
    'Makefile',
    'README.md'
)

foreach ($file in $files) {
    if (Test-Path $file) {
        (Get-Content $file) -replace 'yourusername', $username | Set-Content $file
        Write-Host "  ✓ $file" -ForegroundColor Green
    } else {
        Write-Host "  ✗ $file (не найден)" -ForegroundColor Red
    }
}

Write-Host "`nГенерируем proto файлы..." -ForegroundColor Yellow
make gen

Write-Host "`nИнициализируем Git репозиторий..." -ForegroundColor Yellow
git init
git add .
git commit -m "Initial commit: common proto files"

Write-Host "`n=== Готово! ===" -ForegroundColor Green
Write-Host "`nСледующие шаги:" -ForegroundColor Cyan
Write-Host "1. Создайте репозиторий на GitHub: https://github.com/new" -ForegroundColor White
Write-Host "   Имя: proto-common" -ForegroundColor White
Write-Host "`n2. Добавьте remote и запушьте:" -ForegroundColor White
Write-Host "   git remote add origin git@github.com:$username/proto-common.git" -ForegroundColor Gray
Write-Host "   git branch -M main" -ForegroundColor Gray
Write-Host "   git push -u origin main" -ForegroundColor Gray
Write-Host "`n3. Создайте первый релиз:" -ForegroundColor White
Write-Host "   git tag v0.1.0" -ForegroundColor Gray
Write-Host "   git push origin v0.1.0" -ForegroundColor Gray
Write-Host "`n4. Используйте в своих сервисах:" -ForegroundColor White
Write-Host "   go get github.com/$username/proto-common@v0.1.0" -ForegroundColor Gray

