package main

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type IMetricsDao interface {
	Create(mode *Metrics) error
	GetMetricsDataForEndpoint(route string, fromTime time.Time) ([]Metrics, error)
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

func (dao *MetricsDao) GetMetricsDataForEndpoint(route string, fromTime time.Time) ([]Metrics, error) {
	query := `
        select * from metrics.metrics where route = $1 and createdon >=$2 order by createdon
    `
	metrics := []Metrics{}

	err := dao.db.Select(&metrics, query, route, fromTime)
	if err != nil {
		return nil, err

	}

	return metrics, nil
}
