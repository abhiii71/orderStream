package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/abhiii71/orderStream/graphql/config"
	"github.com/abhiii71/orderStream/graphql/graph"
	"github.com/abhiii71/orderStream/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func main() {

	server, err :=  graph.NewGraphQLServer(config.AccountUrl, config.ProductUrl, config.OrderUrl, config.PaymentUrl, config.RecommenderUrl)
	if err != nil {
		log.Fatal(err)
	}

	serv := handler.New(server.ToExecutableSchema())
	serv.AddTransport(transport.POST{})
	serv.AddTransport(transport.MultipartForm{})

	engine := gin.Default()

	engine.Use(middleware.GinContextToContextMiddlware())

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": " Everything okay",
		})
	})

	engine.POST("/graphql",
		middleware.AuthorizeJWT(),
		gin.WrapH(serv),
	)

	engine.GET("/playground", gin.WrapH(playground.Handler("Playground", "/graphql")))
	log.Fatal(engine.Run(":8080"))

}
