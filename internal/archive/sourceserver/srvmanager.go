package sourceserver

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"
	"mime/multipart"
	"strings"

	"github.com/ARTM2000/archive1/internal/archive/xerrors"
)

func NewSrvManager(config SrvConfig, srvRepo SrvRepository) SrvManager {
	return SrvManager{
		config:        config,
		srvRepository: srvRepo,
	}
}

type newSrvSrcResult struct {
	NewServer *SourceServer
	APIKey    string
}

type SrvConfig struct {
	CorrelationId   string
	StoreMode       string
	DiskStoreConfig DiskStoreConfig
}

type SrvManager struct {
	config        SrvConfig
	srvRepository SrvRepository
}

type StoreManager interface {
	FileStore(srcSrvName string, fileName string, file *multipart.FileHeader, correlationId string) error
	FileRotate(srcSrvName string, fileName string, rotate int, correlationId string) error
	FileStoreValidate(srcSrvName string, fileName string, rotate int) error
}

func (sm *SrvManager) generateAPIKey() (string, error) {
	const apiKeyLength = 64
	apiKeyBytes := make([]byte, apiKeyLength)
	if _, err := rand.Read(apiKeyBytes); err != nil {
		return "", err
	}

	apiKey := base64.RawURLEncoding.EncodeToString(apiKeyBytes)
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

func (sm *SrvManager) AuthorizeSourceServer(srcSrvName string, apiKey string) (*SourceServer, error) {
	srv, err := sm.srvRepository.FindSrvWithName(srcSrvName)
	if err != nil {
		if errors.Is(err, xerrors.ErrRecordNotFound) {
			log.Default().Printf("source server with name '%s' not exists\n", srcSrvName)
			return nil, xerrors.ErrUnauthorized
		}

		log.Default().Printf("[Unhandled] finding source server with name '%s' failed, error: %s", srcSrvName, err.Error())
		return nil, xerrors.ErrUnauthorized
	}

	receivedAPIKeyHashByte := sha256.Sum256([]byte(apiKey))
	receivedAPIKeyHash := hex.EncodeToString(receivedAPIKeyHashByte[:])

	if srv.HashedAPIKey != receivedAPIKeyHash {
		log.Default().Printf(
			"received api key is not valid, receivedHash: '%s' storedHash: '%s'",
			receivedAPIKeyHash,
			srv.HashedAPIKey,
		)
		return nil, xerrors.ErrUnauthorized
	}

	return srv, nil
}

func (sm *SrvManager) getStoreManager() StoreManager {
	switch sm.config.StoreMode {
	case "disk":
		return NewDiskStore(DiskStoreConfig{
			Path: sm.config.DiskStoreConfig.Path,
		})
	default:
		log.Fatalln("store not defined")
		return nil
	}
}

func (sm *SrvManager) RotateFile(srcSrv *SourceServer, rotate int, fileName string, file *multipart.FileHeader) error {
	storeManager := sm.getStoreManager()
	// make sure of final filename
	fnFilename := fileName
	if strings.TrimSpace(fnFilename) == "" {
		fnFilename = file.Filename
	}

	log.Default().Printf(
		"storing file '%s' for source server '%s' with correlationId '%s'\n",
		fnFilename,
		srcSrv.Name,
		sm.config.CorrelationId,
	)

	err := storeManager.FileStoreValidate(srcSrv.Name, fnFilename, rotate)
	if err != nil {
		log.Default().Printf(
			"error in file store, source server name: '%s' correlationId: '%s', error: %s",
			srcSrv.Name,
			sm.config.CorrelationId,
			err.Error(),
		)
		return err
	}

	err = storeManager.FileStore(srcSrv.Name, fnFilename, file, sm.config.CorrelationId)
	if err != nil {
		log.Default().Printf(
			"error in file store, source server name: '%s' correlationId: '%s', error: %s",
			srcSrv.Name,
			sm.config.CorrelationId,
			err.Error(),
		)
		return err
	}

	err = storeManager.FileRotate(srcSrv.Name, fnFilename, rotate, sm.config.CorrelationId)
	if err != nil {
		log.Default().Printf(
			"error in file rotate, source server name: '%s' correlationId: '%s', error: %s",
			srcSrv.Name,
			sm.config.CorrelationId,
			err.Error(),
		)
		return err
	}

	return nil
}
