package sourceserver

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"
	"strings"

	"github.com/ARTM2000/archive1/internal/archive/xerrors"
)

func NewSrvManager(srvRepo SrvRepository) SrvManager {
	return SrvManager{
		srvRepository: srvRepo,
	}
}

type newSrvSrcResult struct {
	NewServer *SrvSchema
	APIKey    string
}

type SrvManager struct {
	srvRepository SrvRepository
}

func (sm *SrvManager) generateAPIKey() (string, error) {
	// Define the length of the API key
	const apiKeyLength = 64

	// Generate a random byte slice with the specified length
	apiKeyBytes := make([]byte, apiKeyLength)
	if _, err := rand.Read(apiKeyBytes); err != nil {
		return "", err
	}

	// Encode the byte slice as a base64 string
	apiKey := base64.RawURLEncoding.EncodeToString(apiKeyBytes)

	// Replace any "+" and "/" characters with alphanumeric characters
	apiKey = strings.Map(func(r rune) rune {
		switch {
		case r == '+':
			return 'A'
		case r == '/':
			return 'B'
		case r == '-':
			return 'x'
		case r == '_':
			return 'X'
		default:
			return r
		}
	}, apiKey)

	// Truncate the API key to the desired length
	apiKey = apiKey[:apiKeyLength]

	return apiKey, nil
}

func (sm *SrvManager) RegisterNewSourceServer(name string) (*newSrvSrcResult, error) {
	existingSrv, err := sm.srvRepository.FindSrvWithName(name)
	if err != nil && !errors.Is(err, xerrors.ErrRecordNotFound) {
		log.Default().Println("[Unhandled] error in finding source server with following name", name, err.Error())
		return nil, xerrors.ErrUnhandled
	}

	if existingSrv != nil {
		log.Default().Printf("source server with following name exists! name: '%s'\n", name)
		return nil, xerrors.ErrSourceServerWithThisNameExists
	}

	newSrvAPIKey, err := sm.generateAPIKey()
	if err != nil {
		log.Default().Println("error in creating api-key for register new server", err.Error())
		return nil, xerrors.ErrUnhandled
	}

	hashedBytes := sha256.Sum256([]byte(newSrvAPIKey))
	newSrvHashedAPIKey := hex.EncodeToString(hashedBytes[:])

	newSrvServer, err := sm.srvRepository.CreateNewSrv(name, newSrvHashedAPIKey)
	if err != nil {
		log.Default().Println("error in creating new source server", err.Error())
		return nil, xerrors.ErrUnhandled
	}

	return &newSrvSrcResult{
		APIKey:    newSrvAPIKey,
		NewServer: newSrvServer,
	}, nil
}
