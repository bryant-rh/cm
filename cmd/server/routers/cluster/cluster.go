package cluster

import (
	"github.com/bryant-rh/cm/cmd/server/global"
	"github.com/bryant-rh/cm/internal/model"
	"github.com/bryant-rh/cm/internal/query"
	"github.com/bryant-rh/cm/pkg/util"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func ClusterRouter(r *gin.RouterGroup) {
	r.GET("/cluster/list", ListCluster)
	r.GET("/cluster/get_id", GetClusterId)
	r.GET("/cluster/label", ListGetLabel)
	r.POST("/cluster/create", CreateCluster)
	r.PUT("/cluster/update", UpdateCluster)
	r.DELETE("/cluster/delete", DeleteCluster)
}

type Cluster struct {
	Name      string `json:"name"`
	ClusterID int64  `json:"cluster_id"`
}

// @BasePath /api/v1
// PingCluster godoc
// @Summary ListCluster
// @Schemes
// @Description List  Cluster
// @Tags ListCluster
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string false "Project Name"
// @Param cluster_name query string false "Cluster Name"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /cluster/list [get]
// @ID ListCluster
func ListCluster(ctx *gin.Context) {
	project_name := ctx.Query("project_name")
	cluster_name := ctx.Query("cluster_name")

	cluster := []model.Cluster{}
	//q := query.Use(global.Config.DB.DB())

	if project_name == "" && cluster_name == "" {
		// err := q.Cluster.WithContext(ctx).Scan(&cluster)
		// if err != nil {
		// 	error_msg := fmt.Sprintf("err: %s ", err)
		// 	util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		// 	return
		// }
		global.Config.DB.DB().
			Raw("SELECT  cluster.id, cluster.cluster_id,cluster.cluster_name,cluster.`describe` ,(select GROUP_CONCAT(label.key, '=',label.value)  FROM  label  where FIND_IN_SET (label.label_id,cluster.labels) )AS labels ,cluster.created_at, cluster.updated_at FROM cluster").
			Scan(&cluster)

	} else if project_name != "" && cluster_name == "" {
		global.Config.DB.DB().
			Raw("SELECT  cluster.id, cluster.cluster_id,cluster.cluster_name,cluster.`describe` ,(select GROUP_CONCAT(label.key, '=',label.value)  FROM  label  where FIND_IN_SET (label.label_id,cluster.labels) )AS labels ,cluster.created_at, cluster.updated_at FROM cluster where cluster.cluster_id IN (SELECT cluster_bind.cluster_id  from cluster_bind WHERE cluster_bind.project_name = ?)", project_name).
			Scan(&cluster)

		// err := q.Cluster.WithContext(ctx).Where(q.Cluster.WithContext(ctx).Columns(q.Cluster.ClusterID).In(q.ClusterBind.WithContext(ctx).Select(q.ClusterBind.ClusterID).Where(q.ClusterBind.ProjectName.Eq(project_name)))).Scan(&cluster)
		// if err != nil {
		// 	error_msg := fmt.Sprintf("err: %s ", err)
		// 	util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		// 	return
		// }

	} else if cluster_name != "" {
		q := query.Use(global.Config.DB.DB()).ClusterBind
		ClusterList, err := q.WithContext(ctx).Select(q.ClusterID).Where(q.ProjectName.Eq(project_name), q.ClusterName.Eq(cluster_name)).First()

		if err != nil {
			if ClusterList == nil {
				error_msg := fmt.Sprintf("The ClusterName: [%s] is not Found! in ProjectName: [%s]", cluster_name, project_name)
				util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
			}
			return
			//ctx.JSON(http.StatusInternalServerError, err)
		}
		global.Config.DB.DB().
			Raw("SELECT  cluster.id, cluster.cluster_id,cluster.cluster_name,cluster.`describe` ,(select GROUP_CONCAT(label.key, '=',label.value)  FROM  label  where FIND_IN_SET (label.label_id,cluster.labels) )AS labels ,cluster.created_at, cluster.updated_at FROM cluster where cluster.cluster_id = (SELECT cluster_bind.cluster_id  from cluster_bind WHERE cluster_bind.project_name = ? and cluster_bind.cluster_name = ?)", project_name, cluster_name).
			Scan(&cluster)

		// err := q.Cluster.WithContext(ctx).Where(q.Cluster.WithContext(ctx).Columns(q.Cluster.ClusterID).Eq(q.ClusterBind.WithContext(ctx).Select(q.ClusterBind.ClusterID).Where(q.ClusterBind.ProjectName.Eq(project_name), q.ClusterBind.ClusterName.Eq(cluster_name)))).Scan(&cluster)
		// if err != nil {
		// 	error_msg := fmt.Sprintf("err: %s ", err)
		// 	util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		// 	return
		// }

	}
	if len(cluster) == 0 {
		util.ReturnMsg(ctx, http.StatusNotFound, "", "The Cluster is not Found!")
		return
	} else {
		util.ReturnMsg(ctx, http.StatusOK, cluster, "success")
	}
}

// @BasePath /api/v1
// PingProject godoc
// @Summary GetClusterId
// @Schemes
// @Description Get ClusterID
// @Tags GetClusterId
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string true "Project Name"
// @Param cluster_name query string true "Cluster Name"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /cluster/get_id [get]
// @ID GetClusterId
func GetClusterId(ctx *gin.Context) {
	project_name := ctx.Query("project_name")
	cluster_name := ctx.Query("cluster_name")

	q := query.Use(global.Config.DB.DB()).ClusterBind
	ClusterList, err := q.WithContext(ctx).Select(q.ClusterID).Where(q.ProjectName.Eq(project_name), q.ClusterName.Eq(cluster_name)).First()

	if err != nil {
		if ClusterList == nil {
			error_msg := fmt.Sprintf("The ClusterName: [%s] is not Found! in ProjectName: [%s]", cluster_name, project_name)
			util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		}
		return
		//ctx.JSON(http.StatusInternalServerError, err)
	}
	util.ReturnMsg(ctx, http.StatusOK, ClusterList.ClusterID, "success")
}

type ClusterId struct {
	ClusterId string `json:"cluster_id"`
}

// @BasePath /api/v1
// PingCluster godoc
// @Summary ListGetLabel
// @Schemes
// @Description List Cluster for label
// @Tags ListGetLabel
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string false "Project Name"
// @Param label_key query string true "Label Key"
// @Param label_value query string true "Label value"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /cluster/label [get]
// @ID ListGetLabel
func ListGetLabel(ctx *gin.Context) {
	project_name := ctx.Query("project_name")
	label_key := ctx.Query("label_key")
	label_value := ctx.Query("label_value")

	q := query.Use(global.Config.DB.DB())
	//判断标签是否存在
	label, err := q.Label.WithContext(ctx).Where(q.Label.Key.Eq(label_key), q.Label.Value.Eq(label_value)).Find()

	if err != nil {
		error_msg := fmt.Sprintf("err: %s ", err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}

	if len(label) == 0 {
		error_msg := fmt.Sprintf("It's no Cluster, for label : [%s=%s]", label_key, label_value)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}

	//clusters := [][]model.Cluster{}
	clusters := []model.Cluster{}

	if project_name == "" {
		for _, v := range label {
			cluster := model.Cluster{}
			global.Config.DB.DB().
				Raw("SELECT  cluster.id, cluster.cluster_id,cluster.cluster_name,cluster.`describe` ,(select GROUP_CONCAT(label.key, '=',label.value) FROM  label  where FIND_IN_SET (label.label_id,cluster.labels) )AS Labels  FROM `cluster` WHERE  FIND_IN_SET  (?,cluster.labels)", v.LabelID).
				Scan(&cluster)
			clusters = append(clusters, cluster)
		}
	} else {
		projectList, err := q.Project.WithContext(ctx).Where(q.Project.ProjectName.Eq(project_name)).First()

		if err != nil {
			if projectList == nil {
				error_msg := fmt.Sprintf("The ProjectName: [%s] is not Found!", project_name)
				util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
			}
			return
		}
		for _, v := range label {
			cluster := model.Cluster{}
			global.Config.DB.DB().
				Raw("SELECT  cluster.id, cluster.cluster_id,cluster.cluster_name,cluster.`describe` ,(select GROUP_CONCAT(label.key, '=',label.value) FROM  label  where FIND_IN_SET (label.label_id,cluster.labels) )AS Labels  FROM `cluster` WHERE  FIND_IN_SET  (?,cluster.labels) and cluster.cluster_id IN (select cluster_bind.cluster_id from cluster_bind where cluster_bind.project_name= ?)", v.LabelID, project_name).
				Scan(&cluster)
			//clusters = append(clusters, cluster)
			if cluster != (model.Cluster{}) {
				clusters = append(clusters, cluster)
			}
		}

	}

	util.ReturnMsg(ctx, http.StatusOK, clusters, "success")

}

type CreateClusterRequestBody struct {
	ProjectName string `json:"project_name" binding:"required"`
	ClusterName string `json:"cluster_name" binding:"required"`
	Description string `json:"description"`
}
type UpdateClusterRequestBody struct {
	ClusterName string `json:"cluster_name"`
	Description string `json:"description"`
}

type ErrorResp struct {
	Msg string `json:"msg"`
}

// @BasePath /api/v1
// PingProject godoc
// @Summary CreateCluster
// @Schemes
// @Description Create Cluster
// @Tags CreateCluster
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ReqeustBody body CreateClusterRequestBody true "Create Cluster"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /cluster/create [post]
// @ID CreateCluster
func CreateCluster(ctx *gin.Context) {
	body := CreateClusterRequestBody{}
	err := ctx.ShouldBind(&body)
	if err != nil {
		error_msg := fmt.Sprintf("err: %s ", err)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return
	}

	q := query.Use(global.Config.DB.DB())
	//判断项目是否存在
	projectID, err := q.Project.WithContext(ctx).Select(q.Project.ProjectID).Where(q.Project.ProjectName.Eq(body.ProjectName)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The ProjectName: [%s] is not exists!", body.ProjectName)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}
	//判断项目下是否已经绑定了该集群
	_, err = q.ClusterBind.WithContext(ctx).Where(q.ClusterBind.ProjectName.Eq(body.ProjectName), q.ClusterBind.ClusterName.Eq(body.ClusterName)).First()

	if !errors.Is(err, gorm.ErrRecordNotFound) {

		error_msg := fmt.Sprintf("The ClusterName: [%s] is already exists! in ProjectName: [%s]", body.ClusterName, body.ProjectName)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return
	}

	//创建集群
	id, err := util.NewIdMgr(5)
	if err != nil {
		fmt.Println(err)
		return
	}

	cluster := model.Cluster{}
	cluster_id := util.Int64toString(id.ID())
	cluster.ClusterName = body.ClusterName
	cluster.Describe = body.Description
	cluster.ClusterID = cluster_id

	ClusterBind := model.ClusterBind{}
	ClusterBind.ClusterID = cluster_id
	ClusterBind.ClusterName = body.ClusterName
	ClusterBind.ProjectName = body.ProjectName
	ClusterBind.ProjectID = projectID.ProjectID

	//qs := query.Use(global.Config.DB)
	q.Transaction(func(tx *query.Query) error {
		if err := tx.Cluster.WithContext(ctx).Create(&cluster); err != nil {
			return err
		}
		if err := tx.ClusterBind.WithContext(ctx).Create(&ClusterBind); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		error_msg := fmt.Sprintf("The ClusterName: [%s] Create Failed! in ProjectName : [%s], err: %s", body.ClusterName, body.ProjectName, err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}
	ok_msg := fmt.Sprintf("The Cluster_Name: [%s] Create successfully! in ProjectName : [%s]", body.ClusterName, body.ProjectName)
	util.ReturnMsg(ctx, http.StatusOK, cluster, ok_msg)
}

// @BasePath /api/v1
// PingCluster godoc
// @Summary UpdateCluster
// @Schemes
// @Description Update Cluster
// @Tags UpdateCluster
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string true "Project_Name"
// @Param cluster_id query string true "Cluster_ID"
// @Param ReqeustBody body UpdateClusterRequestBody true "Update Cluster"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /cluster/update [put]
// @ID UpdateCluster
func UpdateCluster(ctx *gin.Context) {
	project_name := ctx.Query("project_name")
	cluster_id := ctx.Query("cluster_id")

	body := UpdateClusterRequestBody{}
	err := ctx.ShouldBind(&body)
	if err != nil {
		error_msg := fmt.Sprintf("err: %s ", err)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return
	}

	cluster := model.Cluster{}
	cluster.ClusterName = body.ClusterName
	cluster.Describe = body.Description

	ClusterBind := model.ClusterBind{}
	ClusterBind.ClusterName = body.ClusterName

	q := query.Use(global.Config.DB.DB())

	_, err = q.Project.WithContext(ctx).Where(q.Project.ProjectName.Eq(project_name)).First()
	info := gen.ResultInfo{}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The Project_Name: [%s] is not Found!", project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return

	} else {
		if err != nil {
			error_msg := fmt.Sprintf("err: %s ", err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
			return
		}

		q.Transaction(func(tx *query.Query) error {
			//if err := tx.Project.WithContext(ctx).Create(&project); err != nil {
			if info, err = tx.Cluster.WithContext(ctx).Where(tx.Cluster.ClusterID.Eq(cluster_id)).Updates(&cluster); err != nil {
				return err
			}
			if info, err = tx.ClusterBind.WithContext(ctx).Where(tx.ClusterBind.ProjectName.Eq(project_name), tx.ClusterBind.ClusterID.Eq(cluster_id)).Updates(&ClusterBind); err != nil {
				return err
			}
			return nil
		})
	}
	if err != nil {
		error_msg := fmt.Sprintf("err: %s ", err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}

	if info.RowsAffected == 0 {

		error_msg := fmt.Sprintf("The ClusterId: [%s] is not Found! in ProjectName : [%s]", cluster_id, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}

	ok_msg := fmt.Sprintf("The Cluster_id: [%s] Update successfully! in ProjectName : [%s]", cluster_id, project_name)
	util.ReturnMsg(ctx, http.StatusOK, util.Int64toString(info.RowsAffected), ok_msg)

}

// @BasePath /api/v1
// PingCluster godoc
// @Summary DeleteCluster
// @Schemes
// @Description Delete Cluster
// @Tags DeleteCluster
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string true "Project_Name"
// @Param cluster_name query string true "Cluster_Name"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /cluster/delete [delete]
// @ID DeleteCluster
func DeleteCluster(ctx *gin.Context) {
	project_name := ctx.Query("project_name")
	cluster_name := ctx.Query("cluster_name")

	q := query.Use(global.Config.DB.DB())

	_, err := q.Project.WithContext(ctx).Where(q.Project.ProjectName.Eq(project_name)).First()
	info := gen.ResultInfo{}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The Project_Name: [%s] is not Found!", project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return

	} else {
		if err != nil {
			error_msg := fmt.Sprintf("err: %s ", err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
			return
		}

		clusterId, err := q.ClusterBind.WithContext(ctx).Select(q.ClusterBind.ClusterID).Where(q.ClusterBind.ProjectName.Eq(project_name), q.ClusterBind.ClusterName.Eq(cluster_name)).First()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			error_msg := fmt.Sprintf("The ClusterName: [%s] is not Found! in ProjectName: [%s]", cluster_name, project_name)
			util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
			return
		}
		q.Transaction(func(tx *query.Query) error {

			//if err := tx.Project.WithContext(ctx).Create(&project); err != nil {
			if info, err = tx.Cluster.WithContext(ctx).Where(tx.Cluster.ClusterID.Eq(clusterId.ClusterID)).Delete(); err != nil {
				return err
			}
			if info, err = tx.ClusterBind.WithContext(ctx).Where(tx.ClusterBind.ProjectName.Eq(project_name), tx.ClusterBind.ClusterID.Eq(clusterId.ClusterID)).Delete(); err != nil {
				return err
			}
			return nil
		})
	}
	if info.RowsAffected == 0 {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is deleted Failed! in ProjectName: [%s]", cluster_name, project_name)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return
	}

	util.ReturnMsg(ctx, http.StatusOK, util.Int64toString(info.RowsAffected), "Cluster deleted successfully!")

}
