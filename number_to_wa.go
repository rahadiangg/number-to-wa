package numbertowa

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Chat struct {
	Id int `json:"id"`
}

func parse(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("can't decoded incomening update %s", err.Error())
		return nil, err
	}

	return &update, nil
}

func init() {
	functions.HTTP("NumberToWa", NumberToWa)
}

func FindNumberInString(text string) []string {
	regex, err := regexp.Compile(`[0-9]+`)
	if err != nil {
		log.Printf("Error when find number with regex: %s", err.Error())
	}

	result := regex.FindAllString(text, -1)
	return result
}

func NumberToPhone(numbers []string) ([]string, error) {
	var result []string
	for _, value := range numbers {

		// just want to create as function hehe....
		var tenNumber = func(v string) bool {
			if len(v) >= 10 {
				return true
			}
			return false
		}(value)

		if tenNumber {
			if strings.HasPrefix(value, "0") {
				result = append(result, strings.Replace(value, "0", "62", 1))
			} else if strings.HasPrefix(value, "62") {
				result = append(result, value)
			} else {
				result = append(result, "62"+value)
			}
		}
	}

	if len(result) > 0 {
		return result, nil
	}
	return nil, errors.New("no result when parsing number to phone")
}

func PhoneToUri(phones []string) ([]string, error) {
	var result []string
	for index := range phones {
		result = append(result, "wa.me/"+phones[index])
	}

	if len(result) > 0 {
		return result, nil
	}
	return nil, errors.New("no result when parsing phone to URI")

}

func SendTextToTelegram(chatId int, text string) (string, error) {
	log.Printf("Sending %s to chat_id: %d", text, chatId)

	telegramApi := "https://api.telegram.org/bot" + os.Getenv("TOKEN") + "/sendMessage"

	var result_text string

	numbers := FindNumberInString(text)
	log.Println("Number parsed: ", numbers)

	phones, err := NumberToPhone(numbers)
	if err != nil {
		result_text = err.Error()
	}
	log.Println("Phone parsed: ", phones)

	uris, err := PhoneToUri(phones)
	if err != nil {
		result_text = err.Error()
	}

	for index := range uris {
		result_text += fmt.Sprintf("%s\n", uris[index])
	}

	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {result_text},
		},
	)

	if err != nil {
		log.Printf("Error when send text to chat: %s", err.Error())
		return "", err
	}

	defer response.Body.Close()

	bodyBytes, errRead := ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("Error when parsing telegram answer %s", errRead.Error())
		return "", err
	}

	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram response: %s", bodyString)
	return bodyString, nil
}

func NumberToWa(w http.ResponseWriter, r *http.Request) {
	update, err := parse(r)
	if err != nil {
		log.Printf("error parsing update: %s", err.Error())
	}

	telegramResponseBody, err := SendTextToTelegram(update.Message.Chat.Id, update.Message.Text)
	if err != nil {
		log.Printf("Got error %s from telegram, response body us %s", err.Error(), telegramResponseBody)
	} else {
		log.Printf("Successfully send to chat id %d", update.Message.Chat.Id)
	}
}
