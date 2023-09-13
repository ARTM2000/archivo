package sourceserver

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/ARTM2000/archive1/internal/archive/xerrors"
)

const GlobalFileRotateLimit = 100

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

type FileList struct {
	ID        uint32    `json:"id"`
	FileName  string    `json:"filename"`
	Snapshots int       `json:"snapshots"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SnapshotList struct {
	ID        uint32    `json:"id"`
	Name      string    `json:"name"`
	Size      string    `json:"size"`
	ByteSize  int64     `json:"byte_size"`
	Checksum  string    `json:"checksum"`
	CreatedAt time.Time `json:"created_at"`
}

type StoreManager interface {
	FileStore(srcSrvName, fileName string, file *multipart.FileHeader, correlationId string) error
	FileRotate(srcSrvName, fileName string, rotate int, correlationId string) error
	FileStoreValidate(srcSrvName, fileName string, rotate int) error
	FilesList(srcSrvName string) ([]FileList, error)
	SnapshotsList(srcSrvName, filename string) ([]SnapshotList, error)
	ReadSnapshot(srcSrvName, filename, snapshot string) (*[]byte, error)
}

func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(b)/float64(div), "KMGTPE"[exp])
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

func (sm *SrvManager) SourceServersCount() (int64, error) {
	count, err := sm.srvRepository.CountAllSourceServers()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (sm *SrvManager) SourceServerFilesCount() (int64, error) {
	sourceServers, err := sm.srvRepository.AllSourceServers()
	if err != nil {
		return 0, err
	}

	var totalFiles int64 = 0

	storeManager := sm.getStoreManager()
	for _, ss := range *sourceServers {
		filesList, err := storeManager.FilesList(ss.Name)
		if err != nil {
			log.Default().Printf("error in getting filesList for server %s, error: %+v", ss.Name, err)
			continue
		}
		totalFiles += int64(len(filesList))
	}

	return totalFiles, nil
}

func (sm *SrvManager) TotalSnapshotsSize() (string, error) {
	sourceServers, err := sm.srvRepository.AllSourceServers()
	if err != nil {
		return "", err
	}

	var totalSnapshotsSize int64 = 0
	storeManager := sm.getStoreManager()
	for _, ss := range *sourceServers {
		filesList, err := storeManager.FilesList(ss.Name)
		if err != nil {
			log.Default().Printf("error in getting filesList for server %s, error: %+v", ss.Name, err)
			continue
		}

		for _, fl := range filesList {
			snapshots, err := storeManager.SnapshotsList(ss.Name, fl.FileName)
			if err != nil {
				log.Default().Printf("error in getting snapshots for server '%s' on filename '%s', error: %+v", ss.Name, fl.FileName, err)
				continue
			}
			for _, snp := range snapshots {
				totalSnapshotsSize += snp.ByteSize
			}
		}
	}
	log.Default().Println(totalSnapshotsSize)
	return ByteCountDecimal(totalSnapshotsSize), nil
}

func (sm *SrvManager) GetListOfAllSourceServers(option FindAllOption) (*[]SourceServer, int64, error) {
	servers, total, err := sm.srvRepository.FindAllServers(option)

	if err != nil {
		log.Default().Println("[Unhandled] error in finding all source servers", err.Error())
		return nil, 0, xerrors.ErrUnhandled
	}

	return servers, total, nil
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

	if rotate > GlobalFileRotateLimit {
		log.Default().Printf(
			"error in file store, source server name: '%s' correlationId: '%s', error: %s",
			srcSrv.Name,
			sm.config.CorrelationId,
			xerrors.ErrRotateGlobalLimitReached.Error(),
		)
		return xerrors.ErrRotateGlobalLimitReached
	}

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

func (sm *SrvManager) GetListOfSourceServerFiles(srcSrvId uint, options FindAllOption) (*[]FileList, uint32, error) {
	srv, err := sm.srvRepository.FindSrvWithId(srcSrvId)

	if err != nil {
		if errors.Is(err, xerrors.ErrRecordNotFound) {
			log.Default().Printf("source server with ID '%d' not exists\n", srcSrvId)
			return nil, 0, xerrors.ErrRecordNotFound
		}

		log.Default().Printf("[Unhandled] finding source server with ID '%d' failed, error: %s", srcSrvId, err.Error())
		return nil, 0, xerrors.ErrUnhandled
	}

	storeManager := sm.getStoreManager()
	filesList, err := storeManager.FilesList(srv.Name)
	if err != nil {
		if errors.Is(err, xerrors.ErrNoStoreForSourceServer) {
			return nil, 0, err
		}

		log.Default().Printf("[Unhandled] error in finding files for source server by name '%s', error: %s", srv.Name, err)
		return nil, 0, xerrors.ErrUnhandled
	}

	switch options.SortBy {
	case "id":
		sort.Slice(filesList, func(i, j int) bool {
			if options.SortOrder == "DESC" {
				return filesList[i].ID < filesList[j].ID
			}
			return filesList[i].ID > filesList[j].ID
		})
	case "filename":
		sort.Slice(filesList, func(i, j int) bool {
			if options.SortOrder == "DESC" {
				return filesList[i].FileName < filesList[j].FileName
			}
			return filesList[i].FileName > filesList[j].FileName
		})
	case "snapshots":
		sort.Slice(filesList, func(i, j int) bool {
			if options.SortOrder == "DESC" {
				return filesList[i].Snapshots < filesList[j].Snapshots
			}
			return filesList[i].Snapshots > filesList[j].Snapshots
		})
	case "updated_at":
		sort.Slice(filesList, func(i, j int) bool {
			if options.SortOrder == "DESC" {
				return filesList[j].UpdatedAt.After(filesList[i].UpdatedAt)
			}
			return filesList[i].UpdatedAt.After(filesList[j].UpdatedAt)
		})
	default:
		log.Default().Printf("sortBy not defined, sortBy: '%s'", options.SortBy)
	}

	start := options.Start
	end := options.End
	if end > len(filesList) {
		end = len(filesList)
	}
	finalList := filesList[start:end]

	return &finalList, uint32(len(filesList)), nil
}

func (sm *SrvManager) GetListOfFileSnapshotsByFilenameAndSrvId(srcSrvId uint, filename string, options FindAllOption) (*[]SnapshotList, uint32, error) {
	srv, err := sm.srvRepository.FindSrvWithId(srcSrvId)

	if err != nil {
		if errors.Is(err, xerrors.ErrRecordNotFound) {
			log.Default().Printf("source server with ID '%d' not exists\n", srcSrvId)
			return nil, 0, xerrors.ErrRecordNotFound
		}

		log.Default().Printf("[Unhandled] finding source server with ID '%d' failed, error: %s", srcSrvId, err.Error())
		return nil, 0, xerrors.ErrUnhandled
	}

	storeManager := sm.getStoreManager()
	snapshots, err := storeManager.SnapshotsList(srv.Name, filename)

	if err != nil {
		if errors.Is(err, xerrors.ErrNoStoreForSourceServer) {
			log.Default().Printf("no store found for source server '%s' by id '%d'\n", srv.Name, srv.ID)
			return nil, 0, err
		}

		if errors.Is(err, xerrors.ErrNoFileStoredOnSourceServerByThisName) {
			log.Default().Printf("no store found for source server '%s' by id '%d' for filename '%s'\n", srv.Name, srv.ID, filename)
			return nil, 0, err
		}

		log.Default().Printf("[Unhandled] error in finding file snapshots for source server by name '%s' with filename '%s', error: %s", srv.Name, filename, err)
		return nil, 0, xerrors.ErrUnhandled
	}

	switch options.SortBy {
	case "name":
		sort.Slice(snapshots, func(i, j int) bool {
			if options.SortOrder == "DESC" {
				return snapshots[i].Name < snapshots[j].Name
			}
			return snapshots[i].Name > snapshots[j].Name
		})
	case "size":
		sort.Slice(snapshots, func(i, j int) bool {
			if options.SortOrder == "DESC" {
				return snapshots[i].Size < snapshots[j].Size
			}
			return snapshots[i].Size > snapshots[j].Size
		})
	case "created_at":
		sort.Slice(snapshots, func(i, j int) bool {
			if options.SortOrder == "DESC" {
				return snapshots[j].CreatedAt.After(snapshots[i].CreatedAt)
			}
			return snapshots[i].CreatedAt.After(snapshots[j].CreatedAt)
		})
	default:
		log.Default().Printf("sortBy not defined, sortBy: '%s'", options.SortBy)
	}

	start := options.Start
	end := options.End
	if end > len(snapshots) {
		end = len(snapshots)
	}
	finalList := snapshots[start:end]

	return &finalList, uint32(len(snapshots)), nil
}

func (sm *SrvManager) ReadSnapshot(srcSrvId uint, filename, snapshot string) (*[]byte, string, error) {
	srv, err := sm.srvRepository.FindSrvWithId(srcSrvId)

	if err != nil {
		if errors.Is(err, xerrors.ErrRecordNotFound) {
			log.Default().Printf("source server with ID '%d' not exists\n", srcSrvId)
			return nil, "", xerrors.ErrRecordNotFound
		}

		log.Default().Printf("[Unhandled] finding source server with ID '%d' failed, error: %s", srcSrvId, err.Error())
		return nil, "", xerrors.ErrUnhandled
	}

	storeManager := sm.getStoreManager()
	snapshotByte, err := storeManager.ReadSnapshot(srv.Name, filename, snapshot)

	if err != nil {
		if errors.Is(err, xerrors.ErrSnapshotNotFound) {
			log.Default().Println("snapshot not found")
			return nil, "", err
		}

		return nil, "", err
	}

	ext := strings.Split(http.DetectContentType(*snapshotByte), "/")[1]
	if ext == "octet-stream" {
		// in case that mime type not detected, set extension to "txt"
		ext = "txt"
	}
	log.Default().Println("ext >>", ext)
	var finalName string
	if strings.Contains(ext, "plain") {
		finalName = fmt.Sprintf("%d-%s-%s", srcSrvId, strings.ReplaceAll(filename, ".", ""), snapshot)
	} else {
		finalName = fmt.Sprintf("%d-%s-%s.%s", srcSrvId, strings.ReplaceAll(filename, ".", ""), snapshot, ext)
	}

	return snapshotByte, finalName, nil
}
