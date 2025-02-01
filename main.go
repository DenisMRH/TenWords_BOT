// Основной пакет программы
package main

// Импорт необходимых библиотек
import (
	"fmt"
	"log" // Для логирования ошибок
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" // Официальная обёртка Telegram API
	"github.com/joho/godotenv"
)

// Главная функция - точка входа в программу
func main() {
	// Инициализация бота с использованием API токена

	TELEGRAM_BOT_TOKEN := importEnv("token.env", "TELEGRAM_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(TELEGRAM_BOT_TOKEN)
	if err != nil {
		// Если ошибка при создании бота - выводим её и завершаем работу
		log.Panic("Ошибка инициализации бота:", err)
	}

	// Включаем режим отладки (вывод в консоль всех запросов и ответов)
	bot.Debug = true
	// Выводим информацию об успешной авторизации
	log.Printf("Авторизация успешна! Бот %s готов к работе", bot.Self.UserName)

	// Настройка параметров получения обновлений:
	// - Offset = 0 означает получение всех непрочитанных сообщений
	// - Timeout = 60 сек - время ожидания новых обновлений
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Получаем канал обновлений через long-polling
	updates := bot.GetUpdatesChan(updateConfig)

	// Бесконечный цикл обработки входящих обновлений
	for update := range updates {
		// Если обновление не содержит сообщение - пропускаем его
		if update.Message == nil {
			continue
		}

		// Создаем новое сообщение для ответа:
		// - Указываем ID чата, куда отправлять ответ
		// - Используем текст из полученного сообщения
		reply := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			update.Message.Text,
		)

		// Настраиваем цитирование исходного сообщения
		reply.ReplyToMessageID = update.Message.MessageID

		// Отправляем подготовленное сообщение
		if _, err := bot.Send(reply); err != nil {
			// В случае ошибки отправки - логируем и завершаем работу
			log.Panic("Ошибка отправки сообщения:", err)
		}
	}
}

func importEnv(fileName, varName string) (variable string) {
	err := godotenv.Load(fileName)
	if err != nil {
		log.Fatalf("Ошибка импорта файла %v", err)
	}

	variable = os.Getenv(varName)
	if variable == "" {
		log.Fatalf("Переменная %v не найдена.")
	}
	fmt.Println("Переменная ", varName, " из файла ", fileName, " импортирована!")
	return
}
