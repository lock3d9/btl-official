package interfface

import (
	fmt "fmt"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/eiannone/keyboard"
)

type MenuItem struct {
	Title, Description, Warning string
}

type Category struct {
	Title string
	Items []MenuItem
}

var warningShown bool

var restorePointItem = MenuItem{
	Title:       "[РЕКОМЕНДУЕМ] Создать точку восстановления",
	Description: "Создает снимки системных файлов Windows на случай сбоев.",
	Warning:     "",
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func nav(key keyboard.Key, char rune, cursor, max int) int {
	switch {
	case key == keyboard.KeyArrowUp || key == keyboard.KeyArrowLeft || char == 'w' || char == 'W' || char == 'a' || char == 'A' || char == 'ц' || char == 'Ц' || char == 'ф' || char == 'Ф':
		return (cursor - 1 + max) % max
	case key == keyboard.KeyArrowDown || key == keyboard.KeyArrowRight || char == 's' || char == 'S' || char == 'd' || char == 'D' || char == 'ы' || char == 'Ы' || char == 'в' || char == 'В':
		return (cursor + 1) % max
	}
	return cursor
}

func printOptions(opts []string, cursor int) {
	for i, opt := range opts {
		if i == cursor {
			fmt.Printf(" \033[36m➔ %s\033[0m\n", opt)
		} else {
			fmt.Printf("    %s\n", opt)
		}
	}
}

func showStartupWarning() bool {
	opts := []string{"[ Создать точку восстановления сейчас ]", "[ Пропустить и перейти в меню ]"}
	cursor := 0

	for {
		clearScreen()
		fmt.Println("\033[31m========================================================\033[0m")
		fmt.Println("\033[31m                 ВАЖНОЕ ПРЕДУПРЕЖДЕНИЕ                  \033[0m")
		fmt.Println("\033[31m========================================================\033[0m\n")
		fmt.Println("Перед использованием оптимизаций настоятельно рекомендуется")
		fmt.Println("создать точку восстановления Windows.")
		fmt.Println("\nЭто позволит вернуть систему в исходное состояние при сбоях.\n")

		printOptions(opts, cursor)

		char, key, err := keyboard.GetKey()
		if err != nil || key == keyboard.KeyEsc {
			return false
		}
		if key == keyboard.KeyEnter {
			return cursor == 0
		}

		cursor = nav(key, char, cursor, len(opts))
	}
}

func ConfirmAction(item MenuItem) bool {
	opts := []string{"[ Да, продолжить ]", "[ Отмена / Назад ]"}
	cursor := 0

	for {
		clearScreen()
		fmt.Printf("\033[36m=== %s ===\033[0m\n\nОписание:\n%s\n\n", item.Title, item.Description)
		if item.Warning != "" {
			fmt.Printf("\033[31m[ПРЕДУПРЕЖДЕНИЕ]: %s\033[0m\n\n", item.Warning)
		}
		fmt.Println("Вы действительно хотите выполнить это действие?")
		printOptions(opts, cursor)

		char, key, err := keyboard.GetKey()
		if err != nil || key == keyboard.KeyEsc {
			return false
		}
		if key == keyboard.KeyEnter {
			return cursor == 0
		}

		cursor = nav(key, char, cursor, len(opts))
	}
}

func showSubMenu(cat Category) (string, bool) {
	opts := make([]string, len(cat.Items)+1)
	for i, item := range cat.Items {
		opts[i] = item.Title
	}
	opts[len(cat.Items)] = "[ Назад в главное меню ]"

	cursor := 0
	for {
		clearScreen()
		fmt.Printf("\033[36m=== РАЗДЕЛ: %s ===\033[0m\n\nВыберите опцию (W/S — навигация, Enter — выбор, Esc — назад):\n\n", cat.Title)
		printOptions(opts, cursor)

		if cursor < len(cat.Items) {
			fmt.Printf("\n\033[90m----------------------------------------\033[0m\n\033[33mОписание:\033[0m %s\n", cat.Items[cursor].Description)
			if w := cat.Items[cursor].Warning; w != "" {
				fmt.Printf("\033[31m[ПРЕДУПРЕЖДЕНИЕ]: %s\033[0m\n", w)
			}
		}

		char, key, err := keyboard.GetKey()
		if err != nil || key == keyboard.KeyEsc {
			return "", false
		}
		if key == keyboard.KeyEnter {
			if cursor == len(cat.Items) {
				return "", false
			}
			if ConfirmAction(cat.Items[cursor]) {
				return cat.Items[cursor].Title, true
			}
		}

		cursor = nav(key, char, cursor, len(opts))
	}
}

func Showmainmenu(version, buildNum, buildDate string) string {
	if err := keyboard.Open(); err != nil {
		fmt.Printf("[BTL] Ошибка инициализации клавиатуры: %v\n", err)
		os.Exit(1)
	}
	defer keyboard.Close()

	if !warningShown {
		warningShown = true
		if showStartupWarning() {
			if ConfirmAction(restorePointItem) {
				return restorePointItem.Title
			}
		}
	}

	categories := getCategories()
	opts := make([]string, len(categories)+2)
	opts[0] = restorePointItem.Title

	for i, cat := range categories {
		opts[i+1] = "[" + cat.Title + "]"
	}
	opts[len(categories)+1] = "Выход"

	cursor := 0
	for {
		clearScreen()
		figure.NewFigure("// BtL //", "cyberlarge", true).Print()
		fmt.Printf("\033[36m v%s (build %s | data: %s)\033[0m\n\nВыберите раздел (W/S — навигация, Enter — выбор):\n\n", version, buildNum, buildDate)
		printOptions(opts, cursor)

		char, key, err := keyboard.GetKey()
		if err != nil {
			break
		}

		if key == keyboard.KeyEnter {
			if cursor == 0 {
				if ConfirmAction(restorePointItem) {
					return restorePointItem.Title
				}
			} else if cursor == len(categories)+1 {
				return "Выход"
			} else {
				if res, ok := showSubMenu(categories[cursor-1]); ok {
					return res
				}
			}
		}

		cursor = nav(key, char, cursor, len(opts))
	}

	return "Выход"
}

func getCategories() []Category {
	return []Category{
		{
			Title: "Система и Восстановление",
			Items: []MenuItem{
				{"Очистка временных файлов", "Удаляет кэш системы, временные файлы обновлений и корзину.", ""},
				{"Отключение телеметрии", "Отключает службы фонового сбора данных и отправку отчетов в Microsoft.", ""},
			},
		},
		{
			Title: "Оптимизация и Твики",
			Items: []MenuItem{
				{"Удалить встроенный Windows Defender", "Полностью вырезает Защитник Windows из системы.", "Действие необратимо! Ваш ПК останется без встроенной антивирусной защиты."},
				{"Удалить OneDrive", "Удаляет клиент OneDrive и отвязывает синхронизацию системных папок.", "Убедитесь, что важные файлы сохранены на локальном диске."},
				{"Контроль автозагрузки", "Позволяет просмотреть и отключить программы, запускаемые вместе с Windows.", ""},
			},
		},
		{
			Title: "Тестирование и Анализ дисков",
			Items: []MenuItem{
				{"Проверить диск на чтение, запись и скорость", "Проводит тестирование выбранного накопителя путем записи и чтения временных блоков.", "Может занять продолжительное время. Не закрывайте программу во время теста."},
				{"Проверить флешку на муляж", "Проверяет реальную емкость накопителя (защита от поддельных флешек).", "Все данные на накопителе могут быть перезаписаны! Сделайте резервную копию."},
				{"Анализ «тяжелых» файлов", "Сканирует выбранный диск и выводит список самых крупных файлов.", ""},
			},
		},
		{
			Title: "Диагностика и Мониторинг",
			Items: []MenuItem{
				{"Убить зависшие приложения", "Принудительно завершает процессы, которые не отвечают на запросы ОС.", "Несохраненные данные в зависших программах будут утеряны."},
				{"Полная информация о ПК", "Выводит детальные характеристики процессора, ОЗУ, видеокарты и материнской платы.", ""},
			},
		},
	}
}
