package handlers

import (
	//"fmt"
	"io"
)

func InitServer() {

}

func GetPlayerList(stdin io.WriteCloser, token string) {
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
		"playerlist = io.open(\"playerlist\", \"w\")",
		"playerlist:write(pds)",
		"playerlist:close()",
	}

	for _, code := range write_lua_loop {
		stdin.Write([]byte(code + "\n"))
	}
}
