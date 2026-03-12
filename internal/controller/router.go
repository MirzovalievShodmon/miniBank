package controller

import "fmt"

func InitRoutes() {
	fmt.Println()
	fmt.Println("----------------------------------")
	fmt.Println("---Добро пожаловать в miniBank!---")
	for {
		fmt.Println("----------------------------------")
		fmt.Println("---Список функционала:------------")
		fmt.Println("0. Выход")
		fmt.Println("1. Пополнение счёта")
		fmt.Println("2. Снятие со счёта")
		fmt.Println("3. Получение списка всех счетов")
		fmt.Println("4. История операций")
		fmt.Println()
		fmt.Println("---Выберите нужный пункт:")
		cmd := readInput()
		switch cmd {
		case "0":
			fmt.Println("До скорой встречи!)")
			return
		case "1":
			TopUpAccount()
		case "2":
			WithdrawAccount()
		case "3":
			GetAllAccounts()
		case "4":
			GetAllTransactions()
		default:
			fmt.Println("Несуществующая команда, попробуйте еще раз ...")
			fmt.Println()
		}
	}
}
