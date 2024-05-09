package scheduler

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jOpeartion "github.com/weichen-lin/kafka-service/neo4j"
)


type Scheduler struct {
	crontabMap map[string]uuid.UUID
}

var scheduler Scheduler

func Init(driver neo4j.DriverWithContext) {
	crontabs, err := neo4jOpeartion.GetAllUserCrontab(driver)
	if err != nil {
		panic(err)
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	for _, crontab := range crontabs {
		j, err := s.NewJob(
			gocron.DailyJob(
				1,
				gocron.NewAtTimes(
					gocron.NewAtTime(uint(crontab.Hour), 0, 0),

				),
			),
			gocron.NewTask(
				func(a string, b int) {
					// do things
				},
				"hello",
				1,
			),
		)
		if err != nil {
			panic(err)
		}

		scheduler.crontabMap[crontab.UserName] = j.ID()
	}
}

