package label

import (
	"github.com/bryant-rh/cm/cmd/server/global"
	"github.com/bryant-rh/cm/internal/model"
	"github.com/bryant-rh/cm/internal/query"
	"github.com/bryant-rh/cm/pkg/util"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func LabelRouter(r *gin.RouterGroup) {
	r.GET("/label/list", ListClusterLabel)
	r.POST("/label/create", CreateLabel)
	r.DELETE("/label/delete", DeleteLabel)
}

type Cluster struct {
	Name      string `json:"name"`
	ClusterID int64  `json:"cluster_id"`
}

type CluseterLabel struct {
	ProjectName string `json:"project_name"`
	ClusterName string `json:"cluster_name"`
	Labels      string `json:"Labels"`
}

// @BasePath /api/v1
// PingCluster godoc
// @Summary ListClusterLabel
// @Schemes
// @Description List Cluster's label
// @Tags ListClusterLabel
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string true "Project Name"
// @Param cluster_name query string true "Cluster Name"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /label/list [get]
// @ID ListClusterLabel
func ListClusterLabel(ctx *gin.Context) {
	project_name := ctx.Query("project_name")
	cluster_name := ctx.Query("cluster_name")
	q := query.Use(global.Config.DB.DB())
	//判断项目和集群是否存在
	cluster, err := q.ClusterBind.WithContext(ctx).Select(q.ClusterBind.ClusterID).Where(q.ClusterBind.ProjectName.Eq(project_name), q.ClusterBind.ClusterName.Eq(cluster_name)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is not exists! in ProjectName: [%s]", cluster_name, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}
	results := []CluseterLabel{}
	//global.Config.DB.Raw("SELECT  cluster.cluster_name as ClusterName, (select GROUP_CONCAT(label.key, '=',label.value) FROM  label )AS Labels FROM cluster WHERE cluster.cluster_id = ? GROUP BY ID", clusterId.ClusterID).Scan(&results)
	global.Config.DB.DB().
		Raw("SELECT ? as 'ProjectName', ? as 'ClusterName', GROUP_CONCAT(label.key,'=' ,label.value) as 'Labels' FROM  label where FIND_IN_SET (label.label_id,(SELECT cluster.labels FROM cluster WHERE cluster.cluster_id = ?))", project_name, cluster_name, cluster.ClusterID).
		Scan(&results)

	util.ReturnMsg(ctx, http.StatusOK, results, "success")

}

type CreateLabelRequestBody struct {
	ProjectName string `json:"project_name" binding:"required"`
	ClusterName string `json:"cluster_name" binding:"required"`
	LabelKey    string `json:"label_key" binding:"required"`
	LabelValue  string `json:"label_value" binding:"required"`
}
type UpdateClusterRequestBody struct {
	ClusterName string `json:"cluster_name" binding:"required"`
	Description string `json:"description"`
}

type ErrorResp struct {
	Msg string `json:"msg"`
}

// @BasePath /api/v1
// PingProject godoc
// @Summary CreateLabel
// @Schemes
// @Description Create Label
// @Tags CreateLabel
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ReqeustBody body CreateLabelRequestBody true "Create Label"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /label/create [post]
// @ID CreateLabel
func CreateLabel(ctx *gin.Context) {
	body := CreateLabelRequestBody{}
	err := ctx.ShouldBind(&body)
	if err != nil {
		error_msg := fmt.Sprintf("err: %s ", err)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return
	}

	q := query.Use(global.Config.DB.DB())
	//判断项目和集群是否存在
	clusterId, err := q.ClusterBind.WithContext(ctx).Select(q.ClusterBind.ClusterID).Where(q.ClusterBind.ProjectName.Eq(body.ProjectName), q.ClusterBind.ClusterName.Eq(body.ClusterName)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is not exists! in ProjectName: [%s]", body.ClusterName, body.ProjectName)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}
	//判断对应项目的集群下是否已存在此标签

	clusterLabel := ""
	res, err := q.Cluster.WithContext(ctx).Select(q.Cluster.Labels).Where(q.Cluster.ClusterID.Eq(clusterId.ClusterID)).First()
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		clusterLabel = res.Labels
	}

	_, err = q.Label.WithContext(ctx).Where(q.Label.LabelID.FindInSet(clusterLabel), q.Label.Key.Eq(body.LabelKey), q.Label.Value.Eq(body.LabelValue)).First()
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is already exists Label: [%s:%s] in ProjectName: [%s]", body.ClusterName, body.LabelKey, body.LabelValue, body.ProjectName)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return
	}

	//创建集群
	id, err := util.NewIdMgr(5)
	if err != nil {
		fmt.Println(err)
		return
	}

	label := model.Label{}
	label_id := util.Int64toString(id.ID())
	label.Key = body.LabelKey
	label.Value = body.LabelValue
	label.LabelID = label_id

	cluster := model.Cluster{}

	labelSlice := util.StringtoSlice(clusterLabel)
	cluster.Labels = strings.Join(append(labelSlice, label_id), ",")

	q.Transaction(func(tx *query.Query) error {
		if err := tx.Label.WithContext(ctx).Create(&label); err != nil {
			return err
		}
		if _, err = tx.Cluster.WithContext(ctx).Where(tx.Cluster.ClusterID.Eq(clusterId.ClusterID)).Updates(cluster); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		error_msg := fmt.Sprintf("The ClusterName: [%s] Create Label : [%s:%s] Failed! in ProjectName : [%s], err: %s", body.ClusterName, body.LabelKey, body.LabelValue, body.ProjectName, err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}

	results := []CluseterLabel{}
	//global.Config.DB.Raw("SELECT  cluster.cluster_name as ClusterName, (select GROUP_CONCAT(label.key, '=',label.value) FROM  label )AS Labels FROM cluster WHERE cluster.cluster_id = ? GROUP BY ID", clusterId.ClusterID).Scan(&results)
	global.Config.DB.DB().
		Raw("SELECT ? as 'ProjectName', ? as 'ClusterName', GROUP_CONCAT(label.key,'=' ,label.value) as 'Labels' FROM  label where FIND_IN_SET (label.label_id,(SELECT cluster.labels FROM cluster WHERE cluster.cluster_id = ?))", body.ProjectName, body.ClusterName, clusterId.ClusterID).
		Scan(&results)

	ok_msg := fmt.Sprintf("The Cluster_Name: [%s] Create Label : [%s:%s]  successfully! in ProjectName : [%s]", body.ClusterName, body.LabelKey, body.LabelValue, body.ProjectName)
	util.ReturnMsg(ctx, http.StatusOK, results, ok_msg)
	//util.ReturnMsg(ctx, http.StatusOK, label, ok_msg)
}

// @BasePath /api/v1
// PingCluster godoc
// @Summary DeleteLabel
// @Schemes
// @Description Delete Label
// @Tags DeleteLabel
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string true "Project_Name"
// @Param cluster_name query string true "Cluster_Name"
// @Param label_key query string true "Label_Key"
// @Param label_value query string true "Label_value"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /label/delete [delete]
// @ID DeleteLabel
func DeleteLabel(ctx *gin.Context) {
	project_name := ctx.Query("project_name")
	cluster_name := ctx.Query("cluster_name")
	label_key := ctx.Query("label_key")
	label_value := ctx.Query("label_value")

	q := query.Use(global.Config.DB.DB())
	//判断项目和集群是否存在
	clusterId, err := q.ClusterBind.WithContext(ctx).Select(q.ClusterBind.ClusterID).Where(q.ClusterBind.ProjectName.Eq(project_name), q.ClusterBind.ClusterName.Eq(cluster_name)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is not exists! in ProjectName: [%s]", cluster_name, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)

		return
	}
	//判断对应项目的集群下是否已存在此标签
	clusterLabel := ""
	res, err := q.Cluster.WithContext(ctx).Select(q.Cluster.Labels).Where(q.Cluster.ClusterID.Eq(clusterId.ClusterID)).First()
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		clusterLabel = res.Labels
	}

	info, err := q.Label.WithContext(ctx).Where(q.Label.LabelID.FindInSet(clusterLabel), q.Label.Key.Eq(label_key), q.Label.Value.Eq(label_value)).First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is Not Label: [%s:%s] in ProjectName: [%s]", cluster_name, label_key, label_value, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}

	labelSlice := util.StringtoSlice(clusterLabel)

	cluster := model.Cluster{}
	//删除标签
	cluster.Labels = strings.Join(util.SliceDelOne(info.LabelID, labelSlice), ",")
	//更新数据

	res_update := gen.ResultInfo{}
	q.Transaction(func(tx *query.Query) error {

		//if err := tx.Project.WithContext(ctx).Create(&project); err != nil {
		if res_update, err = tx.Label.WithContext(ctx).Where(tx.Label.LabelID.Eq(info.LabelID)).Delete(); err != nil {
			return err
		}
		if res_update, err = tx.Cluster.WithContext(ctx).Where(q.Cluster.ClusterID.Eq(clusterId.ClusterID)).Updates(cluster); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		error_msg := fmt.Sprintf("The ClusterName: [%s] Delete Label : [%s:%s] Failed! in ProjectName : [%s], err: %s", cluster_name, label_key, label_value, project_name, err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}

	ok_msg := fmt.Sprintf("The Cluster_Name: [%s] Delete Label : [%s:%s]  successfully! in ProjectName : [%s]", cluster_name, label_key, label_value, project_name)
	util.ReturnMsg(ctx, http.StatusOK, util.Int64toString(res_update.RowsAffected), ok_msg)

}
