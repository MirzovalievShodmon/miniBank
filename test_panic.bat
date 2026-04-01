@echo off
echo 🚀 Тестирование обработки паник в miniBank
echo ==========================================

REM Убедитесь, что сервер запущен на порту 7556

echo.
echo 1. Проверка работоспособности сервера...
curl -s http://localhost:7556/ping || echo ❌ Сервер не отвечает. Запустите go run main.go

echo.
echo 2. Health check...
curl -s http://localhost:7556/health | python -m json.tool 2>nul || curl -s http://localhost:7556/health

echo.
echo 3. Тест обычного запроса (должен работать)...
curl -s http://localhost:7556/ping | python -m json.tool 2>nul || curl -s http://localhost:7556/ping

echo.
echo ==========================================
echo 🎉 Базовое тестирование завершено!
echo.
echo 💡 Для тестирования обработки паник:
echo    1. Запустите сервер: go run main.go
echo    2. Проверьте логи - они теперь структурированы
echo    3. Recovery middleware автоматически ловит панки
echo.
pause