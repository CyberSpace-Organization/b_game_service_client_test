//@author by. Created in 2022/7/3
package main

import (
	pb "GrpcTestGo/go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

var streamClient pb.RoomClient

const(
	Token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIzMjE2Njc2MTI3ZWM0YTNjYjhmM2ZiM2RjZTA1NGEzNiIsImF1dGhWYWx1ZSI6ImNvbW1vbiIsInVzZXJObyI6IjMyMTY2NzYxMjdlYzRhM2NiOGYzZmIzZGNlMDU0YTM2IiwiZXhwIjoxNjU2OTUxNjc2LCJpYXQiOjE2NTY4NjUyNzZ9.xk5qhutdgzoFRlUBdjjAWDbjv_qzV5yqNaQTsm5hAis"
	Name = "by."
	UserNo = "3216676127ec4a3cb8f3fb3dce054a36"
	UserLevel = "2"
)

//Communicate with the server
func connectToRoom(roomId string){
	stream, err := streamClient.ConnectToTheRoom(context.Background())
	if err!=nil{
		log.Fatalf("Get connect feedback from server via stream: %v",err)
	}
	for n:=0;n<1000;n++{
		var request *pb.ConnectRequest
		if n==0{
			request = &pb.ConnectRequest{IsFirstConnect: true, IsToTerminate: false, Token: Token, RoomCode: roomId, PlayerInfo: &pb.PlayerInfo{}}
		}else if n==999 {
			request = &pb.ConnectRequest{
				IsFirstConnect: false,
				IsToTerminate: true,
				Token: Token,
				RoomCode: roomId,
				PlayerInfo: &pb.PlayerInfo{
					Name:      Name,
					UserNo:    UserNo,
					UserLevel: UserLevel,
					XPosition: int32(n),
					YPosition: int32(n),
					ZPosition: int32(n),
				},
			}
		}else {
			request = &pb.ConnectRequest{
				IsFirstConnect: false,
				IsToTerminate: false,
				Token: "",
				RoomCode: roomId,
				PlayerInfo: &pb.PlayerInfo{
					Name:      Name,
					UserNo:    UserNo,
					UserLevel: UserLevel,
					XPosition: int32(n),
					YPosition: int32(n),
					ZPosition: int32(n),
				},
			}
		}
		err:=stream.Send(request)
		if err != nil{
			log.Fatalf("Connection got stream error! %v",err)
		}
		res,err := stream.Recv()
		if err == io.EOF{
			break
		}
		if err != nil{
			log.Fatalf("Conversation get stream: %v",err)
		}
		log.Println(res.Message)
		log.Println(res.RoomInfo)
		log.Println(res.Players)
	}
}

func main(){
	conn,err := grpc.Dial("localhost:20000",grpc.WithInsecure())
	if err != nil{
		log.Fatal("Connect failed! %v",err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()
	streamClient = pb.NewRoomClient(conn)
	room, err := streamClient.CreateRoom(ctx, &pb.CreateRoomRequest{
		RoomName:   "first room",
		RoomLength: 1000,
		RoomWidth:  1000,
		RoomHeight: 1000,
		NeedPass:   false,
		Password:   "",
	})
	if room.IsSuccess==true{
		connectToRoom(room.RoomCode)
	}

}

