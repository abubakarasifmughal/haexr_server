package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func SignUpWithCode(db *mongo.Database, user *User, code string) bool {
	status := true
	resp := db.Collection("ReferenceInfo").FindOne(context.TODO(), bson.M{"code": code})
	var refer Refer
	signeeReward := 100
	referrerReward := 200
	resp.Decode(&refer)
	if refer.Code == code {
		user.UserWallet.Bonus_cash = signeeReward

		_, err := db.Collection("PersonalDetails").InsertOne(
			context.TODO(), user,
		)
		if err != nil {
			log.Printf(err.Error())
			status = false
		} else {
			log.Printf("Success")
			status = true
			db.Collection("PersonalDetails").UpdateOne(context.TODO(),
				bson.M{"user_uuid": refer.Produce_user_uuid}, bson.M{"$inc": bson.M{"userwallet.bonus_cash": referrerReward}})
		}

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

func GetUserDetails(db *mongo.Database, email string) User {
	resp := db.Collection("PersonalDetails").FindOne(context.TODO(), bson.M{"email": email})
	var userInformation User
	resp.Decode(&userInformation)
	userInformation.Password = ""
	return userInformation
}

func GetUserDetailsUUID(db *mongo.Database, uuid string) User {
	resp := db.Collection("PersonalDetails").FindOne(context.TODO(), bson.M{"user_uuid": uuid})
	var userInformation User
	resp.Decode(&userInformation)
	userInformation.Password = ""
	return userInformation
}

func AddTeam(db *mongo.Database, team *Team) bool {
	status := true
	_, err := db.Collection("Teams").InsertOne(
		context.TODO(), team,
	)
	for i := 0; i < len(team.UsersInTeam); i++ {
		println(team.UsersInTeam[i].User_uuid)
	}
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func AddTeamMember(db *mongo.Database, teamMember User, teamid string) bool {
	status := true
	resp, err := db.Collection("Teams").Find(context.TODO(), bson.M{"teamid": teamid})
	for resp.Next(context.TODO()) {
		var teamTemp Team
		resp.Decode(&teamTemp)
		// teamTemp.UsersInTeam = append(teamTemp.UsersInTeam, teamMember)
		teamTemp.UsersInTeam = append(teamTemp.UsersInTeam, teamMember)
		_, err := db.Collection("Teams").UpdateOne(context.TODO(),
			bson.M{"teamid": teamid}, bson.M{"$set": teamTemp})
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Added user " + teamMember.User_uuid + " to team " + teamid)
		}
	}
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func splice(array []User, nameToRem string) []User {
	index := 0
	for i := 0; i < len(array); i++ {
		if nameToRem == array[i].User_uuid {
			index = i
			break
		}
	}
	return append(array[:index], array[index+1:]...)
}

func DelTeamMember(db *mongo.Database, teamMember *User, teamid string) bool {
	status := true
	resp, err := db.Collection("Teams").Find(context.TODO(), bson.M{"teamid": teamid})
	for resp.Next(context.TODO()) {
		var teamTemp Team
		resp.Decode(&teamTemp)

		teamTemp.UsersInTeam = splice(teamTemp.UsersInTeam, teamMember.User_uuid)
		_, err := db.Collection("Teams").UpdateOne(context.TODO(),
			bson.M{"teamid": teamid}, bson.M{"$set": teamTemp})
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Added user " + teamMember.User_uuid + " to team " + teamid)
		}
	}
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

func GetGame(db *mongo.Database) []Game {
	GamesList := []Game{}
	list, err := db.Collection("GameInformation").Find(
		context.TODO(), bson.M{},
	)
	if err != nil {
		log.Printf(err.Error())

	} else {

		for list.Next(context.TODO()) {
			var game Game
			list.Decode(&game)
			GamesList = append(GamesList, game)
		}

	}
	return GamesList
}

func GetGameInfo(db *mongo.Database, gameid string) Game {
	Game := Game{}
	gameRaw := db.Collection("GameInformation").FindOne(
		context.TODO(), bson.M{"gameid": gameid},
	)
	gameRaw.Decode(&Game)
	return Game
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

func CreateTeams(db *mongo.Database, newTeam Team) bool {
	status := true
	_, err := db.Collection("Teams").InsertOne(
		context.TODO(), newTeam,
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

func GetTeamByName(db *mongo.Database, teamName string) []User {
	res := db.Collection("Teams").FindOne(
		context.TODO(), bson.M{"teamname": teamName},
	)
	var team Team
	res.Decode(&team)
	return team.UsersInTeam
}

func GetTeamsByGameID(db *mongo.Database, gameid string) []Team {
	var teams []Team
	res, err := db.Collection("Teams").Find(
		context.TODO(), bson.M{"gameid": gameid},
	)
	if err != nil {
		return nil
	}
	for res.Next(context.TODO()) {
		var team Team
		res.Decode(&team)
		teams = append(teams, team)
	}
	return teams
}

func GetTeamsWhole(db *mongo.Database) []Team {
	TeamsList := []Team{}
	list, err := db.Collection("Teams").Find(
		context.TODO(), bson.M{},
	)
	if err != nil {
		log.Printf(err.Error())
	} else {
		for list.Next(context.TODO()) {
			var team Team
			list.Decode(&team)
			TeamsList = append(TeamsList, team)
		}
	}
	return TeamsList
}

func AddUserToTeam(db *mongo.Database, user string, team string) bool {
	_, err := db.Collection("Teams").UpdateOne(context.TODO(),
		bson.M{"teamid": team}, bson.M{"$push": bson.M{"usersinteam": user}})
	if err != nil {
		return false
	}
	return true
}

func RemoveUserFromTeam(db *mongo.Database, user string, team string) bool {
	_, err := db.Collection("Teams").UpdateOne(context.TODO(),
		bson.M{"teamid": team}, bson.M{"$pull": bson.M{"usersinteam": user}})
	if err != nil {
		return false
	}
	return true
}

func AddTournament(db *mongo.Database, tournament Tournaments) bool {
	_, err := db.Collection("Tournaments").InsertOne(context.TODO(),
		tournament)
	if err != nil {
		return false
	}
	return true
}

func AddStreamingLinksToTournament(db *mongo.Database, tournament string, steamLink StreamLink) bool {
	_, err := db.Collection("Tournaments").UpdateOne(
		context.TODO(), bson.M{"title": tournament},
		bson.M{"$push": bson.M{"streamlinks": steamLink}},
	)
	if err != nil {
		return false
	}
	return true
}

func AddTeamToTournament(db *mongo.Database, tournament string, team Team) bool {
	_, err := db.Collection("Tournaments").UpdateOne(
		context.TODO(), bson.M{"title": tournament},
		bson.M{"$push": bson.M{"teams": team}},
	)
	if err != nil {
		return false
	}
	return true
}

func GetTournament(db *mongo.Database, tournament string) Tournaments {
	res := db.Collection("Tournaments").FindOne(context.TODO(),
		bson.M{"title": tournament})
	var result Tournaments
	res.Decode(&result)
	return result
}

func GetTournamentByGame(db *mongo.Database, gameid string) []Tournaments {
	res, _ := db.Collection("Tournaments").Find(context.TODO(),
		bson.M{"gameid": gameid})
	tournamentsOfGame := []Tournaments{}

	var result Tournaments
	res.Decode(&result)
	for res.Next(context.TODO()) {
		res.Decode(&result)
		println(result.GameID)
		tournamentsOfGame = append(tournamentsOfGame, result)
	}
	return tournamentsOfGame
}

func GetTournaments(db *mongo.Database) []Tournaments {
	res, err := db.Collection("Tournaments").Find(context.TODO(),
		bson.M{})
	tournaments := []Tournaments{}
	tempHolder := Tournaments{}
	if err != nil {
		return nil
	}
	for res.Next(context.TODO()) {
		res.Decode(&tempHolder)
		tournaments = append(tournaments, tempHolder)
	}
	return tournaments
}

func AddQualifierRoundInTournament(db *mongo.Database, tournament string, qualifier Rounds) bool {
	res, err := db.Collection("Tournaments").UpdateOne(context.TODO(),
		bson.M{"title": tournament}, bson.M{"$push": bson.M{"rounds": qualifier}})
	if err != nil {
		return false
	}

	if res.ModifiedCount == 0 {
		return false
	}

	return true
}

func AddTeamInTournamentGroup(db *mongo.Database, tournament string, qualifier string, group Groups, team Team) Tournaments {
	// key is concatenation of date and time

	res1 := db.Collection("Tournaments").FindOne(context.TODO(), bson.M{"title": tournament, "rounds.qualifiername": qualifier})
	var data Tournaments
	res1.Decode(&data)

	for i := 0; i < len(data.Rounds); i++ {
		if data.Rounds[i].QualifierName == qualifier {
			println("found ")
			currentRound := data.Rounds[i]
			// Qualifier of focus
			// ---------------------------------------------------
			// check if the slot that user wants exists or not
			foundSameSlot := false
			for j := 0; j < len(currentRound.Groups); j++ {
				if currentRound.Groups[j].StartingAtDate == group.StartingAtDate &&
					currentRound.Groups[j].StartingAtTime == group.StartingAtTime {
					foundSameSlot = true

					// if exist then
					if foundSameSlot {
						// check if there is capacity
						if len(data.Rounds[i].Groups[j].Teams) < data.Rounds[i].NumberOfTeamsPerGroup {
							data.Rounds[i].Groups[j].Teams = append(data.Rounds[i].Groups[j].Teams, team)
							// if has capacity then add in same slot
							// write back to database
							db.Collection("Tournaments").UpdateOne(context.TODO(),
								bson.M{"rounds.qualifiername": data.Rounds[i].QualifierName},
								bson.M{"$set": bson.M{"rounds.$.groups": data.Rounds[i].Groups}})
							return Tournaments{}
						} else {
							println("CAPACITY / GROUP REACHED")
							foundSameSlot = false
							continue
						}
					}
					break
				}
			}
			if !foundSameSlot {

				if len(data.Rounds[i].Groups) < data.TotalTeams/data.Rounds[i].NumberOfTeamsPerGroup {
					newGroupWithTeam := Groups{
						GroupID:        "some random id",
						MatchID:        "BGMI #1212",
						StartingAtTime: group.StartingAtTime,
						StartingAtDate: group.StartingAtDate,
						Group:          "Group Name",
						Teams:          []Team{team},
						Results:        []string{},
						RoomID:         "",
						Password:       "",
						Duration:       "45",
						Rounds:         []Match{},
					}
					// here i add new group
					db.Collection("Tournaments").UpdateOne(context.TODO(),
						bson.M{"title": tournament, "rounds.qualifiername": qualifier},
						bson.M{"$push": bson.M{"rounds.$.groups": newGroupWithTeam}})
					foundSameSlot = false
				} else {
					println("No more teams")
				}
			}
			// else make a new one provided slot can be made

			return data
		}
	}

	return Tournaments{}
	/*
		// get all the tournaments to see the dates and time
		res1 := db.Collection("Tournaments").FindOne(context.TODO(), bson.M{"title": tournament, "rounds.qualifiername": qualifier})
		var data Tournaments
		res1.Decode(&data)

		// once got and save now check does the date and time that user is demainding is available pele se ya nai
		for i := 0; i < len(data.Rounds); i++ {
			tempRound := data.Rounds[i]

			res11 := db.Collection("Tournaments").FindOne(context.TODO(), bson.M{"title": tournament, "rounds.qualifiername": qualifier})
			var data2 Tournaments
			res11.Decode(&data2)
			println("----")
			println(len(data2.Rounds[i].Groups))
			println(data2.TotalTeams / data2.Rounds[i].NumberOfTeamsPerGroup)
			println("----")
			if len(data2.Rounds[i].Groups)-2 >= data2.TotalTeams/data2.Rounds[i].NumberOfTeamsPerGroup {
				println("Limit reached cant create more groups")
				return Tournaments{}
			}

			if tempRound.QualifierName == qualifier {
				// println("found")
				// great job now that you have found the qualifier you needed lets search for the slot we're looking for so we can add the user in there
				// check of any group exist or not

				// just in case if groups dont exist then for loop wont execute so I've to manually insert first time
				if len(tempRound.Groups) == 0 {

					println("================")
					println("array was empty adding the first entry")
					println("================")
					newGroupWithTeam := Groups{
						GroupID:        "some random id",
						MatchID:        "BGMI #1212",
						StartingAtTime: group.StartingAtTime,
						StartingAtDate: group.StartingAtDate,
						Group:          "Group Name",
						Teams:          []Team{team},
						Results:        []string{},
						RoomID:         "",
						Password:       "",
						Duration:       "45",
						Rounds:         []Match{},
					}
					// here i add new group
					res2, err := db.Collection("Tournaments").UpdateOne(context.TODO(),
						bson.M{"title": tournament, "rounds.qualifiername": qualifier},
						bson.M{"$push": bson.M{"rounds.$.groups": newGroupWithTeam}})

					if err != nil {
						return Tournaments{}
					}

					if res2.ModifiedCount == 0 {
						return Tournaments{}
					}
					res1 = db.Collection("Tournaments").FindOne(context.TODO(), bson.M{"title": tournament, "rounds.qualifiername": qualifier})
					res1.Decode(&data)
					return data
				} else {
					for j := 0; j < len(tempRound.Groups); j++ {
						if tempRound.Groups[j].StartingAtTime == group.StartingAtTime &&
							tempRound.Groups[j].StartingAtDate == group.StartingAtDate {
							println("Group with same date time found")
							// see if there is capacity in this group or not
							// if teams joined in group are more than number of teams allowed per that group
							// TODO this check doesnt work
							println("Currently we have ")
							println(len(tempRound.Groups[j].Teams))
							if tempRound.NumberOfTeamsPerGroup > len(tempRound.Groups[j].Teams) {
								println("Still has some capacity, adding slot")
								tempRound.Groups[j].Teams = append(tempRound.Groups[j].Teams, team)

								db.Collection("Tournaments").UpdateOne(context.TODO(),
									bson.M{"rounds.qualifiername": tempRound.QualifierName}, bson.M{"$set": bson.M{"rounds.$.groups": tempRound.Groups}})

								// adding a slot in the group
								res1 = db.Collection("Tournaments").FindOne(context.TODO(), bson.M{"title": tournament, "rounds.qualifiername": qualifier})
								res1.Decode(&data)
								return data
							} else {
								// println("Doesnt have any capacity making a new slot")

								newGroupWithTeam := Groups{
									GroupID:        "some random id",
									MatchID:        "BGMI #1212",
									StartingAtTime: group.StartingAtTime,
									StartingAtDate: group.StartingAtDate,
									Group:          "Group Name",
									Teams:          []Team{team},
									Results:        []string{},
									RoomID:         "",
									Password:       "",
									Duration:       "45",
									Rounds:         []Match{},
								}
								// here i add new group
								res2, err := db.Collection("Tournaments").UpdateOne(context.TODO(),
									bson.M{"title": tournament, "rounds.qualifiername": qualifier}, bson.M{"$push": bson.M{"rounds.$.groups": newGroupWithTeam}})
								if err != nil {
									return Tournaments{}
								}

								if res2.ModifiedCount == 0 {
									return Tournaments{}
								}
								// return by adding new team
								res1 = db.Collection("Tournaments").FindOne(context.TODO(), bson.M{"title": tournament, "rounds.qualifiername": qualifier})
								res1.Decode(&data)
								return data
							}
						} else {
							// println("Group doesnt exist, making a new one")
							// so add a new group for this team when no other group for the key is present
							newGroupWithTeam := Groups{
								GroupID:        "some random id",
								MatchID:        "BGMI #1212",
								StartingAtTime: group.StartingAtTime,
								StartingAtDate: group.StartingAtDate,
								Group:          "Group Name",
								Teams:          []Team{team},
								Results:        []string{},
								RoomID:         "",
								Password:       "",
								Duration:       "45",
								Rounds:         []Match{},
							}
							// here i add new group with team in it
							res2, err := db.Collection("Tournaments").UpdateOne(context.TODO(),
								bson.M{"title": tournament, "rounds.qualifiername": qualifier}, bson.M{"$push": bson.M{"rounds.$.groups": newGroupWithTeam}})

							if err != nil {
								return Tournaments{}
							}

							if res2.ModifiedCount == 0 {
								return Tournaments{}
							}
							// and return
							res1 = db.Collection("Tournaments").FindOne(context.TODO(), bson.M{"title": tournament, "rounds.qualifiername": qualifier})
							res1.Decode(&data)
							return data
						}
					}
				}
			}

			// GETTING THE INFORMATION
			print("================================= Round ")
			println(i)
			resinfo := db.Collection("Tournaments").FindOne(context.TODO(), bson.M{"title": tournament, "rounds.qualifiername": qualifier})
			var info Tournaments
			resinfo.Decode(&info)
			println("Tournament informations")
			print("MAX teams ")
			println(info.TotalTeams)
			print("Total # groups currently formed for round are ")
			println(len(info.Rounds[i].Groups))
			println("Total # teams formed/group are ")
			for a := 0; a < len(info.Rounds[i].Groups); a++ {
				print("Group ")
				print(a)
				print(" has ")
				print(len(info.Rounds[i].Groups[a].Teams))
				println(" team(s)")
			}
			print("MAX number of group ")
			println(info.TotalTeams / info.Rounds[i].NumberOfTeamsPerGroup)
			println("=================================")
		}

		res1 = db.Collection("Tournaments").FindOne(context.TODO(), bson.M{"title": tournament, "rounds.qualifiername": qualifier})
		res1.Decode(&data)
		return data
	*/
}
