package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/linus5304/social/internal/store"
)

var usernames = []string{
	"whisperingwillow",
	"cosmicwanderer",
	"silentstorm",
	"fadingember",
	"goldendawn",
	"mysticmoon",
	"crimsondusk",
	"emeraldream",
	"azuresky",
	"midnightsun",
	"starrynight",
	"oceanicdepth",
	"forestwhisper",
	"desertmirage",
	"mountainpeak",
	"riverflow",
	"thundercloud",
	"rainbowbright",
	"fireflyspark",
	"butterflykiss",
	"dragonflight",
	"phoenixrise",
	"unicorndream",
	"mermaidsong",
	"goblingrin",
	"trolltrouble",
	"fairyflight",
	"wizardwisdom",
	"witchyways",
	"ghostlygaze",
	"vampireveil",
	"werewolfhowl",
	"alienarrival",
	"robotreign",
	"cyberninja",
	"pixelpirate",
	"gamergirl",
	"coderkid",
	"hackerhero",
	"geeksquad",
	"bookwormbrain",
	"writer'sworld",
	"artistryavenue",
	"musicmaestro",
	"dancerdream",
	"chef'skiss",
	"baker'sblessing",
	"gardener'sgrace",
	"builder'sbrain",
	"dreamer'sden",
}

var titles = []string{
	"Unveiling the Secrets of",
	"The Ultimate Guide to",
	"Revolutionizing with",
	"5 Surprising Ways to",
	"The Future of: A Bold Prediction",
	"Breaking Down the Barriers:",
	"From Zero to Hero: A Beginner's Guide to",
	"The Little-Known Trick to",
	"Why Everyone's Talking About",
	"The Science Behind",
	"The Dark Side of",
	"The Art of",
	"The Untold Story of",
	"The 10 Commandments of",
	"The Ultimate Showdown:",
	"The Future of Work: Is the Answer?",
	"The Psychology of",
	"The 7 Habits of Highly Effective",
	"The Power of",
	"The Unexpected Benefits of",
}

var content = []string{
	"Discover how artificial intelligence is reshaping industries and our daily lives.",
	"Learn the fundamentals of machine learning and how to build intelligent systems.",
	"Explore how blockchain technology is enhancing security and protecting sensitive data.",
	"Implement these simple techniques to maximize your efficiency and achieve more.",
	"Dive into the future of cloud technology and its impact on businesses and individuals.",
	"Address the challenges of the tech skills gap and explore solutions to bridge the divide.",
	"Start your journey as a web developer with this comprehensive guide to HTML, CSS, and JavaScript.",
	"Discover a secret technique to improve your data analysis skills and gain valuable insights.",
	"Understand the potential of quantum computing and its revolutionary implications.",
	"Explore the technology behind connected devices and how they're transforming our world.",
	"Uncover the negative impacts of social media and how to protect yourself from its harmful effects.",
	"Learn how to communicate clearly, persuasively, and empathetically.",
	"Discover the fascinating history of computing and the pioneers who shaped the digital age.",
	"Follow these essential principles to write clean, efficient, and maintainable code.",
	"Compare and contrast these popular programming languages to determine the best fit for your projects.",
	"Explore the benefits and challenges of remote work and its impact on the future of employment.",
	"Understand the psychological principles that influence user behavior and design better user interfaces.",
	"Adopt these habits to become a more productive and successful software developer.",
	"Learn how mindfulness can improve your focus, reduce stress, and enhance your coding skills.",
	"Discover the surprising advantages of expanding your programming skillset.",
}

var tags = []string{
	"AI", "Machine Learning", "Data Science", "Cybersecurity",
	"Cloud Computing", "DevOps", "Software Engineering",
	"Web Development", "Mobile Development", "Game Development",
	"Blockchain", "Cryptocurrency", "Artificial Intelligence",
	"Internet of Things", "Big Data", "Quantum Computing",
	"Virtual Reality", "Augmented Reality", "Robotics",
	"Biotechnology",
}

var comments = []string{
	"This is a great article! I learned a lot .",
	"I disagree with the author's point .",
	"Can you elaborate on [specific point]?",
	"This is a helpful resource for beginners.",
	"I'm still confused about [specific concept].",
	"I'd like to see more examples of [specific technique].",
	"This article is well-written and easy to understand.",
	"I have a question about [specific question].",
	"I'm looking forward to more articles on this topic.",
	"This is a great starting point for further research.",
	"I'd like to know more about the future of [topic].",
	"This article has changed my perspective on [topic].",
	"I'm impressed by the author's knowledge of [topic].",
	"I'm disappointed that the author didn't address [specific issue].",
	"I'd like to see more practical applications of [topic].",
	"This article is a valuable contribution to the field of [field].",
	"I'm inspired to learn more .",
	"I'd like to see a follow-up article on [related topic].",
	"This article is a must-read for anyone interested in [topic].",
	"I'm looking forward to the author's next article.",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user: ", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "12345",
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
			Content: content[rand.Intn(len(content))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms
}
