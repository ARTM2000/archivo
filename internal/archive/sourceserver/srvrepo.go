package sourceserver

import (
	"errors"
	"log"

	"github.com/ARTM2000/archive1/internal/archive/xerrors"
	"gorm.io/gorm"
)

type SrvSchema struct {
	ID           uint   `gorm:"primaryKey;not null" json:"id"`
	Name         string `gorm:"type:string;not null;unique" json:"name"`
	HashedAPIKey string `gorm:"type:string;not null" json:"-"`
}

func NewSrvRepository(db *gorm.DB) SrvRepository {
	return SrvRepository{
		db: db,
	}
}

type SrvRepository struct {
	db *gorm.DB
}

func (sr *SrvRepository) FindSrvWithName(name string) (*SrvSchema, error) {
	var srv SrvSchema
	dbResult := sr.db.Model(&SrvSchema{}).Where(SrvSchema{Name: name}).First(&srv)

	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			log.Default().Printf("sourceServer with name: '%s' not found\n", name)
			return nil, xerrors.ErrRecordNotFound
		}

		log.Default().Printf("[Unhandled] error in finding source server with name: '%s', error: %s\n", name, dbResult.Error.Error())
		return nil, xerrors.ErrUnhandled
	}

	return &srv, nil
}

func (sr *SrvRepository) CreateNewSrv(name string, hashedAPIKey string) (*SrvSchema, error) {
	var newSrv = SrvSchema{
		Name:         name,
		HashedAPIKey: hashedAPIKey,
	}
	dbResult := sr.db.Model(&SrvSchema{}).Create(&newSrv)

	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrDuplicatedKey) {
			log.Default().Printf("error in creating new source server %+v\n", newSrv)
			return nil, xerrors.ErrDuplicateViolation
		}

		log.Default().Printf("[Unhandled] error in creating new source server. error: %s\n", dbResult.Error.Error())
		return nil, xerrors.ErrUnhandled
	}

	return &newSrv, nil
}
