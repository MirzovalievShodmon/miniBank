#!/bin/bash

# Тестовый скрипт для проверки системы аутентификации miniBank
# Убедитесь, что сервер запущен на порту 7556

BASE_URL="http://localhost:7556"
EMAIL="test@example.com"
PASSWORD="secure123"

echo "🚀 Тестирование системы аутентификации miniBank"
echo "=================================================="

# 1. Health check
echo "1. Проверка работоспособности сервера..."
curl -s "$BASE_URL/ping" && echo " ✅ Ping OK"

echo -e "\n2. Health check..."
curl -s "$BASE_URL/health" | python -m json.tool

# 2. Регистрация пользователя
echo -e "\n3. Регистрация нового пользователя..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\":\"$EMAIL\",
    \"password\":\"$PASSWORD\",
    \"first_name\":\"Тестовый\",
    \"last_name\":\"Пользователь\"
  }")

echo "Ответ регистрации:"
echo "$REGISTER_RESPONSE" | python -m json.tool

# 3. Вход в систему
echo -e "\n4. Вход в систему..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\":\"$EMAIL\",
    \"password\":\"$PASSWORD\"
  }")

echo "Ответ входа:"
echo "$LOGIN_RESPONSE" | python -m json.tool

# Извлекаем JWT токен
TOKEN=$(echo "$LOGIN_RESPONSE" | python -c "import sys, json; print(json.load(sys.stdin)['token'])" 2>/dev/null)

if [ -z "$TOKEN" ] || [ "$TOKEN" == "None" ]; then
    echo "❌ Не удалось получить JWT токен"
    exit 1
fi

echo "✅ JWT токен получен: ${TOKEN:0:20}..."

# 4. Проверка профиля
echo -e "\n5. Получение профиля пользователя..."
PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/profile" \
  -H "Authorization: Bearer $TOKEN")

echo "Профиль пользователя:"
echo "$PROFILE_RESPONSE" | python -m json.tool

# 5. Создание счета
echo -e "\n6. Создание нового счета..."
ACCOUNT_RESPONSE=$(curl -s -X POST "$BASE_URL/accounts" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\":\"Тестовый счет\"
  }")

echo "Созданный счет:"
echo "$ACCOUNT_RESPONSE" | python -m json.tool

# Извлекаем ID созданного счета
ACCOUNT_ID=$(echo "$ACCOUNT_RESPONSE" | python -c "import sys, json; data=json.load(sys.stdin); print(data['id'] if 'id' in data else '')" 2>/dev/null)

if [ -n "$ACCOUNT_ID" ] && [ "$ACCOUNT_ID" != "None" ]; then
    echo "✅ Счет создан с ID: $ACCOUNT_ID"
    
    # 6. Пополнение счета
    echo -e "\n7. Пополнение счета на 10000 рублей..."
    TOPUP_RESPONSE=$(curl -s -X POST "$BASE_URL/accounts/$ACCOUNT_ID/top-up" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{\"amount\":10000}")
    
    echo "Результат пополнения:"
    echo "$TOPUP_RESPONSE" | python -m json.tool
    
    # 7. Получение списка счетов
    echo -e "\n8. Получение списка счетов..."
    ACCOUNTS_RESPONSE=$(curl -s -X GET "$BASE_URL/accounts" \
      -H "Authorization: Bearer $TOKEN")
    
    echo "Список счетов пользователя:"
    echo "$ACCOUNTS_RESPONSE" | python -m json.tool
    
    # 8. История транзакций
    echo -e "\n9. История транзакций счета..."
    TRANSACTIONS_RESPONSE=$(curl -s -X GET "$BASE_URL/accounts/$ACCOUNT_ID/transactions" \
      -H "Authorization: Bearer $TOKEN")
    
    echo "Транзакции:"
    echo "$TRANSACTIONS_RESPONSE" | python -m json.tool
    
else
    echo "❌ Не удалось создать счет"
fi

# 9. Проверка защищенного endpoint'а без токена
echo -e "\n10. Проверка защиты endpoint'ов (запрос без токена)..."
UNAUTH_RESPONSE=$(curl -s -X GET "$BASE_URL/accounts")
echo "Ответ без авторизации:"
echo "$UNAUTH_RESPONSE" | python -m json.tool

echo -e "\n=================================================="
echo "🎉 Тестирование завершено!"