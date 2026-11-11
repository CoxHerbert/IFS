package dao

import (
	"baize/app/cms/models"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestArticleSelectSQLCanBeNamedBound(t *testing.T) {
	query := &models.ArticleDQL{}
	sqlText := articleDaoImpl.selectSql + articleDaoImpl.fromSql + articleDaoImpl.buildWhere(query)

	if _, _, err := sqlx.Named(sqlText, query); err != nil {
		t.Fatalf("article select SQL cannot be named bound: %v", err)
	}
}
