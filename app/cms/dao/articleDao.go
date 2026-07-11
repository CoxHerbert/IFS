package dao

import (
	"baize/app/cms/models"
	"baize/app/common/datasource"
	"baize/app/constant/constants"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var articleDaoImpl *articleDao

func init() {
	articleDaoImpl = &articleDao{
		selectSql: `select article_id, title, slug, summary, category, cover_url, content, status, sort,
			ifnull(date_format(publish_time, '%Y-%m-%d %H:%i:%s'), '') as publish_time,
			create_by, create_time, update_by, update_time `,
		fromSql: ` from cms_article`,
	}
}

type articleDao struct {
	selectSql string
	fromSql   string
}

func GetArticleDao() *articleDao {
	return articleDaoImpl
}

func (dao *articleDao) SelectList(query *models.ArticleDQL) (list []*models.ArticleVo, total *int64) {
	whereSql := dao.buildWhere(query)
	countRow, err := datasource.GetMasterDb().NamedQuery(constants.MysqlCount+dao.fromSql+whereSql, query)
	if err != nil { panic(err) }
	total = new(int64)
	if countRow.Next() { countRow.Scan(total) }
	defer countRow.Close()
	list = make([]*models.ArticleVo, 0, query.Size)
	if *total > query.Offset {
		orderSql := " order by sort desc, publish_time desc, create_time desc"
		if query.Limit != "" { orderSql += query.Limit }
		rows, err := datasource.GetMasterDb().NamedQuery(dao.selectSql+dao.fromSql+whereSql+orderSql, query)
		if err != nil { panic(err) }
		for rows.Next() {
			item := new(models.ArticleVo)
			if err := rows.StructScan(item); err != nil { panic(err) }
			list = append(list, item)
		}
		defer rows.Close()
	}
	return
}

func (dao *articleDao) SelectById(articleId int64) *models.ArticleVo {
	item := new(models.ArticleVo)
	err := datasource.GetMasterDb().Get(item, dao.selectSql+dao.fromSql+" where article_id = ?", articleId)
	if err == sql.ErrNoRows { return nil }
	if err != nil { panic(err) }
	return item
}

func (dao *articleDao) SelectBySlug(slug string) *models.ArticleVo {
	item := new(models.ArticleVo)
	err := datasource.GetMasterDb().Get(item, dao.selectSql+dao.fromSql+" where slug = ? and status = '0'", slug)
	if err == sql.ErrNoRows { return nil }
	if err != nil { panic(err) }
	return item
}

func (dao *articleDao) Insert(item *models.ArticleDML) {
	_, err := datasource.GetMasterDb().NamedExec(`insert into cms_article(
		article_id, title, slug, summary, category, cover_url, content, status, sort, publish_time, create_by, create_time, update_by, update_time
	) values (
		:article_id, :title, :slug, :summary, :category, :cover_url, :content, :status, :sort,
		nullif(:publish_time, ''), :create_by, now(), :update_by, now()
	)`, item)
	if err != nil { panic(err) }
}

func (dao *articleDao) Update(item *models.ArticleDML) {
	_, err := datasource.GetMasterDb().NamedExec(`update cms_article set
		title=:title, slug=:slug, summary=:summary, category=:category, cover_url=:cover_url, content=:content,
		status=:status, sort=:sort, publish_time=nullif(:publish_time, ''), update_by=:update_by, update_time=now()
		where article_id=:article_id`, item)
	if err != nil { panic(err) }
}

func (dao *articleDao) DeleteByIds(articleIds []int64) {
	if len(articleIds) == 0 { return }
	query, args, err := sqlx.In("delete from cms_article where article_id in (?)", articleIds)
	if err != nil { panic(err) }
	_, err = datasource.GetMasterDb().Exec(query, args...)
	if err != nil { panic(err) }
}

func (dao *articleDao) buildWhere(query *models.ArticleDQL) string {
	whereSql := ""
	if query.Title != "" { whereSql += " AND title like concat('%', :title, '%')" }
	if query.Category != "" { whereSql += " AND category = :category" }
	if query.Status != "" { whereSql += " AND status = :status" }
	if query.Keyword != "" {
		whereSql += " AND (title like concat('%', :keyword, '%') OR summary like concat('%', :keyword, '%') OR content like concat('%', :keyword, '%'))"
	}
	if whereSql != "" { whereSql = " where " + whereSql[5:] }
	return whereSql
}
