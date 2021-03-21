package internals

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

/*
"/status"
"/aboutme"
"/cards"
"/debts"
"/rules"
"/search"
"/aboutMyPayment"
*/

func Process(adapter *Adapter, cmd string, isMain bool, update tgbotapi.Update) string {
	var response string

	id := update.Message.From.ID

	if _, err := getRowByTgId(adapter, id); err != nil {
		return err.Error()
	}

	if isMain {
		return "Постарайтесь не использовать это в общем чате"
	}

	switch cmd {
	case "status":
		response = status(adapter)
	case "aboutme":
		response = aboutme(adapter, id)
	case "aboutMyPayment":
		response = aboutMyPayment(adapter, id)
	default:
		response = "Функция не реализована"
	}

	return response
}

func status(adapter *Adapter) string {
	var out string

	for i := 7; i <= 10; i++ {
		valuesRange := fmt.Sprintf("A%v:B%v", i, i)
		resVal := parseRow(adapter.GetValues(valuesRange))
		out += fmt.Sprintf("%v - %v\n", resVal[0], resVal[1])
	}

	return out
}

func aboutme(adapter *Adapter, id int) string {
	return ""
}

func aboutMyPayment(adapter *Adapter, id int) string {
	row, _ := getRowByTgId(adapter, id)
	valuesRange := fmt.Sprintf("Взносы!D%v:N%v", row, row)
	resVal := parseRow(adapter.GetValues(valuesRange))
	titlesRange := "Взносы!D1:N1"
	resTitle := parseRow(adapter.GetValues(titlesRange))

	out := "Статистика по платежам\n"

	for i, _ := range resVal {
		if resTitle[i] == "" {
			break
		}
		if resVal[i] == "!" {
			continue
		}
		out += fmt.Sprintf("%v - %v рублей\n", resTitle[i], resVal[i])
	}

	return out
}

func getRowByTgId(adapter *Adapter, tgId int) (int, error) {
	res := adapter.GetValues("Участники!A2:B250")

	var rowId int

	for idx, row := range res {
		if len(row) == 1 {
			return 0, fmt.Errorf("ID нет в табличке")
		}
		val, _ := strconv.Atoi(fmt.Sprintf("%v", row[1]))
		if val == tgId {
			log.Printf("Idx=%v row=%s row[0]=%s row[1]=%s", idx, row, row[0], row[1])
			rowId = idx
			break
		}
	}

	log.Println(res[rowId][1])

	return rowId + 2, nil
}

/*
//Парсим результат если нужен столбец
func parseCollumn(values [][]interface{}) []string {
	r := []rune(fmt.Sprintf("%v", values))
	r = r[2 : len(r)-2]
	var out string = string(r)
	out = strings.ReplaceAll(out, "] [", ";")

	return strings.Split(out, ";")
}
*/

//Парсим результат если нужна строка
func parseRow(values [][]interface{}) []string {
	var out string
	for _, row := range values[0] {
		out += fmt.Sprintf("%v;", row)
	}
	return strings.Split(out, ";")

}
