package stdao

import "gorm.io/gorm"

func forkDB(db *gorm.DB) (newDB *gorm.DB, err error) {
	plugins := map[string]gorm.Plugin{}
	for k, v := range db.Config.Plugins {
		plugins[k] = v
	}
	newDB, err = gorm.Open(db.Dialector, &gorm.Config{
		SkipDefaultTransaction:                   db.Config.SkipDefaultTransaction,
		NamingStrategy:                           db.Config.NamingStrategy,
		FullSaveAssociations:                     db.Config.FullSaveAssociations,
		Logger:                                   db.Config.Logger,
		NowFunc:                                  db.Config.NowFunc,
		DryRun:                                   db.Config.DryRun,
		PrepareStmt:                              db.Config.PrepareStmt,
		DisableAutomaticPing:                     db.Config.DisableAutomaticPing,
		DisableForeignKeyConstraintWhenMigrating: db.Config.DisableForeignKeyConstraintWhenMigrating,
		IgnoreRelationshipsWhenMigrating:         db.Config.IgnoreRelationshipsWhenMigrating,
		DisableNestedTransaction:                 db.Config.DisableNestedTransaction,
		AllowGlobalUpdate:                        db.Config.AllowGlobalUpdate,
		QueryFields:                              db.Config.QueryFields,
		CreateBatchSize:                          db.Config.CreateBatchSize,
		ClauseBuilders:                           db.Config.ClauseBuilders,
		ConnPool:                                 db.Config.ConnPool,
		Plugins:                                  plugins,
	})
	return
}
