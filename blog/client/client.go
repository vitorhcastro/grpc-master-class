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

	doCreate(c)
}

func doCreate(c blogpb.BlogServiceClient) {
	log.Println("Creating a blog")

	blog := &blogpb.Blog{
		AuthorId: "Vitor",
		Title:    "My first blog",
		Content:  "Content of the first blog",
	}

	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	log.Printf("Blog has been created: %v", res)
}
