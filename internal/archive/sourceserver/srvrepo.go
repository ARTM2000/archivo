package sourceserver

import (
	"errors"
	"log"
	"time"

	"github.com/ARTM2000/archive1/internal/archive/xerrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SourceServer struct {
	ID           uint      `gorm:"primaryKey;not null" json:"id"`
	Name         string    `gorm:"type:string;not null;unique" json:"name"`
	HashedAPIKey string    `gorm:"type:string;not null" json:"-"`
	CreatedAt    time.Time `gorm:"autoUpdateTime:milli" json:"created_at"`
}

func NewSrvRepository(db *gorm.DB) SrvRepository {
	return SrvRepository{
		db: db,
	}
}

type SrvRepository struct {
	db *gorm.DB
}

func (sr *SrvRepository) FindSrvWithName(name string) (*SourceServer, error) {
	var srv SourceServer
	dbResult := sr.db.Model(&SourceServer{}).Where(SourceServer{Name: name}).First(&srv)

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

func (sr *SrvRepository) FindSrvWithId(id uint) (*SourceServer, error) {
	var srv SourceServer
	dbResult := sr.db.Model(&SourceServer{}).Where(SourceServer{ID: id}).First(&srv)

	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			log.Default().Printf("sourceServer with ID: '%d' not found\n", id)
			return nil, xerrors.ErrRecordNotFound
		}

		log.Default().Printf("[Unhandled] error in finding source server with id: '%d', error: %s\n", id, dbResult.Error.Error())
		return nil, xerrors.ErrUnhandled
	}

	return &srv, nil
}

func (sr *SrvRepository) CreateNewSrv(name string, hashedAPIKey string) (*SourceServer, error) {
	var newSrv = SourceServer{
		Name:         name,
		HashedAPIKey: hashedAPIKey,
	}
	dbResult := sr.db.Model(&SourceServer{}).Create(&newSrv)

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

type FindAllOption struct {
	SortBy    string
	SortOrder string
	Start     int
	End       int
}

func (sr *SrvRepository) FindAllServers(option FindAllOption) (*[]SourceServer, int64, error) {
	var srvs []SourceServer
	var DESC bool
	if option.SortOrder == "ASC" {
		DESC = false
	} else {
		DESC = true
	}
	dbResult := sr.db.Model(&SourceServer{}).Order(clause.OrderByColumn{Column: clause.Column{Name: option.SortBy}, Desc: DESC}).Offset(option.Start).Limit(option.End).Find(&srvs)

	if dbResult.Error != nil {
		log.Default().Printf("[Unhandled] error in finding source servers, error: %s\n", dbResult.Error.Error())
		return nil, 0, xerrors.ErrUnhandled
	}

	var total int64
	dbResult = sr.db.Model(&SourceServer{}).Count(&total)
	if dbResult.Error != nil {
		log.Default().Printf("[Unhandled] error in counting source servers, error: %s\n", dbResult.Error.Error())
		return nil, 0, xerrors.ErrUnhandled
	}

	return &srvs, total, nil
}

func (sr *SrvRepository) CountAllSourceServers() (int64, error) {
	var total int64
	dbResult := sr.db.Model(&SourceServer{}).Count(&total)
	if dbResult.Error != nil {
		log.Default().Println("[Unhandled] error in counting all source servers", dbResult.Error)
		return 0, xerrors.ErrUnhandled
	}

	return total, nil
}

func (sr *SrvRepository) AllSourceServers() (*[]SourceServer, error) {
	var sourceServers []SourceServer
	dbResult := sr.db.Model(&SourceServer{}).Find(&sourceServers)
	if dbResult.Error != nil {
		log.Default().Println("[Unhandled] error in finding all source servers without pagination", dbResult.Error)
		return nil, xerrors.ErrUnhandled
	}

	return &sourceServers, nil
}
