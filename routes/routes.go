package routes

import (
    "github.com/gin-gonic/gin"
    "issue-server/controllers"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.POST("/issue", controllers.CreateIssue)
    r.GET("/issues", controllers.ListIssues)
    r.GET("/issue/:id", controllers.GetIssue)
    r.PATCH("/issue/:id", controllers.UpdateIssue)

    return r
}
