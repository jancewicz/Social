package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/jancewicz/social/internal/store"
)

var usernames = []string{
	"bob", "thomas", "paul", "alice", "jane", "johnny",
	"mike", "susan", "dave", "carla", "rachel", "tony",
	"linda", "mark", "kimberly", "chris", "emma", "leo",
	"jack", "sophie", "max", "anna", "peter", "lucy",
	"gregory", "nancy", "oliver", "harry", "george", "tina",
	"frank", "sam", "danny", "sara", "louis", "kevin",
	"brenda", "charlie", "julia", "robert", "maggie", "vicky",
	"derek", "stacy", "greg", "victoria", "mario", "josh",
	"harvey", "clara",
}

var titles = []string{
	"Learning Go Basics",
	"Understanding Goroutines",
	"Mastering Interfaces in Go",
	"Building REST APIs with Go",
	"Go Concurrency Made Easy",
	"Error Handling in Go",
	"Introduction to Go Modules",
	"Working with JSON in Go",
	"Writing Unit Tests in Go",
	"Using Context in Go",
	"Building CLI Tools with Go",
	"Exploring Go Slices",
	"Effective Go Logging",
	"Go Channels Deep Dive",
	"Creating HTTP Servers in Go",
	"Deploying Go Apps to the Cloud",
	"Optimizing Go Code for Performance",
	"Go Structs and Methods Explained",
	"Reading and Writing Files in Go",
	"Getting Started with Websockets in Go",
}

var contents = []string{
	"Go is a powerful language for building fast and scalable applications.",
	"Goroutines make it easy to handle concurrency in Go.",
	"Interfaces in Go allow for flexible and reusable code.",
	"Learn how to build robust REST APIs with Go's net/http package.",
	"Go concurrency helps you handle multiple tasks efficiently.",
	"Error handling in Go ensures your program is reliable.",
	"Go modules simplify dependency management for your projects.",
	"JSON encoding and decoding are straightforward with Go's standard library.",
	"Unit testing in Go is simple and encourages clean code.",
	"Context in Go helps manage deadlines and cancellations.",
	"Go is a great choice for creating command-line tools.",
	"Slices in Go provide powerful ways to work with collections.",
	"Logging is essential for debugging and monitoring Go applications.",
	"Channels are a key feature for communication between goroutines.",
	"Build scalable HTTP servers with Go's lightweight frameworks.",
	"Deploying Go applications to the cloud is seamless with Docker.",
	"Go's performance can be optimized using profiling tools.",
	"Structs in Go are the backbone of defining complex data models.",
	"File operations in Go are simple with the io and os packages.",
	"Websockets in Go enable real-time communication with minimal effort.",
}

var tags = []string{
	"golang", "programming", "webdev", "backend", "concurrency",
	"restapi", "json", "cli", "testing", "cloud",
	"devops", "performance", "logging", "tutorial", "opensource",
	"frameworks", "microservices", "databases", "websockets", "deployment",
}

var commentsArr = []string{
	"Great post! Learned a lot about Go concurrency.",
	"I've been struggling with Go channels, this really helped.",
	"Fantastic explanation of Go interfaces, very clear!",
	"Thanks for the tips on testing in Go, I’ll use them in my next project.",
	"Your example of working with JSON in Go is super helpful.",
	"Go modules are such a game changer for managing dependencies, thanks for the info!",
	"Could you write a follow-up post on deploying Go apps to AWS?",
	"I had no idea Go made error handling this easy, thanks for sharing.",
	"Your article on structs and methods in Go cleared up a lot of confusion.",
	"Really enjoyed the section on building REST APIs with Go, looking forward to more posts!",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error on creating user: ", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error on creating post: ", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error on creating comment: ", err)
			return
		}
	}

	log.Panicln("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123456",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: titles[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		comments[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: commentsArr[rand.Intn(len(comments))],
		}
	}

	return comments
}
