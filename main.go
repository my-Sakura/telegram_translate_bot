package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type result struct {
	From         string `json: "from"`
	To           string `json: "to"`
	Trans_result []test `json: "trans_result"`
}

type test struct {
	Src string `json: "src"`
	Dst string `json: "dst"`
}

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

var bot *tgbotapi.BotAPI

func init() {
	//new telegramBot
	token := os.Getenv("token")

	bot, _ = tgbotapi.NewBotAPI(token)
	bot.Debug = true
	//new update
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// link := "https://translate-one.vercel.app/"
	link := "https://e1b5b78dfd47.ngrok.io/"
	bot.SetWebhook(tgbotapi.NewWebhook(link + token))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var update tgbotapi.Update

	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &update); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	var target result
	//baidu appid
	appid := os.Getenv("appid")
	q := update.Message.Text
	salt := 123
	secret := os.Getenv("secret")
	data := []byte(appid + q + strconv.Itoa(salt) + secret)
	sign := fmt.Sprintf("%x", md5.Sum(data))
	url := fmt.Sprintf("https://fanyi-api.baidu.com/api/trans/vip/translate?q=%s&from=%s&to=%s&salt=%d&appid=%s&sign=%s", q, "auto", "en", salt, appid, sign)
	result, err := httpGet(url)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(result, &target)
	if err != nil {
		fmt.Println(err)
	}
	msg.Text = target.Trans_result[0].Dst
	if _, err = bot.Send(msg); err != nil {
		fmt.Println(err)
	}
}

func main() {
	http.HandleFunc("/1306747225:AAEpFms8-OJCb2VcuMt04AW5xR6SUVz9wCM", Handler)
	http.ListenAndServe(":1000", nil)
}
