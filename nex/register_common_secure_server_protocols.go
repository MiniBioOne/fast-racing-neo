package nex

import (
	"crypto/rand"
	"encoding/binary"

	"github.com/PretendoNetwork/fast-racing-neo/globals"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	commonmatchmaking "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making"
	commonmatchmakingext "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making-ext"
	commonmatchmakeextension "github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension"
	commonnattraversal "github.com/PretendoNetwork/nex-protocols-common-go/v2/nat-traversal"
	commonranking "github.com/PretendoNetwork/nex-protocols-common-go/v2/ranking"
	commonsecure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	commonutility "github.com/PretendoNetwork/nex-protocols-common-go/v2/utility"
	matchmaking "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	matchmakingext "github.com/PretendoNetwork/nex-protocols-go/v2/match-making-ext"
	match_making_types "github.com/PretendoNetwork/nex-protocols-go/v2/match-making/types"
	matchmakeextension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	nattraversal "github.com/PretendoNetwork/nex-protocols-go/v2/nat-traversal"
	ranking "github.com/PretendoNetwork/nex-protocols-go/v2/ranking"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	utility "github.com/PretendoNetwork/nex-protocols-go/v2/utility"
)

func CreateReportDBRecord(_ types.PID, _ types.UInt32, _ types.QBuffer) error {
	return nil
}

func cleanupMatchmakeSessionSearchCriteriasHandler(searchCriterias types.List[match_making_types.MatchmakeSessionSearchCriteria]) {

}

// thank you mr. Trace Pretendo for the knowledge
func generateNEXUniqueIDHandler() uint64 {
	var uniqueID uint64

	err := binary.Read(rand.Reader, binary.BigEndian, &uniqueID)
	if err != nil {
		globals.Logger.Error(err.Error())
	}

	return uniqueID
}

func onAfterAutoMatchmakeWithParamPostpone(_ nex.PacketInterface, _ match_making_types.AutoMatchmakeParam) {
	globals.MatchmakingManager.Mutex.Lock()

	_, err := globals.MatchmakingManager.Database.Exec(`UPDATE matchmaking.matchmake_sessions SET open_participation=true WHERE game_mode=12`)
	if err != nil {
		globals.Logger.Error(err.Error())
	}

	globals.MatchmakingManager.Mutex.Unlock()
}

func registerCommonSecureServerProtocols() {
	secureProtocol := secure.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	commonSecureProtocol := commonsecure.NewCommonProtocol(secureProtocol)

	commonSecureProtocol.CreateReportDBRecord = CreateReportDBRecord

	natTraversalProtocol := nattraversal.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(natTraversalProtocol)
	commonnattraversal.NewCommonProtocol(natTraversalProtocol)

	matchMakingProtocol := matchmaking.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchMakingProtocol)
	commonMatchMakingProtocol := commonmatchmaking.NewCommonProtocol(matchMakingProtocol)
	commonMatchMakingProtocol.SetManager(globals.MatchmakingManager)

	matchMakingExtProtocol := matchmakingext.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchMakingExtProtocol)
	commonMatchMakingExtProtocol := commonmatchmakingext.NewCommonProtocol(matchMakingExtProtocol)
	commonMatchMakingExtProtocol.SetManager(globals.MatchmakingManager)

	matchmakeExtensionProtocol := matchmakeextension.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchmakeExtensionProtocol)
	commonMatchmakeExtensionProtocol := commonmatchmakeextension.NewCommonProtocol(matchmakeExtensionProtocol)
	commonMatchmakeExtensionProtocol.CleanupMatchmakeSessionSearchCriterias = cleanupMatchmakeSessionSearchCriteriasHandler
	commonMatchmakeExtensionProtocol.OnAfterAutoMatchmakeWithParamPostpone = onAfterAutoMatchmakeWithParamPostpone
	commonMatchmakeExtensionProtocol.SetManager(globals.MatchmakingManager)

	utilityProtocol := utility.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(utilityProtocol)
	commonUtilityProtocol := commonutility.NewCommonProtocol(utilityProtocol)
	commonUtilityProtocol.GenerateNEXUniqueID = generateNEXUniqueIDHandler

	rankingProtocol := ranking.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(rankingProtocol)
	commonranking.NewCommonProtocol(rankingProtocol)
}
