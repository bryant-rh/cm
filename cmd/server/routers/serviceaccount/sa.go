package serviceaccount

import (
	"github.com/bryant-rh/cm/cmd/server/global"
	"github.com/bryant-rh/cm/internal/model"
	"github.com/bryant-rh/cm/internal/query"
	"github.com/bryant-rh/cm/pkg/util"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func SaRouter(r *gin.RouterGroup) {
	r.GET("/sa/list", ListSa)
	r.GET("/sa/listns", ListNs)
	r.GET("/sa/gettoken", Gettoken)
	r.POST("/sa/create", CreateSa)
	r.PUT("/sa/update", UpdateSa)
	r.DELETE("/sa/delete", DeleteSa)
	r.POST("/sa/addns", AddNameSpace)
	r.DELETE("/sa/delns", DeleteNs)
}

type Cluster struct {
	Name      string `json:"name"`
	ClusterID int64  `json:"cluster_id"`
}

type SaId struct {
	SaId string `json:"sa_id"`
}

type NameSpace struct {
	NsName string `json:"ns_name"`
}

type SaRes struct {
	ID        int32     `json:"id"`
	SaID      string    `json:"sa_id"`
	SaName    string    `json:"sa_name"`
	SaToken   string    `json:"sa_token"`
	NameSpace string    `json:"namespace"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// @BasePath /api/v1
// PingCluster godoc
// @Summary ListSa
// @Schemes
// @Description List All ServiceAccount
// @Tags ListSa
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string true "Project Name"
// @Param cluster_name query string true "Cluster Name"
// @Param sa_name query string false "Sa Name"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /sa/list [get]
// @ID ListSa
func ListSa(ctx *gin.Context) {
	project_name := ctx.Query("project_name")
	cluster_name := ctx.Query("cluster_name")
	sa_name := ctx.Query("sa_name")
	q := query.Use(global.Config.DB.DB())
	//判断项目和集群是否存在
	cluster, err := q.ClusterBind.WithContext(ctx).Select(q.ClusterBind.ClusterID).Where(q.ClusterBind.ProjectName.Eq(project_name), q.ClusterBind.ClusterName.Eq(cluster_name)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is not exists! in ProjectName: [%s]", cluster_name, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}
	//results := []*model.Serviceaccount{}
	results := []SaRes{}
	said_res := []SaId{}
	if sa_name == "" {
		err = q.Namespace.WithContext(ctx).Distinct(q.Namespace.SaID).Where(q.Namespace.ClusterID.Eq(cluster.ClusterID)).Scan(&said_res)
	} else {
		err = q.Namespace.WithContext(ctx).Distinct(q.Namespace.SaID).Where(q.Namespace.ClusterID.Eq(cluster.ClusterID), q.Namespace.SaName.Eq(sa_name)).Scan(&said_res)

	}
	if len(said_res) != 0 {
		saIDs := []string{}
		for _, v := range said_res {
			saIDs = append(saIDs, v.SaId)
		}
		saIdSlice := util.SlicetoString(saIDs)
		//results, err = q.Serviceaccount.WithContext(ctx).Where(q.Serviceaccount.SaID.FindInSet(saIdSlice)).Find()
		global.Config.DB.DB().
			Raw("SELECT serviceaccount.id, serviceaccount.sa_id, serviceaccount.sa_name,serviceaccount.sa_token, (select GROUP_CONCAT(namespace.ns_name) from namespace WHERE namespace.sa_id=serviceaccount.sa_id ) AS NameSpace ,serviceaccount.created_at,serviceaccount.updated_at from serviceaccount where FIND_IN_SET (serviceaccount.sa_id ,?)", saIdSlice).
			Scan(&results)

	}

	if err != nil {
		error_msg := fmt.Sprintf("The ServiceAccount Find Failed!, err: %s", err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}
	if len(results) == 0 {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is not Found  ServiceAccount! in ProjectName : [%s]", cluster_name, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)

	} else {
		ok_msg := fmt.Sprintf("The ClusterName: [%s] Find  ServiceAccount Success! in ProjectName : [%s]", cluster_name, project_name)
		util.ReturnMsg(ctx, http.StatusOK, results, ok_msg)

	}

}

// @BasePath /api/v1
// PingCluster godoc
// @Summary ListNs
// @Schemes
// @Description List All ServiceAccount's NameSpace
// @Tags ListNs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string true "Project Name"
// @Param cluster_name query string true "Cluster Name"
// @Param sa_name query string true "Sa_Name"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /sa/listns [get]
// @ID ListNs
func ListNs(ctx *gin.Context) {
	project_name := ctx.Query("project_name")
	cluster_name := ctx.Query("cluster_name")
	sa_name := ctx.Query("sa_name")

	q := query.Use(global.Config.DB.DB())
	//判断项目和集群是否存在
	cluster, err := q.ClusterBind.WithContext(ctx).Select(q.ClusterBind.ClusterID).Where(q.ClusterBind.ProjectName.Eq(project_name), q.ClusterBind.ClusterName.Eq(cluster_name)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is not exists! in ProjectName: [%s]", cluster_name, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}
	ns_res := []NameSpace{}
	err = q.Namespace.WithContext(ctx).Select(q.Namespace.NsName).Where(q.Namespace.SaName.Eq(sa_name), q.Namespace.ClusterID.Eq(cluster.ClusterID)).Scan(&ns_res)

	if err != nil {
		error_msg := fmt.Sprintf("The ServiceAccount:[%s] Find NameSpace Failed!, err: %s", sa_name, err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}
	if len(ns_res) == 0 {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is not Found NameSpace, For ServiceAccount: [%s] in ProjectName : [%s]", cluster_name, sa_name, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}

	ok_msg := fmt.Sprintf("The ClusterName: [%s] Find Namespace Succeed! For  ServiceAccount : [%s] Success! in ProjectName : [%s]", cluster_name, sa_name, project_name)
	util.ReturnMsg(ctx, http.StatusOK, ns_res, ok_msg)

}

// @BasePath /api/v1
// PingCluster godoc
// @Summary Gettoken
// @Schemes
// @Description Get SaToken For NameSpace
// @Tags Gettoken
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string true "Project Name"
// @Param cluster_name query string true "Cluster Name"
// @Param ns_name query string true "Ns_Name"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /sa/gettoken [get]
// @ID Gettoken
func Gettoken(ctx *gin.Context) {
	project_name := ctx.Query("project_name")
	cluster_name := ctx.Query("cluster_name")
	ns_name := ctx.Query("ns_name")

	q := query.Use(global.Config.DB.DB())
	//判断项目和集群是否存在
	cluster, err := q.ClusterBind.WithContext(ctx).Select(q.ClusterBind.ClusterID).Where(q.ClusterBind.ProjectName.Eq(project_name), q.ClusterBind.ClusterName.Eq(cluster_name)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is not exists! in ProjectName: [%s]", cluster_name, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}
    SaId, err := q.Namespace.WithContext(ctx).Select(q.Namespace.SaID).Where(q.Namespace.NsName.Eq(ns_name), q.Namespace.ClusterID.Eq(cluster.ClusterID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			error_msg := fmt.Sprintf("The ClusterName: [%s], NameSpace: [%s] is not Found  ServiceAccount ,in ProjectName : [%s]", cluster_name, ns_name, project_name)
			util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)

			return
		} else {
			error_msg := fmt.Sprintf("The NameSpace:[%s] Find  ServiceAccount Failed!, err: %s", ns_name, err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
			return
		}
	}
	token, err := q.Serviceaccount.WithContext(ctx).Where(q.Serviceaccount.SaID.Eq(SaId.SaID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			error_msg := fmt.Sprintf("The ClusterName: [%s], NameSpace: [%s] is not Found  ServiceAccount ,in ProjectName : [%s]", cluster_name, ns_name, project_name)
			util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)

			return
		} else {
			error_msg := fmt.Sprintf("The NameSpace:[%s] Find  ServiceAccount Failed!, err: %s", ns_name, err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
			return
		}
	}

	util.ReturnMsg(ctx, http.StatusOK, token.SaToken, "success")

}

type CreateSaRequestBody struct {
	ProjectName string `json:"project_name" binding:"required"`
	ClusterName string `json:"cluster_name" binding:"required"`
	SaName      string `json:"sa_name" binding:"required"`
	SaToken     string `json:"sa_token" binding:"required"`
	Namespace   string `json:"namespace" binding:"required"`
}

type ErrorResp struct {
	Msg string `json:"msg"`
}

// @BasePath /api/v1
// PingProject godoc
// @Summary CreateSa
// @Schemes
// @Description Create ServiceAccount
// @Tags CreateSa
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ReqeustBody body CreateSaRequestBody true "Create Sa"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /sa/create [post]
// @ID CreateSa
func CreateSa(ctx *gin.Context) {
	body := CreateSaRequestBody{}
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

	//生成id
	id, err := util.NewIdMgr(5)
	if err != nil {
		fmt.Println(err)
		return
	}

	//判断对应项目的集群下是否已存在此sa
	result := []SaId{}
	err = q.Serviceaccount.WithContext(ctx).Select(q.Serviceaccount.SaID).Where(q.Serviceaccount.SaName.Eq(body.SaName)).Scan(&result)
	if err != nil {
		error_msg := fmt.Sprintf("The ServiceAccount: [%s] Find Failed!, err: %s", body.SaName, err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}

	sa_id := ""
	if len(result) != 0 {
		saIDs := []string{}
		for _, v := range result {
			saIDs = append(saIDs, v.SaId)
		}
		saIdSlice := util.SlicetoString(saIDs)

		fmt.Printf("saIdslice: %s", saIdSlice)

		//查看当前集群绑定的sa
		//sa_res := model.Serviceaccount{}
		said_res := SaId{}
		err := q.Namespace.WithContext(ctx).Distinct(q.Namespace.SaID).Where(q.Namespace.SaName.Eq(body.SaName), q.Namespace.ClusterID.Eq(clusterId.ClusterID)).Scan(&said_res)
		if err != nil {
			error_msg := fmt.Sprintf("The SaId: [%s] Find Failed!, err: %s", body.SaName, err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
			return
		}

		if said_res != (SaId{}) {
			error_msg := fmt.Sprintf("The ClusterName: [%s] is already Exist ServiceAccount: [%s], in ProjectName : [%s]", body.ClusterName, body.SaName, body.ProjectName)
			util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
			return
			// global.Config.DB.DB().Raw("SELECT * FROM `serviceaccount` WHERE  FIND_IN_SET  (?,?)", said_res.SaId, saIdSlice).Scan(&sa_res)
			// if sa_res != (model.Serviceaccount{}) {
			// 	sa_id = said_res.SaId
			// }

		} else {
			sa_id = util.Int64toString(id.ID())

		}

	} else {
		sa_id = util.Int64toString(id.ID())

	}
	//创建sa
	sa := model.Serviceaccount{}
	sa.SaID = sa_id
	sa.SaName = body.SaName
	sa.SaToken = body.SaToken

	err = q.Serviceaccount.WithContext(ctx).Create(&sa)
	if err != nil {
		error_msg := fmt.Sprintf("The ServiceAccount: [%s] Create Failed!, err: %s", body.SaName, err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}

	ns_slice := []model.Namespace{}
	ns := util.StringtoSlice(body.Namespace)
	for _, v := range ns {
		_, err = q.Namespace.WithContext(ctx).Where(q.Namespace.SaName.Eq(body.SaName), q.Namespace.ClusterID.Eq(clusterId.ClusterID), q.Namespace.NsName.Eq(v)).First()

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {

				ns := model.Namespace{}
				ns_id := util.Int64toString(id.ID())
				ns.NsID = ns_id
				ns.NsName = v
				ns.SaID = sa_id
				ns.SaName = body.SaName
				ns.ClusterID = clusterId.ClusterID

				err = q.Namespace.WithContext(ctx).Create(&ns)
				if err == nil {
					ns_slice = append(ns_slice, ns)
				}
			} else {
				error_msg := fmt.Sprintf("The ClusterName: [%s] Create ServiceAccount : [%s] for NameSpace: [%s] Failed! in ProjectName : [%s], err: %s", body.ClusterName, body.SaName, body.Namespace, body.ProjectName, err)
				util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
				return
			}
		}

	}

	if len(ns_slice) == 0 {
		//error_msg := fmt.Sprintf("The ClusterName: [%s] Create ServiceAccount : [%s] for NameSpace: [%s] already Exist! in ProjectName : [%s]", body.ClusterName, body.SaName, body.Namespace, body.ProjectName)
		error_msg := fmt.Sprintf("The ClusterName: [%s] Create ServiceAccount : [%s] for NameSpace: [%s] Failed! in ProjectName : [%s], err: %s", body.ClusterName, body.SaName, body.Namespace, body.ProjectName, err)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return

	} else {
		res := []SaRes{}
		global.Config.DB.DB().
			Raw("SELECT serviceaccount.id,serviceaccount.sa_id,serviceaccount.sa_name,serviceaccount.sa_token, (select GROUP_CONCAT(namespace.ns_name) from namespace WHERE namespace.sa_id =? and namespace.cluster_id = ? ) AS NameSpace ,serviceaccount.created_at,serviceaccount.updated_at from serviceaccount WHERE  serviceaccount.sa_id = ?", sa_id, clusterId.ClusterID, sa_id).
			Scan(&res)
		ok_msg := fmt.Sprintf("The ClusterName: [%s] Create ServiceAccount : [%s] for NameSpace: [%s]  successfully! in ProjectName : [%s]", body.ClusterName, body.SaName, body.Namespace, body.ProjectName)
		util.ReturnMsg(ctx, http.StatusOK, res, ok_msg)
		return

	}
}

type UpdateSaRequestBody struct {
	SaToken string `json:"sa_token" binding:"required"`
}

// @BasePath /api/v1
// PingCluster godoc
// @Summary UpdateSaToken
// @Schemes
// @Description Update ServiceAccount SaToken
// @Tags UpdateSaToken
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string true "Project_Name"
// @Param cluster_name query string true "Cluster_Name"
// @Param sa_name query string true "Sa_Name"
// @Param ReqeustBody body UpdateSaRequestBody true "Update ServiceAccount Token"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /sa/update [put]
// @ID UpdateSa
func UpdateSa(ctx *gin.Context) {
	project_name := ctx.Query("project_name")
	cluster_name := ctx.Query("cluster_name")
	sa_name := ctx.Query("sa_name")

	body := UpdateSaRequestBody{}
	err := ctx.ShouldBind(&body)
	if err != nil {
		error_msg := fmt.Sprintf("err: %s ", err)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return
	}

	q := query.Use(global.Config.DB.DB())
	//判断项目和集群是否存在
	cluster, err := q.ClusterBind.WithContext(ctx).Select(q.ClusterBind.ClusterID).Where(q.ClusterBind.ProjectName.Eq(project_name), q.ClusterBind.ClusterName.Eq(cluster_name)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is not exists! in ProjectName: [%s]", cluster_name, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}
	said_res := SaId{}
	err = q.Namespace.WithContext(ctx).Distinct(q.Namespace.SaID).Where(q.Namespace.SaName.Eq(sa_name), q.Namespace.ClusterID.Eq(cluster.ClusterID)).Scan(&said_res)

	if err != nil {
		error_msg := fmt.Sprintf("The SaId: [%s] Find Failed!, err: %s", sa_name, err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}
	if said_res == (SaId{}) {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is not Found  ServiceAccount: [%s] in ProjectName : [%s]", cluster_name, sa_name, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}
	sa := model.Serviceaccount{}
	sa.SaToken = body.SaToken

	info, err := q.Serviceaccount.WithContext(ctx).Where(q.Serviceaccount.SaID.Eq(said_res.SaId)).Updates(&sa)

	if err != nil {
		error_msg := fmt.Sprintf("The ClusterName: [%s] Update ServiceAccount : [%s] Failed! in ProjectName : [%s], err: %s", cluster_name, sa_name, project_name, err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}

	if info.RowsAffected == 0 {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is not Found  ServiceAccount: [%s] in ProjectName : [%s]", cluster_name, sa_name, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}

	ok_msg := fmt.Sprintf("The ClusterName: [%s] Update ServiceAccount : [%s] Success! in ProjectName : [%s]", cluster_name, sa_name, project_name)
	util.ReturnMsg(ctx, http.StatusOK, util.Int64toString(info.RowsAffected), ok_msg)

}

// @BasePath /api/v1
// PingCluster godoc
// @Summary DeleteSa
// @Schemes
// @Description Delete ServiceAccount
// @Tags DeleteSa
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string true "Project_Name"
// @Param cluster_name query string true "Cluster_Name"
// @Param sa_name query string true "sa_name"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /sa/delete [delete]
// @ID DeleteSa
func DeleteSa(ctx *gin.Context) {
	project_name := ctx.Query("project_name")
	cluster_name := ctx.Query("cluster_name")
	sa_name := ctx.Query("sa_name")

	q := query.Use(global.Config.DB.DB())
	//判断项目和集群是否存在
	clusterId, err := q.ClusterBind.WithContext(ctx).Select(q.ClusterBind.ClusterID).Where(q.ClusterBind.ProjectName.Eq(project_name), q.ClusterBind.ClusterName.Eq(cluster_name)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is not exists! in ProjectName: [%s]", cluster_name, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}
	//判断对应项目的集群下是否已存在此sa
	//clusterLabel := ""
	//_, err = q.Namespace.WithContext(ctx).Select(q.Cluster.Label).Where(q.Namespace.ClusterID.Eq(clusterId.ClusterID), q.Namespace.SaName.Eq(sa_name)).First()
	said_res := SaId{}
	err = q.Namespace.WithContext(ctx).Distinct(q.Namespace.SaID).Where(q.Namespace.ClusterID.Eq(clusterId.ClusterID), q.Namespace.SaName.Eq(sa_name)).Scan(&said_res)
	if err != nil {
		error_msg := fmt.Sprintf("The SaId: [%s] Find Failed!, err: %s", sa_name, err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}
	if said_res == (SaId{}) {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is Not ServiceAccount: [%s] in ProjectName: [%s]", cluster_name, sa_name, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}

	res_delete := gen.ResultInfo{}
	q.Transaction(func(tx *query.Query) error {

		if res_delete, err = tx.Serviceaccount.WithContext(ctx).Where(tx.Serviceaccount.SaID.Eq(said_res.SaId)).Delete(); err != nil {
			return err
		}
		if res_delete, err = tx.Namespace.WithContext(ctx).Where(q.Namespace.ClusterID.Eq(clusterId.ClusterID), q.Namespace.SaID.Eq(said_res.SaId)).Delete(); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		error_msg := fmt.Sprintf("The ClusterName: [%s] Delete ServiceAccount : [%s] Failed! in ProjectName : [%s], err: %s", cluster_name, sa_name, project_name, err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}

	ok_msg := fmt.Sprintf("The Cluster_Name: [%s] Delete ServiceAccount : [%s]  successfully! in ProjectName : [%s]", cluster_name, sa_name, project_name)
	util.ReturnMsg(ctx, http.StatusOK, util.Int64toString(res_delete.RowsAffected), ok_msg)

}

type AddSaNsRequestBody struct {
	ProjectName string `json:"project_name" binding:"required"`
	ClusterName string `json:"cluster_name" binding:"required"`
	SaName      string `json:"sa_name" binding:"required"`
	Namespace   string `json:"namespace" binding:"required"`
}

// @BasePath /api/v1
// PingProject godoc
// @Summary AddNameSpace
// @Schemes
// @Description ADD ServiceAccount NameSpace
// @Tags AddNameSpace
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ReqeustBody body AddSaNsRequestBody true "Add Sa NameSpace"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /sa/addns [post]
// @ID AddNameSpace
func AddNameSpace(ctx *gin.Context) {
	body := AddSaNsRequestBody{}
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

	//生成id
	id, err := util.NewIdMgr(5)
	if err != nil {
		fmt.Println(err)
		return
	}

	//判断对应项目的集群下是否存在此sa
	result := []SaId{}
	err = q.Serviceaccount.WithContext(ctx).Select(q.Serviceaccount.SaID).Where(q.Serviceaccount.SaName.Eq(body.SaName)).Scan(&result)
	if err != nil {
		error_msg := fmt.Sprintf("The ServiceAccount: [%s] Find Failed!, err: %s", body.SaName, err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}

	sa_id := ""
	if len(result) != 0 {
		saIDs := []string{}
		for _, v := range result {
			saIDs = append(saIDs, v.SaId)
		}
		saIdSlice := util.SlicetoString(saIDs)
		//查看当前集群绑定的sa
		sa_res := model.Serviceaccount{}
		said_res := SaId{}
		err := q.Namespace.WithContext(ctx).Distinct(q.Namespace.SaID).Where(q.Namespace.SaName.Eq(body.SaName), q.Namespace.ClusterID.Eq(clusterId.ClusterID)).Scan(&said_res)
		if err != nil {
			error_msg := fmt.Sprintf("The SaId: [%s] Find Failed!, err: %s", body.SaName, err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
			return
		}
		if said_res != (SaId{}) {
			global.Config.DB.DB().
				Raw("SELECT * FROM `serviceaccount` WHERE  FIND_IN_SET  (?,?)", said_res.SaId, saIdSlice).
				Scan(&sa_res)

			if sa_res != (model.Serviceaccount{}) {
				sa_id = said_res.SaId
			}

		} else {
			error_msg := fmt.Sprintf("The ClusterName: [%s] is Not ServiceAccount: [%s] in ProjectName: [%s]", body.ClusterName, body.SaName, body.ProjectName)
			util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
			return

		}

	} else {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is Not ServiceAccount: [%s] in ProjectName: [%s]", body.ClusterName, body.SaName, body.ProjectName)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}

	ns_slice := []model.Namespace{}
	ns := util.StringtoSlice(body.Namespace)
	for _, v := range ns {
		_, err = q.Namespace.WithContext(ctx).Where(q.Namespace.SaName.Eq(body.SaName), q.Namespace.ClusterID.Eq(clusterId.ClusterID), q.Namespace.NsName.Eq(v)).First()

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {

				ns := model.Namespace{}
				ns_id := util.Int64toString(id.ID())
				ns.NsID = ns_id
				ns.NsName = v
				ns.SaID = sa_id
				ns.SaName = body.SaName
				ns.ClusterID = clusterId.ClusterID

				err = q.Namespace.WithContext(ctx).Create(&ns)
				if err == nil {
					ns_slice = append(ns_slice, ns)
				}
			} else {
				error_msg := fmt.Sprintf("The ClusterName: [%s] Create ServiceAccount : [%s] for NameSpace: [%s] Failed! in ProjectName : [%s], err: %s", body.ClusterName, body.SaName, body.Namespace, body.ProjectName, err)
				util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
				return
			}
		}

	}

	if len(ns_slice) == 0 {
		ok_msg := fmt.Sprintf("The ClusterName: [%s] Create ServiceAccount : [%s] for NameSpace: [%s] already Exist! in ProjectName : [%s]", body.ClusterName, body.SaName, body.Namespace, body.ProjectName)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", ok_msg)
		return

	} else {
		res := []SaRes{}
		ok_msg := fmt.Sprintf("The ClusterName: [%s] Create ServiceAccount : [%s] for NameSpace: [%s]  successfully! in ProjectName : [%s]", body.ClusterName, body.SaName, body.Namespace, body.ProjectName)
		global.Config.DB.DB().
			Raw("SELECT serviceaccount.id,serviceaccount.sa_id,serviceaccount.sa_name,serviceaccount.sa_token, (select GROUP_CONCAT(namespace.ns_name) from namespace WHERE namespace.sa_name =? and namespace.cluster_id = ? ) AS NameSpace ,serviceaccount.created_at,serviceaccount.updated_at from serviceaccount WHERE  serviceaccount.sa_name = ?", body.SaName, clusterId.ClusterID, body.SaName).
			Scan(&res)
		util.ReturnMsg(ctx, http.StatusOK, res, ok_msg)
		return
	}
	//ok_msg := fmt.Sprintf("The ClusterName: [%s] Create ServiceAccount : [%s] for NameSpace: [%s]  successfully! in ProjectName : [%s]", body.ClusterName, body.SaName, body.Namespace, body.ProjectName)
	//util.ReturnMsg(ctx, http.StatusOK, ns_slice, ok_msg)
}

// @BasePath /api/v1
// PingProject godoc
// @Summary DeleteNs
// @Schemes
// @Description Delete ServiceAccount NameSpace
// @Tags DeleteNs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string true "Project_Name"
// @Param cluster_name query string true "Cluster_Name"
// @Param sa_name query string true "Sa_Name"
// @Param namespace query string true "NameSpace"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /sa/delns [delete]
// @ID DeleteNs
func DeleteNs(ctx *gin.Context) {
	project_name := ctx.Query("project_name")
	cluster_name := ctx.Query("cluster_name")
	sa_name := ctx.Query("sa_name")
	namespace := ctx.Query("namespace")

	q := query.Use(global.Config.DB.DB())
	//判断项目和集群是否存在
	clusterId, err := q.ClusterBind.WithContext(ctx).Select(q.ClusterBind.ClusterID).Where(q.ClusterBind.ProjectName.Eq(project_name), q.ClusterBind.ClusterName.Eq(cluster_name)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is not exists! in ProjectName: [%s]", cluster_name, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)

		return
	}

	//判断对应项目的集群下是否存在此sa
	result := []SaId{}
	err = q.Serviceaccount.WithContext(ctx).Select(q.Serviceaccount.SaID).Where(q.Serviceaccount.SaName.Eq(sa_name)).Scan(&result)
	if err != nil {
		error_msg := fmt.Sprintf("The ServiceAccount: [%s] Find Failed!, err: %s", sa_name, err)
		util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		return
	}

	if len(result) != 0 {

		//查看当前集群绑定的sa
		said_res := SaId{}
		err := q.Namespace.WithContext(ctx).Distinct(q.Namespace.SaID).Where(q.Namespace.SaName.Eq(sa_name), q.Namespace.ClusterID.Eq(clusterId.ClusterID)).Scan(&said_res)
		if err != nil {
			error_msg := fmt.Sprintf("The SaId: [%s] Find Failed!, err: %s", sa_name, err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
			return
		}
		if said_res == (SaId{}) {
			error_msg := fmt.Sprintf("The ClusterName: [%s] is Not ServiceAccount: [%s] in ProjectName: [%s]", cluster_name, sa_name, project_name)
			util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
			return

		}

	} else {
		error_msg := fmt.Sprintf("The ClusterName: [%s] is Not ServiceAccount: [%s] in ProjectName: [%s]", cluster_name, sa_name, project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return
	}

	ns := util.StringtoSlice(namespace)
	res_delete := gen.ResultInfo{}
	for _, v := range ns {
		ns_id, err := q.Namespace.WithContext(ctx).Select(q.Namespace.NsID).Where(q.Namespace.SaName.Eq(sa_name), q.Namespace.ClusterID.Eq(clusterId.ClusterID), q.Namespace.NsName.Eq(v)).First()

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			res_delete, err = q.Namespace.WithContext(ctx).Where(q.Namespace.NsID.Eq(ns_id.NsID)).Delete()
			if err != nil {
				error_msg := fmt.Sprintf("The ClusterName: [%s] 下的ServiceAccount : [%s] for NameSpace: [%s] Delete Failed! in ProjectName : [%s], err: %s", cluster_name, sa_name, v, project_name, err)
				util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
				return

			}
		}

	}

	if res_delete.RowsAffected == 0 {
		error_msg := fmt.Sprintf("The ClusterName: [%s] Create ServiceAccount : [%s] for NameSpace: [%s] Delete Failed! in ProjectName : [%s]", cluster_name, sa_name, namespace, project_name)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return

	} else {
		ok_msg := fmt.Sprintf("The ClusterName: [%s] Create ServiceAccount : [%s] for NameSpace: [%s]  successfully! in ProjectName : [%s]", cluster_name, sa_name, namespace, project_name)
		util.ReturnMsg(ctx, http.StatusOK, util.Int64toString(res_delete.RowsAffected), ok_msg)
		return
	}
}
