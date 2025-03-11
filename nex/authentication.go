package nex

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PretendoNetwork/fast-racing-neo/globals"
	
)

func StartAuthenticationServer() {
	globals.AuthenticationServer = nex.NewPRUDPServer()
	globals.AuthenticationServer.ByteStreamSettings.UseStructureHeader = true

	globals.AuthenticationEndpoint = nex.NewPRUDPEndPoint(1)
	globals.AuthenticationEndpoint.ServerAccount = globals.AuthenticationServerAccount
	globals.AuthenticationEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.AuthenticationEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.AuthenticationServer.BindPRUDPEndPoint(globals.AuthenticationEndpoint)

	globals.AuthenticationServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(3, 8, 15))
	globals.AuthenticationServer.AccessKey = "811aa39f"

	globals.AuthenticationEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("==Fast Racing NEO - Auth==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID)
		fmt.Printf("Method ID: %#v\n", request.MethodID)
		fmt.Println("===============")
	})

	globals.AuthenticationEndpoint.OnError(func(err *nex.Error) {
		globals.Logger.Errorf("Auth: %v", err)
	})

	registerCommonAuthenticationServerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_FAST_RACING_NEO_AUTHENTICATION_SERVER_PORT"))

	globals.AuthenticationServer.Listen(port)
}
