package main

import "github.com/jmoiron/sqlx"

type IMetricsDao interface {
	Create(mode *Metrics) error
}

type MetricsDao struct {
	db *sqlx.DB
}

func newMetricsDao(db *sqlx.DB) IMetricsDao {
	return &MetricsDao{db: db}
}

func (dao *MetricsDao) Create(model *Metrics) error {
	query := `insert into metrics.metrics( id , route , method , responsetime , statuscode , createdon ) 
    values (  :id , :route , :method , :responsetime , :statuscode , :createdon)
    `
	_, err := dao.db.NamedExec(query, model)
	return err
}
