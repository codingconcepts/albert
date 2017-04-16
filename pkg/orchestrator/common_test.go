package orchestrator

import (
	"time"

	"github.com/codingconcepts/albert/pkg/model"
)

var (
	gatherTimeout  = model.ConfigDuration{Duration: time.Second}
	gatherChanSize = 10
	applications   = Applications{
		Application{
			Name:       "notepad",
			Schedule:   "1 2 3 4 5 6",
			Percentage: 0.75,
		},
	}
	config = `
	{
		"gatherTimeout": "1s",
		"gatherChanSize": 10,
		
		"applications": [
			{
				"name": "notepad",
				"schedule": "1 2 3 4 5 6",
				"percentage": 0.75
			}
		]
	}`
)
