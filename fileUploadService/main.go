package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
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

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables file")
	}

	//ctx := context.Background()
	if err := waitForHost("mydbinstance.c1cnaivzlk0f.us-east-1.rds.amazonaws.com", "5432"); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Connection established")

	DB_USERNAME := os.Getenv("DB_USERNAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOSTNAME := os.Getenv("DB_HOSTNAME")
	DB_PORT := os.Getenv("DB_PORT")
	DATABASE := os.Getenv("DATABASE")

	db, err := helpers.Open(helpers.Config{
		Username: DB_USERNAME,
		Password: DB_PASSWORD,
		Hostname: DB_HOSTNAME,
		Port:     DB_PORT,
		Database: DATABASE,
	})

	if err != nil {
		log.Println(err)
		return
	}

	db.AutoMigrate(&structures.Errors{}, &structures.Users{}, structures.UserActions{})

	defer db.Close()

	svr := server.Server{
		Db: db,
	}

	app.Get("/ping", svr.HealthCheck)
	app.Post("/file_upload/upload_error", svr.UploadError)
	app.Get("/file_upload/get_latest_unsolved", svr.GetUnresolvedJourneys)
	app.Post("/file_upload/update_error_state", svr.UpdateErrorState)
	app.Post("file_upload/update_final_state", svr.UpdateFinalState)
	app.Post("/generateDocument", svr.GenerateDocument)
	app.Post("/getDashboard", svr.getDashboard)
	app.Post("/file_upload/user_action", svr.InsertUserActions)
	app.Post("/file_upload/delete_user_action", svr.DeleteUserAction)
	app.Post("/file_upload/validate_user_action", svr.ValidateUserAction)
	app.Post("/editing/error", svr.GetRawErrorDocs)
	app.Post("/preDocEdit/getLatestErrorRaw", svr.GetLatestRawError)
	app.Post("/editing/images", svr.GetImagesFromS3)
	app.Post("/auth/login", svr.Login)
	app.Post("/auth/register", svr.Register)
	app.Get("/oauth2/login", svr.HandleLogin)
	app.Get("/oauth2/callback", svr.HandleCallback)

	fmt.Println("Routing established!!")

	if helpers.IsLambda() {
		fiberLambda = fiberadapter.New(app)
		lambda.Start(Handler)
	} else {
		fmt.Println("Starting server locally!!")
		err = app.Listen(":8080")

		if err != nil {
			fmt.Println("An error occured while starting the server : ", err)
		}
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
