package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/qinains/fastergoding"
)

type User struct {
	User_uuid             string
	Email                 string
	Password              string
	Fname                 string
	Lname                 string
	Telephone             string
	Address               string
	Country               string
	ProfileImage          string
	PreferredGames        []Game
	UserWallet            Wallet
	UsersGamesInformation []GameInformationOfUser
}

type Team struct {
	TeamID      string
	TeamName    string
	GameID      string
	UsersInTeam []string
}

type GameInformationOfUser struct {
	GameID     string
	Total_time string
	IGN        string // in game name
	ID         string // in game id
	TeamId     string // which game
}

type Game struct {
	GameID       string
	GameName     string
	GameTeamType []string
	GameLogo     string
}

// Composed in the user
type Wallet struct {
	Wallet_id    string
	Deposit_cash int
	Winning_cash int
	Bonus_cash   int
}

type Transaction struct {
	Transaction_id string
	Wallet_id      string
	Source         string
	Timestamp      string
}

type Refer struct {
	Refer_id          string
	Produce_user_uuid string // who generated this reference
	Validity          string
	Timestamp         string
}

const connectionString = "mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false"
const currentDB = "haexrdb"
const Success = 200
const NotAcceptable = 406

func main() {
	fastergoding.Run()
	var ServerOK = false
	log.Print("> Starting the Haexr Servers...")
	server := fiber.New()

	log.Print("> Server Loaded")

	log.Print("> Connecting to Databases...")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Print("> Connection Failed")
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Print("> Connection attempt failed, Disconnecting.")
		}
	}()

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Print("> Cannot Ping")
	} else {
		ServerOK = true
	}

	if ServerOK {
		fmt.Println("Successfully connected and pinged.")
	}

	// Root API
	server.Get("/", func(c *fiber.Ctx) error {
		if ServerOK {
			return c.Render("./index.html", nil, "")
		}
		return c.SendString("System not OK ...")
	})

	// -----------------------------------------------------------------

	// User SignUp API
	server.Post("/register", func(c *fiber.Ctx) error {
		userData := &User{}
		json.Unmarshal(c.Body(), userData)
		if SignUpUser(client.Database(currentDB), userData) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})
	// User unregister API
	server.Post("/unregister", func(c *fiber.Ctx) error {
		userData := &User{}
		json.Unmarshal(c.Body(), userData)
		if UnRegUser(client.Database(currentDB), userData) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Post("/login", func(c *fiber.Ctx) error {
		userData := &User{}
		json.Unmarshal(c.Body(), userData)
		fmt.Println(userData.Email)
		fmt.Println(userData.Password)
		if FindUser(client.Database(currentDB), userData) {
			// Token will be given
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Post("/logout", func(c *fiber.Ctx) error {
		userData := &User{}
		json.Unmarshal(c.Body(), userData)
		fmt.Println(userData.Email)
		fmt.Println(userData.Password)
		if FindUser(client.Database(currentDB), userData) {
			// Token will be revoked
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Post("/updateuser", func(c *fiber.Ctx) error {
		userData := &User{}
		json.Unmarshal(c.Body(), userData)
		fmt.Println(userData.Email)
		fmt.Println(userData.Password)
		if UpdateUser(client.Database(currentDB), userData) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	// Add Team
	server.Post("/addteam", func(c *fiber.Ctx) error {
		teamData := &Team{}
		json.Unmarshal(c.Body(), teamData)
		if AddTeam(client.Database(currentDB), teamData) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	// Add Team Member
	server.Post("/addteammember/:teamid", func(c *fiber.Ctx) error {
		newMemberData := &User{}
		json.Unmarshal(c.Body(), newMemberData)
		if AddTeamMember(client.Database(currentDB), newMemberData, c.Params("teamid")) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	// Add Users GameInformation
	server.Post("/addgameinformationofuser", func(c *fiber.Ctx) error {
		gameInformationOfUser := &GameInformationOfUser{}
		json.Unmarshal(c.Body(), gameInformationOfUser)
		if AddUsersGameInfo(client.Database(currentDB), gameInformationOfUser) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	// Add Users Game
	server.Post("/addgame", func(c *fiber.Ctx) error {
		gameInfo := &Game{}
		json.Unmarshal(c.Body(), gameInfo)
		if AddGame(client.Database(currentDB), gameInfo) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	// Add Users Game
	server.Post("/addtransaction", func(c *fiber.Ctx) error {
		transactionInfo := &Transaction{}
		json.Unmarshal(c.Body(), transactionInfo)
		if addTransaction(client.Database(currentDB), transactionInfo) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Post("/addreference", func(c *fiber.Ctx) error {
		referenceInfo := &Refer{}
		json.Unmarshal(c.Body(), referenceInfo)
		if addReference(client.Database(currentDB), referenceInfo) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Listen(":3000")
}

func SignUpUser(db *mongo.Database, user *User) bool {
	status := true
	_, err := db.Collection("PersonalDetails").InsertOne(
		context.TODO(), user,
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func UnRegUser(db *mongo.Database, user *User) bool {
	status := true
	log.Println("Unregister User")
	_, err := db.Collection("PersonalDetails").DeleteOne(
		context.TODO(), &fiber.Map{
			"email":    user.Email,
			"password": user.Password,
		},
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func UpdateUser(db *mongo.Database, user *User) bool {
	status := true
	_, err := db.Collection("PersonalDetails").UpdateOne(context.TODO(), bson.M{
		"email": user.Email,
	}, bson.M{"$set": user}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func FindUser(db *mongo.Database, user *User) bool {
	status := true
	resp, err := db.Collection("PersonalDetails").
		Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	for resp.Next(context.TODO()) {
		var userTemp User
		err := resp.Decode(&userTemp)
		if err != nil {
			log.Fatal(err)
		}
		if userTemp.Password == user.Password && user.Email == userTemp.Email {
			status = true
			return status
		}
	}
	status = false
	return status
}

func AddTeam(db *mongo.Database, team *Team) bool {
	status := true
	_, err := db.Collection("Teams").InsertOne(
		context.TODO(), team,
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func AddTeamMember(db *mongo.Database, teamMember *User, teamid string) bool {
	status := true

	// find the teaminformation extract the team members in array then add new member and then put it back
	db.Collection("Teams").FindOne(context.TODO())

	_, err := db.Collection("Teams").UpdateOne(context.TODO(),
		bson.M{"teamid": teamid})
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func AddUsersGameInfo(db *mongo.Database, gameInformationOfUser *GameInformationOfUser) bool {
	status := true
	_, err := db.Collection("UsersGameInformation").InsertOne(
		context.TODO(), gameInformationOfUser,
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func AddGame(db *mongo.Database, gameInfo *Game) bool {
	status := true
	_, err := db.Collection("GameInformation").InsertOne(
		context.TODO(), gameInfo,
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func addTransaction(db *mongo.Database, transactionInfo *Transaction) bool {
	status := true
	_, err := db.Collection("TransactionInfo").InsertOne(
		context.TODO(), transactionInfo,
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func addReference(db *mongo.Database, reference *Refer) bool {
	status := true
	_, err := db.Collection("ReferenceInfo").InsertOne(
		context.TODO(), reference,
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}
