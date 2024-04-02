package main

import (
	"bufio"
	"context"
	"time"

	//"crypto/bcrypt"
	"encoding/json"
	"env-mngr/backend/db"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	//"golang.org/x/crypto/bcrypt"
)

type entry struct {
	Key   string
	Entry string
}

var secretKey = []byte("secret-key")

func main() {
	app := fiber.New()
	app.Use(cors.New())

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://brad:12345678@localhost:5432/env-manager-v1")
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close(ctx)
	queries := db.New(conn)
	app.Get("/", func(c *fiber.Ctx) error {
		log.Info("/")
		return c.SendString("Hello, World!")
	})
	app.Post("/api/signin", func(c *fiber.Ctx) error {
		log.Info("/api/signin")
		headersCopy := c.GetReqHeaders()
		authHeaderVal := headersCopy["authorization"]

		if authHeaderVal[0] != "" {
			ok, err := authenticate(c.FormValue("email"), c.FormValue("password"), queries, ctx)
			if err != nil {

				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "msg": err.Error()})
			}
			if ok {
				err := verifyToken(authHeaderVal[0])
				if err != nil {

					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "msg": err.Error()})
				}

				return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"error": false, "msg": "user authenticated"})
			}
		} else {
			ok, err := authenticate(c.FormValue("email"), c.FormValue("password"), queries, ctx)
			if err != nil {

				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "msg": err.Error()})
			}
			if ok {
				token, err := createToken(c.FormValue("email"))
				if err != nil {

					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "msg": err.Error()})
				}
				return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"error": false, "token": token, "msg": "user authenticated"})

			}
		}
		return nil

	})
	app.Post("/api/create-user", func(c *fiber.Ctx) error {
		log.Info("/api/create-user")
		password := []byte(c.FormValue("password"))
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		insertedUser, err := queries.CreateUser(ctx, db.CreateUserParams{Password: string(hashedPassword), Email: c.FormValue("email")})

		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "msg": err.Error()})
		}
		if insertedUser.Email == c.FormValue("email") {
			log.Error(err)
			return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"error": false, "msg": "user created successfully"})
		}
		return nil
	})
	app.Post("/api/upload", func(c *fiber.Ctx) error {
		log.Info("/api/upload")
		// var envContainer []entry
		envContainer := make(map[string]string)
		file, err := c.FormFile("envFile")
		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "msg": err.Error()})
		}
		dat, err := file.Open()
		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "msg": err.Error()})
		}
		defer dat.Close()
		scanner := bufio.NewScanner(dat)
		for scanner.Scan() {
			fmt.Printf("env var: %s\n", scanner.Text())
			slicedEnvEntryLine := strings.Split(scanner.Text(), "=")
			if len(slicedEnvEntryLine) > 1 && string(slicedEnvEntryLine[0][0]) != "#" {
				envContainer[slicedEnvEntryLine[0]] = slicedEnvEntryLine[1]
				// envContainer = append(envContainer, entry{Key: slicedEnvEntryLine[0], Entry: slicedEnvEntryLine[1]})
			}
		}
		log.Info(envContainer)
		if len(envContainer) != 0 {
			jsonEnvBucket, err := json.Marshal(envContainer)
			log.Info(string(jsonEnvBucket))
			if err != nil {
				log.Error(err)
				return c.JSON(fiber.Map{"error": true, "msg": err})
			}
			insertedBucket, err := queries.CreateBucket(ctx, db.CreateBucketParams{Name: "Some bucket", Envvariables: jsonEnvBucket})
			if err != nil {

				return c.JSON(fiber.Map{"error": true, "msg": err})
			}

			return c.JSON(fiber.Map{"error": false, "msg": nil, "info": fmt.Sprintf("added %d variables for bucket %s ", len(envContainer), insertedBucket.Name)})
		} else {
			return c.JSON(fiber.Map{"error": true, "msg": "no environment variables added"})
		}
	})
	app.Post("/api/generateEnvFile", func(c *fiber.Ctx) error {
		bucketId, err := strconv.Atoi(c.FormValue("bucketId"))
		envContainer := make(map[string]string)

		if err != nil {
			log.Error(err)
		}
		bucket, err := queries.GetBucket(ctx, int32(bucketId))
		if err := json.Unmarshal(bucket.Envvariables, &envContainer); err != nil {
			log.Error(err)
		} else {
			file, err := os.Create("/tmp/env-manager/.env")
			if err != nil {
				return c.JSON(fiber.Map{"error": true, "msg": err})
			}
			defer file.Close()
			for key, entry := range envContainer {
				line := fmt.Sprintf("%s=%s", key, entry)
				_, err := file.WriteString(line)

				if err != nil {
					return c.JSON(fiber.Map{"error": true, "msg": err})
				}

			}
			return c.Download("/tmp/env-manager/.env")

		}

		log.Info(envContainer)
		if err != nil {
			return c.JSON(fiber.Map{"error": true, "msg": err})
		}
		return c.JSON(fiber.Map{"error": false, "msg": nil, "info": fmt.Sprintf("found")})
	})
	app.Listen(":3002")

}
func authenticate(email string, password string, queries *db.Queries, ctx context.Context) (bool, error) {
	user, err := queries.GetUser(ctx, email)
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {

		log.Error(err)
		return false, err
	} else {
		return true, nil
	}

}
func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
