package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

func InitServer() {

}

func WritePlayerList(stdin io.WriteCloser, token string) {
	/*
		// Valid fields from ClienTable
		ping      21
		userid    KU_Yy3t5oy_
		friend    false
		playerage 43
		userflags 0
		name      420
		admin     true
		steamid
		prefab   wx78
		colour
	*/
	write_lua_loop := []string{
		"pds = \"\"",
		"clients = TheNet:GetClientTable()",
		"for idx,_ in pairs(clients) do if (clients[idx].ping ~= nil) then pds = pds .. clients[idx].name .. \"\\n\" end end",
		"playerlist = io.open(\"starvewrap_playerlist\", \"w\")",
		"playerlist:write(pds)",
		"playerlist:close()",
	}

	for _, code := range write_lua_loop {
		stdin.Write([]byte(code + "\n"))
	}
}

func GetNumPlayers(dir string) int {
	var numP int
	file := dir + "/starvewrap_playerlist"
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Could not parse %v, %v\n", file, err)
	}
	numP = len(strings.Split(string(contents), "\n")) - 1
	//fmt.Printf("%v", numP)
	return numP
}
