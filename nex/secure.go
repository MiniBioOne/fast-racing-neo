package nex

import (
	"fmt"
	"os"
	"strconv"

	commonglobals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"

	"github.com/PretendoNetwork/fast-racing-neo/globals"
	"github.com/PretendoNetwork/nex-go/v2"
)

func StartSecureServer() {
	globals.SecureServer = nex.NewPRUDPServer()
	globals.SecureServer.ByteStreamSettings.UseStructureHeader = true

	globals.SecureEndpoint = nex.NewPRUDPEndPoint(1)
	globals.SecureEndpoint.IsSecureEndPoint = true
	globals.SecureEndpoint.ServerAccount = globals.SecureServerAccount
	globals.SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.SecureServer.BindPRUDPEndPoint(globals.SecureEndpoint)

	globals.SecureServer.LibraryVersions.SetDefault(nex.NewLibraryVersion(3, 9, 19))
	globals.SecureServer.AccessKey = "811aa39f"

	globals.SecureEndpoint.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("==Fast Racing NEO - Secure==")
		fmt.Printf("Protocol ID: %d\n", request.ProtocolID)
		fmt.Printf("Method ID: %d\n", request.MethodID)
		fmt.Println("===============")
	})

	globals.SecureEndpoint.OnError(func(err *nex.Error) {
		globals.Logger.Errorf("Secure: %v", err)
	})

	globals.MatchmakingManager = commonglobals.NewMatchmakingManager(globals.SecureEndpoint, globals.Postgres)

	registerCommonSecureServerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_FAST_RACING_NEO_SECURE_SERVER_PORT"))

	globals.SecureServer.Listen(port)
}
