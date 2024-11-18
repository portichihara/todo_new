package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var db *gorm.DB

type User struct {
    gorm.Model
    Username string `json:"username"`
}

type Todo struct {
    gorm.Model
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      string `json:"status"`
    Priority    int    `json:"priority"`
    UserID      uint   `json:"user_id"`
}

func main() {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    var err error
    for i := 0; i < 10; i++ {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            break
        }
        log.Printf("Failed to connect to database, retrying in 5 seconds... (%d/10)\n", i+1)
        time.Sleep(5 * time.Second)
    }
    if err != nil {
        log.Fatal("Failed to connect to database", err)
    }

    db.AutoMigrate(&User{}, &Todo{})

    createInitialUser()

    r := setupRouter()

    r.LoadHTMLFiles("templates/index.html")

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    r.Run(":" + port)
}

func createInitialUser() {
    user := User{
        Username: "defaultuser",
    }
    result := db.FirstOrCreate(&user, User{Username: "defaultuser"})
    if result.Error != nil {
        log.Fatalf("Failed to create initial user: %v", result.Error)
    }
}

func setupRouter() *gin.Engine {
    r := gin.Default()

    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", nil)
    })

    api := r.Group("/api/v1")
    {
        api.POST("/users", createUser)
        api.POST("/todos", createTodo)
        api.GET("/todos", listTodos)
        api.PUT("/todos/:id", updateTodo)
        api.DELETE("/todos/:id", deleteTodo)
    }

    return r
}

func createUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result := db.Create(&user)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusCreated, user)
}

func createTodo(c *gin.Context) {
    var todo Todo
    if err := c.ShouldBindJSON(&todo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if todo.UserID == 0 {
        var user User
        result := db.First(&user, "username = ?", "defaultuser")
        if result.Error != nil {
            log.Fatalf("Failed to find user: %v", result.Error)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
            return
        }
        todo.UserID = user.ID
    }

    todo.Status = "未着手"
    result := db.Create(&todo)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusCreated, todo)
}

func listTodos(c *gin.Context) {
    var todos []Todo
    title := c.Query("title")
    query := db.Model(&Todo{})
    if title != "" {
        query = query.Where("title LIKE ?", "%"+title+"%")
    }

    result := query.Find(&todos)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"todos": todos})
}

func updateTodo(c *gin.Context) {
    var todo Todo
    id := c.Param("id")

    result := db.First(&todo, id)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    if result.RowsAffected == 0 {
        c.Status(http.StatusNotFound)
        return
    }

    if err := c.ShouldBindJSON(&todo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result = db.Save(&todo)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}

func deleteTodo(c *gin.Context) {
    var todo Todo
    id := c.Param("id")

    result := db.First(&todo, id)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    if result.RowsAffected == 0 {
        c.Status(http.StatusNotFound)
        return
    }

    result = db.Delete(&todo)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}