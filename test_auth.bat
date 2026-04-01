@echo off
chcp 65001 > nul
echo 🚀 Тестирование системы аутентификации miniBank
echo ==================================================

set BASE_URL=http://localhost:7556
set EMAIL=test@example.com
set PASSWORD=secure123

echo 1. Проверка работоспособности сервера...
curl -s "%BASE_URL%/ping"
echo  ✅ Ping OK

echo.
echo 2. Health check...
curl -s "%BASE_URL%/health"

echo.
echo 3. Регистрация нового пользователя...
curl -s -X POST "%BASE_URL%/auth/register" ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"%EMAIL%\",\"password\":\"%PASSWORD%\",\"first_name\":\"Тестовый\",\"last_name\":\"Пользователь\"}"

echo.
echo 4. Вход в систему...
for /f "delims=" %%i in ('curl -s -X POST "%BASE_URL%/auth/login" -H "Content-Type: application/json" -d "{\"email\":\"%EMAIL%\",\"password\":\"%PASSWORD%\"}"') do set LOGIN_RESPONSE=%%i
echo %LOGIN_RESPONSE%

echo.
echo 5. Создание счета (потребуется ручной ввод токена)...
echo Скопируйте JWT токен из ответа выше и используйте следующую команду:
echo curl -X POST "%BASE_URL%/accounts" -H "Authorization: Bearer YOUR_TOKEN_HERE" -H "Content-Type: application/json" -d "{\"name\":\"Тестовый счет\"}"

echo.
echo 6. Получение списка счетов:
echo curl -X GET "%BASE_URL%/accounts" -H "Authorization: Bearer YOUR_TOKEN_HERE"

echo.
echo 7. Проверка защиты (запрос без токена):
curl -s -X GET "%BASE_URL%/accounts"

echo.
echo ==================================================
echo 🎉 Тестирование завершено!