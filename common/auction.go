package common

import (
	"fmt"
	"strconv"
	"time"
)

type Auction struct {
	Id                    int
	Name                  string
	StartingBet           float32
	CurrentBet            float32
	RemainingTime         time.Duration
	StartingTime          time.Time
	CurrentBestContestant string
}

func PrettifyAuction(auction Auction) string {
	prettifiedAuction := ""

	prettifiedAuction += "id: " + strconv.Itoa(auction.Id)
	prettifiedAuction += " | name: " + auction.Name
	prettifiedAuction += " | current bet: " + fmt.Sprintf("%.2f", auction.CurrentBet)
	prettifiedAuction += " | best contestant : " + auction.CurrentBestContestant
	prettifiedAuction += " | remaining time: " + auction.RemainingTime.String()

	return prettifiedAuction
}
