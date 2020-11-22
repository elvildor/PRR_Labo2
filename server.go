package main

import (
	"bufio"
	"common"
	"lamport"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

type client chan<- string // an outgoing message channel

type requestClientStruct struct {
	cli client
	req common.Request
	address string
}

type notificationStruct struct {
	typeOfNotification string
	message string
	auctionId string
}

var (
	clients                    = make(map[client]bool) // all connected clients
	usernames                  = make(map[client]string)
	existingUsername           = make(map[string]bool)
	idAuction                  = 0
	auctions                   = make(map[client][]common.Auction)
	clientsNotifyNew           = make(map[client]bool)
	clientsNotifyAuctionUpdate = make(map[common.Auction]map[client]bool)

	entering      = make(chan client)
	leaving       = make(chan client)
	messages      = make(chan string) // all incoming client messages
	logs          = make(chan string)
	logsFatal     = make(chan error)
	requests      = make(chan requestClientStruct)
	notifications = make(chan notificationStruct)
)

func main()  {
	go logging()
	listener, err := net.Listen("tcp", common.ServerAddress)
	if err != nil {
		logsFatal <- err
	}
	go broadcaster()
	go notify()
	logs <- "broadcast created"
	for {
		conn, err := listener.Accept()
		if err != nil {
			logsFatal <- err
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	for {
		select {
		case request := <-requests:
			requestToClient, errTreatment := treatmentRequest(request.req, request.cli)
			if errTreatment != nil {
				logsFatal <- errors.New(request.address + " " + errTreatment.Error())
			}
			request.cli <- requestToClient
		case msg := <-messages: // broadcaster <- handleConn
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg // clientwriter (handleConn) <- broadcaster
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli) // delete client

			delete(existingUsername, usernames[cli]) //delete username
			delete(usernames, cli)

			delete(auctions, cli) // delete auctions
			close(cli)
		}

		// update remaining time
		for _, auctionsClient := range auctions {
			for idx, auction := range auctionsClient {
				auction.RemainingTime = auction.RemainingTime - (time.Now().Sub(auction.StartingTime))

				if auction.RemainingTime <= 0 {
					auctionsClient = append(auctionsClient[:idx], auctionsClient[idx+1:]...)
					if auction.CurrentBestContestant == "" {
						notifications <- notificationStruct{
							typeOfNotification: "expired",
							message: fmt.Sprintf(
								"auction %v with name \"%s\" could not have been sold. Time expired",
								auction.Id,
								auction.Name,
								),
						}
					} else {
						notifications <- notificationStruct{
							typeOfNotification: "sold",
							message: fmt.Sprintf(
								"auction %v with name \"%s\" had been sold to the best bidder %s for %f",
								auction.Id,
								auction.Name,
								auction.CurrentBestContestant,
								auction.CurrentBet,
								),
						}
					}
				} else
				{
					auctionsClient[idx] = auction
				}
			}
		}
	}
}

func notify() {
	for true {
		select {
		case notification := <-notifications:
			switch notification.typeOfNotification {
			case "add":
				for cli := range clientsNotifyNew {
					if clientsNotifyNew[cli] {
						cli <- notification.message
					}
				}
				break
			case "update":
				auction, err := getAuction(notification.auctionId)
				if err != nil {
					logsFatal <- err
					break
				}
				for cli := range clientsNotifyAuctionUpdate[auction] {
					if clientsNotifyAuctionUpdate[auction][cli] {
						cli <- notification.message
					}
				}
				break

			default:
				logs <- "unknown type"
				break
			}
		}
	}
}

func logging() {
	for {
		select {
		case logMsg := <-logs:
			log.Println(logMsg)
		case logError := <-logsFatal:
			log.Println("[ALERT]: " + logError.Error())
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // channel 'client' mais utilisé ici dans les 2 sens
	who := conn.RemoteAddr().String()

	go func() { // clientwriter
		for msg := range ch { // clientwriter <- broadcaster, handleConn
			logs <- who + " " + msg
			_, _ = fmt.Fprintln(conn, msg) // netcat Client <- clientwriter
		}
	}()

	logs <- who + " has arrived" // broadcaster <- handleConn
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() { // handleConn <- netcat client
		//logs <- who + ": " + input.Text()
		var tramJSON common.Request
		_ = json.Unmarshal(input.Bytes(), &tramJSON)
		requests <- requestClientStruct{req: tramJSON, cli: ch, address: who}
	}

	leaving <- ch
	logs <- who + " has left" // broadcaster <- handleConn
	_ = conn.Close()
}


func treatmentRequest(request common.Request, cli client) (string, error) {
	switch request.Header {
	case "CONNECT":
		var nameCli string
		var ok bool
		// this avoid a panic kernel
		if x, found := request.Body["Username"]; found {
			if nameCli, ok = x.(string); !ok {
				//do whatever you want to handle errors - this means this wasn't a string
				return "", errors.New("not a string")
			}
		} else {
			//handle error - the map didn't contain this key
			return "", errors.New("doesn't have a username")
		}
		logs <- "CONNECT " + nameCli

		if existingUsername[nameCli] {
			return "", errors.New("username " + nameCli + " already exist")
		}

		usernames[cli] = nameCli
		existingUsername[nameCli] = true

		return "ok", nil
	case "ADD":
		tmpAuctionJSON, _ := json.Marshal(request.Body["Auction"])
		var auction common.Auction
		err := json.Unmarshal(tmpAuctionJSON, &auction)
		if err != nil {
			return "", err
		}
		logs <- "ADD " + auction.Name
		auction.Id = idAuction
		idAuction++
		auction.StartingTime = time.Now()
		auction.CurrentBet = auction.StartingBet
		auctions[cli] = append(auctions[cli], auction)

		notifications <- notificationStruct{
			typeOfNotification: "add",
			message: fmt.Sprintf("auction %v has been added", auction.Id),
			auctionId: strconv.Itoa(auction.Id),
		}

		return "ok", nil
	case "LIST":
		logs <- "LIST"
		allAuctions := getAllAuctions()
		auctionsJSON, err := json.Marshal(allAuctions)
		return string(auctionsJSON), err

	case "SELECT":
		var id string
		var ok bool
		// this avoid a panic kernel
		if x, found := request.Body["Id"]; found {
			if id, ok = x.(string); !ok {
				//do whatever you want to handle errors - this means this wasn't a string
				return "", errors.New("not a string")
			}
		} else {
			//handle error - the map didn't contain this key
			return "", errors.New("doesn't have an id")
		}

		logs <- "SELECT " + id

		auction, err := getAuction(id)

		if err == nil {
			auctionJSON, err := json.Marshal(auction)
			return string(auctionJSON), err
		}

		return "", errors.New("id not found")

	case "RAISE":
		var id string
		var raise float32
		var ok bool
		// this avoid a panic kernel
		if x, found := request.Body["Id"]; found {
			if id, ok = x.(string); !ok {
				//do whatever you want to handle errors - this means this wasn't a string
				return "", errors.New("not a string")
			}
		} else {
			//handle error - the map didn't contain this key
			return "", errors.New("doesn't have an id")
		}

		if x, found := request.Body["Raise"]; found {
			var tmp float64
			if tmp, ok = x.(float64); !ok {
				//do whatever you want to handle errors - this means this wasn't a string
				return "", errors.New("not a float")
			}
			raise = float32(tmp)
		} else {
			//handle error - the map didn't contain this key
			return "", errors.New("doesn't have an value to raise")
		}

		logs <- "RAISE " + id + " to " + fmt.Sprintf("%f", raise)

		auction, err := getAuction(id)

		if err != nil {
			return "", err
		}

		if auction.CurrentBet < raise {
			auction.CurrentBet = raise
			auction.CurrentBestContestant = usernames[cli]
			updateAuction(auction)
			notifications <- notificationStruct{
				typeOfNotification: "update",
				message: fmt.Sprintf("auction %v has been raised to %f", auction.Id, raise),
				auctionId: strconv.Itoa(auction.Id),
			}

			return "ok", nil
		} else {
			return "", errors.New("the raise is inferior to the current bet")
		}

	case "NOTIFY":
		var id string
		var newAuctionNotification bool
		var addNotification bool
		var ok bool

		// this avoid a panic kernel
		if x, found := request.Body["Id"]; found {
			if id, ok = x.(string); !ok {
				//do whatever you want to handle errors - this means this wasn't a string
				return "", errors.New("not a string")
			}
		} else {
			//handle error - the map didn't contain this key
			return "", errors.New("doesn't have an id")
		}

		// this avoid a panic kernel
		if x, found := request.Body["New"]; found {
			if newAuctionNotification, ok = x.(bool); !ok {
				//do whatever you want to handle errors - this means this wasn't a string
				return "", errors.New("not a bool")
			}
		} else {
			//handle error - the map didn't contain this key
			newAuctionNotification = false
		}

		// this avoid a panic kernel
		if x, found := request.Body["AddNotification"]; found {
			if addNotification, ok = x.(bool); !ok {
				//do whatever you want to handle errors - this means this wasn't a string
				return "", errors.New("not a bool")
			}
		} else {
			addNotification = false
		}

		if newAuctionNotification {
			if addNotification {
				clientsNotifyNew[cli] = true
			} else {
				clientsNotifyNew[cli] = false
			}
		}

		if id != "" {
			auction, err := getAuction(id)
			if err != nil {
				return "", errors.New("id not found")
			}

			if clientsNotifyAuctionUpdate[auction] == nil {
				clientsNotifyAuctionUpdate[auction] = make(map[client]bool)
			}

			if addNotification {
				clientsNotifyAuctionUpdate[auction][cli] = true
			} else {
				clientsNotifyAuctionUpdate[auction][cli] = false
			}
		}

		return "ok", nil
	default:
		return "", errors.New("command '" + request.Header + "' doesn't exist")
	}
}

func getAllAuctions() []common.Auction {
	var allAuctions []common.Auction

	for elementClient := range auctions {
		allAuctions = append(allAuctions, auctions[elementClient]...)
	}

	return allAuctions
}

func getAuction(id string) (common.Auction, error){
	for _, auction := range getAllAuctions() {
		if strconv.Itoa(auction.Id) == id {
			return auction, nil
		}
	}

	return common.Auction{}, errors.New("auction not found")
}

func updateAuction(auction common.Auction) {

	for elementClient := range auctions {
		for idx := range auctions[elementClient] {
			if auctions[elementClient][idx].Id == auction.Id {
				auctions[elementClient][idx].CurrentBet = auction.CurrentBet
				auctions[elementClient][idx].RemainingTime = auction.RemainingTime
				auctions[elementClient][idx].CurrentBestContestant = auction.CurrentBestContestant
			}
		}
	}
}





func callLamport(command string, args []string) {
	switch command {
		case "REQ":
			processorTimestamp,_ := strconv.ParseUint(args[0],10,64)
			processorId,_ := strconv.Atoi(args[1])
			lamport.REQ(processorId,lamport.Clock{Timestamp: processorTimestamp})
		break

		case "ACK":
			processorTimestamp,_ := strconv.ParseUint(args[0],10,64)
			processorId,_ := strconv.Atoi(args[1])
			lamport.ACK(processorId,lamport.Clock{Timestamp: processorTimestamp})
		break

		case "REL":
			processorTimestamp,_ := strconv.ParseUint(args[0],10,64)
			processorId,_ := strconv.Atoi(args[1])
			lamport.REL(processorId,lamport.Clock{Timestamp: processorTimestamp})
		break
	}

	common.Debug("command lamport ended")
}