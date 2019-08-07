package main

import (
	"context"
	"log"

	"github.com/vitorhcastro/grpc-master-class/blog/blogpb"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Hello! I'm a client.")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	res := doCreate(c)

	doRead(c, "doesNotExist")
	doRead(c, res.GetBlog().GetId())
}

func doCreate(c blogpb.BlogServiceClient) *blogpb.CreateBlogResponse {
	log.Println("Creating a blog")

	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Vitor",
			Title:    "My first blog",
			Content:  "Content of the first blog",
		},
	}

	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	log.Printf("Blog has been created: %v", res)

	return res
}

func doRead(c blogpb.BlogServiceClient, blogID string) {
	log.Println("Reading a blog")

	req := &blogpb.ReadBlogRequest{
		BlogId: blogID,
	}

	res, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		log.Printf("Error happened while reading: %v", err)
		return
	}
	log.Printf("Blog was read: %v", res)
}
