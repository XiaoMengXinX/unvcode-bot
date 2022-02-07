package api

import (
	"encoding/json"
	"fmt"
	"github.com/XiaoMengXinX/go-font"
	"github.com/XiaoMengXinX/go-unvcode"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"net/http"
)

type response struct {
	Msg    string `json:"text"`
	ChatID int64  `json:"chat_id"`
	Method string `json:"method"`
}

var unv *unvcode.Unv

func init() {
	unv, _ = unvcode.New(font.Font)
}

func UnvBot(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	var update tgbotapi.Update

	err := json.Unmarshal(body, &update)
	if err != nil {
		log.Println(err)
		return
	}

	if update.Message.Text != "" {
		text := fmt.Sprintln(unv.Unvcode(update.Message.Text))

		data := response{
			Msg:    text,
			Method: "sendMessage",
			ChatID: update.Message.Chat.ID,
		}
		msg, _ := json.Marshal(data)

		log.Printf("Response %s", string(msg))

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, string(msg))
	}
}
