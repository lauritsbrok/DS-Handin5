package main

import (
	proto "AuctionSystem/grpc"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedAuctionSystemServer // Necessary
	name                                   string
	port                                   int
}

const initialPort = 8000

var highestBid int64 = 0
var isBidCompleted bool = false
var highestBidder int64

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + initialPort

	// Get the port from the command line when the server is run
	flag.Parse()

	// Create a server struct
	server := &Server{
		name: "serverName",
		port: int(ownPort),
	}

	// Start the server
	go startServer(server)

	go CloseAuction()

	go RandomCrash()

	// Keep the server running until it is manually quit
	for {

	}
}

func startServer(server *Server) {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", server.port)

	// Register the grpc server and serve its listener
	proto.RegisterAuctionSystemServer(grpcServer, server)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

func (c *Server) AskForBid(ctx context.Context, in *proto.AskForHighestBid) (*proto.HighestBid, error) {
	log.Printf("Client with ID %d asked for the highest bid\n", in.ClientId)
	output := ""
	if highestBid == 0 && !isBidCompleted {
		output = "No bid has been made yet"
	} else if isBidCompleted {
		output = fmt.Sprintf("The action is over. The winner was client %s with a bid of %s", strconv.FormatInt(highestBidder, 10), strconv.FormatInt(highestBid, 10))
	} else {
		output = fmt.Sprintf("The highest bid is %s", strconv.FormatInt(highestBid, 10))
	}
	return &proto.HighestBid{HighestBid: output}, nil
}

func (c *Server) MakeABid(ctx context.Context, in *proto.MakeBid) (*proto.Ack, error) {
	log.Printf("Client with ID %d makes a bid for %d \n", in.ClientId, in.Bid)
	message := ""
	if isBidCompleted {
		message = "The bid is over. Type 'result' to see who won."
	} else if highestBid < in.Bid {
		highestBid = in.Bid
		highestBidder = in.ClientId
		message = fmt.Sprintf("Success: The new highest bid is %s", strconv.FormatInt(highestBid, 10))
	} else {
		message = fmt.Sprintf("Failed: Your bid of %s was not high enough. The highest bid is %s", strconv.FormatInt(in.Bid, 10), strconv.FormatInt(highestBid, 10))
	}
	return &proto.Ack{Message: message}, nil
}

func CloseAuction() {
	time.Sleep(45 * time.Second)
	isBidCompleted = true
	log.Printf("Auction is now closed")
}

func RandomCrash() {
	time.Sleep(20 * time.Second)
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(5-1+1) + 1
	if num == 2 || num == 4 {
		log.Print("Oh no, i'm crashing. Good thing there is other replicas.")
		os.Exit(0)
	}
}
