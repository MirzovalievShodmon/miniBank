package controller

import (
	"fmt"

	"github.com/MirzovalievShodmon/miniBank.git/internal/service"
)

func GetAllTransactions() {
	transactions, err := service.GetAllTransactions()
	if err != nil {
		fmt.Printf("Ошибка при транзакции: %s\n", err.Error())
		return
	}

	fmt.Println("===== Список транзакций =====")

	if len(transactions) == 0 {
		fmt.Println("Список транзакций - пуст!")
		fmt.Println("=============================")
		return
	}

	for _, t := range transactions {
		fmt.Printf("ID: %d, Сумма: %2.f, Счет: %d, Тип: %s,Время создания: %v\n",
			t.ID, t.Amount, t.AccountID, t.Type, t.CreatedAt)
	}
	fmt.Println("=============================")
}
