package routers

import (
	"github.com/bryant-rh/cm/cmd/server/routers/cluster"
	"github.com/bryant-rh/cm/cmd/server/routers/label"
	"github.com/bryant-rh/cm/cmd/server/routers/project"
	"github.com/bryant-rh/cm/cmd/server/routers/serviceaccount"
	"github.com/bryant-rh/cm/cmd/server/routers/user"
	"github.com/bryant-rh/cm/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func NewRooter(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	user.UserRouter(v1)
	v1.Use(middleware.JWTAuthMiddleware())
	{
		project.ProjectRouter(v1)
		cluster.ClusterRouter(v1)
		label.LabelRouter(v1)
		serviceaccount.SaRouter(v1)

	}

}
