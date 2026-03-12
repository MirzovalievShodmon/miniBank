package controller

import (
	"fmt"
	"strconv"

	"github.com/MirzovalievShodmon/miniBank.git/internal/service"
)

func GetAllAccounts() {
	accounts, err := service.GetAllAccounts()
	if err != nil {
		fmt.Printf("Ошибка при получении списка счетов: %s\n", err.Error())
		fmt.Println("==================================================")
		fmt.Println()
		return
	}

	fmt.Println("===== Список счетов =====")

	if len(accounts) == 0 {
		fmt.Println("В базе данных пока нет ни одного открытого счета.")
		fmt.Println("==================================================")
		fmt.Println()
		return
	}

	for _, a := range accounts {
		fmt.Printf("ID: %d, Баланс: %.2f, Владелец: %s\n", a.ID, a.Balance, a.Owner)
	}
	fmt.Println("==================================================")
	fmt.Println()
}

func TopUpAccount() {
	fmt.Println("===== Пополнение счета =====")
	fmt.Println("Введите ID счета, который хотите пополнить")
	idStr := readInput()
	accountID, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Ошибка: ID должен быть числом.")
		fmt.Println("==================================================")
		fmt.Println()
		return
	}

	if accountID <= 0 {
		fmt.Println("ID не можеть быть меньше или равным нулю")
		fmt.Println("==================================================")
		fmt.Println()
		return
	}

	fmt.Println("Введите сумму на которую хотите пополнить:")
	amountStr := readInput()
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("Ошибка: Неверный формат суммы. Используйте цифры")
		fmt.Println("==================================================")
		fmt.Println()
		return
	}

	if amount <= 0 {
		fmt.Println("Сумма не можеть быть меньше или равной нулю")
		fmt.Println("==================================================")
		fmt.Println()
		return
	}

	err = service.TopUpAccount(accountID, amount)
	if err != nil {
		fmt.Printf("Ошибка при пополнении счета: %s\n", err.Error())
		fmt.Println("==================================================")
		fmt.Println()
		return
	}

	fmt.Printf("Указанный счет успешно пополнен на %.2f сомони\n", amount)
	fmt.Println("==================================================")
	fmt.Println()
}

func WithdrawAccount() {
	fmt.Println("===== Снятие со счета =====")
	fmt.Println("Введите ID счета с которого хотите снять деньги:")
	idStr := readInput()
	accountID, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Ошибка: ID должен быть числом")
		fmt.Println("==================================================")
		fmt.Println()
		return
	}

	if accountID <= 0 {
		fmt.Println("ID не можеть быть меньше или равным нулю")
		fmt.Println("==================================================")
		fmt.Println()
		return
	}

	fmt.Println("Введите сумму которую хотите снять:")
	amountStr := readInput()
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("Ошибка: Неверный формат суммы. Используйте цифры")
		fmt.Println("==================================================")
		fmt.Println()
		return
	}

	if amount <= 0 {
		fmt.Println("Сумма не может быть меньше или равной нулю")
		fmt.Println("==================================================")
		fmt.Println()
		return
	}

	err = service.WithdrawAccount(accountID, amount)
	if err != nil {
		fmt.Printf("Ошибка при снятии денег со счета: %s\n", err.Error())
		fmt.Println("==================================================")
		fmt.Println()
		return
	}

	fmt.Printf("С указанного счета было успешно снято %.2f сомони\n", amount)
	fmt.Println("==================================================")
	fmt.Println()
}
