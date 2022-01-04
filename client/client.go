package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	pb "example.com/increment"
	"google.golang.org/grpc"
)

func main() {

	fe := frontend{
		ctx:      context.Background(),
		repch:    make(chan pb.IncrementReply),
		replicas: make(map[int32]pb.IncrementServiceClient),
	}

	for i := 5000; i < 5003; i++ {
		address := fmt.Sprintf(":%v", i)
		// log.Printf("Dialing %v", address)
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		//defer conn.Close()
		fe.replicas[int32(i)] = pb.NewIncrementServiceClient(conn)
	}

	fmt.Print("\n=== Welcome to the distributed incrementor ===\n\n╒═════════════════ COMMANDS ══════════════════╕\n│ Write the value you want your incrementor   │\n│ to have. You cannot set the incrementor to  │\n│ a lower value than it already is.           │\n│ Only positive integers are accepted.        │\n└─────────────────────────────────────────────┘\n\n")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		input, err := strconv.ParseInt(command, 10, 32)
		if err != nil {
			fmt.Printf("%v is not an integer.\n", command)
		} else {
			fmt.Println(fe.increment(int32(input)))
		}
	}
}
