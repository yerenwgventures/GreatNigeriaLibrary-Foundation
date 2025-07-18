package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Info("Starting Great Nigeria React Frontend Server")

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		logger.Warn("Error loading .env file, using environment variables")
	}

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))




	// Serve static assets (e.g., /static/js, /static/css)
        router.Static("/static", "./great-nigeria-frontend/build/static")

        // Serve root HTML directly
       //  router.StaticFile("/", "./great-nigeria-frontend/build/index.html")



	// Serve other static files from the React build directory
	router.StaticFile("/asset-manifest.json", "./great-nigeria-frontend/build/asset-manifest.json")
	router.StaticFile("/manifest.json", "./great-nigeria-frontend/build/manifest.json")
	router.StaticFile("/favicon.ico", "./great-nigeria-frontend/build/favicon.ico")
	router.StaticFile("/logo192.png", "./great-nigeria-frontend/build/logo192.png")
	router.StaticFile("/logo512.png", "./great-nigeria-frontend/build/logo512.png")

	// Add a simple test route to see if routing is working
	router.GET("/test-route", func(c *gin.Context) {
		c.String(http.StatusOK, "Test route is working!")
	})

	// Serve the React app at the root path and other paths
	router.GET("/", serveReactApp)
	router.GET("/app", serveReactApp)
	router.GET("/react-app", serveReactApp)
	router.GET("/book-viewer", serveReactApp)

	// API routes with mock data based on the original API structure
	api := router.Group("/api")
	{
		// Health check endpoint
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now().Format(time.RFC3339),
			})
		})

		// Books API
		api.GET("/books", func(c *gin.Context) {
			c.JSON(http.StatusOK, []gin.H{
				{
					"id":          1,
					"title":       "Great Nigeria – Awakening the Giant",
					"author":      "Great Nigeria Network",
					"description": "A comprehensive manifesto that diagnoses the root causes of Nigeria's challenges and calls for unified citizen action to transform the nation.",
					"coverImage":  "/static/img/book1-cover.jpg",
					"published":   true,
					"createdAt":   time.Now().AddDate(0, -1, 0),
					"updatedAt":   time.Now(),
				},
				{
					"id":          2,
					"title":       "Great Nigeria – The Masterplan",
					"author":      "Great Nigeria Network",
					"description": "A detailed implementation plan for transforming Nigeria through citizen action and institutional reform.",
					"coverImage":  "/static/img/book2-cover.jpg",
					"published":   true,
					"createdAt":   time.Now().AddDate(0, -1, 0),
					"updatedAt":   time.Now(),
				},
			})
		})

		// Get book by ID
		api.GET("/books/:id", func(c *gin.Context) {
			id := c.Param("id")
			if id == "1" {
				c.JSON(http.StatusOK, gin.H{
					"id":          1,
					"title":       "Great Nigeria – Awakening the Giant",
					"author":      "Great Nigeria Network",
					"description": "A comprehensive manifesto that diagnoses the root causes of Nigeria's challenges and calls for unified citizen action to transform the nation.",
					"coverImage":  "/static/img/book1-cover.jpg",
					"published":   true,
					"createdAt":   time.Now().AddDate(0, -1, 0),
					"updatedAt":   time.Now(),
				})
			} else if id == "2" {
				c.JSON(http.StatusOK, gin.H{
					"id":          2,
					"title":       "Great Nigeria – The Masterplan",
					"author":      "Great Nigeria Network",
					"description": "A detailed implementation plan for transforming Nigeria through citizen action and institutional reform.",
					"coverImage":  "/static/img/book2-cover.jpg",
					"published":   true,
					"createdAt":   time.Now().AddDate(0, -1, 0),
					"updatedAt":   time.Now(),
				})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			}
		})

		// Get book chapters
		api.GET("/books/:id/chapters", func(c *gin.Context) {
			id := c.Param("id")
			if id == "1" {
				c.JSON(http.StatusOK, []gin.H{
					{
						"id":          1,
						"bookId":      1,
						"title":       "The Bleeding Giant & Ghosts of the Past",
						"number":      1,
						"description": "This chapter vividly illustrates Nigeria's vast, squandered potential and begins unearthing the deep historical roots that continue to bind the nation.",
						"published":   true,
					},
					{
						"id":          2,
						"bookId":      1,
						"title":       "Governance Failures & Institutional Decay",
						"number":      2,
						"description": "An analysis of the systemic governance challenges that have prevented Nigeria from realizing its potential.",
						"published":   true,
					},
				})
			} else {
				c.JSON(http.StatusOK, []gin.H{})
			}
		})

		// Auth API
		api.POST("/auth/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"token": "mock-jwt-token",
				"user": gin.H{
					"id":    1,
					"name":  "Demo User",
					"email": "demo@example.com",
					"role":  "user",
				},
			})
		})

		api.POST("/auth/register", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"token": "mock-jwt-token",
				"user": gin.H{
					"id":    1,
					"name":  "Demo User",
					"email": "demo@example.com",
					"role":  "user",
				},
			})
		})

		api.GET("/auth/me", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"id":    1,
				"name":  "Demo User",
				"email": "demo@example.com",
				"role":  "user",
			})
		})

		// Celebrate Nigeria API
		api.GET("/celebrate/featured", func(c *gin.Context) {
			c.JSON(http.StatusOK, []gin.H{
				{
					"id":          "1", // String ID to match what React expects
					"type":        "person",
					"name":        "Wole Soyinka",
					"slug":        "wole-soyinka",
					"description": "Nobel Prize-winning playwright, poet, and essayist",
					"image":       "/static/img/celebrate/wole-soyinka.jpg",
					"votes":       120,
				},
				{
					"id":          "2", // String ID to match what React expects
					"type":        "place",
					"name":        "Osun Sacred Grove",
					"slug":        "osun-sacred-grove",
					"description": "UNESCO World Heritage site and home to the goddess of fertility Osun",
					"image":       "/static/img/celebrate/osun-grove.jpg",
					"votes":       85,
				},
				{
					"id":          "3", // String ID to match what React expects
					"type":        "innovation",
					"name":        "Fintech Revolution",
					"slug":        "fintech-revolution",
					"description": "Nigeria's pioneering role in Africa's financial technology transformation",
					"image":       "/static/img/celebrate/fintech.jpg",
					"votes":       95,
				},
				{
					"id":          "4", // String ID to match what React expects
					"type":        "culture",
					"name":        "Nollywood",
					"slug":        "nollywood",
					"description": "The world's second-largest film industry by volume",
					"image":       "/static/img/celebrate/nollywood.jpg",
					"votes":       110,
				},
			})
		})

		// Get entries by type
		api.GET("/celebrate/types/:type", func(c *gin.Context) {
			entryType := c.Param("type")

			if entryType == "person" {
				c.JSON(http.StatusOK, []gin.H{
					{
						"id":          "1",
						"type":        "person",
						"name":        "Wole Soyinka",
						"slug":        "wole-soyinka",
						"description": "Nobel Prize-winning playwright, poet, and essayist",
						"image":       "/static/img/celebrate/wole-soyinka.jpg",
						"votes":       120,
					},
					{
						"id":          "5",
						"type":        "person",
						"name":        "Chimamanda Ngozi Adichie",
						"slug":        "chimamanda-adichie",
						"description": "Award-winning author and feminist advocate",
						"image":       "/static/img/celebrate/chimamanda.jpg",
						"votes":       115,
					},
					{
						"id":          "6",
						"type":        "person",
						"name":        "Aliko Dangote",
						"slug":        "aliko-dangote",
						"description": "Africa's richest man and industrial magnate",
						"image":       "/static/img/celebrate/dangote.jpg",
						"votes":       90,
					},
				})
			} else if entryType == "place" {
				c.JSON(http.StatusOK, []gin.H{
					{
						"id":          "2",
						"type":        "place",
						"name":        "Osun Sacred Grove",
						"slug":        "osun-sacred-grove",
						"description": "UNESCO World Heritage site and home to the goddess of fertility Osun",
						"image":       "/static/img/celebrate/osun-grove.jpg",
						"votes":       85,
					},
					{
						"id":          "7",
						"type":        "place",
						"name":        "Yankari Game Reserve",
						"slug":        "yankari-game-reserve",
						"description": "Nigeria's richest wildlife oasis",
						"image":       "/static/img/celebrate/yankari.jpg",
						"votes":       75,
					},
				})
			} else if entryType == "innovation" {
				c.JSON(http.StatusOK, []gin.H{
					{
						"id":          "3",
						"type":        "innovation",
						"name":        "Fintech Revolution",
						"slug":        "fintech-revolution",
						"description": "Nigeria's pioneering role in Africa's financial technology transformation",
						"image":       "/static/img/celebrate/fintech.jpg",
						"votes":       95,
					},
				})
			} else if entryType == "culture" {
				c.JSON(http.StatusOK, []gin.H{
					{
						"id":          "4",
						"type":        "culture",
						"name":        "Nollywood",
						"slug":        "nollywood",
						"description": "The world's second-largest film industry by volume",
						"image":       "/static/img/celebrate/nollywood.jpg",
						"votes":       110,
					},
				})
			} else {
				c.JSON(http.StatusOK, []gin.H{})
			}
		})

		// Get celebrate entry by type and slug
		api.GET("/celebrate/:type/:slug", func(c *gin.Context) {
			entryType := c.Param("type")
			slug := c.Param("slug")

			if entryType == "person" && slug == "wole-soyinka" {
				c.JSON(http.StatusOK, gin.H{
					"id":           "1",
					"type":         "person",
					"name":         "Wole Soyinka",
					"slug":         "wole-soyinka",
					"description":  "Nobel Prize-winning playwright, poet, and essayist",
					"full_content": "Wole Soyinka is a Nigerian playwright, novelist, poet, and essayist who won the Nobel Prize for Literature in 1986. Born on July 13, 1934, in Abeokuta, Nigeria, Soyinka has been a strong critic of successive Nigerian governments, especially the country's many military dictators, as well as other political tyrannies, including the Mugabe regime in Zimbabwe. Much of his writing has been concerned with the oppressive boot and the irrelevance of the color of the foot that wears it.\n\nSoyinka was educated at the University of Leeds and has held fellowships and professorships at several universities worldwide, including Yale, Cornell, and Oxford. His works often blend traditional Yoruba folklore with Western literary traditions, creating a unique style that has influenced generations of African writers.\n\nHis notable works include the plays 'Death and the King's Horseman' and 'A Dance of the Forests,' the novel 'The Interpreters,' and the autobiographical work 'Aké: The Years of Childhood.' Through his writing and activism, Soyinka has consistently advocated for human rights, democracy, and social justice in Nigeria and across Africa.",
					"image":        "/static/img/celebrate/wole-soyinka.jpg",
					"votes":        120,
					"created_at":   time.Now().AddDate(0, -2, 0),
					"related_entries": []gin.H{
						{
							"id":    "5",
							"type":  "person",
							"name":  "Chimamanda Ngozi Adichie",
							"slug":  "chimamanda-adichie",
							"image": "/static/img/celebrate/chimamanda.jpg",
						},
						{
							"id":    "4",
							"type":  "culture",
							"name":  "Nollywood",
							"slug":  "nollywood",
							"image": "/static/img/celebrate/nollywood.jpg",
						},
					},
				})
			} else if entryType == "place" && slug == "osun-sacred-grove" {
				c.JSON(http.StatusOK, gin.H{
					"id":           "2",
					"type":         "place",
					"name":         "Osun Sacred Grove",
					"slug":         "osun-sacred-grove",
					"description":  "UNESCO World Heritage site and home to the goddess of fertility Osun",
					"full_content": "The Osun Sacred Grove is a dense forest situated on the outskirts of Osogbo, the capital city of Osun State in southwestern Nigeria. This sacred grove is one of the last remnants of primary high forest in southern Nigeria and is regarded as the abode of the goddess of fertility Osun, one of the pantheon of Yoruba gods.\n\nThe grove, which is now a UNESCO World Heritage Site, contains sanctuaries, shrines, sculptures, and art works erected in honor of Osun and other Yoruba deities. Many of the sculptures were created by Austrian artist Susanne Wenger (who later became a Yoruba priestess) and her New Sacred Art movement in the 1950s and 1960s.\n\nThe annual Osun-Osogbo festival, which attracts thousands of Osun worshippers, tourists, and spectators from all over the world, takes place in the grove. During this festival, the Arugba (a virgin maiden) carries a calabash containing sacrificial materials to the Osun River, accompanied by a large procession of people singing, dancing, and praying for blessings.\n\nBeyond its cultural and spiritual significance, the Osun Sacred Grove is also an important biodiversity conservation site, hosting diverse flora and fauna, including endangered species. It represents a remarkable example of a cultural landscape that illustrates the adaptation of traditional beliefs and practices to environmental challenges.",
					"image":        "/static/img/celebrate/osun-grove.jpg",
					"votes":        85,
					"created_at":   time.Now().AddDate(0, -3, 0),
					"related_entries": []gin.H{
						{
							"id":    "7",
							"type":  "place",
							"name":  "Yankari Game Reserve",
							"slug":  "yankari-game-reserve",
							"image": "/static/img/celebrate/yankari.jpg",
						},
					},
				})
			} else if entryType == "innovation" && slug == "fintech-revolution" {
				c.JSON(http.StatusOK, gin.H{
					"id":           "3",
					"type":         "innovation",
					"name":         "Fintech Revolution",
					"slug":         "fintech-revolution",
					"description":  "Nigeria's pioneering role in Africa's financial technology transformation",
					"full_content": "Nigeria has emerged as the epicenter of Africa's fintech revolution, with Lagos often referred to as the continent's 'Silicon Valley.' This transformation has been driven by a combination of factors: a large unbanked population, high mobile phone penetration, a youthful tech-savvy demographic, and regulatory support for financial innovation.\n\nCompanies like Paystack (acquired by Stripe for over $200 million), Flutterwave (valued at over $1 billion), and Interswitch have led the charge, creating solutions that address uniquely African challenges while meeting global standards for financial technology. These platforms have revolutionized payment processing, making digital transactions accessible to millions of Nigerians previously excluded from the formal financial system.\n\nBeyond payments, Nigerian fintech innovations have expanded into digital banking, lending, wealth management, and cryptocurrency. Digital banks like Kuda and Carbon are challenging traditional banking models, while platforms like PiggyVest and Cowrywise are democratizing access to savings and investment opportunities.\n\nThe impact of Nigeria's fintech revolution extends beyond convenience—it's driving financial inclusion, reducing corruption through transparent digital transactions, creating thousands of high-skilled jobs, and attracting significant foreign investment. The sector has become a bright spot in Nigeria's economy, demonstrating how technological innovation can address developmental challenges while creating economic opportunities.",
					"image":        "/static/img/celebrate/fintech.jpg",
					"votes":        95,
					"created_at":   time.Now().AddDate(0, -1, 0),
				})
			} else if entryType == "culture" && slug == "nollywood" {
				c.JSON(http.StatusOK, gin.H{
					"id":           "4",
					"type":         "culture",
					"name":         "Nollywood",
					"slug":         "nollywood",
					"description":  "The world's second-largest film industry by volume",
					"full_content": "Nollywood, Nigeria's film industry, has grown from humble beginnings in the early 1990s to become the world's second-largest film industry by volume, producing approximately 2,500 films annually. What began with the success of Kenneth Nnebue's 'Living in Bondage' (1992)—a direct-to-video release that sold over a million copies—has evolved into a global cultural phenomenon that generates an estimated $1 billion annually and employs over one million Nigerians.\n\nNollywood's success stems from its authentic storytelling that resonates with African audiences and the diaspora. The films often address relevant social issues, family dynamics, cultural traditions, and contemporary challenges facing Nigerian society. This cultural authenticity has helped Nollywood films find audiences across Africa and beyond, becoming one of Nigeria's most significant cultural exports.\n\nThe industry has evolved significantly over the years. While early Nollywood was characterized by low-budget productions with quick turnaround times, today's industry includes a growing segment of high-production-value films. Directors like Kemi Adetiba ('King of Boys'), Kunle Afolayan ('Citation'), and Genevieve Nnaji ('Lionheart'—Netflix's first Nigerian original film) are creating content that meets international standards while maintaining distinctly Nigerian narratives.\n\nStreaming platforms like Netflix, Amazon Prime, and local services like IrokoTV have further expanded Nollywood's global reach, bringing Nigerian stories to international audiences and creating new revenue streams for filmmakers. The industry continues to innovate, with recent expansions into series production, animation, and documentaries.\n\nNollywood represents more than entertainment—it's a powerful vehicle for cultural diplomacy, shaping global perceptions of Nigeria and Africa while creating economic opportunities and preserving Nigerian stories for future generations.",
					"image":        "/static/img/celebrate/nollywood.jpg",
					"votes":        110,
					"created_at":   time.Now().AddDate(0, -1, -15),
					"related_entries": []gin.H{
						{
							"id":    "1",
							"type":  "person",
							"name":  "Wole Soyinka",
							"slug":  "wole-soyinka",
							"image": "/static/img/celebrate/wole-soyinka.jpg",
						},
					},
				})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
			}
		})

		// Vote for a celebrate entry
		api.POST("/celebrate/:id/vote", func(c *gin.Context) {
			id := c.Param("id")
			var requestBody struct {
				Direction string `json:"direction"` // "up" or "down"
			}

			if err := c.BindJSON(&requestBody); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			// In a real implementation, we would update the vote count in the database
			// For now, we'll just return a success response with a mock updated vote count
			var newVotes int
			if requestBody.Direction == "up" {
				newVotes = 121 // Increment from 120
			} else {
				newVotes = 119 // Decrement from 120
			}

			c.JSON(http.StatusOK, gin.H{
				"id":    id,
				"votes": newVotes,
			})
		})

		// Forum API
		api.GET("/forum/categories", func(c *gin.Context) {
			c.JSON(http.StatusOK, []gin.H{
				{
					"id":           "1", // String ID to match what React expects
					"name":         "General Discussion",
					"description":  "General discussions about Nigeria's transformation",
					"topics_count": 24, // Changed to match React's expected property name
				},
				{
					"id":           "2", // String ID to match what React expects
					"name":         "Book Discussions",
					"description":  "Discussions related to the Great Nigeria books",
					"topics_count": 15, // Changed to match React's expected property name
				},
			})
		})

		// Get topics by category
		api.GET("/forum/categories/:id/topics", func(c *gin.Context) {
			categoryId := c.Param("id")
			if categoryId == "1" || categoryId == "all" {
				c.JSON(http.StatusOK, []gin.H{
					{
						"id":            "1", // String ID to match what React expects
						"category_id":   "1", // Changed to match React's expected property name
						"title":         "How can we improve Nigeria's education system?",
						"content":       "I believe education is the foundation for national development. What practical steps can we take to improve Nigeria's education system?",
						"replies_count": 5,                            // Changed to match React's expected property name
						"views_count":   120,                          // Changed to match React's expected property name
						"created_at":    time.Now().AddDate(0, 0, -5), // Changed to match React's expected property name
						"author": gin.H{ // Nested author object to match React's expected structure
							"id":   "1",
							"name": "Demo User",
						},
					},
					{
						"id":            "2", // String ID to match what React expects
						"category_id":   "1", // Changed to match React's expected property name
						"title":         "Infrastructure development priorities",
						"content":       "Which infrastructure projects should be prioritized for Nigeria's development?",
						"replies_count": 8,                            // Changed to match React's expected property name
						"views_count":   95,                           // Changed to match React's expected property name
						"created_at":    time.Now().AddDate(0, 0, -3), // Changed to match React's expected property name
						"author": gin.H{ // Nested author object to match React's expected structure
							"id":   "2",
							"name": "Jane Smith",
						},
					},
				})
			} else if categoryId == "2" {
				c.JSON(http.StatusOK, []gin.H{
					{
						"id":            "3", // String ID to match what React expects
						"category_id":   "2", // Changed to match React's expected property name
						"title":         "Book 1 Chapter 3 Discussion: Economic Transformation",
						"content":       "What are your thoughts on the economic transformation strategies outlined in Chapter 3?",
						"replies_count": 12,                           // Changed to match React's expected property name
						"views_count":   150,                          // Changed to match React's expected property name
						"created_at":    time.Now().AddDate(0, 0, -7), // Changed to match React's expected property name
						"author": gin.H{ // Nested author object to match React's expected structure
							"id":   "3",
							"name": "John Doe",
						},
					},
					{
						"id":            "4", // String ID to match what React expects
						"category_id":   "2", // Changed to match React's expected property name
						"title":         "Implementing the ideas from Book 2",
						"content":       "How can we start implementing some of the practical ideas from Book 2 in our local communities?",
						"replies_count": 7,                            // Changed to match React's expected property name
						"views_count":   85,                           // Changed to match React's expected property name
						"created_at":    time.Now().AddDate(0, 0, -2), // Changed to match React's expected property name
						"author": gin.H{ // Nested author object to match React's expected structure
							"id":   "4",
							"name": "Sarah Johnson",
						},
					},
				})
			} else {
				c.JSON(http.StatusOK, []gin.H{})
			}
		})

		// Get topic by ID
		api.GET("/forum/topics/:id", func(c *gin.Context) {
			topicId := c.Param("id")
			if topicId == "1" {
				c.JSON(http.StatusOK, gin.H{
					"id":            "1",
					"category_id":   "1",
					"title":         "How can we improve Nigeria's education system?",
					"content":       "I believe education is the foundation for national development. What practical steps can we take to improve Nigeria's education system?",
					"replies_count": 5,
					"views_count":   120,
					"created_at":    time.Now().AddDate(0, 0, -5),
					"author": gin.H{
						"id":   "1",
						"name": "Demo User",
					},
					"replies": []gin.H{
						{
							"id":         "101",
							"topic_id":   "1",
							"content":    "We need to invest more in teacher training and development. Quality teachers are the backbone of any education system.",
							"created_at": time.Now().AddDate(0, 0, -4),
							"votes":      15,
							"author": gin.H{
								"id":   "2",
								"name": "Jane Smith",
							},
						},
						{
							"id":         "102",
							"topic_id":   "1",
							"content":    "Infrastructure is also critical. Many schools lack basic facilities like proper classrooms, libraries, and laboratories.",
							"created_at": time.Now().AddDate(0, 0, -3),
							"votes":      8,
							"author": gin.H{
								"id":   "3",
								"name": "John Doe",
							},
						},
						{
							"id":         "103",
							"topic_id":   "1",
							"content":    "Curriculum reform is needed to focus more on critical thinking and practical skills rather than rote memorization.",
							"created_at": time.Now().AddDate(0, 0, -2),
							"votes":      12,
							"author": gin.H{
								"id":   "4",
								"name": "Sarah Johnson",
							},
						},
					},
				})
			} else if topicId == "2" {
				c.JSON(http.StatusOK, gin.H{
					"id":            "2",
					"category_id":   "1",
					"title":         "Infrastructure development priorities",
					"content":       "Which infrastructure projects should be prioritized for Nigeria's development?",
					"replies_count": 8,
					"views_count":   95,
					"created_at":    time.Now().AddDate(0, 0, -3),
					"author": gin.H{
						"id":   "2",
						"name": "Jane Smith",
					},
					"replies": []gin.H{
						{
							"id":         "201",
							"topic_id":   "2",
							"content":    "Power infrastructure should be the top priority. Reliable electricity is fundamental to all other development.",
							"created_at": time.Now().AddDate(0, 0, -2),
							"votes":      20,
							"author": gin.H{
								"id":   "1",
								"name": "Demo User",
							},
						},
						{
							"id":         "202",
							"topic_id":   "2",
							"content":    "Transportation networks are equally important. We need better roads, railways, and ports to facilitate trade and movement.",
							"created_at": time.Now().AddDate(0, 0, -1),
							"votes":      15,
							"author": gin.H{
								"id":   "3",
								"name": "John Doe",
							},
						},
					},
				})
			} else if topicId == "3" {
				c.JSON(http.StatusOK, gin.H{
					"id":            "3",
					"category_id":   "2",
					"title":         "Book 1 Chapter 3 Discussion: Economic Transformation",
					"content":       "What are your thoughts on the economic transformation strategies outlined in Chapter 3?",
					"replies_count": 12,
					"views_count":   150,
					"created_at":    time.Now().AddDate(0, 0, -7),
					"author": gin.H{
						"id":   "3",
						"name": "John Doe",
					},
					"replies": []gin.H{
						{
							"id":         "301",
							"topic_id":   "3",
							"content":    "The emphasis on diversification away from oil dependence is crucial. We've been talking about this for decades but haven't made enough progress.",
							"created_at": time.Now().AddDate(0, 0, -6),
							"votes":      18,
							"author": gin.H{
								"id":   "1",
								"name": "Demo User",
							},
						},
						{
							"id":         "302",
							"topic_id":   "3",
							"content":    "I particularly liked the section on developing the agricultural value chain. This could create millions of jobs and improve food security.",
							"created_at": time.Now().AddDate(0, 0, -5),
							"votes":      22,
							"author": gin.H{
								"id":   "2",
								"name": "Jane Smith",
							},
						},
					},
				})
			} else if topicId == "4" {
				c.JSON(http.StatusOK, gin.H{
					"id":            "4",
					"category_id":   "2",
					"title":         "Implementing the ideas from Book 2",
					"content":       "How can we start implementing some of the practical ideas from Book 2 in our local communities?",
					"replies_count": 7,
					"views_count":   85,
					"created_at":    time.Now().AddDate(0, 0, -2),
					"author": gin.H{
						"id":   "4",
						"name": "Sarah Johnson",
					},
					"replies": []gin.H{
						{
							"id":         "401",
							"topic_id":   "4",
							"content":    "Community action cells as described in Chapter 5 could be a good starting point. We could form small groups focused on specific local issues.",
							"created_at": time.Now().AddDate(0, 0, -1),
							"votes":      10,
							"author": gin.H{
								"id":   "3",
								"name": "John Doe",
							},
						},
					},
				})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
			}
		})

		// Resource API
		api.GET("/resources/categories", func(c *gin.Context) {
			c.JSON(http.StatusOK, []gin.H{
				{
					"id":          1,
					"name":        "Educational Materials",
					"description": "Resources for learning and education",
				},
				{
					"id":          2,
					"name":        "Community Development",
					"description": "Resources for community development projects",
				},
			})
		})
	}

	// Add routes for React app client-side routing
	// router.GET("/book-viewer", serveReactApp)
	router.GET("/book-viewer/*path", serveReactApp)
	router.GET("/books", serveReactApp)
	router.GET("/community", serveReactApp)
	router.GET("/community/*path", serveReactApp)
	router.GET("/celebrate", serveReactApp)
	router.GET("/celebrate/*path", serveReactApp)
	router.GET("/resources", serveReactApp)
	router.GET("/about", serveReactApp)
	router.GET("/contact", serveReactApp)
	router.GET("/login", serveReactApp)
	router.GET("/register", serveReactApp)
	router.GET("/profile", serveReactApp)

	// Add a catch-all route for the React app
	router.NoRoute(func(c *gin.Context) {
		// For all paths that don't start with /api/ or /static/, serve the React app
		if !strings.HasPrefix(c.Request.URL.Path, "/api/") && !strings.HasPrefix(c.Request.URL.Path, "/static/") {
			serveReactApp(c)
		}
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Print a message to the console
	fmt.Println("Starting React Frontend server on port", port)
	logger.Infof("Server running on port %s", port)

	if err := router.Run(fmt.Sprintf("0.0.0.0:%s", port)); err != nil {
		// Print the error to the console
		fmt.Println("Failed to start server:", err)
		logger.WithError(err).Fatal("Failed to start server")
	}
}

// serveReactApp serves the React app's index.html file
func serveReactApp(c *gin.Context) {
	// Always serve the index.html file for client-side routing
	c.File("./great-nigeria-frontend/build/index.html")
}