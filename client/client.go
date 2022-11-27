package main

import (
	proto "AuctionSystem/grpc"
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	id         int
	portNumber int
}

var (
	serverPort = flag.Int("sPort", 8000, "server port number (should match the port used for the server)")
)

const initialPort = 8000

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + initialPort
	// Parse the flags to get the port for the client
	flag.Parse()

	// Create a client
	client := &Client{
		id:         int(ownPort),
		portNumber: int(ownPort),
	}

	// Wait for the client (user) to ask for the time
	go askAuction(client)

	for {

	}
}

func connectToServer() (proto.AuctionSystemClient, error) {
	// Dial the server at the specified port.
	*serverPort = *serverPort + 1
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *serverPort)
	} else {
		log.Printf("Connected to the server at port %d\n", *serverPort)
	}
	return proto.NewAuctionSystemClient(conn), nil
}

func askAuction(client *Client) {
	// Connect to the server
	serverConnection1, _ := connectToServer()
	serverConnection2, _ := connectToServer()
	serverConnection3, _ := connectToServer()

	serverConnections := []proto.AuctionSystemClient{serverConnection1, serverConnection2, serverConnection3}

	serverConnection1 = serverConnections[0]

	// Wait for input in the client terminal
	scanner := bufio.NewScanner(os.Stdin)
	log.Printf("Type 'result' to get the highest bid. Type 'bid' to make a bid.")
	for scanner.Scan() {
		input := scanner.Text()

		if input == "result" {
			for i := 0; i < 3; i++ {
				highestBidReturnMessage, err := serverConnections[i].AskForBid(context.Background(), &proto.AskForHighestBid{
					ClientId: int64(client.id),
				})

				if err != nil {
					log.Printf("Server %s is currently down, skipping ...", strconv.Itoa(i+1))
				} else {
					log.Printf("Server %s says: %s\n", strconv.Itoa(i+1), highestBidReturnMessage.HighestBid)
				}
			}
		} else if input == "bid" {
			log.Printf("How much would you like to bid? Enter a valid positive number:")
			for scanner.Scan() {
				newbid, err := strconv.ParseInt(scanner.Text(), 10, 64)
				if err != nil {
					log.Printf("Wrong input: Not a number!")
					break
				}

				for i := 0; i < 3; i++ {

					makeBidReturnMessage, err := serverConnections[i].MakeABid(context.Background(), &proto.MakeBid{
						ClientId: int64(client.id),
						Bid:      newbid,
					})

					if err != nil {
						log.Printf("Server %s is currently down, skipping ...", strconv.Itoa(i+1))
					} else {
						log.Printf("Server %s says: %s", strconv.Itoa(i+1), makeBidReturnMessage.Message)
					}
				}
				break
			}
		} else {
			log.Printf("Invalid command. Can you type 'result' or 'bid'")
		}

	}
}
