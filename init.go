package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PretendoNetwork/fast-racing-neo/globals"
	pb "github.com/PretendoNetwork/grpc-go/account"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func init() {
	globals.Logger = plogger.NewLogger()

	var err error

	err = godotenv.Load()
	if err != nil {
		globals.Logger.Warning("Error loading .env file")
	}

	authenticationServerPort := os.Getenv("PN_FAST_RACING_NEO_AUTHENTICATION_SERVER_PORT")
	aesKey := os.Getenv("PN_FRIENDS_CONFIG_AES_KEY")
	secureServerHost := os.Getenv("PN_FAST_RACING_NEO_SECURE_SERVER_HOST")
	secureServerPort := os.Getenv("PN_FAST_RACING_NEO_SECURE_SERVER_PORT")
	accountGRPCHost := os.Getenv("PN_FAST_RACING_NEO_ACCOUNT_GRPC_HOST")
	accountGRPCPort := os.Getenv("PN_FAST_RACING_NEO_ACCOUNT_GRPC_PORT")
	accountGRPCAPIKey := os.Getenv("PN_FAST_RACING_NEO_ACCOUNT_GRPC_API_KEY")

	kerberosPassword := make([]byte, 0x10)
	_, err = rand.Read(kerberosPassword)
	if err != nil {
		globals.Logger.Error("Error generating Kerberos password")
		os.Exit(0)
	}

	globals.KerberosPassword = string(kerberosPassword)

	globals.InitAccounts()

	if strings.TrimSpace(authenticationServerPort) == "" {
		globals.Logger.Error("PN_FAST_RACING_NEO_AUTHENTICATION_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(authenticationServerPort); err != nil {
		globals.Logger.Errorf("PN_FAST_RACING_NEO_AUTHENTICATION_SERVER_PORT is not a valid port. Expected 0-65535, got %s", authenticationServerPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_FAST_RACING_NEO_AUTHENTICATION_SERVER_PORT is not a valid port. Expected 0-65535, got %s", authenticationServerPort)
		os.Exit(0)
	}

	if strings.TrimSpace(secureServerHost) == "" {
		globals.Logger.Error("PN_FAST_RACING_NEO_SECURE_SERVER_HOST environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(secureServerPort) == "" {
		globals.Logger.Error("PN_FAST_RACING_NEO_SECURE_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(secureServerPort); err != nil {
		globals.Logger.Errorf("PN_FAST_RACING_NEO_SECURE_SERVER_PORT is not a valid port. Expected 0-65535, got %s", secureServerPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_FAST_RACING_NEO_SECURE_SERVER_PORT is not a valid port. Expected 0-65535, got %s", secureServerPort)
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCHost) == "" {
		globals.Logger.Error("PN_FAST_RACING_NEO_ACCOUNT_GRPC_HOST environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCPort) == "" {
		globals.Logger.Error("PN_FAST_RACING_NEO_ACCOUNT_GRPC_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(accountGRPCPort); err != nil {
		globals.Logger.Errorf("PN_FAST_RACING_NEO_ACCOUNT_GRPC_PORT is not a valid port. Expected 0-65535, got %s", accountGRPCPort)
		os.Exit(0)
	} else if port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_FAST_RACING_NEO_ACCOUNT_GRPC_PORT is not a valid port. Expected 0-65535, got %s", accountGRPCPort)
		os.Exit(0)
	}

	if strings.TrimSpace(aesKey) == "" {
		globals.Logger.Error("PN_FRIENDS_CONFIG_AES_KEY environment variable not set")
		os.Exit(0)
	} else {
		globals.AESKey, err = hex.DecodeString(os.Getenv("PN_FRIENDS_CONFIG_AES_KEY"))
		if err != nil {
			globals.Logger.Criticalf("Failed to decode AES key: %v", err)
			os.Exit(0)
		}
	}

	if strings.TrimSpace(accountGRPCAPIKey) == "" {
		globals.Logger.Warning("Insecure gRPC server detected. PN_FAST_RACING_NEO_ACCOUNT_GRPC_API_KEY environment variable not set")
	}

	globals.GRPCAccountClientConnection, err = grpc.NewClient(fmt.Sprintf("dns:%s:%s", accountGRPCHost, accountGRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		globals.Logger.Criticalf("Failed to connect to account gRPC server: %v", err)
		os.Exit(0)
	}

	globals.GRPCAccountClient = pb.NewAccountClient(globals.GRPCAccountClientConnection)
	globals.GRPCAccountCommonMetadata = metadata.Pairs(
		"X-API-Key", accountGRPCAPIKey,
	)
	globals.Postgres, err = sql.Open("postgres", os.Getenv("PN_FAST_RACING_NEO_POSTGRES_URI"))
	if err != nil {
		globals.Logger.Critical(err.Error())
	}
	globals.Logger.Success("Connected to Postgres!")
}
