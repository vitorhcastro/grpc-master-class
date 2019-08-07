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

	blog := &blogpb.Blog{
		AuthorId: "Vitor",
		Title:    "My first blog",
		Content:  "Content of the first blog",
	}

	res := doCreate(c, blog)

	doRead(c, "doesNotExist")
	doRead(c, res.GetBlog().GetId())

	newBlog := &blogpb.Blog{
		Id:       res.GetBlog().GetId(),
		AuthorId: "Vitor Castro",
		Title:    "A better title for a blog",
		Content:  "Content of the first blog with a little more content",
	}

	doUpdate(c, newBlog)

	doDelete(c, "doesNotExist")
	doDelete(c, res.GetBlog().GetId())
}

func doCreate(c blogpb.BlogServiceClient, blog *blogpb.Blog) *blogpb.CreateBlogResponse {
	log.Println("Creating a blog")

	req := &blogpb.CreateBlogRequest{
		Blog: blog,
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

func doUpdate(c blogpb.BlogServiceClient, blog *blogpb.Blog) {
	log.Println("Updating a blog")

	req := &blogpb.UpdateBlogRequest{
		Blog: blog,
	}

	res, err := c.UpdateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	log.Printf("Blog has been updated: %v", res)
}

func doDelete(c blogpb.BlogServiceClient, blogID string) {
	log.Println("Deleting a blog")

	req := &blogpb.DeleteBlogRequest{
		BlogId: blogID,
	}

	res, err := c.DeleteBlog(context.Background(), req)
	if err != nil {
		log.Printf("Error happened while deleting: %v", err)
		return
	}
	log.Printf("Blog was deleted: %v", res)
}
