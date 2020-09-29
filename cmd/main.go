package main

import (
	"flag"
	armamods "steam/workshop/armamods"
)

func main() {
	steamCMDDirArg := flag.String("steamcmddir", "", "Steam cmd path")
	loginArg := flag.String("login", "", "Steam login")
	passwordArg := flag.String("password", "", "Steam password")
	armaDirArg := flag.String("armadir", "", "ArmA 3 Dir")
	modListArg := flag.String("addonlist", "", "Addon file")
	flag.Parse()

	armamods.DownloadMods(*steamCMDDirArg, *loginArg, *passwordArg, *armaDirArg, *modListArg)
}
