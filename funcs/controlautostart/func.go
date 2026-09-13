package autostart

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/eiannone/keyboard"
)

type AutoStartItem struct {
	Type   string `json:"Type"`
	Status string `json:"Status"`
	Name   string `json:"Name"`
	Path   string `json:"Path"`
}

func ManageAutostart() {
	fmt.Println("\n[BTL] Анализ автозагрузки системы...")

	psScript := `
	[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

	$regPaths = @(
		"Registry::HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run",
		"Registry::HKEY_LOCAL_MACHINE\Software\Microsoft\Windows\CurrentVersion\Run"
	)

	$keepList = @(
		"security", "defender", "antivirus", "avast", "kaspersky", "eset", "malwarebytes",
		"realtek", "audio", "nvidia", "amd", "intel", "bluetooth", "wi-fi", "network",
		"ctfmon", "windows", "microsoft", "touchpad", "synaptics", "vgtray"
	)

	$disableList = @(
		"update", "assistant", "helper", "download", "manager", "adobe", "browser",
		"chrome", "firefox", "opera", "yandex", "discord", "spotify", "epic",
		"skype", "teams", "slack", "cortana", "ccleaner", "torrent", "utorrent",
		"bittorrent", "steelseries", "gg", "epicgames", "onedrive"
	)

	$results = [System.Collections.Generic.List[PSObject]]::new()

	foreach ($path in $regPaths) {
		if (Test-Path $path) {
			$props = Get-ItemProperty -Path $path
			$props.PSObject.Properties | Where-Object { $_.Name -notin @('PSPath','PSParentPath','PSChildName','PSProvider') -and $_.Value -ne $null } | ForEach-Object {
				$name = $_.Name
				$val = $_.Value
				$status = "ПО ЖЕЛАНИЮ"

				foreach ($item in $disableList) {
					if ($name -match $item -or $val -match $item) {
						$status = "СОВЕТУЕМ ОТКЛЮЧИТЬ"
						break
					}
				}

				if ($status -eq "ПО ЖЕЛАНИЮ") {
					foreach ($item in $keepList) {
						if ($name -match $item -or $val -match $item) {
							$status = "ЛУЧШЕ ОСТАВИТЬ"
							break
						}
					}
				}

				$results.Add([PSCustomObject]@{
					Type   = "Реестр"
					Status = $status
					Name   = $name
					Path   = $val
				})
			}
		}
	}

	$startupFolders = @(
		"$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Startup",
		"$env:ProgramData\Microsoft\Windows\Start Menu\Programs\Startup"
	)

	foreach ($folder in $startupFolders) {
		if (Test-Path $folder) {
			Get-ChildItem -Path $folder -File | ForEach-Object {
				$results.Add([PSCustomObject]@{
					Type   = "Папка"
					Status = "СОВЕТУЕМ ОТКЛЮЧИТЬ"
					Name   = $_.Name
					Path   = $_.FullName
				})
			}
		}
	}

	$results | ConvertTo-Json -Compress
	`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	outputBytes, err := cmd.Output()
	if err != nil {
		fmt.Printf("\n[BTL] Ошибка при чтении автозагрузки: %v\n", err)
		return
	}

	var items []AutoStartItem
	if err := json.Unmarshal(outputBytes, &items); err != nil {
		var singleItem AutoStartItem
		if errSingle := json.Unmarshal(outputBytes, &singleItem); errSingle == nil {
			items = append(items, singleItem)
		}
	}

	fmt.Println("\n==================================================")
	fmt.Println("\033[32m[BTL] АНАЛИЗ АВТОЗАГРУЗКИ\033[0m")
	fmt.Println("==================================================")

	if len(items) == 0 {
		fmt.Println(" Автозагрузка пуста.")
	} else {
		for _, item := range items {
			shortPath := item.Path
			if len(shortPath) > 150 {
				shortPath = "..." + shortPath[len(shortPath)-47:]
			}

			switch item.Status {
			case "ЛУЧШЕ ОСТАВИТЬ":
				fmt.Printf(" \033[32m[ОСТАВИТЬ]    %-15s | %s\033[0m\n", item.Name, shortPath)
			case "СОВЕТУЕМ ОТКЛЮЧИТЬ":
				fmt.Printf(" \033[31m[ВЫКЛЮЧИТЬ]   %-15s | %s\033[0m\n", item.Name, shortPath)
			default:
				fmt.Printf(" \033[33m[ПО ЖЕЛАНИЮ]  %-15s | %s\033[0m\n", item.Name, shortPath)
			}
		}
	}

	fmt.Println("==================================================")
	fmt.Println("[BTL] Действие выполнено успешно!")
	fmt.Print("\nНажмите Enter для возврата в меню...")

	if err := keyboard.Open(); err == nil {
		defer keyboard.Close()
		for {
			_, key, err := keyboard.GetKey()
			if err != nil || key == keyboard.KeyEnter {
				break
			}
		}
	}
}
