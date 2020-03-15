package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/zserge/lorca"
)

const (
	minimumWaitTime = 2000
	maximumWaitTime = 5000
)

type config struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func main() {
	dat, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Cannot load config file: %s", err)
	}
	cfg := &config{}
	err = json.Unmarshal(dat, cfg)
	if err != nil {
		log.Fatalf("Cannot unmarshal config file content into config struct: %s", err)
	}
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

		time.Sleep(time.Duration(rand.Intn(maximumWaitTime-minimumWaitTime)+minimumWaitTime) * time.Millisecond)

		ui.Eval(fmt.Sprintf(`
			document.querySelector("#log_email").value = "%s";
			document.querySelector("#log_password").value = "%s";
			document.querySelector('input[type="submit"]').click();
		`, cfg.Login, cfg.Password))

		for !waitForElement(ui, `'#session_button'`) {
		}

		time.Sleep(time.Duration(rand.Intn(maximumWaitTime-minimumWaitTime)+minimumWaitTime) * time.Millisecond)

		ui.Eval(`
			document.querySelector('#session_button').click();
		`)

		for !waitForElement(ui, `'#start_session_button'`) && !waitForElement(ui, `'#continue_session_button'`) {
		}

		time.Sleep(time.Duration(rand.Intn(maximumWaitTime-minimumWaitTime)+minimumWaitTime) * time.Millisecond)

		ui.Eval(getMainScript())
	}()

	<-ui.Done()
}

func waitForElement(ui lorca.UI, querySelector string) bool {
	return ui.Eval(fmt.Sprintf("!!document.querySelector(%s)", querySelector)).Bool()
}

func getMainScript() string {
	return `
	const wait = (amount = 0) =>
  new Promise(resolve => setTimeout(resolve, amount));

const randomBreak = (min, max) => wait(Math.floor(Math.random() * max) + min);

const writeAnswerToInput = async response => {
  const answerInput = document.querySelector('#answer');
  const checkButton = document.querySelector('#check');
  const knowNewButton = document.querySelector('#know_new');
  const dontKnowNewButton = document.querySelector('#dont_know_new');
  const skipButton = document.querySelector('#skip');
  if (
    knowNewButton &&
    dontKnowNewButton &&
    document.querySelector('#new_word_form').style.display !== 'none'
  ) {
    await randomBreak(300, 600);
    if (Math.random() >= 0.5) {
      knowNewButton.click();
    } else {
      dontKnowNewButton.click();
    }
    await randomBreak(300, 600);
    skipButton.click();
    return;
  }
  answerInput.focus();

  let offset = 0;
  await randomBreak(400, 750);
  while (offset < response.word.length) {
    answerInput.value += response.word[offset];
    await randomBreak(100, 150);
    offset++;
  }
  checkButton.click();
};

const clickNextWord = async () => {
  const nextWordButton = document.querySelector('#nextword');
  nextWordButton.click();
};

const oldXHROpen = window.XMLHttpRequest.prototype.open;
window.XMLHttpRequest.prototype.open = function(method, url) {
  this.addEventListener('load', function() {
    if (url.includes('generate_next_word.php')) {
      const json = JSON.parse(this.responseText);
      if (json.summary) {
        window.closeWindow();
      } else {
        writeAnswerToInput(json);
      }
    } else if (url.includes('save_answer.php')) {
      clickNextWord();
    }
  });
  return oldXHROpen.apply(this, arguments);
};

if (
  document.querySelector('#start_session_button') &&
  document.querySelector('#start_session_page').style.display === 'block'
) {
  document.querySelector('#start_session_button').click();
} else {
  document.querySelector('#continue_session_button').click();
}

		`
}
