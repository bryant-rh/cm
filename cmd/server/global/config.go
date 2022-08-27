package global

import (
	"github.com/bryant-rh/cm/internal/model"
	"github.com/bryant-rh/cm/internal/query"
	"context"
	"errors"

	"github.com/kunlun-qilian/confmysql"
	"github.com/kunlun-qilian/confserver"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	confserver.SetServiceName("cm-server", "..")
	confserver.ConfP(&Config)

	// migrate 数据库 命令
	confserver.AddCommand(Config.DB.Commands()...)
	InitUser(context.Background(), Config.Default_UserName, Config.Default_PassWord)
}

// 需要migrate 到数据库的信息
var AutoMigrateModelList = []interface{}{
	model.User{},
	model.Cluster{},
	model.Project{},
	model.Label{},
	model.ClusterBind{},
	model.Namespace{},
	model.Serviceaccount{},
}

var Config = struct {
	DB     *confmysql.MySQL
	Server *confserver.Server

	TestEnvStr       string `env:""`
	Default_UserName string
	Default_PassWord string
}{
	Server: &confserver.Server{
		Mode: "debug",
	},

	DB: &confmysql.MySQL{
		DSN: "root:a8EyHxVuaZeS9J@tcp(127.0.0.1:3306)/cluster_mgr?charset=utf8mb4&parseTime=True&loc=Local",
		AutoMigrateConfig: &confmysql.AutoMigrateConfig{
			Models:    AutoMigrateModelList,
			ModelPath: "./internal/model",
			QueryPath: "./internal/query",
		},
		Config: &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	},
	TestEnvStr:       "global.config",
	Default_UserName: "admin",
	Default_PassWord: "123456",
}

func InitUser(ctx context.Context, username, password string) error {
	if username == "" && password == "" {
		logrus.Warningf("No default User: [%s/%s]", username, password)
		return nil
	}
	q := query.Use(Config.DB.DB()).User
	_, err := q.WithContext(ctx).Where(q.Username.Eq(username)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user := model.User{}
		//加密密码
		hashPwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user.Username = username

		user.Password = string(hashPwd)

		err = q.WithContext(ctx).Create(&user)
		if err != nil {
			return err
		}
		logrus.Infof("Create default User: [%s/%s]", username, password)
	}
	return nil
}
