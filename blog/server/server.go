package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vitorhcastro/grpc-master-class/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct{}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	log.Printf("Create blog request: %v", req)

	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprint("Cannot convert to OID"),
		)
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	log.Printf("Read blog request: %v", req)

	blogID := req.GetBlogId()

	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprint("Cannot parse ID"),
		)
	}

	data := &blogItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err = res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find log with ID '%v': %v", blogID, err),
		)
	}

	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Hello! I'm a server.")

	log.Println("Opening MongoDB connection")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to start mongo client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to mongodb: %v", err)
	}

	collection = client.Database("mydb").Collection("blog")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		log.Println("Starting server...")
		if err := s.Serve(lis); err == nil {
			log.Fatalf("failed to serve : %v", err)
		}
	}()

	// Wait for CTRL + C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	log.Println("Stopping the server")
	s.Stop()
	log.Println("Closing the listener")
	lis.Close()
	log.Println("Closing MongoDB connection")
	client.Disconnect(ctx)
	log.Println("End of program")
}
