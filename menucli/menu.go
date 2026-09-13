package interfface

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/common-nighthawk/go-figure"
	"github.com/eiannone/keyboard"
)

type MenuItem struct {
	Title       string
	Description string
	Warning     string
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func ConfirmAction(item MenuItem) bool {
	cursor := 0
	options := []string{"[ Да, продолжить ]", "[ Отмена / Назад ]"}

	for {
		clearScreen()

		fmt.Printf("\033[36m=== %s ===\033[0m\n\n", item.Title)
		fmt.Printf("Описание:\n%s\n\n", item.Description)

		if item.Warning != "" {
			fmt.Printf("\033[31m[ПРЕДУПРЕЖДЕНИЕ]: %s\033[0m\n\n", item.Warning)
		}

		fmt.Println("Вы действительно хотите выполнить это действие?")

		for i, opt := range options {
			if i == cursor {
				fmt.Printf(" \033[36m➔ %s\033[0m\n", opt)
			} else {
				fmt.Printf("    %s\n", opt)
			}
		}

		char, key, err := keyboard.GetKey()
		if err != nil {
			return false
		}

		if key == keyboard.KeyArrowUp || key == keyboard.KeyArrowLeft || char == 'w' || char == 'W' || char == 'a' || char == 'A' || char == 'ц' || char == 'Ц' || char == 'ф' || char == 'Ф' {
			if cursor > 0 {
				cursor--
			} else {
				cursor = len(options) - 1
			}
		} else if key == keyboard.KeyArrowDown || key == keyboard.KeyArrowRight || char == 's' || char == 'S' || char == 'd' || char == 'D' || char == 'ы' || char == 'Ы' || char == 'в' || char == 'В' {
			if cursor < len(options)-1 {
				cursor++
			} else {
				cursor = 0
			}
		} else if key == keyboard.KeyEnter {
			return cursor == 0
		} else if key == keyboard.KeyEsc {
			return false
		}
	}
}

func Showmainmenu(version string, buildNum string, buildDate string) string {
	items := []MenuItem{
		{
			Title:       "[РЕКОМЕНДУЕМ] Создать точку восстановления",
			Description: "Создает снимки системных файлов Windows на случай сбоев.",
			Warning:     "",
		},
		{
			Title:       "Проверить диск на чтение, запись и скорость",
			Description: "Проводит тестирование выбранного накопителя путем записи и чтения временных блоков.",
			Warning:     "Может занять продолжительное время. Не закрывайте программу во время теста.",
		},
		{
			Title:       "Проверить флешку на муляж",
			Description: "Проверяет реальную емкость накопителя (защита от поддельных флешек).",
			Warning:     "Все данные на накопителе могут быть перезаписаны! Сделайте резервную копию.",
		},
		{
			Title:       "Очистка временных файлов",
			Description: "Удаляет кэш системы, временные файлы обновлений и корзину.",
			Warning:     "",
		},
		{
			Title:       "Отключение телеметрии",
			Description: "Отключает службы фонового сбора данных и отправку отчетов в Microsoft.",
			Warning:     "",
		},
		{
			Title:       "Удалить встроенный Windows Defender",
			Description: "Полностью вырезает Защитник Windows из системы.",
			Warning:     "Действие необратимо! Ваш ПК останется без встроенной антивирусной защиты.",
		},
		{
			Title:       "Удалить OneDrive",
			Description: "Удаляет клиент OneDrive и отвязывает синхронизацию системных папок.",
			Warning:     "Убедитесь, что важные файлы сохранены на локальном диске.",
		},
		{
			Title:       "Анализ «тяжелых» файлов",
			Description: "Сканирует выбранный диск и выводит список самых крупных файлов.",
			Warning:     "",
		},
		{
			Title:       "Убить зависшие приложения",
			Description: "Принудительно завершает процессы, которые не отвечают на запросы ОС.",
			Warning:     "Несохраненные данные в зависших программах будут утеряны.",
		},
		{
			Title:       "Контроль автозагрузки",
			Description: "Позволяет просмотреть и отключить программы, запускаемые вместе с Windows.",
			Warning:     "",
		},
		{
			Title:       "Полная информация о ПК",
			Description: "Выводит детальные характеристики процессора, ОЗУ, видеокарты и материнской платы.",
			Warning:     "",
		},
		{
			Title:       "Выход",
			Description: "Завершение работы программы.",
			Warning:     "",
		},
	}

	cursor := 0

	err := keyboard.Open()
	if err != nil {
		fmt.Printf("[BTL] Ошибка инициализации клавиатуры: %v\n", err)
		os.Exit(1)
	}
	defer keyboard.Close()

	for {
		clearScreen()

		logoapp := figure.NewFigure("// BtL //", "cyberlarge", true)
		logoapp.Print()

		fmt.Printf("\033[36m v%s (build %s | data: %s)\033[0m\n", version, buildNum, buildDate)
		fmt.Println("\nВыберите действие (W/S — навигация, Enter — выбор):\n")

		for i, item := range items {
			if i == cursor {
				fmt.Printf(" \033[36m➔ %s\033[0m\n", item.Title)
			} else {
				fmt.Printf("   %s\n", item.Title)
			}
		}

		char, key, err := keyboard.GetKey()
		if err != nil {
			break
		}

		if key == keyboard.KeyArrowUp || char == 'w' || char == 'W' || char == 'ц' || char == 'Ц' {
			if cursor > 0 {
				cursor--
			} else {
				cursor = len(items) - 1
			}
		} else if key == keyboard.KeyArrowDown || char == 's' || char == 'S' || char == 'ы' || char == 'Ы' {
			if cursor < len(items)-1 {
				cursor++
			} else {
				cursor = 0
			}
		} else if key == keyboard.KeyEnter {
			selectedItem := items[cursor]

			if selectedItem.Title == "Выход" {
				return selectedItem.Title
			}

			if ConfirmAction(selectedItem) {
				return selectedItem.Title
			}
		}
	}

	return "Выход"
}
