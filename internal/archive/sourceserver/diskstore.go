package sourceserver

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ARTM2000/archive1/internal/archive/xerrors"
)

func NewDiskStore(config DiskStoreConfig) *DiskStore {
	return &DiskStore{
		Config: config,
	}
}

type DiskStoreConfig struct {
	Path string
}

type DiskStore struct {
	Config DiskStoreConfig
}

func (ds *DiskStore) FileStore(srcSrvName string, fileName string, file *multipart.FileHeader, correlationId string) error {
	// check if directory exist
	storePath := path.Join(ds.Config.Path, srcSrvName, fileName)
	if spData, err := os.Stat(storePath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(storePath, os.ModePerm); err != nil {
				log.Default().Printf(
					"error in creating store path for source server '%s' filename '%s' correlationId '%s'",
					srcSrvName,
					fileName,
					correlationId,
				)
				log.Default().Println(err.Error())
				return xerrors.ErrUnableToCreateStoreDirectory
			}
		} else if !spData.IsDir() {
			log.Default().Printf(
				"error in creating store path for source server '%s' filename '%s' correlationId '%s' store path exists and is not a directory\n",
				srcSrvName,
				fileName,
				correlationId,
			)
			return xerrors.ErrStorePathExistButNotADirectory
		}
	}

	// create new file snapshot name
	now := time.Now()
	fileSnapshotName := strings.Replace(now.Format("20060102150405.000"), ".", "", 1)

	// store file to desire path
	snapshotPath := path.Join(storePath, fileSnapshotName)
	f, err := os.OpenFile(snapshotPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Default().Printf(
			"error in creating store path for source server '%s' filename '%s' correlationId '%s'. error: %s\n",
			srcSrvName,
			fileName,
			correlationId,
			err.Error(),
		)
		return err
	}
	defer f.Close()

	oFile, err := file.Open()
	if err != nil {
		log.Default().Printf(
			"error in creating store path for source server '%s' filename '%s' correlationId '%s'. error: %s\n",
			srcSrvName,
			fileName,
			correlationId,
			err.Error(),
		)
		return err
	}

	_, err = io.Copy(f, oFile)
	if err != nil {
		log.Default().Printf(
			"error in creating store path for source server '%s' filename '%s' correlationId '%s'. error: %s\n",
			srcSrvName,
			fileName,
			correlationId,
			err.Error(),
		)
		return err
	}

	return nil
}

func (ds *DiskStore) FileRotate(srcSrvName string, fileName string, rotate uint64, correlationId string) error {


	return nil
}
