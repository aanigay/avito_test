package router

import (
	"api/internal/segment"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(segmentsHandler *segment.Handler) {
	r = gin.Default()

	r.Use(CORS())

	r.POST("/createSegment", segmentsHandler.CreateSegment)
	r.POST("/createSegmentsWithUsers", segmentsHandler.CreateSegmentWithUsers)
	r.DELETE("/deleteSegment", segmentsHandler.DeleteSegment)
	r.POST("updateSegments", segmentsHandler.UpdateSegments)
	r.GET("/getSegmentsByUserId/:userId", segmentsHandler.GetSegmentsByUserId)
	r.GET("/report/:date", segmentsHandler.GetReports)
	r.GET("/download/:filepath", segmentsHandler.DownloadFile)
}

func Start(addr string) error {
	return r.Run(addr)
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
