package project

import (
	"github.com/bryant-rh/cm/cmd/server/global"
	"github.com/bryant-rh/cm/pkg/util"
	"errors"
	"fmt"
	"net/http"

	"github.com/bryant-rh/cm/internal/model"
	"github.com/bryant-rh/cm/internal/query"

	"github.com/gin-gonic/gin"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func ProjectRouter(r *gin.RouterGroup) {
	r.GET("/project/list", ListProject)
	r.GET("/project/get_id", GetProjectId)
	r.POST("/project/create", CreateProject)
	r.PUT("/project/update", UpdateProject)
	r.DELETE("/project/delete", DeleteProject)
}

type Project struct {
	Name      string `json:"name"`
	ProjectID int64  `json:"project_id"`
}

// @BasePath /api/v1
// PingProject godoc
// @Summary ListProject
// @Schemes
// @Description List Project
// @Tags ListProject
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string false "Project Name"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /project/list [get]
// @ID ListProject
func ListProject(ctx *gin.Context) {
	project_name := ctx.Query("project_name")

	result := []model.Project{}
	q := query.Use(global.Config.DB.DB()).Project
	if project_name == "" {
		err := q.WithContext(ctx).Scan(&result)
		if err != nil {
			error_msg := fmt.Sprintf("err: %s ", err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, result, error_msg)
			return
		}
	} else {
		err := q.WithContext(ctx).Where(q.ProjectName.Eq(project_name)).Scan(&result)
		if err != nil {
			error_msg := fmt.Sprintf("err: %s ", err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, result, error_msg)
			return
		}

	}
	if len(result) == 0 {
		msg := fmt.Sprintf("The ProjectName: [%s] is not Found!", project_name)
		util.ReturnMsg(ctx, http.StatusNotFound, result, msg)
	} else {
		msg := "success"
		util.ReturnMsg(ctx, http.StatusOK, result, msg)
	}
}

// @BasePath /api/v1
// PingProject godoc
// @Summary GetProjectId
// @Schemes
// @Description Get ProjectID
// @Tags GetProjectId
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string true "Project Name"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /project/get_id [get]
// @ID GetProjectId
func GetProjectId(ctx *gin.Context) {
	project_name := ctx.Query("project_name")

	q := query.Use(global.Config.DB.DB()).Project
	projectList, err := q.WithContext(ctx).Select(q.ProjectID).Where(q.ProjectName.Eq(project_name)).First()

	if err != nil {
		if projectList == nil {
			error_msg := fmt.Sprintf("The ProjectName: [%s] is not Found!", project_name)
			util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		}
		return
	}
	util.ReturnMsg(ctx, http.StatusOK, projectList.ProjectID, "success")
}

type CreateProjectRequestBody struct {
	ProjectName string `json:"project_name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateProjectRequestBody struct {
	ProjectName string `json:"project_name"`
	Description string `json:"description"`
}

type ErrorResp struct {
	Msg string `json:"msg"`
}

type SuccessResp struct {
	Msg string `json:"msg"`
}

// @BasePath /api/v1
// PingProject godoc
// @Summary CreateProject
// @Schemes
// @Description Create Project
// @Tags CreateProject
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ReqeustBody body CreateProjectRequestBody true "Create Project"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /project/create [post]
// @ID CreateProject
func CreateProject(ctx *gin.Context) {
	body := CreateProjectRequestBody{}
	err := ctx.ShouldBind(&body)
	if err != nil {
		error_msg := fmt.Sprintf("err: %s ", err)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return
	}

	q := query.Use(global.Config.DB.DB()).Project
	_, err = q.WithContext(ctx).Where(q.ProjectName.Eq(body.ProjectName)).First()

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The ProjectName: [%s] already exists!", body.ProjectName)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)

		return
	}

	// // Generate a snowflake ID.
	// id, err := util.SnowId()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	id, err := util.NewIdMgr(5)
	if err != nil {
		fmt.Println(err)
		return
	}
	Project_id := util.Int64toString(id.ID())

	project := model.Project{}
	project.ProjectName = body.ProjectName
	project.Describe = body.Description
	project.ProjectID = Project_id

	err = q.WithContext(ctx).Create(&project)

	//ok_msg := fmt.Sprintf("The ProjectName: [%s] Create successfully!", body.Name)

	if err != nil {
		msg := fmt.Sprintf("The ProjectName: [%s] Create Failed!, err: %s", body.ProjectName, err)
		util.ReturnMsg(ctx, http.StatusNotFound, "", msg)
	} else {
		ok_msg := fmt.Sprintf("The Project_Name: [%s] Create successfully!", body.ProjectName)
		util.ReturnMsg(ctx, http.StatusOK, project, ok_msg)

	}

}

// @BasePath /api/v1
// PingProject godoc
// @Summary UpdateProject
// @Schemes
// @Description Update Project
// @Tags UpdateProject
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_id query string true "Project_ID"
// @Param ReqeustBody body CreateProjectRequestBody true "Update Project"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /project/update [put]
// @ID UpdateProject
func UpdateProject(ctx *gin.Context) {
	project_id := ctx.Query("project_id")

	body := UpdateProjectRequestBody{}
	err := ctx.ShouldBind(&body)
	if err != nil {
		error_msg := fmt.Sprintf("err: %s ", err)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)

		return
	}

	project := model.Project{}
	project.ProjectName = body.ProjectName
	project.Describe = body.Description

	ClusterBind := model.ClusterBind{}
	ClusterBind.ProjectName = body.ProjectName

	q := query.Use(global.Config.DB.DB())

	_, err = q.ClusterBind.WithContext(ctx).Select(q.ClusterBind.ClusterID).Where(q.ClusterBind.ProjectID.Eq(project_id)).First()
	info := gen.ResultInfo{}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		info, err = q.Project.WithContext(ctx).Where(q.Project.ProjectID.Eq(project_id)).Updates(&project)

	} else {
		if err != nil {
			error_msg := fmt.Sprintf("err: %s ", err)
			util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)

			return
		}

		q.Transaction(func(tx *query.Query) error {
			//if err := tx.Project.WithContext(ctx).Create(&project); err != nil {
			if info, err = tx.Project.WithContext(ctx).Where(tx.Project.ProjectID.Eq(project_id)).Updates(&project); err != nil {
				return err
			}
			if info, err = tx.ClusterBind.WithContext(ctx).Where(tx.ClusterBind.ProjectID.Eq(project_id)).Updates(&ClusterBind); err != nil {
				return err
			}
			return nil
		})
	}
	if err != nil {
		error_msg := fmt.Sprintf("err: %s ", err)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return
	}

	if info.RowsAffected == 0 {
		msg := "The ProjectID is not Found!"
		util.ReturnMsg(ctx, http.StatusNotFound, "", msg)

		return
	}

	ok_msg := fmt.Sprintf("The Project_id: [%s] Update successfully!", project_id)
	util.ReturnMsg(ctx, http.StatusOK, util.Int64toString(info.RowsAffected), ok_msg)

}

// @BasePath /api/v1
// PingProject godoc
// @Summary DeleteProject
// @Schemes
// @Description Delete Project
// @Tags DeleteProject
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param project_name query string true "Project_Name"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /project/delete [delete]
// @ID DeleteProject
func DeleteProject(ctx *gin.Context) {
	project_name := ctx.Query("project_name")

	q := query.Use(global.Config.DB.DB())

	_, err := q.ClusterBind.WithContext(ctx).Select(q.ClusterBind.ClusterID).Where(q.ClusterBind.ProjectName.Eq(project_name)).First()

	info := gen.ResultInfo{}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if info, err = q.Project.WithContext(ctx).Where(q.Project.ProjectName.Eq(project_name)).Delete(); err != nil {
			if err != nil {
				error_msg := fmt.Sprintf("err: %s ", err)
				util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
				return
			}
		}

		if info.RowsAffected == 0 {
			error_msg := fmt.Sprintf("The Project_Name: [%s] is not Found!", project_name)
			util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
			return
		}

		ok_msg := fmt.Sprintf("The Project_Name: [%s] deleted successfully!!", project_name)
		util.ReturnMsg(ctx, http.StatusOK, util.Int64toString(info.RowsAffected), ok_msg)
		return

	} else {
		error_msg := fmt.Sprintf("The Project_Name: [%s] 存在绑定的集群，无法删除!", project_name)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return

	}

}
