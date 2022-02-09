package main

import (
	"blog/blogpb"
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog client")
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50052", opts)
	if err != nil {
		log.Fatalf("could not connect : %v", err)
	}
	defer cc.Close()
	c := blogpb.NewBlogServiceClient(cc)
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Ihksan",
		Title:    "Hello Word",
		Content:  "its my first blog content",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Fatalf("Unexpected error : %v", err)
	}
	fmt.Printf("Blog has been created : %v \n", createBlogRes)

	// read blog
	fmt.Println("Read blog")

	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "asdasd"})
	if err2 != nil {
		fmt.Printf("Error happened while reading: %v \n", err2)
	}

	blogId := createBlogRes.GetBlog().GetId()
	readBlogReq := &blogpb.ReadBlogRequest{
		BlogId: blogId}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(),
		readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happended while reading: %v \n", readBlogErr)
	}

	fmt.Printf("Blog was read: %v \n", readBlogRes)

	// update Blog
	newBlog := &blogpb.Blog{
		Id:       blogId,
		AuthorId: "Change author",
		Title:    "Word breaker (edited)",
		Content:  "its sec blog content",
	}

	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: newBlog,
	})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf("Blog was updated: %v \n", updateRes)

	// delete Blog
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: blogId,
	})
	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", deleteErr)
	}
	fmt.Printf("Blog was deleted : %v", deleteRes)

	// list blog
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling GreetManyTImes RPC : %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Println(res.GetBlog())
	}
}
