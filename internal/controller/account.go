package controller

import (
	"fmt"

	"github.com/MirzovalievShodmon/miniBank.git/internal/service"
)

func GetAllAccounts() {
	accounts, err := service.GetAllAccounts()
	if err != nil {
		fmt.Printf("Ошибка при получении списка счетов: %s\n", err.Error())
		return
	}

	fmt.Println("===== Список счетов =====")
	for _, a := range accounts {
		fmt.Printf("ID: %d, Баланс: %2.f, Владелец: %s\n", a.ID, a.Balance, a.Owner)
	}
	fmt.Println("=========================")
}

func TopUpAccount() {
	fmt.Println("===== Пополнение счета =====")
	fmt.Println("Введите ID счета, который хотите пополнить")
	var accountID int
	_, err := fmt.Scan(&accountID)
	if err != nil {
		fmt.Printf("Неправильная форма ID. Ошибка: %s\n", err.Error())
		return
	}

	if accountID <= 0 {
		fmt.Println("ID не можеть быть меньше или равным нулю")
		return
	}

	fmt.Println("Введите сумму на которую хотите пополнить:")
	var amount float64
	_, err = fmt.Scan(&amount)
	if err != nil {
		fmt.Printf("Неправильный формат суммы. Ошибка: %s\n", err.Error())
		return
	}

	if amount <= 0 {
		fmt.Println("Сумма не можеть быть меньше или равной нулю")
		return
	}

	err = service.TopUpAccount(accountID, amount)
	if err != nil {
		fmt.Printf("Ошибка при пополнении счета. Ошибка: %s\n", err.Error())
		return
	}

	fmt.Printf("Указанный счет успешно пополнен на %.2f сомони\n", amount)
	fmt.Println("===========================")
}

func WithdrawAccount() {
	fmt.Println("===== Снятие со счета =====")
	fmt.Println("Введите ID счета с которого хотите снять деньги:")
	var accountID int
	_, err := fmt.Scan(&accountID)
	if err != nil {
		fmt.Printf("Неправильный формат ID. Ошибка: %s\n", err.Error())
		return
	}

	if accountID <= 0 {
		fmt.Println("ID не можеть быть меньше или равным нулю")
		return
	}

	fmt.Println("Введите сумму которую хотите снять:")
	var amount float64
	_, err = fmt.Scan(&amount)
	if err != nil {
		fmt.Printf("Неправильный формат суммы. Ошибка: %s\n", err.Error())
		return
	}

	if amount <= 0 {
		fmt.Println("Сумма не может быть меньше или равной нулю")
		return
	}

	err = service.WithdrawAccount(accountID, amount)
	if err != nil {
		fmt.Printf("Ошибка при снятии денег со счета. Ошибка: %s\n", err.Error())
		return
	}

	fmt.Printf("С указанного счета было успешно снято %.2f сомони\n", amount)
	fmt.Println("===========================")
}
