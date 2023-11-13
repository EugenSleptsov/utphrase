package main

import (
	"bufio"
	"fmt"
	botapi "github.com/EugenSleptsov/utphrase/api/telegram"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	token := os.Getenv("UTPHRASE_TOKEN")

	bot, err := botapi.NewBot(token)
	if err != nil {
		log.Panic(err)
	}

	err = bot.SetCommandList(botapi.CommandPhrase, botapi.CommandSlap)
	if err != nil {
		log.Print(err)
	}

	updates := bot.GetUpdateChannel(30)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.IsCommand() {
				command := update.Message.Command()
				switch command {
				case "add":
					commandAdd(bot, update)
				case "phrase", "fraza", "frazochka":
					commandPhrase(bot, update)
				case "slap":
					commandSlap(bot, update)
				}
			} else {
				switch update.Message.Text {
				case "!фраза", "!фразочка", "!фраza", "!fraza", "!frazochka", "!frazo4ka", "/фраза":
					commandPhrase(bot, update)
				}
			}

		}
	}
}

// commandAdd добавляет новую фразочку
func commandAdd(bot *botapi.Bot, update botapi.Update) {
	if len(update.Message.CommandArguments()) == 0 {
		bot.Reply(update.Message.Chat.ID, update.Message.MessageID, "Текст не задан")
		return
	}

	phrase := update.Message.CommandArguments()

	// Define the file path
	filePath := "data/phrases.txt"
	createEmptyFileIfNotExists(filePath)

	// Open the file for appending (create it if it doesn't exist)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close() // Ensure the file is closed when the function returns

	// Write the new phrase to the file, followed by a newline
	_, err = file.WriteString(phrase + "\n")
	if err != nil {
		log.Panic(err)
	}

	// Send a reply to confirm the phrase addition
	replyText := "Добавлена новая фраза: " + phrase
	bot.Reply(update.Message.Chat.ID, update.Message.MessageID, replyText)
}

// commandPhrase отвечает случайной фразочкой
func commandPhrase(bot *botapi.Bot, update botapi.Update) {
	// Read all phrases from the file
	phrases, err := readPhrasesFromFile("data/phrases.txt")
	if err != nil {
		log.Panic(err)
	}

	if len(phrases) == 0 {
		bot.Reply(update.Message.Chat.ID, update.Message.MessageID, "Ни одной фразочки не было добавлено")
		return
	}

	// Select a random phrase from the list
	randPhrase := phrases[rand.Intn(len(phrases))]

	// Send the random phrase as a reply
	bot.Reply(update.Message.Chat.ID, update.Message.MessageID, randPhrase)
}

// commandSlap слапает пользователя, если указан, иначе выбирает случайного
func commandSlap(bot *botapi.Bot, update botapi.Update) {
	user := update.Message.CommandArguments()
	if len(user) == 0 {
		// Select a random user or username from a slice
		randomUsers := []string{"@daishi424", "DenDr", "@OceanDrive80x", "Dutch", "lord-z", "Systez", "@magic_frontier"} // Add your user list here
		randomUser := getRandomUser(randomUsers)
		user = randomUser
	}

	// @todo сделать в файле
	formatList := []string{
		"slaps %s around a bit with a large trout",
		"slaps %s around a bit with a big smelly trout",
		"%s couldn't handle the pressure and took their own life.",
		"%s was slapped into another dimension.",
	}

	// getting random format
	randFormat := formatList[rand.Intn(len(formatList))]

	bot.Say(update.Message.Chat.ID, fmt.Sprintf(randFormat, user), false)
}

// readPhrasesFromFile reads all phrases from a file and returns them as a slice
func readPhrasesFromFile(filePath string) ([]string, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return []string{}, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var phrases []string

	for scanner.Scan() {
		phrases = append(phrases, scanner.Text())
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return phrases, nil
}

func createEmptyFileIfNotExists(filePath string) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			log.Panic(err)
		}
		file.Close()
	}
}

// getRandomUser returns a random user from the given list
func getRandomUser(users []string) string {
	if len(users) == 0 {
		return ""
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(users))
	return users[randomIndex]
}
