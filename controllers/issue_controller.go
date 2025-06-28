package controllers

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "issue-server/models"
)

var users = []models.User{
    {ID: 1, Name: "김개발"},
    {ID: 2, Name: "이디자인"},
    {ID: 3, Name: "박기획"},
}

var issues []models.Issue
var nextIssueID uint = 1

func findUserByID(id uint) *models.User {
    for i := range users {
        if users[i].ID == id {
            return &users[i]
        }
    }
    return nil
}

func isValidStatus(s string) bool {
    switch s {
    case "PENDING", "IN_PROGRESS", "COMPLETED", "CANCELLED":
        return true
    default:
        return false
    }
}

func CreateIssue(c *gin.Context) {
    var req struct {
        Title       string `json:"title" binding:"required"`
        Description string `json:"description"`
        UserID      *uint  `json:"userId"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "필수 파라미터 누락 또는 잘못된 형식", "code": 400})
        return
    }

    var assignedUser *models.User
    if req.UserID != nil {
        assignedUser = findUserByID(*req.UserID)
        if assignedUser == nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "존재하지 않는 사용자", "code": 400})
            return
        }
    }

    status := "PENDING"
    if assignedUser != nil {
        status = "IN_PROGRESS"
    }

    issue := models.Issue{
        ID:          nextIssueID,
        Title:       req.Title,
        Description: req.Description,
        Status:      status,
        User:        assignedUser,
        CreatedAt:   time.Now().UTC(),
        UpdatedAt:   time.Now().UTC(),
    }
    nextIssueID++
    issues = append(issues, issue)

    c.JSON(http.StatusCreated, issue)
}

func ListIssues(c *gin.Context) {
    statusFilter := c.Query("status")
    if statusFilter != "" && !isValidStatus(statusFilter) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 상태값", "code": 400})
        return
    }

    var result []models.Issue
    for _, is := range issues {
        if statusFilter == "" || is.Status == statusFilter {
            result = append(result, is)
        }
    }
    c.JSON(http.StatusOK, gin.H{"issues": result})
}

func GetIssue(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID", "code": 400})
        return
    }
    for _, is := range issues {
        if is.ID == uint(id) {
            c.JSON(http.StatusOK, is)
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "이슈를 찾을 수 없습니다", "code": 404})
}

func UpdateIssue(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID", "code": 400})
        return
    }

    var req struct {
        Title       *string `json:"title"`
        Description *string `json:"description"`
        Status      *string `json:"status"`
        UserID      *uint   `json:"userId"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 요청", "code": 400})
        return
    }

    for i := range issues {
        if issues[i].ID == uint(id) {
            is := &issues[i]

            // 5.3 ‑ COMPLETED or CANCELLED → immutable
            if is.Status == "COMPLETED" || is.Status == "CANCELLED" {
                c.JSON(http.StatusForbidden, gin.H{"error": "완료/취소된 이슈는 수정할 수 없습니다", "code": 403})
                return
            }

            // userId changes
            if req.UserID != nil {
                if *req.UserID == 0 {
                    // 담당자 제거
                    is.User = nil
                    is.Status = "PENDING"
                } else {
                    user := findUserByID(*req.UserID)
                    if user == nil {
                        c.JSON(http.StatusBadRequest, gin.H{"error": "존재하지 않는 사용자", "code": 400})
                        return
                    }
                    is.User = user
                    if is.Status == "PENDING" && req.Status == nil {
                        is.Status = "IN_PROGRESS"
                    }
                }
            }

            // status change
            if req.Status != nil {
                if !isValidStatus(*req.Status) {
                    c.JSON(http.StatusBadRequest, gin.H{"error": "유효하지 않은 상태", "code": 400})
                    return
                }
                if (*req.Status == "IN_PROGRESS" || *req.Status == "COMPLETED") && is.User == nil {
                    c.JSON(http.StatusBadRequest, gin.H{"error": "담당자 없이 해당 상태로 변경 불가", "code": 400})
                    return
                }
                is.Status = *req.Status
            }

            if req.Title != nil {
                is.Title = *req.Title
            }
            if req.Description != nil {
                is.Description = *req.Description
            }

            is.UpdatedAt = time.Now().UTC()
            c.JSON(http.StatusOK, *is)
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "이슈를 찾을 수 없습니다", "code": 404})
}
