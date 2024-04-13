package job

import (
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"math/rand"

	"github.com/robfig/cron"
	"gorm.io/gorm"
)

// timeouts example: @exmple
const (
	MAX_LOVE_METER = 20
	MIN_LOVE_METER = 0

	MIN_SATIETY = 0

	HUNGER_TIMEOUT = "@every 10s"
	LOVE_TIMEOUT   = "@every 10s"
	SEX_TIMEOUT    = "@every 5s"
)

type Job struct {
	c  *cron.Cron
	db *gorm.DB
}

func NewJob(c *cron.Cron, db *gorm.DB) *Job {
	return &Job{
		c:  c,
		db: db,
	}
}

func (j *Job) Run() {
	j.petJobs()

	j.c.Start()
}

func (j *Job) petJobs() {
	// hunger func
	j.c.AddFunc(HUNGER_TIMEOUT, func() {
		var pets []domain.Pet

		if err := j.db.Find(&pets).Error; err != nil {
			return
		}

		for _, pet := range pets {
			if pet.Satiety <= 0 {
				if err := j.db.Where("id = ?", pet.ID).Delete(&domain.Pet{}).Commit().Error; err != nil {
					return
				}
			}

			pet.Satiety = pet.Satiety - 5
			if err := j.db.Save(&pet).Error; err != nil {
				return
			}
		}
	})

	// love meter func
	j.c.AddFunc(LOVE_TIMEOUT, func() {
		var pets []domain.Pet

		if err := j.db.Find(&pets).Error; err != nil {
			return
		}

		for _, pet := range pets {
			if pet.LoveMeter == MAX_LOVE_METER {
				continue
			}

			if pet.Satiety > 20 {
				pet.LoveMeter = pet.LoveMeter + 0.5
			}

			if err := j.db.Save(&pet).Error; err != nil {
				return
			}
		}
	})

	// sex func
	j.c.AddFunc(SEX_TIMEOUT, func() {
		var pets []domain.Pet

		if err := j.db.Find(&pets).Error; err != nil {
			return
		}

		var petMale, petFemale domain.Pet

		if err := j.db.Where("sex = ? and love_meter = ?", types.SEX_MALE, MAX_LOVE_METER).Find(&petMale).Error; err != nil {
			return
		}

		if err := j.db.Where("sex = ? and love_meter = ?", types.SEX_FEMALE, MAX_LOVE_METER).Find(&petFemale).Error; err != nil {
			return
		}

		if petMale.Type == petFemale.Type && petMale.LoveMeter == MAX_LOVE_METER && petFemale.LoveMeter == MAX_LOVE_METER {
			if err := j.db.Create(&domain.Pet{
				Type:      petMale.Type,
				Sex:       randomSex(),
				Satiety:   100,
				LoveMeter: 0,
				Cost:      petMale.Cost,
			}).Error; err != nil {
				return
			}

			petMale.LoveMeter = MIN_LOVE_METER
			petFemale.LoveMeter = MIN_LOVE_METER

			if err := j.db.Save(&petMale).Error; err != nil {
				return
			}

			if err := j.db.Save(&petFemale).Error; err != nil {
				return
			}
		}
	})
}

func randomSex() string {
	return []string{"male", "female"}[rand.Intn(2)]
}
