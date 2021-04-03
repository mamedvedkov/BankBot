package internals

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

/*
"/status" +
"/aboutme"
"/cards" +
"/debts" +
"/rules" +
"/search"
"/aboutMyPayment" +
*/

func Process(repo *Repo, cmd string, isMain bool, update tgbotapi.Update) string {
	var response string

	id := update.Message.From.ID

	if _, err := getRowByTgId(repo, id); err != nil {
		return err.Error()
	}

	if isMain {
		return "Постарайтесь не использовать это в общем чате"
	}

	switch cmd {
	case "status":
		response = status(repo)
	case "aboutme":
		response = aboutme(repo, id)
	case "aboutMyPayment":
		response = aboutMyPayment(repo, id)
	case "cards":
		response = cardHolders(repo)
	case "rules":
		response = showRules()
	case "debts":
		response = debts(repo)

	default:
		response = "Функция не реализована"
	}

	return response
}

func status(repo *Repo) string {
	var out string

	for i := 7; i <= 10; i++ {
		valuesRange := fmt.Sprintf("A%v:B%v", i, i)
		resVal := parseRow(repo.GetValues(valuesRange))
		out += fmt.Sprintf("%v - %v\n", resVal[0], resVal[1])
	}

	return out
}

func aboutme(adapter *Repo, id int) string {
	return ""
}

func aboutMyPayment(repo *Repo, id int) string {
	row, _ := getRowByTgId(repo, id)
	valuesRange := fmt.Sprintf("Взносы!D%v:N%v", row, row)
	resVal := parseRow(repo.GetValues(valuesRange))
	titlesRange := "Взносы!D1:N1"
	resTitle := parseRow(repo.GetValues(titlesRange))

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

func cardHolders(repo *Repo) string {
	cardHoldersRange := "Держатели!A2:B25"
	res := repo.GetValues(cardHoldersRange)
	number := res[0][0]
	sum := res[0][1]
	city := res[2][0]
	bankName := res[3][0]
	systemName := res[4][0]
	link := res[5][0]

	result := fmt.Sprintf("номер карты:\t%s\n", number)
	result += fmt.Sprintf("сумма на карте:\t%s\n", sum)
	result += fmt.Sprintf("город:\t%s\n", city)
	result += fmt.Sprintf("карта:\t%s %s", bankName, systemName)

	if link != "" {
		result += fmt.Sprintf("\nПополнить без комиссии можно по ссылке:\n%s\n", link)
	}

	number = res[17][0]
	sum = res[17][1]
	city = res[19][0]
	bankName = res[20][0]
	systemName = res[21][0]

	if number != "" {
		result += fmt.Sprintf("\nномер карты:\t%s\n", number)
		result += fmt.Sprintf("сумма на карте:\t%s\n", sum)
		result += fmt.Sprintf("город:\t%s\n", city)
		result += fmt.Sprintf("карта:\t%s %s", bankName, systemName)
	}

	return result
}

func showRules() string {
	rulesLink := "заглушка для ссылки на правила"
	return fmt.Sprintf("Правила работы кассы\n%s", rulesLink)
}

func debts(repo *Repo) string {
	debtsRange := "Займы!A2:C100"
	res := repo.GetValues(debtsRange)
	result := ""

	for _, row := range res {
		if row[1] == "0" || row[1] == "" {
			continue
		}
		result += fmt.Sprintf("%v: %v\tдо %v\n", row[0], row[1], row[2])
	}

	return result
}

func getRowByTgId(repo *Repo, tgId int) (int, error) {
	res := repo.GetValues("Участники!A2:B250")

	var rowId int

	for idx, row := range res {
		if len(row) == 1 {
			return 0, fmt.Errorf("ID нет в табличке")
		}
		val, _ := strconv.Atoi(fmt.Sprintf("%v", row[1]))
		if val == tgId {
			log.Printf("Idx=%v row=%s row[0]=%s row[1]=%s\n", idx, row, row[0], row[1])
			rowId = idx
			break
		}
	}

	log.Println(res[rowId][1])

	return rowId + 2, nil
}

//Парсим результат если нужен столбец
func parseCollumn(values [][]interface{}) []string {
	r := []rune(fmt.Sprintf("%v", values))
	r = r[2 : len(r)-2]
	var out string = string(r)
	out = strings.ReplaceAll(out, "] [", ";")

	return strings.Split(out, ";")
}

//Парсим результат если нужна строка
func parseRow(values [][]interface{}) []string {
	var out string
	for _, row := range values[0] {
		out += fmt.Sprintf("%v;", row)
	}
	return strings.Split(out, ";")

}
