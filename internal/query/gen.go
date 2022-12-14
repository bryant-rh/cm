// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

func Use(db *gorm.DB) *Query {
	return &Query{
		db:             db,
		Cluster:        newCluster(db),
		ClusterBind:    newClusterBind(db),
		Label:          newLabel(db),
		Namespace:      newNamespace(db),
		Project:        newProject(db),
		Serviceaccount: newServiceaccount(db),
		User:           newUser(db),
	}
}

type Query struct {
	db *gorm.DB

	Cluster        cluster
	ClusterBind    clusterBind
	Label          label
	Namespace      namespace
	Project        project
	Serviceaccount serviceaccount
	User           user
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:             db,
		Cluster:        q.Cluster.clone(db),
		ClusterBind:    q.ClusterBind.clone(db),
		Label:          q.Label.clone(db),
		Namespace:      q.Namespace.clone(db),
		Project:        q.Project.clone(db),
		Serviceaccount: q.Serviceaccount.clone(db),
		User:           q.User.clone(db),
	}
}

type queryCtx struct {
	Cluster        *clusterDo
	ClusterBind    *clusterBindDo
	Label          *labelDo
	Namespace      *namespaceDo
	Project        *projectDo
	Serviceaccount *serviceaccountDo
	User           *userDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Cluster:        q.Cluster.WithContext(ctx),
		ClusterBind:    q.ClusterBind.WithContext(ctx),
		Label:          q.Label.WithContext(ctx),
		Namespace:      q.Namespace.WithContext(ctx),
		Project:        q.Project.WithContext(ctx),
		Serviceaccount: q.Serviceaccount.WithContext(ctx),
		User:           q.User.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	return &QueryTx{q.clone(q.db.Begin(opts...))}
}

type QueryTx struct{ *Query }

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
