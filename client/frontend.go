package main

import (
	"context"
	"fmt"
	"time"

	pb "example.com/increment"
)

type frontend struct {
	replicas map[int32]pb.IncrementServiceClient
	ctx      context.Context
	repch    chan pb.IncrementReply
}

func (fe *frontend) incrementReplica(req *pb.IncrementRequest, c pb.IncrementServiceClient) {
	rep, err := c.Increment(fe.ctx, req)
	if err != nil {
		fe.repch <- pb.IncrementReply{Success: false, ValueBefore: -1}
	} else {
		fe.repch <- *rep
	}
}

func (fe *frontend) increment(value int32) string {
	if value <= 0 {
		return fmt.Sprint("You can't set a negative value.")
	}

	req := &pb.IncrementRequest{
		Value: value,
	}
	var rep pb.IncrementReply

incrementreplicas:
	for k, v := range fe.replicas {
		go fe.incrementReplica(req, v)

		start := time.Now()
	fivesecondcheck:
		for start.Add(5 * time.Second).After(time.Now()) {
			select {
			case reply := <-fe.repch:
				if reply.ValueBefore < 0 {
					delete(fe.replicas, k)
					continue incrementreplicas
				}
				rep = reply
				break fivesecondcheck
			default:
			}
		}
	}

	if rep.Success {
		return fmt.Sprintf("Cool! Incremented to %v - was %v before.", value, rep.ValueBefore)
	} else {
		return fmt.Sprintf("Uh oh! You cannot increment to %v, since the value is already %v.", value, rep.ValueBefore)
	}

}
