package migrationutil

import (
	"github.com/drone/drone/server/helper"
	"github.com/drone/drone/shared/model"
	"github.com/russross/meddler"
)

type Migrate_201410301440 struct{}

var CreateBuildsTable = &Migrate_201410301440{}

func (m *Migrate_201410301440) Revision() int64 {
	return 201410301440
}

func (m *Migrate_201410301440) Up(mg *MigrationDriver) error {
	t := mg.T

	if _, err := mg.CreateTable("builds", []string{
		t.Pk("build_id"),
		t.Integer("commit_id"),
		t.Integer("build_index"),
		t.String("build_name"),
		t.String("build_status"),
		t.Integer("build_started"),
		t.Integer("build_finished"),
		t.Integer("build_duration"),
		t.Bool("build_allow_fail"),
		t.Blob("build_output"),
		t.Integer("build_created"),
		t.Integer("build_updated"),
	}); err != nil {
		return err
	}

	if _, err := mg.AddIndex("builds", []string{"build_id", "commit_id"}, "unique"); err != nil {
		return err
	}
	if _, err := mg.AddIndex("builds", []string{"build_index", "commit_id"}, "unique"); err != nil {
		return err
	}

	tx := mg.Tx

	var dst []*model.Commit
	if err := meddler.QueryAll(tx, &dst, "SELECT * FROM commits"); err != nil {
		return err
	}

	for _, commit := range dst {
		build := &model.Build{}
		build.CommitID = commit.ID
		build.Index = 1
		build.Name = "Build 1"
		build.Duration = commit.Duration
		build.Status = commit.Status
		build.Started = commit.Started
		build.Finished = commit.Finished
		build.Created = commit.Created
		build.Updated = commit.Updated

		var out string
		if err := tx.QueryRow(helper.Rebind("SELECT output_raw FROM output WHERE commit_id = ?"), commit.ID).Scan(&out); err != nil {
			return err
		}

		build.Output = out

		if err := meddler.Insert(tx, "builds", build); err != nil {
			return err
		}
	}

	_, err := mg.DropTable("output")
	return err
}

func (m *Migrate_201410301440) Down(mg *MigrationDriver) error {
	t := mg.T

	if _, err := mg.DropTable("builds"); err != nil {
		return err
	}

	if _, err := mg.CreateTable("output", []string{
		t.Pk("output_id"),
		t.Integer("commit_id"),
		t.Blob("output_raw"),
	}); err != nil {
		return err
	}

	_, err := mg.AddIndex("output", []string{"commit_id"}, "unique")
	return err
}
