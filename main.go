package main

import (
	checkdiskerrors "btl/funcs/checkdisk"
	checkflashff "btl/funcs/checkflashfake"
	"btl/funcs/checkheavyfiles"
	"btl/funcs/cleantemp"
	autostart "btl/funcs/controlautostart"
	"btl/funcs/createbackup"
	deleteonedriv "btl/funcs/deleteonedrive"
	"btl/funcs/disabledefender"
	sysinfo "btl/funcs/fullinformationpc"
	"btl/funcs/killtasks"
	interfface "btl/menucli"

	"fmt"
	"os"
	"os/exec"

	"golang.org/x/sys/windows"
)

var (
	Version   = "00.1"
	BuildNum  = "0"
	BuildDate = "unknown"
)

func isAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	return err == nil && member
}

func runAsAdmin() {
	elevateMe := exec.Command("cmd.exe", "/c", "powershell", "Start-Process", "\""+os.Args[0]+"\"", "-Verb", "RunAs")
	err := elevateMe.Run()
	if err != nil {
		fmt.Println("[BTL] Не удалось запросить права администратора:", err)
	}
}
func main() {
	for {
		if !isAdmin() {
			fmt.Println("[BTL] Программе требуются права администратора. Делаем запрос...")
			runAsAdmin()
			os.Exit(0)
		}
		choice := interfface.Showmainmenu(Version, BuildNum, BuildDate)
		switch choice {
		case "Выход":
			fmt.Println("[BTL] Выходим...")
			os.Exit(0)
		case "Проверить диск на чтение, запись и скорость":
			checkdiskerrors.Ccheckdsk()
		case "Полная информация о ПК":
			sysinfo.ShowSystemInfo()
		case "Проверить флешку на муляж":
			checkflashff.H2testwCheck()
		case "Убить зависшие приложения":
			killtasks.KillFrozenApps()
		case "Анализ «тяжелых» файлов":
			checkheavyfiles.FindHeavyFiles()
		case "Контроль автозагрузки":
			autostart.ManageAutostart()
		case "[РЕКОМЕНДУЕМ] Создать точку восстановления":
			createbackup.CreateRestorePoint()
		case "Очистка временных файлов":
			cleantemp.CleanTempFiles()
		case "Удалить встроенный Windows Defender":
			disabledefender.DisableDefender()
		case "Удалить OneDrive":
			deleteonedriv.RemoveOneDrive()
		default:
			fmt.Println("[BTL] Выполняется действие...")
		}
		fmt.Printf("[BTL] Версия: v%s (build #%s | %s)\n", Version, BuildNum, BuildDate)
	}
}
