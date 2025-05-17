package stat

import (
	"gorm.io/datatypes"
	"links-service/pkg/db"
	"time"
)

type StatRepository struct {
	Database *db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		Database: db,
	}
}

func (repo *StatRepository) AddClick(linkId uint) {
	var stat Stat
	dateNow := datatypes.Date(time.Now())
	repo.Database.Find(&stat, "link_id = ? and date = ?", linkId, dateNow)
	if stat.ID == 0 {
		repo.Database.Create(&Stat{
			LinkId: linkId,
			Clicks: 1,
			Date:   dateNow,
		})
	} else {
		stat.Clicks++
		repo.Database.Save(stat)
	}
}
