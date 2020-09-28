package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type armaMod struct {
	workshopID string
	targetDir  string
}

func main() {
	steamcmd := flag.String("steamcmd", "", "Steam cmd path")
	login := flag.String("login", "", "Steam login")
	password := flag.String("password", "", "Steam password")
	armaDir := flag.String("armadir", "", "ArmA 3 Dir")
	modList := flag.String("addonlist", "", "Addon file")
	flag.Parse()

	println(*steamcmd)
	println(*login)
	println(*password)
	println(*armaDir)
	println(*modList)

	mods := getFilesList(*modList)
	for _, mod := range mods {
		downloadMod(*steamcmd, *login, *password, mod.workshopID, *armaDir+mod.targetDir)
	}
}

func getFilesList(modlistFile string) []armaMod {
	file, err := os.Open(modlistFile)

	if err != nil {
		panic("failed to open")
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	modList := []armaMod{}

	for _, line := range text {
		lineParts := strings.Split(line, ";")

		if strings.HasPrefix(lineParts[0], "#") {
			continue
		}

		mod := armaMod{
			workshopID: lineParts[0],
			targetDir:  lineParts[1],
		}

		modList = append(modList, mod)
	}

	return modList
}

func downloadMod(steamcmd string, login string, password string, workshopID string, targetDir string) {
	loginArg := "+login $login $password"
	loginArg = strings.Replace(loginArg, "$login", login, 1)
	loginArg = strings.Replace(loginArg, "$password", password, 1)
	downloadArg := "+workshop_download_item 107410 $workshopID +quit"
	downloadArg = strings.Replace(downloadArg, "$workshopID", workshopID, 1)

	commandExec := exec.Command(steamcmd, loginArg, downloadArg, "+quit")
	stderr, _ := commandExec.StdoutPipe()
	commandExec.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	commandExec.Wait()
}
