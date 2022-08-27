// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/bryant-rh/cm/internal/model"
)

func newServiceaccount(db *gorm.DB) serviceaccount {
	_serviceaccount := serviceaccount{}

	_serviceaccount.serviceaccountDo.UseDB(db)
	_serviceaccount.serviceaccountDo.UseModel(&model.Serviceaccount{})

	tableName := _serviceaccount.serviceaccountDo.TableName()
	_serviceaccount.ALL = field.NewField(tableName, "*")
	_serviceaccount.ID = field.NewInt32(tableName, "id")
	_serviceaccount.SaID = field.NewString(tableName, "sa_id")
	_serviceaccount.SaName = field.NewString(tableName, "sa_name")
	_serviceaccount.SaToken = field.NewString(tableName, "sa_token")
	_serviceaccount.CreatedAt = field.NewTime(tableName, "created_at")
	_serviceaccount.UpdatedAt = field.NewTime(tableName, "updated_at")

	_serviceaccount.fillFieldMap()

	return _serviceaccount
}

type serviceaccount struct {
	serviceaccountDo serviceaccountDo

	ALL       field.Field
	ID        field.Int32
	SaID      field.String
	SaName    field.String
	SaToken   field.String
	CreatedAt field.Time
	UpdatedAt field.Time

	fieldMap map[string]field.Expr
}

func (s serviceaccount) Table(newTableName string) *serviceaccount {
	s.serviceaccountDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s serviceaccount) As(alias string) *serviceaccount {
	s.serviceaccountDo.DO = *(s.serviceaccountDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *serviceaccount) updateTableName(table string) *serviceaccount {
	s.ALL = field.NewField(table, "*")
	s.ID = field.NewInt32(table, "id")
	s.SaID = field.NewString(table, "sa_id")
	s.SaName = field.NewString(table, "sa_name")
	s.SaToken = field.NewString(table, "sa_token")
	s.CreatedAt = field.NewTime(table, "created_at")
	s.UpdatedAt = field.NewTime(table, "updated_at")

	s.fillFieldMap()

	return s
}

func (s *serviceaccount) WithContext(ctx context.Context) *serviceaccountDo {
	return s.serviceaccountDo.WithContext(ctx)
}

func (s serviceaccount) TableName() string { return s.serviceaccountDo.TableName() }

func (s serviceaccount) Alias() string { return s.serviceaccountDo.Alias() }

func (s *serviceaccount) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *serviceaccount) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 6)
	s.fieldMap["id"] = s.ID
	s.fieldMap["sa_id"] = s.SaID
	s.fieldMap["sa_name"] = s.SaName
	s.fieldMap["sa_token"] = s.SaToken
	s.fieldMap["created_at"] = s.CreatedAt
	s.fieldMap["updated_at"] = s.UpdatedAt
}

func (s serviceaccount) clone(db *gorm.DB) serviceaccount {
	s.serviceaccountDo.ReplaceDB(db)
	return s
}

type serviceaccountDo struct{ gen.DO }

func (s serviceaccountDo) Debug() *serviceaccountDo {
	return s.withDO(s.DO.Debug())
}

func (s serviceaccountDo) WithContext(ctx context.Context) *serviceaccountDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s serviceaccountDo) ReadDB() *serviceaccountDo {
	return s.Clauses(dbresolver.Read)
}

func (s serviceaccountDo) WriteDB() *serviceaccountDo {
	return s.Clauses(dbresolver.Write)
}

func (s serviceaccountDo) Clauses(conds ...clause.Expression) *serviceaccountDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s serviceaccountDo) Returning(value interface{}, columns ...string) *serviceaccountDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s serviceaccountDo) Not(conds ...gen.Condition) *serviceaccountDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s serviceaccountDo) Or(conds ...gen.Condition) *serviceaccountDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s serviceaccountDo) Select(conds ...field.Expr) *serviceaccountDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s serviceaccountDo) Where(conds ...gen.Condition) *serviceaccountDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s serviceaccountDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *serviceaccountDo {
	return s.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (s serviceaccountDo) Order(conds ...field.Expr) *serviceaccountDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s serviceaccountDo) Distinct(cols ...field.Expr) *serviceaccountDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s serviceaccountDo) Omit(cols ...field.Expr) *serviceaccountDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s serviceaccountDo) Join(table schema.Tabler, on ...field.Expr) *serviceaccountDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s serviceaccountDo) LeftJoin(table schema.Tabler, on ...field.Expr) *serviceaccountDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s serviceaccountDo) RightJoin(table schema.Tabler, on ...field.Expr) *serviceaccountDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s serviceaccountDo) Group(cols ...field.Expr) *serviceaccountDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s serviceaccountDo) Having(conds ...gen.Condition) *serviceaccountDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s serviceaccountDo) Limit(limit int) *serviceaccountDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s serviceaccountDo) Offset(offset int) *serviceaccountDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s serviceaccountDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *serviceaccountDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s serviceaccountDo) Unscoped() *serviceaccountDo {
	return s.withDO(s.DO.Unscoped())
}

func (s serviceaccountDo) Create(values ...*model.Serviceaccount) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s serviceaccountDo) CreateInBatches(values []*model.Serviceaccount, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s serviceaccountDo) Save(values ...*model.Serviceaccount) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s serviceaccountDo) First() (*model.Serviceaccount, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Serviceaccount), nil
	}
}

func (s serviceaccountDo) Take() (*model.Serviceaccount, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Serviceaccount), nil
	}
}

func (s serviceaccountDo) Last() (*model.Serviceaccount, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Serviceaccount), nil
	}
}

func (s serviceaccountDo) Find() ([]*model.Serviceaccount, error) {
	result, err := s.DO.Find()
	return result.([]*model.Serviceaccount), err
}

func (s serviceaccountDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Serviceaccount, err error) {
	buf := make([]*model.Serviceaccount, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s serviceaccountDo) FindInBatches(result *[]*model.Serviceaccount, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s serviceaccountDo) Attrs(attrs ...field.AssignExpr) *serviceaccountDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s serviceaccountDo) Assign(attrs ...field.AssignExpr) *serviceaccountDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s serviceaccountDo) Joins(fields ...field.RelationField) *serviceaccountDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s serviceaccountDo) Preload(fields ...field.RelationField) *serviceaccountDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s serviceaccountDo) FirstOrInit() (*model.Serviceaccount, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Serviceaccount), nil
	}
}

func (s serviceaccountDo) FirstOrCreate() (*model.Serviceaccount, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Serviceaccount), nil
	}
}

func (s serviceaccountDo) FindByPage(offset int, limit int) (result []*model.Serviceaccount, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s serviceaccountDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s serviceaccountDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s *serviceaccountDo) withDO(do gen.Dao) *serviceaccountDo {
	s.DO = *do.(*gen.DO)
	return s
}
