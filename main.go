package main

import (
	"api-golang/auth"
	"api-golang/middleware"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"reflect"
	// package used to read the .env file
	_ "github.com/lib/pq" // postgres golang driver

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

type Student struct {
	Student_id      uint64 `json:"student_id" binding:"required"`
	Student_name    string `json:"student_name" binding:"required"`
	Student_age     uint8  `json:"student_age" binding:"required"`
	Student_address string `json:"student_address" binding:"required"`
	Student_phone   string `json:"student_phone" binding:"required"`
}

func rowToStruct(rows *sql.Rows, dest interface{}) error {
	destv := reflect.ValueOf(dest).Elem()
	args := make([]interface{}, destv.Type().Elem().NumField())
	for rows.Next() {
		rowp := reflect.New(destv.Type().Elem())
		rowv := rowp.Elem()
		for i := 0; i < rowv.NumField(); i++ {
			args[i] = rowv.Field(i).Addr().Interface()
		}
		if err := rows.Scan(args...); err != nil {
			return err
		}
		destv.Set(reflect.Append(destv, rowv))
	}
	return nil
}

func postHandler(c *gin.Context, db *gorm.DB) {
	var student Student

	if err := c.Bind(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := db.Create(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data Student Berhasil ditambahkan",
		"data":    student,
	})

}

func getAllHandler(c *gin.Context, db *gorm.DB) {
	var students []Student
	if err := db.Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if len(students) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data tidak ada",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data Student",
		"data":    students,
	})
}

func getHandler(c *gin.Context, db *gorm.DB) {
	var student Student
	studentId := c.Param("student_id")
	if db.Find(&student, "student_id =?", studentId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data tidak ada",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data Student",
		"data":    student,
	})
}

func putHandler(c *gin.Context, db *gorm.DB) {
	var reqStudent Student
	studentId := c.Param("student_id")

	if db.Find(&reqStudent, "student_id = ?", studentId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data tidak ada",
		})
		return
	}

	if err := c.Bind(&reqStudent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := db.Model(&reqStudent).Where("student_id = ?", studentId).Updates(&reqStudent).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data Student Berhasil diupdate",
		"data":    reqStudent,
	})
}

func deleteHandler(c *gin.Context, db *gorm.DB) {
	var student Student
	studentId := c.Param("student_id")
	if db.Find(&student, "student_id =?", studentId).RecordNotFound() {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data tidak ada",
		})
		return
	}
	if err := db.Delete(&student, "student_id =?", studentId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data berhasil dihapus",
	})
}

func setupRouter() *gin.Engine {
	conn := "postgres://postgres:Haruhiko_123@127.0.0.1:5432/api_golang?sslmode=disable"
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatal("Error load env")
	}
	conn = os.Getenv("POSTGRES_URL")
	DB, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
	Migrate(DB)
	r := gin.Default()
	r.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "Welcome",
		})
	})
	r.POST("/login", auth.LoginHandler)
	v1 := r.Group("api/v1")
	v1.POST("/student", func(context *gin.Context) {
		postHandler(context, DB)
	})
	//
	v1.GET("/students", middleware.AuthValid, func(context *gin.Context) {
		getAllHandler(context, DB)
	})

	v1.GET("/student/:student_id", func(context *gin.Context) {
		getHandler(context, DB)
	})

	v1.PUT("/student/:student_id", func(context *gin.Context) {
		putHandler(context, DB)
	})

	v1.DELETE("/student/:student_id", func(context *gin.Context) {
		deleteHandler(context, DB)
	})

	return r
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Student{})

	data := Student{}
	if db.Find(&data).RecordNotFound() {
		fmt.Println("===Menjalankan seeder user===")
		seederUser(db)
	}
}

func seederUser(db *gorm.DB) {
	data := Student{
		Student_id:      1,
		Student_name:    "Jovi",
		Student_address: "Bandung",
		Student_age:     21,
		Student_phone:   "081212",
	}
	db.Create(&data)
}

func main() {
	r := setupRouter()
	r.Run(":666")

}
