package autostart

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/eiannone/keyboard"
	"golang.org/x/text/encoding/charmap"
)

func ManageAutostart() {
	fmt.Println("\n[BTL] Анализ автозагрузки системы...")

	psCommand := `
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

				[PSCustomObject]@{
					Type   = "Реестр"
					Status = $status
					Name   = $name
					Path   = $val
				}
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
				[PSCustomObject]@{
					Type   = "Папка"
					Status = "СОВЕТУЕМ ОТКЛЮЧИТЬ"
					Name   = $_.Name
					Path   = $_.FullName
				}
			}
		}
	}
	`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCommand)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("\n[BTL] Ошибка при чтении автозагрузки: %v\n", err)
		return
	}

	decoder := charmap.CodePage866.NewDecoder()
	outputStr, err := decoder.String(string(outputBytes))
	if err != nil {
		outputStr = string(outputBytes)
	}

	fmt.Println("\n==================================================")
	fmt.Println("\033[32m[BTL] АНАЛИЗ АВТОЗАГРУЗКИ И РЕКОМЕНДАЦИИ\033[0m")
	fmt.Println("==================================================")

	trimmedOutput := strings.TrimSpace(outputStr)
	if trimmedOutput == "" {
		fmt.Println(" Автозагрузка пуста.")
	} else {
		lines := strings.Split(outputStr, "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && !strings.HasPrefix(trimmed, "---") && !strings.HasPrefix(trimmed, "ComputerName") {
				if strings.Contains(trimmed, "ЛУЧШЕ ОСТАВИТЬ") {
					fmt.Printf(" \033[32m[ОСТАВИТЬ]    %s\033[0m\n", trimmed)
				} else if strings.Contains(trimmed, "СОВЕТУЕМ ОТКЛЮЧИТЬ") {
					fmt.Printf(" \033[31m[ВЫКЛЮЧИТЬ]   %s\033[0m\n", trimmed)
				} else {
					fmt.Printf(" \033[33m[ПО ЖЕЛАНИЮ]  %s\033[0m\n", trimmed)
				}
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
