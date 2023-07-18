package sourceserver

import (
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/ARTM2000/archive1/internal/archive/xerrors"
)

type metaData struct {
	Rotate int `json:"rotate"`
}

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

func (ds *DiskStore) FileStoreValidate(srcSrvName string, fileName string, rotate int) error {
	storePath := path.Join(ds.Config.Path, srcSrvName, fileName)
	if _, err := os.Stat(storePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	metaFilePath := path.Join(storePath, ".archive1.meta")
	metaF, err := os.Open(metaFilePath)
	if err != nil {
		log.Default().Println("error in opening meta file, error: ", err.Error())
		return err
	}
	defer metaF.Close()

	metaFileBytes, err := io.ReadAll(metaF)
	if err != nil {
		log.Default().Println("error in reading meta file, error: ", err.Error())
		return err
	}

	var metaData metaData
	if err := json.Unmarshal(metaFileBytes, &metaData); err != nil {
		log.Default().Println("error in parsing meta file, error: ", err.Error())
		return err
	}

	if metaData.Rotate > rotate {
		log.Default().Printf("received rotate file count is lower than previous one and can not be processed")
		return xerrors.ErrFileRotateCountIsLowerThanPreviousOne
	}

	return nil
}

func (ds *DiskStore) FileRotate(srcSrvName string, fileName string, rotate int, correlationId string) error {
	// get files count from file store path

	// the path should exists and then this method to be called.
	// so we prevent store path folder checks
	storePath := path.Join(ds.Config.Path, srcSrvName, fileName)

	ents, err := os.ReadDir(storePath)
	if err != nil {
		log.Default().Println("error in read directory of file store. error: ", err.Error())
		return err
	}

	var fileSnapshotNames []string
	for _, ent := range ents {
		if ent.Name() != ".archive1.meta" {
			fileSnapshotNames = append(fileSnapshotNames, ent.Name())
		}
	}

	// if rotate meta file not found, create it
	metaFilePath := path.Join(storePath, ".archive1.meta")
	metaF, err := os.Create(metaFilePath)
	if err != nil {
		log.Default().Println("error in writing meta file. error: ", err.Error())
		return err
	}
	defer metaF.Close()

	mData := metaData{Rotate: rotate}
	jsonMetaData, _ := json.Marshal(mData)
	_, err = metaF.Write(jsonMetaData)
	if err != nil {
		log.Default().Println("error in write default meta data to file, error: ", err.Error())
		return err
	}

	// sort files by date
	sort.Strings(fileSnapshotNames)
	log.Default().Printf(
		"sorted filename '%s' snapshots for source server '%s' with correlationId '%s'. %+v",
		fileName,
		srcSrvName,
		correlationId,
		fileSnapshotNames,
	)

	// delete required file(s)
	if len(fileSnapshotNames) < rotate {
		log.Default().Printf(
			"sorted filename '%s' snapshots for source server '%s' with correlationId '%s'. no need to file rotate",
			fileName,
			srcSrvName,
			correlationId,
		)
		return nil
	}
	filesForDelete := fileSnapshotNames[:(len(fileSnapshotNames) - rotate)]
	log.Default().Printf(
		"filename '%s' snapshots for source server '%s' with correlationId '%s'. going to delete files: %+v",
		fileName,
		srcSrvName,
		correlationId,
		filesForDelete,
	)

	for _, ffdName := range filesForDelete {
		err := os.Remove(path.Join(storePath, ffdName))
		if err != nil {
			log.Default().Printf(
				"error in deleting file for rotation, filename: '%s', source server: '%s', correlationId: '%s', snapshotName: '%s', error: %s",
				fileName,
				srcSrvName,
				correlationId,
				ffdName,
				err.Error(),
			)
			return err
		}
	}

	log.Default().Printf(
		"file rotation completed. source server: '%s', filename: '%s', correlationId: '%s', rotate: '%d'",
		srcSrvName,
		fileName,
		correlationId,
		rotate,
	)

	return nil
}

func (ds *DiskStore) FilesList(srcSrvName string) ([]FileList, error) {

	// check that is there any directory for requested source server or not
	srcSrvStorePath := path.Join(ds.Config.Path, srcSrvName)
	ents, err := os.ReadDir(srcSrvStorePath)

	if err != nil {
		log.Default().Println("error in reading source server store directory, error:", err.Error())
		if os.IsNotExist(err) {
			return nil, xerrors.ErrNoStoreForSourceServer
		}
		return nil, xerrors.ErrUnhandled
	}

	var filenamesList []string
	for _, ent := range ents {
		if ent.IsDir() {
			filenamesList = append(filenamesList, ent.Name())
		}
	}

	// sort files by their
	sort.Strings(filenamesList)
	var filesList []FileList
	for i, filename := range filenamesList {
		dirName := path.Join(srcSrvStorePath, filename)
		fInfo, _ := os.Stat(dirName)
		snapshots, _ := os.ReadDir(dirName)

		// (-1) is to exclude `.archive1.meta` metadata file
		fileNameSnapshotCounts := len(snapshots) - 1

		filesList = append(filesList, FileList{
			ID:        uint32(i + 1),
			FileName:  filename,
			Snapshots: fileNameSnapshotCounts,
			UpdatedAt: fInfo.ModTime(),
		})
	}

	return filesList, nil
}
