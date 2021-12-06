package router

import (
	"github.com/gin-gonic/gin"
	"lab-platform/handler/check"
	"lab-platform/router/middleware"
	"net/http"
)

// Load loads the middleware's, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middleware's
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// pprof router
	//pprof.Register(g)

	//clusterGroup := g.Group("/qbus/clusters")
	//{
	//	clusterGroup.POST("/addcluster", cluster.AddCluster)
	//	clusterGroup.GET("/deletecluster", cluster.DeleteCluster)
	//	clusterGroup.GET("/listall", cluster.ListAllCluster)
	//	clusterGroup.GET("/getclusterdetail", cluster.GetClusterDetail)
	//	clusterGroup.GET("/getclusterdiskinfo", cluster.GetClusterDiskInfo)
	//}

	//topicGroup := g.Group("/qbus/topics")
	//{
	//	topicGroup.POST("/addtopic", topic.Create)
	//	topicGroup.GET("/deletetopic", topic.Delete)
	//	topicGroup.GET("/gettopicdetail", topic.GetTopicDetail)
	//	topicGroup.POST("/addpartition", topic.AddPartition)
	//	topicGroup.POST("/altertopic", topic.AlterTopic)
	//	topicGroup.POST("/settopicmaxsize", topic.SetTopicMaxSize)
	//	topicGroup.GET("/getoveragetopiclist", topic.GetOverRageTopicList)
	//	topicGroup.GET("/getconsumerandproducerlist", topic.GetCPListOfTopic)
	//	topicGroup.POST("/registerproducerandconsumer", topic.RegisterProduceAndConsumer)
	//}
	//
	// The check check handlers
	systemCheck := g.Group("/check")
	{
		systemCheck.GET("/health", check.HealthCheck)
		systemCheck.GET("/disk", check.DiskCheck)
		systemCheck.GET("/cpu", check.CPUCheck)
		systemCheck.GET("/ram", check.RAMCheck)
	}
	//
	//consumerGroup := g.Group("/qbus/consumer")
	//{
	//	consumerGroup.POST("/getconsumerdetail", consumer.GetConsumerDetail)
	//	consumerGroup.POST("/resetconsumeroffset", consumer.ResetConsumerDetail)
	//}

	return g
}

