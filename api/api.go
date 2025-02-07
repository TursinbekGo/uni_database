package api

import (
	_ "app/api/docs"
	"app/api/handler"
	"app/config"
	"app/pkg/logger"
	"app/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewApi(r *gin.Engine, cfg *config.Config, storage storage.StorageI, logger logger.LoggerI) {

	// @securityDefinitions.apikey ApiKeyAuth
	// @in header
	// @name Authorization

	handler := handler.NewHandler(cfg, storage, logger)

	r.Use(customCORSMiddleware())
	// v1 := r.Group("/v1")
	// v1.Use(handler.AuthMiddleware())

	r.POST("/login", handler.Login)
	r.POST("/register", handler.Register)

	//uploading
	// r.Static("/uploads", "./uploads")
	// r.GET("/download/:filename", handler.DownloadFileHandler)
	// r.POST("/uploadd", handler.UploadHandler)
	// r.GET("/images", handler.ListImagesHandler)
	// r.GET("/image/:filename", handler.GetImageHandler)

	//SisAmin
	r.POST("/admin", handler.CreateAdmin)
	r.GET("/admin/:id", handler.GetByIdAdmin)
	r.GET("/admin", handler.GetListAdmin)
	r.PUT("/admin/:id", handler.UpdateAdmin)
	r.DELETE("/admin/:id", handler.DeleteAdmin)

	// User
	// r.POST("/user", handler.CreateUser)
	r.GET("/user/:id", handler.GetByIdUser)
	r.GET("/user", handler.GetListUser)
	r.PUT("/user/:id", handler.UpdateUser)
	r.DELETE("/user/:id", handler.DeleteUser)
	r.GET("/get_user_activity_counts", handler.GetUserActivityCounts)
	r.GET("/top_contributors", handler.GetTopContributors)
	r.GET("/users/scores", handler.GetUserScores)
	r.GET("/users/{user_id}/rank", handler.GetUserRank)
	r.GET("/users/statistics", handler.GetUserStatistics)

	// Course
	r.POST("/course", handler.CreateCourse)
	r.GET("/course/:id", handler.GetByIdCourse)
	r.GET("/course", handler.GetListCourse)
	r.PUT("/course/:id", handler.UpdateCourse)
	r.DELETE("/course/:id", handler.DeleteCourse)
	r.GET("/courses_by_semester_id", handler.GetListCoursesBySemesterId)

	// Semester
	r.POST("/semester", handler.CreateSemester)
	r.GET("/semester/:id", handler.GetByIdSemester)
	r.GET("/semester", handler.GetListSemester)
	r.PUT("/semester/:id", handler.UpdateSemester)
	r.DELETE("/semester/:id", handler.DeleteSemester)
	// Like
	r.POST("/like", handler.CreateLike)
	r.GET("/like/:id", handler.GetByIdLike)
	r.GET("/like", handler.GetListLike)
	r.PUT("/like/:id", handler.UpdateLike)
	r.DELETE("/like/:id", handler.DeleteLike)
	// Download
	r.POST("/download", handler.CreateDownload)
	r.GET("/download/:id", handler.GetByIdDownload)
	r.GET("/download", handler.GetListDownload)
	r.PUT("/download/:id", handler.UpdateDownload)
	r.DELETE("/download/:id", handler.DeleteDownload)
	// Publication
	r.POST("/publication", handler.CreatePublication)
	r.GET("/publication/:id", handler.GetByIdPublication)
	r.GET("/publication", handler.GetListPublication)
	r.PUT("/publication/:id", handler.UpdatePublication)
	r.DELETE("/publication/:id", handler.DeletePublication)
	r.GET("/get_publication_stats", handler.GetPublicationStats)
	r.GET("/publications/tags", handler.GetPublicationsByTag)

	// Notification
	r.POST("/notification", handler.CreateNotification)
	r.GET("/notification/:id", handler.GetByIdNotification)
	r.GET("/notification", handler.GetListNotification)
	r.PUT("/notification/:id", handler.UpdateNotification)
	r.DELETE("/notification/:id", handler.DeleteNotification)
	///////////////////////////////////////////////////

	// Worker
	// r.POST("/worker", handler.AuthMiddleware(), handler.CreateWorker)
	// r.GET("/worker/:id", handler.AuthMiddleware(), handler.GetByIdWorker)
	// r.GET("/worker", handler.AuthMiddleware(), handler.GetListWorker)
	// r.PUT("/worker/:id", handler.AuthMiddleware(), handler.UpdateWorker)
	// r.DELETE("/worker/:id", handler.AuthMiddleware(), handler.DeleteWorker)
	// r.PATCH("/worker/:id", handler.AuthMiddleware(), handler.PatchWorker)
	//Login
	// r.POST("/login", handler.Login)

	// r.GET("/petrol_history", handler.AuthMiddleware(), handler.GetListPetrolHistory)
	// r.GET("/petrol_history_by_car_id", handler.AuthMiddleware(), handler.GetListPetrolHistoryByCarId)
	// r.GET("/users_by_company_id", handler.AuthMiddleware(), handler.GetListUserByCompanyId)

	// // r.POST("/upload", handler.UploadPhoto)   TO FIrebase
	// r.Static("/uploads", "./uploads")
	r.POST("/uploadd", handler.UploadHandler)
	r.POST("/upload_profile", handler.UploadHandlerProfile)
	r.GET("/images", handler.ListImagesHandler)
	r.GET("/image/:filename", handler.GetImageHandler)
	r.GET("/profile_image/:filename", handler.GetProfileImageHandler)
	// r.GET("/petrol_history_sum_remaining", handler.AuthMiddleware(), handler.GetSumRemainingPetrolByMonth)
	// r.GET("/petrol_history_sum_remaining_by_year", handler.AuthMiddleware(), handler.GetSumRemainingPetrolByYear)
	r.GET("/file_download/:filename", handler.FileDownloadFileHandler)
	// r.GET("/petrol_history/by_car", handler.GetListPetrolHistoryByCarID)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

func customCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS, HEAD")
		c.Header("Access-Control-Allow-Headers", "Platform-Id, Content-Type, Accesp-Encoding, Authorization, Cache-Control")
		c.Header("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
