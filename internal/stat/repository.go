package stat

import (
	"gorm.io/datatypes"
	"links-service/pkg/db"
	req "links-service/pkg/request"
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
func (repo *StatRepository) GetStats(params req.StatParams) []GetStatResponse {
	var stats []GetStatResponse
	var selectQuery string
	switch params.By {
	case req.GroupByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case req.GroupByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	}

	repo.Database.Table("stats").
		Select(selectQuery).
		Where("date BETWEEN ? and ?", params.From, params.To).
		Group("period").
		Order("period").
		Find(&stats)

	return stats
}
