package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/spf13/viper"
	"github.com/zserge/lorca"
)

func waitForElement(ui lorca.UI, querySelector string) bool {
	return ui.Eval(fmt.Sprintf("!!document.querySelector(%s)", querySelector)).Bool()
}

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

func main() {
	// Create UI with basic HTML passed via data URI
	ui, err := lorca.New("https://instaling.pl/teacher.php?page=login", "", 800, 600)
	if err != nil {
		log.Fatal(err)
	}
	ui.Bind("closeWindow", func() {
		ui.Close()
	})
	defer ui.Close()

	go func() {
		for !waitForElement(ui, `'form[action="./teacher.php?page=teacherActions"'`) {
		}

		time.Sleep(time.Duration(rand.Intn(10000-2000)+2000) * time.Millisecond)

		ui.Eval(fmt.Sprintf(`
			document.querySelector("#log_email").value = "%s";
			document.querySelector("#log_password").value = "%s";
			document.querySelector('input[type="submit"]').click();
		`, viper.GetString("login"), viper.GetString("password")))

		for !waitForElement(ui, `'#session_button'`) {
		}

		time.Sleep(time.Duration(rand.Intn(10000-2000)+2000) * time.Millisecond)

		ui.Eval(`
			document.querySelector('#session_button').click();
		`)

		for !waitForElement(ui, `'#start_session_button'`) && !waitForElement(ui, `'#continue_session_button'`) {
		}

		time.Sleep(time.Duration(rand.Intn(10000-2000)+2000) * time.Millisecond)

		ui.Eval(`
			console.log("called")
			const wait = (amount = 0) =>
			new Promise(resolve => setTimeout(resolve, amount));

			const startWriting = async response => {
			const answerInput = document.querySelector("#answer");
			const checkButton = document.querySelector("#check");
			answerInput.focus();
			let offset = 0;

			const time = Math.floor(Math.random() * 500) + 100;
			await wait(time);

			while (offset < response.word.length) {
				answerInput.value += response.word[offset];
				const time = Math.floor(Math.random() * 150) + 100;
				await wait(time);
				offset++;
			}

			checkButton.click();
			};

			const nextWord = async () => {
			const nextWordButton = document.querySelector("#nextword");
			nextWordButton.click();
			};

			const oldXHROpen = window.XMLHttpRequest.prototype.open;
			window.XMLHttpRequest.prototype.open = function(
			method,
			url,
			async,
			user,
			password
			) {
			this.addEventListener("load", function() {
				if (url.includes("generate_next_word.php")) {
						const json = JSON.parse(this.responseText);
						if(json.summary) {
							window.closeWindow();
						} else {
							startWriting(json);
						}
				} else if (url.includes("save_answer.php")) {
				nextWord();
				}
			});

			return oldXHROpen.apply(this, arguments);
			};

			if(document.querySelector('#start_session_button') && document.querySelector("#start_session_page").style.display === "block") {
				document.querySelector('#start_session_button').click();
			} else {
				document.querySelector('#continue_session_button').click();
			}
		`)
	}()

	<-ui.Done()
}
