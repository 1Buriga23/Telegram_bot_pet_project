package telegram

import (
	"context"
	"fmt"
	"gracia-bot/internal/models"
	"gracia-bot/internal/repository"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleUpdate(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		if update.CallbackQuery.Message != nil {
			b.HandleCallback(update.CallbackQuery)
		}
		return
	}

	if update.Message != nil {
		b.HandleMessage(update.Message)
		return
	}
}

func (b *Bot) ShowAdminPanel(chatID int64) {

	message := tgbotapi.NewMessage(
		chatID,
		"⚙️ Админ-панель",
	)

	message.ReplyMarkup = AdminPanelMenu()

	_, err := b.API.Send(message)

	if err != nil {
		log.Println(err)
	}
}

func (b *Bot) EditMenu(
	chatID int64,
	messageID int,
	text string,
	keyboard tgbotapi.InlineKeyboardMarkup,
) {
	edit := tgbotapi.NewEditMessageText(
		chatID,
		messageID,
		text,
	)

	edit.ReplyMarkup = &keyboard

	_, err := b.API.Send(edit)
	if err != nil {
		log.Println(err)
	}
}

func (b *Bot) ShowMainMenu(
	chatID int64,
	text string,
) {

	categories, err := repository.GetCategories(
		context.Background(),
		b.DB,
	)

	if err != nil {
		log.Println(err)
		return
	}

	message := tgbotapi.NewMessage(
		chatID,
		text,
	)

	message.ReplyMarkup = MainMenu(categories)

	msg, err := b.API.Send(message)

	if err != nil {
		log.Println(err)
		return
	}

	session, exists := Sessions[chatID]

	if !exists {
		session = &UserSession{}
	}

	session.MenuMessageID = msg.MessageID

	Sessions[chatID] = session
}

func (b *Bot) ShowSearchResults(
	chatID int64,
	items []models.Item,
) {

	if len(items) == 0 {

		message := tgbotapi.NewMessage(
			chatID,
			"🔎 Ничего не найдено",
		)

		message.ReplyMarkup = BackToMenu()

		_, err := b.API.Send(message)

		if err != nil {
			log.Println(err)
		}

		return
	}

	message := tgbotapi.NewMessage(
		chatID,
		"🔎 Результаты поиска:",
	)

	message.ReplyMarkup = ItemsMenu(items)

	_, err := b.API.Send(message)

	if err != nil {
		log.Println(err)
	}
}

func (b *Bot) HandleCallback(callback *tgbotapi.CallbackQuery) {
	data := callback.Data
	if strings.HasPrefix(data, "admin:") {
		isAdmin, err := repository.IsAdmin(
			context.Background(),
			b.DB,
			callback.From.ID,
		)

		if err != nil {
			log.Println(err)
			return
		}

		if !isAdmin {
			_, err := b.API.Request(
				tgbotapi.NewCallback(
					callback.ID,
					"Нет доступа",
				),
			)

			if err != nil {
				log.Println(err)
			}

			return
		}
	}
	switch {

	case data == "admin:add_category":
		chatID := callback.Message.Chat.ID
		messageID := callback.Message.MessageID

		// Получаем или создаём сессию пользователя.
		session, exists := Sessions[chatID]

		if !exists {
			session = &UserSession{}
		}

		// Включаем режим добавления категории.
		session.AdminMode = AdminModeAddCategory
		session.SearchMode = false
		session.FromSearch = false

		Sessions[chatID] = session

		// Убираем индикатор загрузки с кнопки.
		_, err := b.API.Request(
			tgbotapi.NewCallback(callback.ID, ""),
		)

		if err != nil {
			log.Println(err)
		}

		// Редактируем существующее сообщение админ-панели.

		b.EditMenu(
			chatID,
			messageID,
			"Введите название категории:",
			tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: make(
					[][]tgbotapi.InlineKeyboardButton,
					0,
				),
			},
		)

	case data == "admin:add_video":
		chatID := callback.Message.Chat.ID
		messageID := callback.Message.MessageID

		// Получаем категории из базы данных.
		categories, err := repository.GetCategories(
			context.Background(),
			b.DB,
		)

		if err != nil {
			log.Println(err)
			return
		}

		// Если категорий нет, показываем админ-панель.
		if len(categories) == 0 {
			_, err := b.API.Request(
				tgbotapi.NewCallback(callback.ID, ""),
			)
			if err != nil {
				log.Println(err)
			}

			b.EditMenu(
				chatID,
				messageID,
				"Сначала создайте категорию.",
				AdminPanelMenu(),
			)
			return
		}

		// Получаем текущую сессию.
		session, exists := Sessions[chatID]
		if !exists {
			session = &UserSession{}
		}

		// Переходим к выбору категории для видео.
		session.AdminMode = AdminModeSelectVideoCategory
		session.MenuMessageID = messageID
		session.SelectedItemID = 0

		Sessions[chatID] = session

		// Убираем индикатор загрузки с кнопки.
		_, err = b.API.Request(
			tgbotapi.NewCallback(callback.ID, ""),
		)
		if err != nil {
			log.Println(err)
		}

		// Редактируем сообщение админки.
		b.EditMenu(
			chatID,
			messageID,
			"🎥 Выберите категорию, в которой находится нужный элемент:",
			AdminVideoCategoriesMenu(categories),
		)

	case strings.HasPrefix(data, "admin:video_category:"):
		chatID := callback.Message.Chat.ID
		messageID := callback.Message.MessageID

		// Извлекаем ID выбранной категории.
		categoryIDText := strings.TrimPrefix(
			data,
			"admin:video_category:",
		)

		categoryID, err := strconv.ParseInt(
			categoryIDText,
			10,
			64,
		)
		if err != nil {
			log.Println(err)
			return
		}

		// Проверяем текущий сценарий.
		session, exists := Sessions[chatID]

		if !exists || session.AdminMode != AdminModeSelectVideoCategory {
			_, err := b.API.Request(
				tgbotapi.NewCallback(
					callback.ID,
					"Сценарий завершён",
				),
			)
			if err != nil {
				log.Println(err)
			}
			return
		}

		// Получаем элементы выбранной категории.
		items, err := repository.GetItemsByCategoryID(
			context.Background(),
			b.DB,
			categoryID,
		)
		if err != nil {
			log.Println(err)
			return
		}

		// Запоминаем выбранную категорию.
		session.SelectedCategoryID = categoryID
		session.AdminMode = AdminModeSelectVideoItem
		session.MenuMessageID = messageID

		Sessions[chatID] = session

		// Подтверждаем нажатие кнопки.
		_, err = b.API.Request(
			tgbotapi.NewCallback(callback.ID, ""),
		)
		if err != nil {
			log.Println(err)
		}

		// Показываем элементы этой категории.
		text := "Выберите элемент, к которому нужно добавить видео:"

		if len(items) == 0 {
			text = "В этой категории пока нет элементов. Сначала добавьте элемент."
		}

		b.EditMenu(
			chatID,
			messageID,
			text,
			AdminVideoItemsMenu(items),
		)

	case strings.HasPrefix(data, "admin:video_item:"):
		chatID := callback.Message.Chat.ID
		messageID := callback.Message.MessageID

		// Получаем ID выбранного элемента.
		itemIDText := strings.TrimPrefix(
			data,
			"admin:video_item:",
		)

		itemID, err := strconv.ParseInt(itemIDText, 10, 64)
		if err != nil {
			log.Println(err)
			return
		}

		// Проверяем, что пользователь выбирает элемент
		// в сценарии добавления видео.
		session, exists := Sessions[chatID]

		if !exists || session.AdminMode != AdminModeSelectVideoItem {
			_, err := b.API.Request(
				tgbotapi.NewCallback(
					callback.ID,
					"Сначала выберите категорию",
				),
			)
			if err != nil {
				log.Println(err)
			}
			return
		}

		// Получаем элемент из базы.
		item, err := repository.GetItemByID(
			context.Background(),
			b.DB,
			itemID,
		)
		if err != nil {
			log.Println(err)
			return
		}

		// Проверяем, что элемент относится к выбранной категории.
		if item.CategoryID != session.SelectedCategoryID {
			_, err := b.API.Request(
				tgbotapi.NewCallback(
					callback.ID,
					"Элемент не относится к выбранной категории",
				),
			)
			if err != nil {
				log.Println(err)
			}
			return
		}

		// Запоминаем элемент, к которому добавляем видео.
		session.SelectedItemID = item.ID
		session.AdminMode = AdminModeAddVideo
		session.MenuMessageID = messageID

		Sessions[chatID] = session

		// Подтверждаем нажатие.
		_, err = b.API.Request(
			tgbotapi.NewCallback(callback.ID, ""),
		)
		if err != nil {
			log.Println(err)
		}

		// Просим отправить видео в том же сообщении.
		b.EditMenu(
			chatID,
			messageID,
			fmt.Sprintf(
				"🎥 Элемент: «%s»\n\nОтправьте видео для этого элемента.",
				item.Name,
			),
			tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: make(
					[][]tgbotapi.InlineKeyboardButton,
					0,
				),
			},
		)

	case strings.HasPrefix(data, "category:"):
		b.DeleteVideos(callback.Message.Chat.ID)
		categoryID := strings.TrimPrefix(data, "category:")
		b.HandleCategory(callback, categoryID)

	case strings.HasPrefix(data, "item:"):
		b.DeleteVideos(callback.Message.Chat.ID)
		itemID := strings.TrimPrefix(data, "item:")
		b.HandleItem(callback, itemID)

	case strings.HasPrefix(data, "admin:item_category:"):
		chatID := callback.Message.Chat.ID
		messageID := callback.Message.MessageID

		// Получаем ID выбранной категории из callback.
		categoryIDText := strings.TrimPrefix(
			data,
			"admin:item_category:",
		)

		categoryID, err := strconv.ParseInt(
			categoryIDText,
			10,
			64,
		)
		if err != nil {
			log.Println(err)
			return
		}

		// Получаем текущую сессию.
		session, exists := Sessions[chatID]
		if !exists || session.AdminMode != AdminModeAddItem {
			_, err := b.API.Request(
				tgbotapi.NewCallback(
					callback.ID,
					"Сценарий завершён",
				),
			)
			if err != nil {
				log.Println(err)
			}
			return
		}

		// Запоминаем выбранную категорию.
		session.SelectedCategoryID = categoryID
		session.MenuMessageID = messageID

		Sessions[chatID] = session

		// Убираем индикатор загрузки кнопки.
		_, err = b.API.Request(
			tgbotapi.NewCallback(callback.ID, ""),
		)
		if err != nil {
			log.Println(err)
		}

		// Редактируем существующее сообщение.
		b.EditMenu(
			chatID,
			messageID,
			"Введите название элемента:",
			tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: make(
					[][]tgbotapi.InlineKeyboardButton,
					0,
				),
			},
		)

	case data == "menu:main":
		chatID := callback.Message.Chat.ID
		messageID := callback.Message.MessageID
		b.DeleteVideos(callback.Message.Chat.ID)

		session, exists := Sessions[chatID]

		if !exists {
			session = &UserSession{}
		}

		session.SearchMode = false
		session.FromSearch = false

		Sessions[chatID] = session

		categories, _ := repository.GetCategories(
			context.Background(),
			b.DB,
		)

		b.EditMenu(
			chatID,
			messageID,
			"🤸 Gracia Bot 🤸\n\n\nВыбери каталог:",
			MainMenu(categories),
		)

	case data == "admin:add_more_video":
		chatID := callback.Message.Chat.ID
		messageID := callback.Message.MessageID

		session, exists := Sessions[chatID]
		if !exists || session.SelectedItemID == 0 {
			_, err := b.API.Request(
				tgbotapi.NewCallback(callback.ID, "Сначала создайте элемент"),
			)
			if err != nil {
				log.Println(err)
			}
			return
		}

		session.AdminMode = AdminModeAddVideo
		session.MenuMessageID = messageID
		Sessions[chatID] = session

		_, err := b.API.Request(
			tgbotapi.NewCallback(callback.ID, ""),
		)
		if err != nil {
			log.Println(err)
		}

		b.EditMenu(
			chatID,
			messageID,
			"🎥 Отправьте следующее видео:",
			AddVideoActionsMenu(),
		)

	case data == "admin:finish_videos":
		chatID := callback.Message.Chat.ID
		messageID := callback.Message.MessageID

		session, exists := Sessions[chatID]
		if exists {
			session.AdminMode = AdminModeNone
			session.SelectedItemID = 0
			session.SelectedCategoryID = 0
			Sessions[chatID] = session
		}

		_, err := b.API.Request(
			tgbotapi.NewCallback(callback.ID, ""),
		)
		if err != nil {
			log.Println(err)
		}

		b.EditMenu(
			chatID,
			messageID,
			"✅ Добавление видео завершено!",
			AdminPanelMenu(),
		)

	case data == "admin:add_item":
		chatID := callback.Message.Chat.ID
		messageID := callback.Message.MessageID

		// Получаем категории из PostgreSQL.
		categories, err := repository.GetCategories(
			context.Background(),
			b.DB,
		)

		if err != nil {
			log.Println(err)
			return
		}

		// Убираем индикатор ожидания на кнопке.
		_, err = b.API.Request(
			tgbotapi.NewCallback(callback.ID, ""),
		)
		if err != nil {
			log.Println(err)
		}

		// Если категорий нет, сначала нужно создать категорию.
		if len(categories) == 0 {
			b.EditMenu(
				chatID,
				messageID,
				"Сначала добавьте категорию.",
				AdminPanelMenu(),
			)
			return
		}

		// Получаем или создаём сессию.
		session, exists := Sessions[chatID]
		if !exists {
			session = &UserSession{}
		}

		session.AdminMode = AdminModeAddItem
		session.MenuMessageID = messageID
		session.SearchMode = false
		session.FromSearch = false

		Sessions[chatID] = session

		// Показываем категории в том же сообщении.
		b.EditMenu(
			chatID,
			messageID,
			"Выберите категорию для нового элемента:",
			AdminSelectCategoryMenu(categories),
		)

	case data == "search":
		chatID := callback.Message.Chat.ID
		messageID := callback.Message.MessageID

		session, exists := Sessions[chatID]

		if !exists {
			session = &UserSession{}
		}

		session.SearchMode = true
		session.FromSearch = true

		Sessions[chatID] = session

		b.EditMenu(
			chatID,
			messageID,
			"🔎 Введите название элемента:",
			BackToMenu(),
		)

	}
}

func (b *Bot) HandleMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		chatID := message.Chat.ID
		b.ClearAll(chatID)
		switch message.Command() {

		case "start":

			textHello := fmt.Sprintf("Привет, %s!\n\n\n🤸 Gracia Bot 🤸\n\n\nВыбери каталог:", message.From.FirstName)

			if message.From.FirstName == "" {
				textHello = "Привет!\n\n\n🤸 Gracia Bot 🤸\n\n\nВыбери каталог:"
			}

			b.ShowMainMenu(
				chatID,
				textHello,
			)
		case "admin":

			isAdmin, err := repository.IsAdmin(
				context.Background(),
				b.DB,
				message.From.ID,
			)

			if err != nil {
				log.Println(err)
				return
			}

			if isAdmin {

				b.ShowAdminPanel(message.Chat.ID)

				return
			}

			session, exists := Sessions[message.Chat.ID]

			if !exists {
				session = &UserSession{}
			}

			session.WaitingAdminPassword = true

			Sessions[message.Chat.ID] = session

			msg := tgbotapi.NewMessage(
				message.Chat.ID,
				"Введите пароль администратора:",
			)

			b.API.Send(msg)

			return

		default:
			return
		}
	}

	if message.Video != nil {
		chatID := message.Chat.ID

		session, exists := Sessions[chatID]

		// Видео принимаем только в режиме его добавления.
		if !exists || session.AdminMode != AdminModeAddVideo {
			return
		}

		b.DeleteChatMessage(
			chatID,
			message.MessageID,
		)

		videoName := strings.TrimSpace(message.Caption)

		if videoName == "" {
			videoName = "Видео"
		}

		err := repository.CreateVideo(
			context.Background(),
			b.DB,
			videoName,
			message.Video.FileID,
			session.SelectedItemID,
		)

		if err != nil {
			log.Println(err)

			b.EditMenu(
				chatID,
				session.MenuMessageID,
				"❌ Не удалось сохранить видео. Отправьте его ещё раз.",
				tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: make(
						[][]tgbotapi.InlineKeyboardButton,
						0,
					),
				},
			)

			return
		}

		messageID := session.MenuMessageID

		// Видео сохранено. Ждём решения администратора.
		session.AdminMode = AdminModeNone

		// Важно: SelectedItemID пока НЕ обнуляем.
		// Он нужен, если администратор захочет добавить ещё видео.
		Sessions[chatID] = session

		b.EditMenu(
			chatID,
			messageID,
			"✅ Видео успешно добавлено!\n\nДобавить ещё одно видео?",
			AddVideoActionsMenu(),
		)

		return

	}

	session, exists := Sessions[message.Chat.ID]

	if exists && session.WaitingAdminPassword {

		if message.Text == GetAdminPassword() {

			err := repository.AddAdmin(
				context.Background(),
				b.DB,
				message.From.ID,
			)

			if err != nil {
				log.Println(err)
				return
			}

			session.WaitingAdminPassword = false
			Sessions[message.Chat.ID] = session

			msg := tgbotapi.NewMessage(
				message.Chat.ID,
				"✅ Админ доступ получен",
			)

			b.API.Send(msg)

			b.ShowAdminPanel(message.Chat.ID)

		} else {

			msg := tgbotapi.NewMessage(
				message.Chat.ID,
				"❌ Неверный пароль",
			)

			b.API.Send(msg)
		}

		return
	}

	if exists && session.AdminMode == AdminModeAddCategory {
		b.DeleteChatMessage(
			message.Chat.ID,
			message.MessageID,
		)
		chatID := message.Chat.ID

		// Повторно проверяем права перед изменением базы.
		isAdmin, err := repository.IsAdmin(
			context.Background(),
			b.DB,
			message.From.ID,
		)

		if err != nil {
			log.Println(err)
			return
		}

		if !isAdmin {
			session.AdminMode = AdminModeNone
			Sessions[chatID] = session

			msg := tgbotapi.NewMessage(
				chatID,
				"❌ У вас нет доступа к этому действию.",
			)
			if _, err := b.API.Send(msg); err != nil {
				log.Println(err)
			}
			return
		}

		name := strings.TrimSpace(message.Text)

		if name == "" {
			msg := tgbotapi.NewMessage(
				chatID,
				"Название не может быть пустым. Введите название категории:",
			)
			if _, err := b.API.Send(msg); err != nil {
				log.Println(err)
			}
			return
		}

		err = repository.CreateCategory(
			context.Background(),
			b.DB,
			name,
		)

		if err != nil {
			log.Println(err)

			msg := tgbotapi.NewMessage(
				chatID,
				"❌ Не удалось добавить категорию. Попробуйте ещё раз.",
			)
			if _, err := b.API.Send(msg); err != nil {
				log.Println(err)
			}
			return
		}

		session.AdminMode = AdminModeNone
		Sessions[chatID] = session

		msg := tgbotapi.NewMessage(
			chatID,
			fmt.Sprintf("✅ Категория «%s» добавлена!", name),
		)
		if _, err := b.API.Send(msg); err != nil {
			log.Println(err)
		}

		return
	}

	if exists && session.AdminMode == AdminModeAddItem {
		b.DeleteChatMessage(
			message.Chat.ID,
			message.MessageID,
		)
		chatID := message.Chat.ID

		// Получаем название элемента.
		name := strings.TrimSpace(message.Text)

		if name == "" {
			b.EditMenu(
				chatID,
				session.MenuMessageID,
				"Введите название элемента:",
				tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: make(
						[][]tgbotapi.InlineKeyboardButton,
						0,
					),
				},
			)
			return
		}

		// Сохраняем элемент в PostgreSQL.
		itemID, err := repository.CreateItem(
			context.Background(),
			b.DB,
			name,
			session.SelectedCategoryID,
		)

		if err != nil {
			log.Println(err)

			b.EditMenu(
				chatID,
				session.MenuMessageID,
				"❌ Не удалось добавить элемент. Попробуйте ещё раз:",
				tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: make(
						[][]tgbotapi.InlineKeyboardButton,
						0,
					),
				},
			)

			return
		}

		messageID := session.MenuMessageID

		// Элемент создан. Теперь ожидаем видео.
		session.SelectedItemID = itemID
		session.AdminMode = AdminModeAddVideo

		Sessions[chatID] = session

		b.EditMenu(
			chatID,
			messageID,
			fmt.Sprintf(
				"✅ Элемент «%s» добавлен!\n\n🎥 Теперь отправьте одно видео для этого элемента.\nПри желании добавьте подпись к видео.",
				name,
			),
			tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: make(
					[][]tgbotapi.InlineKeyboardButton,
					0,
				),
			},
		)

		return

	}

	if exists && session.SearchMode {
		b.DeleteChatMessage(
			message.Chat.ID,
			message.MessageID,
		)
		query := message.Text

		if query == "" {
			return
		}

		items, err := repository.SearchItems(
			context.Background(),
			b.DB,
			query,
		)

		if err != nil {
			log.Println(err)
			return
		}

		session.SearchMode = false
		Sessions[message.Chat.ID] = session

		b.ShowSearchResults(
			message.Chat.ID,
			items,
		)

		return
	}
}

func (b *Bot) HandleCategory(callback *tgbotapi.CallbackQuery, categoryID string) {

	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	id, err := strconv.ParseInt(
		categoryID,
		10,
		64,
	)

	if err != nil {
		log.Println(err)
		return
	}

	items, err := repository.GetItemsByCategoryID(
		context.Background(),
		b.DB,
		id,
	)

	if err != nil {
		log.Println(err)
		return
	}

	category, err := repository.GetCategoryByID(
		context.Background(),
		b.DB,
		id,
	)

	if err != nil {
		log.Println(err)
		return
	}

	text := fmt.Sprintf("%s\n\nВыберите элемент:", category.Name)

	if category.Name == "Выступления" {
		text = fmt.Sprintf("%s\n\nВыберите номер:", category.Name)
	}

	b.EditMenu(
		chatID,
		messageID,
		text,
		ItemsMenu(items),
	)
}

func (b *Bot) HandleItem(callback *tgbotapi.CallbackQuery, itemID string) {
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	id, err := strconv.ParseInt(
		itemID,
		10,
		64,
	)

	if err != nil {
		log.Println(err)
		return
	}

	item, err := repository.GetItemByID(
		context.Background(),
		b.DB,
		id,
	)

	text := fmt.Sprintf("%s", item.Name)

	videos, err := repository.GetVideosByItemID(
		context.Background(),
		b.DB,
		id,
	)

	if err != nil {
		log.Println(err)
		return
	}

	session, exists := Sessions[chatID]

	keyboard := ItemMenu(item.CategoryID)

	if exists && session.FromSearch {
		keyboard = SearchItemMenu()

		session.FromSearch = false
		Sessions[chatID] = session
	}

	b.EditMenu(
		chatID,
		messageID,
		text,
		keyboard,
	)

	b.SendVideos(
		chatID,
		videos,
	)
}

func (b *Bot) SendVideos(chatID int64, videos []models.Video) {
	if len(videos) == 0 {
		return
	}

	var media []interface{}

	for _, video := range videos {
		item := tgbotapi.NewInputMediaVideo(tgbotapi.FileID(video.FileID))

		item.Caption = video.Name

		media = append(media, item)
	}

	message := tgbotapi.NewMediaGroup(
		chatID,
		media,
	)

	messages, err := b.API.SendMediaGroup(message)

	if err != nil {
		fmt.Println(err)
		return
	}

	var messageIDs []int

	for _, message := range messages {
		messageIDs = append(messageIDs, message.MessageID)
	}

	session, exists := Sessions[chatID]

	if !exists {
		session = &UserSession{}
	}

	session.VideoMessages = messageIDs

	Sessions[chatID] = session
}

func (b *Bot) Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := b.API.GetUpdatesChan(updateConfig)

	for update := range updates {
		b.HandleUpdate(update)
	}
}
