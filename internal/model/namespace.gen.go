// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameNamespace = "namespace"

// Namespace mapped from table <namespace>
type Namespace struct {
	ID        int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	NsID      string    `gorm:"column:ns_id;not null;default:0" json:"ns_id"`
	NsName    string    `gorm:"column:ns_name;not null" json:"ns_name"`
	ClusterID string    `gorm:"column:cluster_id;not null;default:0" json:"cluster_id"`
	SaID      string    `gorm:"column:sa_id;not null;default:0" json:"sa_id"`
	SaName    string    `gorm:"column:sa_name;not null" json:"sa_name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName Namespace's table name
func (*Namespace) TableName() string {
	return TableNameNamespace
}
