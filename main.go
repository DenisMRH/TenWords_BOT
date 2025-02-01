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

	
	err := godotenv.Load("C:\Users\denis\Documents\go\TelegramBOT\token.env")
		if err != nil {
		log.Fatalf("Ошибка загрузки token.env файла %v", err)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("Токен не найден в token.env")
	}
	fmt.Println("Токен успешно загружен")

	bot, err := tgbotapi.NewBotAPI(token)
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
