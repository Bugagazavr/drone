package fixtures

import (
	"database/sql"

	"github.com/drone/drone/shared/model"
	"github.com/russross/meddler"
)

func LoadBuilds(db *sql.DB) {
	meddler.Insert(db, "builds", &model.Build{
		Index:     1,
		CommitID:  1,
		Name:      "Build 1",
		Status:    "Success",
		Started:   1398065345,
		Finished:  1398069999,
		Duration:  854,
		AllowFail: false,
		Output:    "sample console output",
		Created:   1398065343,
		Updated:   1398065344,
	})

	meddler.Insert(db, "builds", &model.Build{
		Index:     2,
		CommitID:  1,
		Name:      "Build 2",
		Status:    "Success",
		Started:   1398065345,
		Finished:  1398069999,
		Duration:  854,
		AllowFail: false,
		Output:    "sample console output.....",
		Created:   1398065343,
		Updated:   1398065344,
	})

	meddler.Insert(db, "builds", &model.Build{
		Index:     1,
		CommitID:  2,
		Name:      "Build 1",
		Status:    "Success",
		Started:   1398065345,
		Finished:  1398069999,
		Duration:  854,
		AllowFail: false,
		Output:    "sample console output",
		Created:   1398065343,
		Updated:   1398065344,
	})

	meddler.Insert(db, "builds", &model.Build{
		Index:     1,
		CommitID:  3,
		Name:      "Build 1",
		Status:    "Success",
		Started:   1398065345,
		Finished:  1398069999,
		Duration:  854,
		AllowFail: false,
		Output:    "sample console output",
		Created:   1398065343,
		Updated:   1398065344,
	})

	meddler.Insert(db, "builds", &model.Build{
		Index:     1,
		CommitID:  4,
		Name:      "Build 1",
		Status:    "Success",
		Started:   1398065345,
		Finished:  1398069999,
		Duration:  854,
		AllowFail: false,
		Output:    "sample console output",
		Created:   1398065343,
		Updated:   1398065344,
	})

	meddler.Insert(db, "builds", &model.Build{
		Index:     1,
		CommitID:  5,
		Name:      "Build 1",
		Status:    "Started",
		Started:   1398065345,
		Finished:  1398069999,
		Duration:  854,
		AllowFail: false,
		Output:    "sample console output",
		Created:   1398065343,
		Updated:   1398065344,
	})
}
