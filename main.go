package armamods

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type armaMod struct {
	workshopID string
	targetDir  string
}

var armaSteamID = "107410"

// DownloadMods ...
func DownloadMods(steamCMDDir string, login string, password string, armaDir string, modListFile string) {
	steamCMDDir = strings.TrimRight(steamCMDDir, "/")
	steamCMDDir = strings.TrimRight(steamCMDDir, "\\")
	armaDir = strings.TrimRight(armaDir, "/")
	armaDir = strings.TrimRight(armaDir, "\\")

	mods := getFilesList(modListFile)
	for _, mod := range mods {
		modDir, err := downloadMod(steamCMDDir, login, password, mod.workshopID)
		if err != nil {
			panic(err)
		}

		err = createLink(modDir, armaDir+mod.targetDir)
		if err != nil {
			panic(err)
		}
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

func downloadMod(steamCMDDir string, login string, password string, workshopID string) (downloadDir string, err error) {
	loginArg := "+login $login $password"
	loginArg = strings.Replace(loginArg, "$login", login, 1)
	loginArg = strings.Replace(loginArg, "$password", password, 1)
	downloadArg := "+workshop_download_item 107410 $workshopID"
	downloadArg = strings.Replace(downloadArg, "$workshopID", workshopID, 1)

	commandExec := exec.Command(steamCMDDir+"/steamcmd.exe", loginArg, downloadArg, "+quit")
	stderr, err := commandExec.StdoutPipe()
	if err != nil {
		return "", err
	}

	commandExec.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Print(m + " ")
	}
	commandExec.Wait()

	return steamCMDDir + "/steamapps/workshop/content/" + armaSteamID + "/" + workshopID, nil
}

func createLink(from string, to string) error {
	if _, err := os.Stat(to); os.IsNotExist(err) {
		err = os.MkdirAll(to, os.ModeDir)
		if err != nil {
			return err
		}
	}

	//TODO: nested files
	fileInfos, err := ioutil.ReadDir(from)

	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		targetPath := to + "/" + fileInfo.Name()
		sourcePath := from + "/" + fileInfo.Name()

		if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
			continue
		}

		err = os.Link(sourcePath, targetPath)

		if err != nil {
			return err
		}
	}

	return nil
}
