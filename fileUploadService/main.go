package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pkg/errors"
	helpers "github.com/saiteja111997/throttle_backend/pkg/helper"
	"github.com/saiteja111997/throttle_backend/pkg/server"
	"github.com/saiteja111997/throttle_backend/pkg/structures"
)

var fiberLambda *fiberadapter.FiberLambda

func main() {
	app := fiber.New()

	// Use the CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	//ctx := context.Background()
	if err := waitForHost("mydbinstance.c1cnaivzlk0f.us-east-1.rds.amazonaws.com", "5432"); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Connection established")

	db, err := helpers.Open(helpers.Config{
		Username: "postgres",
		Password: "Saiteja11",
		Hostname: "mydbinstance.c1cnaivzlk0f.us-east-1.rds.amazonaws.com",
		Port:     "5432",
		Database: "postgres",
	})

	if err != nil {
		log.Println(err)
		return
	}

	db.AutoMigrate(&structures.Errors{}, &structures.User{}, structures.UserActions{})

	defer db.Close()

	svr := server.Server{
		Db: db,
	}

	app.Get("/ping", svr.HealthCheck)
	app.Post("/file_upload/upload_error", svr.UploadError)
	app.Post("/generate_document", svr.GenerateDocument)
	app.Post("/file_upload/user_action", svr.InsertUserActions)
	app.Post("/editing/error", svr.GetRawErrorDocs)
	app.Post("/editing/images", svr.GetImagesFromS3)

	fmt.Println("Routing established!!")

	if helpers.IsLambda() {
		fiberLambda = fiberadapter.New(app)
		lambda.Start(Handler)
	} else {
		app.Listen(":3000")
	}

}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Proxy the request to the Fiber app and get the response
	response, err := fiberLambda.ProxyWithContext(ctx, request)

	response.Headers = make(map[string]string)

	// Add CORS headers to the response
	response.Headers["Access-Control-Allow-Origin"] = "*"
	response.Headers["Access-Control-Allow-Methods"] = "GET,POST,PUT,DELETE"
	response.Headers["Access-Control-Allow-Headers"] = "Origin, Content-Type, Accept"
	response.Headers["Access-Control-Allow-Credentials"] = "true"

	return response, err
}

func waitForHost(host, port string) error {
	timeOut := time.Second

	if host == "" {
		return errors.Errorf("unable to connect to %v:%v", host, port)
	}

	for i := 0; i < 60; i++ {
		fmt.Printf("waiting for %v:%v ...\n", host, port)
		conn, err := net.DialTimeout("tcp", host+":"+port, timeOut)
		if err == nil {
			fmt.Println("done!")
			conn.Close()
			return nil
		}

		time.Sleep(time.Second)
	}

	return errors.Errorf("timeout attempting to connect to %v:%v", host, port)
}
