package internals

import (
	"fmt"
	"github.com/mamedvedkov/BankBot/internals/repository"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

/*
"/status" +
"/aboutme" +
"/cards" +
"/debts" +
"/rules" +
"/search" +
"/aboutMyPayment" +
*/

func Process(repo *repository.Repo, cmd string, isMain bool, update tgbotapi.Update) string {
	var response string

	id := update.Message.From.ID

	// TODO: обработка новых пользователей, а не отфутболивание
	if _, err := GetRowByTgId(repo, id); err != nil {
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
	case "search":
		response = searchForUser(repo, update.Message.Text)
	default:
		response = "Функция не реализована"
	}

	return response
}

func status(repo *repository.Repo) string {
	var out string

	valuesRange := "A7:B11"
	res := repo.GetSummary(valuesRange)

	for _, row := range res {
		out += fmt.Sprintf("`%-10v- %10v`\n",
			row[0], row[1])
	}

	return out
}

func aboutMyPayment(repo *repository.Repo, id int) string {
	row, _ := GetRowByTgId(repo, id)
	valuesRange := fmt.Sprintf("Взносы!D%v:N%v", row, row)
	resVal := parseRow(repo.GetPayments(valuesRange))
	titlesRange := "Взносы!D1:N1"
	resTitle := parseRow(repo.GetPayments(titlesRange))

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

func cardHolders(repo *repository.Repo) string {
	cardHoldersRange := "Держатели!A2:B25"
	res := repo.GetHolders(cardHoldersRange)
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

func debts(repo *repository.Repo) string {
	debtsRange := "Займы!A2:C100"
	res := repo.GetLoans(debtsRange)
	userRange := "Участники!A2:B100"
	users := repo.GetMembers(userRange)

	delay := []string{"Просрочки"}
	thisWeek := []string{"Платежи на этой неделе"}
	otherDeb := []string{"Остальные задолжности"}

	for _, row := range res {
		if row[1] == "0" || row[1] == "" {
			continue
		}
		rawDate := strings.ReplaceAll(fmt.Sprintf("%v", row[2]), ".", "-")
		date, err := time.Parse("02-01-2006", rawDate)
		if err != nil {
			return err.Error()
		}
		if isDelay(date) {
			userRow, _ := getRowByName(repo, fmt.Sprintf("%v", row[0]))
			userId := users[userRow-2][1]
			delay = append(delay, fmt.Sprintf("[%s](tg://user?id=%v): %v\tдо %v", row[0], userId, row[1], row[2]))
			continue
		}
		if isThisWeek(date) {
			userRow, _ := getRowByName(repo, fmt.Sprintf("%v", row[0]))
			userId := users[userRow-2][1]
			thisWeek = append(thisWeek, fmt.Sprintf("[%s](tg://user?id=%v): %v\tдо %v", row[0], userId, row[1], row[2]))
			continue
		}
		otherDeb = append(otherDeb, fmt.Sprintf("%v: %v\tдо %v", row[0], row[1], row[2]))
	}

	var result string

	if len(delay) != 1 {
		for _, val := range delay {
			result += val
			result += "\n"
		}
		result += "\n"
	}

	if len(thisWeek) != 1 {
		for _, val := range thisWeek {
			result += val
			result += "\n"
		}
		result += "\n"
	}

	if len(otherDeb) != 1 {
		for _, val := range otherDeb {
			result += val
			result += "\n"
		}
	}

	return result
}

func isThisWeek(date time.Time) bool {
	_, now := time.Now().ISOWeek()
	_, debtWeek := date.ISOWeek()
	return now == debtWeek
}

func isDelay(date time.Time) bool {
	return (date.Add(time.Hour * 23).Add(time.Minute * 59)).Before(time.Now())
}

func searchForUser(repo *repository.Repo, text string) string {
	splitedText := strings.Split(text, "/search ")

	if len(splitedText) == 1 {
		return "Стоит указать кого вы ищите"
	}

	text = splitedText[1]

	if text == "" {
		return "Стоит указать кого вы ищите"
	}

	rowNum, err := getRowByName(repo, text)
	if err != nil {
		return err.Error()
	}

	return searchByRow(repo, rowNum)
}

func aboutme(repo *repository.Repo, id int) string {
	rowNum, err := GetRowByTgId(repo, id)
	if err != nil {
		return err.Error()
	}

	return searchByRow(repo, rowNum)
}

func searchByRow(repo *repository.Repo, rowNum int) string {
	valuesRange := fmt.Sprintf("Участники!A%v:M%v", rowNum, rowNum)
	resVal := parseRow(repo.GetMembers(valuesRange))
	titlesRange := "Участники!A1:M1"
	resTitle := parseRow(repo.GetMembers(titlesRange))

	out := "`Информация о пользователе`\n\n"

	for i, _ := range resVal {
		if resTitle[i] == "" {
			break
		}
		if i == 1 || i == 2 || resVal[i] == "" {
			continue
		}
		if i == 7 && resVal[i] == "0" {
			continue
		}
		if i == 11 || i == 12 {
			temp := ""
			if resVal[i] == "1" {
				temp = "оплачены"
			} else {
				temp = "неоплачены"
			}
			if i == 11 {
				resTitle[i] = "За 3 месяца"
			}
			if i == 12 {
				resTitle[i] = "За этот месяц"
			}
			out += fmt.Sprintf("`%v:%s`\n",
				resTitle[i], formatOutStringToLenghtAddSpacesLeft(24-len([]rune(resTitle[i])), temp))
			continue
		}

		out += fmt.Sprintf("`%v:%s`\n",
			resTitle[i], formatOutStringToLenghtAddSpacesLeft(24-len([]rune(resTitle[i])), resVal[i]))
	}

	return out
}

func formatOutStringToLenghtAddSpacesLeft(length int, str string) string {
	length -= len([]rune(str))
	for i := 0; i < length; i++ {
		str = " " + str
	}
	return str
}

func getRowByName(repo *repository.Repo, name string) (int, error) {
	res := repo.GetMembers("Участники!A2:A250")
	name = strings.ReplaceAll(name, "ё", "е")
	name = strings.ToLower(name)

	var rowId int
	var counter int
	var tempText string

	r, err := regexp.Compile(fmt.Sprintf(`(%s){1}`, name))
	if err != nil {
		return 0, err
	}

	for idx, row := range res {
		tempText = fmt.Sprintf("%v", row[0])
		tempText = strings.ReplaceAll(tempText, "ё", "е")
		tempText = strings.ToLower(tempText)
		found := r.MatchString(tempText)
		if found {
			rowId = idx + 2
			counter++
		}
	}

	if counter == 0 {
		return 0, fmt.Errorf("поиск не дал результатов, попробуйте другой запрос")
	}

	if counter > 1 {
		return 0, fmt.Errorf("слишком много совпадений, попробуйте другой запрос")
	}

	return rowId, nil
}

func GetRowByTgId(repo *repository.Repo, tgId int) (int, error) {
	res := repo.GetMembers("Участники!A2:B250")

	var rowId int

	for idx, row := range res {
		if len(row) == 1 {
			return 0, fmt.Errorf("ID нет в табличке")
		}
		val, _ := strconv.Atoi(fmt.Sprintf("%v", row[1]))
		if val == tgId {
			//log.Printf("Idx=%v row=%s row[0]=%s row[1]=%s\n", idx, row, row[0], row[1])
			rowId = idx
			break
		}
	}

	//log.Println(res[rowId][1])

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
