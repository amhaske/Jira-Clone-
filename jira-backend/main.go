package main

import (
	"fmt"
	"jira-backend/dbutils"
	"jira-backend/handlers"
	"jira-backend/utils"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error
var config *viper.Viper

func corsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

func setConfigInContext(c *gin.Context) {
	c.Set("config", config)
	c.Next()
}

func main() {

	// Migrate the schema
	dbutils.InitializeConn()

	config, err = utils.ConfigReader()
	if err != nil {
		fmt.Println("Cannot Read from the Config", err)
		os.Exit(1)
	}

	r := gin.Default()
	r.Use(setConfigInContext)
	r.Use(cors.Default())

	project := r.Group("/api/project")
	{
		project.POST("/list", handlers.ListProjects)
		project.POST("/info", handlers.GetProjectInfo)
		project.POST("/delete", handlers.DeleteProject)
		project.POST("/create", handlers.CreateProject)
		project.POST("/members", handlers.ListMembers)
	}

	sprint := r.Group("/api/sprint")
	{
		sprint.POST("/list", handlers.ListSprints)
		sprint.POST("/info", handlers.GetSprintInfo)
		sprint.POST("/create", handlers.CreateSprint)
		sprint.POST("/delete", handlers.DeleteSprint)
		sprint.POST("/update", handlers.UpdateSprint)
	}

	issue := r.Group("/api/issue")
	{
		issue.POST("/info", handlers.GetIssueInfo)
		issue.POST("/create", handlers.CreateIssue)
		issue.POST("/delete", handlers.DeleteIssue)
		issue.POST("/update", handlers.UpdateIssue)
		issue.POST("/list", handlers.ListIssuesForSprint)
		issue.POST("/move", handlers.UpdateIssueStatus)
	}

	ip_address := fmt.Sprintf("%s:%d", config.GetString("server.ip_address"), config.GetInt("server.port"))
	r.Run(ip_address)

}
