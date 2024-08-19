package handlers

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	structs "github.com/tapmahtec/TNL_bot"
	"github.com/tapmahtec/TNL_bot/service"
)

type Bot struct {
	Session  *discordgo.Session
	servises *service.Service
}

func NewBot(token string, allowedChannelID string, services *service.Service) (*Bot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %v", err)
	}

	// Включаем необходимые интенты
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	bot := &Bot{
		Session:  dg,
		servises: services,
	}

	// Обработчик сообщений
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Проверяем, что сообщение не от бота
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Проверяем, что сообщение из разрешенного канала
		if m.ChannelID != allowedChannelID {
			fmt.Printf("Message ignored because it's not from the allowed channel: %s\n", m.ChannelID)
			return
		}
		parts := strings.Split(m.Content, " ")
		command := strings.ToLower(parts[0])

		// Обрабатываем команды
		switch command {
		case "!add":
			if len(parts) < 3 {
				s.ChannelMessageSend(m.ChannelID, "Введте !add [Системное имя активности] [тэг] [тэг]")
				return
			}
			activity, err := bot.servises.GetActivityBySid(parts[1])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Введте !add [Системное имя активности] [тэг] [тэг]... Ошибка: "+err.Error())
				return
			}
			if activity.Id == 0 {
				s.ChannelMessageSend(m.ChannelID, "Введте !add [Системное имя активности] [тэг] [тэг]... Ошибка: не верное имя активности")
				return
			}
			playersNames := parts[2:]

			var players []structs.Players
			for _, name := range playersNames {
				player, err := bot.servises.GetPlayerByName(name)
				if err != nil || player.Id == 0 {
					_, err = bot.servises.AddPlayer(structs.Players{Name: name})
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "Введите !add [Игроки]... Ошибка: "+err.Error())
						return
					}
					player, err = bot.servises.GetPlayerByName(name)
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "Введите !add [Игроки]... Ошибка: "+err.Error())
						return
					}
				}
				players = append(players, player)
			}

			for _, p := range players {
				bot.servises.AddPlayerActivity(p, activity)
				bot.servises.UpdatePlayerScore(p.Id, p.Score+activity.Score)
			}
			bot.WriteFullUsersScore(dg, m)
			s.ChannelMessageSend(m.ChannelID, "Выполнено ")

		case "!add_activity":
			var activity structs.Activities
			if len(parts) < 2 {
				s.ChannelMessageSend(m.ChannelID, "Введте !add_activity [Название] [Системное имя (на английском)] [Количество очков (число)]")
				return
			}
			activity.Name = parts[1]
			activity.Sid = parts[2]
			score, err := strconv.Atoi(parts[3])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Вы не верно ввели счет")
				return
			}
			activity.Score = score

			_, err = bot.servises.CreateActivity(activity)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Введте \"Название\", \"Системное имя\" (на английском), \"Количество очков\" (число)"+err.Error())
				return
			}

			s.ChannelMessageSend(m.ChannelID, "Активность успешно создана")
		case "!delete_activity":
			if len(parts) < 2 {
				s.ChannelMessageSend(m.ChannelID, "Введте !delete_activity [системное имя]")
				return
			}
			sid := parts[1]
			err := bot.servises.DeleteActivityBySid(sid)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Введте !delete_activity [системное имя]. Ошибка: "+err.Error())
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Активность успешно удалена")
		case "!list_activity":

			act, err := bot.servises.GetActivities()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Ошибка получения списка активностей")
				return
			}
			var message string
			for _, a := range act {
				message += fmt.Sprintf("Название: %s,\t Системное имя: %s, \t Количество очков: %d\n", a.Name, a.Sid, a.Score)
			}
			s.ChannelMessageSend(m.ChannelID, message)

		case "!top":
			var players []structs.Players
			if len(parts) < 2 {
				players, err = bot.servises.GetTopPlayers(100)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Ошибка в получении игроков"+err.Error())
					return
				}
			} else {
				limit, err := strconv.Atoi(parts[1])
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Введите число игроков для топа")
					return
				}
				players, err = bot.servises.GetTopPlayers(limit)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Ошибка в получении игроков")
					return
				}
			}

			var message string

			for _, p := range players {
				message += fmt.Sprintf("%s: %d очков\n", p.Name, p.Score)
			}

			s.ChannelMessageSend(m.ChannelID, message)
			bot.WriteFullUsersScore(dg, m)
		}

	})

	return bot, nil
}

func (b *Bot) Start() error {
	err := b.Session.Open()
	if err != nil {
		return fmt.Errorf("error opening Discord session: %v", err)
	}

	// Ожидаем сигнала завершения
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Закрываем соединение с Discord
	b.Session.Close()
	return nil
}

func (b *Bot) WriteFullUsersScore(dg *discordgo.Session, m *discordgo.MessageCreate) error {
	users, err := b.servises.GetTopPlayers(10000)
	if err != nil {
		return err
	}
	for _, user := range users {
		name := user.Name
		score := strconv.Itoa(user.Score)
		err = writeUserScore(name, score, dg, m)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func writeUserScore(savedUser string, points string, dg *discordgo.Session, m *discordgo.MessageCreate) error {
	// Получи информацию о пользователе
	userID := extractUserID(savedUser)
	if userID == "" {
		return errors.New("не удалось извлечь ID пользователя из строки")
	}

	member, err := dg.GuildMember(os.Getenv("GUILD_ID"), userID)
	if err != nil {
		return err
	}

	// Текущий никнейм пользователя
	currentNickname := member.Nick
	if currentNickname == "" {
		currentNickname = member.User.Username
	}

	re := regexp.MustCompile(`\[\d+\]\s*`)
	newNickname := re.ReplaceAllString(currentNickname, "")

	newNickname = "[" + points + "] " + newNickname

	// Измени никнейм пользователя
	err = dg.GuildMemberNickname(m.GuildID, userID, newNickname)
	if err != nil {
		return err
	}
	return nil
}

func extractUserID(savedUser string) string {
	// Удали символ '@' и верни ID
	return strings.TrimPrefix(strings.TrimSuffix(savedUser, ">"), "<@")
}
