package job

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"math/rand"
	"time"

	"github.com/robfig/cron"
	"gorm.io/gorm"
)

// timeouts example: @exmple
const (
	MAX_LOVE_METER = 20
	MIN_LOVE_METER = 0

	MIN_SATIETY    = 0
	HUNGRY_SATIETY = 20

	HUNGER_TIMEOUT = "@every 10s"
	LOVE_TIMEOUT   = "@every 10s"
	SEX_TIMEOUT    = "@every 5s"

	SEX_MALE   = "male"
	SEX_FEMALE = "female"
)

type Job struct {
	c       *cron.Cron
	db      *gorm.DB
	wsConns map[uint][]*websocket.Conn
}

type webSocketLoveData struct {
	Male   domain.Pet `json:"male"`
	Female domain.Pet `json:"female"`
	Child  domain.Pet `json:"child"`
}

func NewJob(c *cron.Cron, db *gorm.DB, wsConns map[uint][]*websocket.Conn) *Job {
	return &Job{
		c:       c,
		db:      db,
		wsConns: wsConns,
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

				var child domain.User
				if err := j.db.First(&child, "id = ?", pet.UserID).Error; err != nil {
					return
				}
				var parent domain.User
				if err := j.db.First(&parent, "email = ?", child.Email).Error; err != nil {
					return
				}

				j.broadcastStructToUserById(child.ID, types.WebSocketMessage{
					Event: "pet.death",
					Data:  pet,
				})
				j.broadcastStructToUserById(parent.ID, types.WebSocketMessage{
					Event: "pet.death",
					Data:  pet,
				})
			}

			if pet.Satiety <= HUNGRY_SATIETY {
				j.broadcastStructToUserById(pet.UserID, types.WebSocketMessage{
					Event: "pet.hungry",
					Data:  pet,
				})
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

		var users []domain.User
		if err := j.db.Find(&users); err != nil {
			return
		}
		for _, user := range users {
			var petMale, petFemale domain.Pet

			if err := j.db.Where("sex = ? and love_meter = ? and user_id = ?", SEX_MALE, MAX_LOVE_METER, user.ID).Find(&petMale).Error; err != nil {
				return
			}

			if err := j.db.Where("sex = ? and love_meter = ? and user_id = ?", SEX_FEMALE, MAX_LOVE_METER, user.ID).Find(&petFemale).Error; err != nil {
				return
			}

			if petMale.Type == petFemale.Type && petMale.LoveMeter == MAX_LOVE_METER && petFemale.LoveMeter == MAX_LOVE_METER {
				child := domain.Pet{
					Type:      petMale.Type,
					Sex:       randomSex(),
					Satiety:   100,
					LoveMeter: 0,
					Cost:      petMale.Cost,
				}
				if err := j.db.Create(&child).Error; err != nil {
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

				j.broadcastStructToUserById(user.ID, types.WebSocketMessage{
					Event: "pet.love",
					Data: webSocketLoveData{
						Male:   petMale,
						Female: petFemale,
						Child:  child,
					},
				})
			}
		}
	})
}

func (j *Job) broadcastStructToUserById(id uint, msg interface{}) {
	rawM, err := json.Marshal(msg)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to parse json for websocket"))
		return
	}
	for _, ws := range j.wsConns[id] {
		ws.WriteMessage(websocket.TextMessage, rawM)
	}
	if wsMsg, ok := msg.(types.WebSocketMessage); ok {
		data, err := json.Marshal(wsMsg.Data)

		if err != nil {
			logrus.Error(errors.Wrap(err, "failed to parse json for websocket"))
			return
		}
		dbMsg := domain.Message{
			CreatedAt: time.Now(),
			Event:     wsMsg.Event,
			Data:      string(data),
			IsRead:    false,
			UserID:    id,
		}
		if err := j.db.Save(&dbMsg).Error; err != nil {
			logrus.Error(errors.Wrap(err, "failed to save messages"))
			return
		}
	}
}

func randomSex() string {
	return []string{"male", "female"}[rand.Intn(2)]
}
