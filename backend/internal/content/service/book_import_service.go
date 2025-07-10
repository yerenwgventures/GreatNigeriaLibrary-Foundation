package service

import (
        "encoding/json"
        "log"
        "strconv"
        "strings"
        "time"

        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/models"
        "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/content/repository"
)

// BookImportService defines the interface for book import operations
type BookImportService interface {
        ImportBook1() error
}

// BookImportServiceImpl implements the BookImportService interface
type BookImportServiceImpl struct {
        bookRepo repository.BookRepository
}

// NewBookImportService creates a new book import service
func NewBookImportService(bookRepo repository.BookRepository) BookImportService {
        return &BookImportServiceImpl{
                bookRepo: bookRepo,
        }
}

// createInteractiveElements creates interactive elements for a section
func (s *BookImportServiceImpl) createInteractiveElements(bookID, sectionID uint, sectionTitle string) error {
        // Create appropriate interactive elements based on section title
        var elements []models.InteractiveElement
        
        // Check section title to determine what type of interactive elements to create
        if strings.Contains(strings.ToLower(sectionTitle), "forum topic") {
                // Create discussion prompts
                discussionContent := models.DiscussionPromptContent{
                        Topic:         "Community Perspectives on " + sectionTitle,
                        InitialPrompt: "Share your thoughts and experiences related to this topic. How does this issue manifest in your local community?",
                        SupportingPoints: []string{
                                "Consider both historical and contemporary examples",
                                "Reflect on potential solutions or interventions",
                                "Think about your personal role in addressing this challenge",
                        },
                        Guidelines: "Please keep responses respectful and constructive. The goal is to build understanding and identify potential paths forward.",
                }
                
                contentJSON, err := json.Marshal(discussionContent)
                if err != nil {
                        return err
                }
                
                elements = append(elements, models.InteractiveElement{
                        SectionID:      sectionID,
                        Position:       1,
                        Type:           models.DiscussionPromptType,
                        Title:          "Join the Discussion",
                        Description:    "Share your perspectives and experiences related to the topics in this section.",
                        Content:        string(contentJSON),
                        CompletionType: "no-check",
                        PointsValue:    10,
                        RequiredStatus: false,
                })
        } else if strings.Contains(strings.ToLower(sectionTitle), "actionable step") {
                // Create reflection element
                reflectionContent := models.ReflectionContent{
                        Prompt:           "Take a moment to reflect on how you can apply these action steps in your own context.",
                        GuidingQuestions: []string{
                                "Which of these action steps resonates most strongly with you, and why?",
                                "What specific challenges might you face when implementing these steps in your local context?",
                                "What resources or support would help you successfully implement these actions?",
                        },
                        MinResponseLength: 100,
                        SharingOptions:    []string{"private", "peers", "public"},
                }
                
                contentJSON, err := json.Marshal(reflectionContent)
                if err != nil {
                        return err
                }
                
                elements = append(elements, models.InteractiveElement{
                        SectionID:      sectionID,
                        Position:       1,
                        Type:           models.ReflectionType,
                        Title:          "Personal Reflection",
                        Description:    "Reflect on how these action steps apply to your situation.",
                        Content:        string(contentJSON),
                        CompletionType: "self-check",
                        PointsValue:    15,
                        RequiredStatus: false,
                })
                
                // Create call to action
                ctaContent := models.CallToActionContent{
                        ActionType:  "click",
                        Text:        "Ready to put these ideas into action? Join an implementation group focused on this area.",
                        ButtonText:  "Find Implementation Groups",
                        URL:         "/community/implementation-groups",
                        TrackingID:  "action_" + strconv.Itoa(int(sectionID)),
                }
                
                ctaJSON, err := json.Marshal(ctaContent)
                if err != nil {
                        return err
                }
                
                elements = append(elements, models.InteractiveElement{
                        SectionID:      sectionID,
                        Position:       2,
                        Type:           models.CallToActionType,
                        Title:          "Take Action Now",
                        Description:    "Connect with others who are implementing these solutions.",
                        Content:        string(ctaJSON),
                        CompletionType: "no-check",
                        PointsValue:    5,
                        RequiredStatus: false,
                })
        } else if strings.Contains(strings.ToLower(sectionTitle), "parasitic governance") || strings.Contains(strings.ToLower(sectionTitle), "the bleeding giant") {
                // Create a poll for these key concept sections
                pollContent := models.PollContent{
                        Question: "Which aspect of this section do you believe is most critical to address for Nigeria's transformation?",
                        Options: []models.PollOption{
                                {ID: "1", Text: "Leadership accountability and transparency"},
                                {ID: "2", Text: "Citizen engagement and participation"},
                                {ID: "3", Text: "Institutional reform and capacity building"},
                                {ID: "4", Text: "Education and public awareness"},
                                {ID: "5", Text: "Other (please specify in comments)"},
                        },
                        AllowMultiple: false,
                        ShowResults:   "after-vote",
                        AllowComments: true,
                }
                
                pollJSON, err := json.Marshal(pollContent)
                if err != nil {
                        return err
                }
                
                elements = append(elements, models.InteractiveElement{
                        SectionID:      sectionID,
                        Position:       1,
                        Type:           models.PollType,
                        Title:          "Community Priorities",
                        Description:    "Help us understand what aspects of these challenges you believe are most important to address.",
                        Content:        string(pollJSON),
                        CompletionType: "no-check",
                        PointsValue:    5,
                        RequiredStatus: false,
                })
        }
        
        // Create quiz for section to test understanding (general for all sections)
        if !strings.Contains(strings.ToLower(sectionTitle), "forum topic") && !strings.Contains(strings.ToLower(sectionTitle), "actionable step") {
                quizContent := models.QuizContent{
                        Questions: []models.QuizQuestion{
                                {
                                        ID:           1,
                                        QuestionText: "Based on the content in this section, which of the following best describes the Nigerian paradox?",
                                        QuestionType: "multiple-choice",
                                        Options: []models.QuizOption{
                                                {ID: "a", Text: "A country with limited natural resources but efficient governance"},
                                                {ID: "b", Text: "A nation with abundant resources and potential alongside devastating underperformance"},
                                                {ID: "c", Text: "A developing country with steady but slow progress over decades"},
                                                {ID: "d", Text: "A homogeneous society with unified national identity"},
                                        },
                                        CorrectAnswer: "b",
                                        Explanation:   "The Nigerian paradox refers to the contrast between the country's vast potential (human and natural resources) and its persistent developmental challenges.",
                                        Difficulty:    "medium",
                                },
                                {
                                        ID:           2,
                                        QuestionText: "According to the text, what is a key factor that distinguishes parasitic governance from simple corruption?",
                                        QuestionType: "multiple-choice",
                                        Options: []models.QuizOption{
                                                {ID: "a", Text: "It affects only specific sectors rather than the entire government"},
                                                {ID: "b", Text: "It involves foreign rather than domestic actors"},
                                                {ID: "c", Text: "It represents isolated incidents rather than systematic patterns"},
                                                {ID: "d", Text: "It transforms institutions themselves into extraction mechanisms by design"},
                                        },
                                        CorrectAnswer: "d",
                                        Explanation:   "Parasitic governance involves the systematic redirection of state institutions to serve extraction rather than public service - a structural rather than incidental problem.",
                                        Difficulty:    "hard",
                                },
                        },
                        Randomize: true,
                        PassScore: 75,
                }
                
                quizJSON, err := json.Marshal(quizContent)
                if err != nil {
                        return err
                }
                
                elements = append(elements, models.InteractiveElement{
                        SectionID:      sectionID,
                        Position:       99, // At the end of the section
                        Type:           models.QuizType,
                        Title:          "Check Your Understanding",
                        Description:    "Test your comprehension of key concepts from this section.",
                        Content:        string(quizJSON),
                        CompletionType: "graded",
                        PointsValue:    20,
                        RequiredStatus: false,
                })
        }
        
        // Save all interactive elements
        for _, element := range elements {
                err := s.bookRepo.CreateInteractiveElement(&element)
                if err != nil {
                        log.Printf("Error creating interactive element: %v", err)
                        return err
                }
        }
        
        return nil
}

// ImportBook1 imports content for Book 1 based on the TOC
func (s *BookImportServiceImpl) ImportBook1() error {
        // Create a new book
        book := &models.Book{
                Title:       "Great Nigeria – Awakening the Giant: A Call to Urgent United Citizen Action",
                Author:      "Great Nigeria Network",
                Description: "The first book in the Great Nigeria series, serving as a manifesto that focuses on diagnosing the 'Why' – the root causes and systemic nature of Nigeria's crisis.",
                Published:   true,
                CoverImage:  "/static/img/nigeria-landscape.svg",
                CreatedAt:   time.Now(),
                UpdatedAt:   time.Now(),
        }

        // Save the book
        err := s.bookRepo.CreateBook(book)
        if err != nil {
                log.Printf("Error creating book: %v", err)
                return err
        }

        // Create front matter
        frontMatter := &models.BookFrontMatter{
                BookID:      book.ID,
                Introduction: `
# Introduction

## The Nigerian Imperative

I've lived through Nigeria's peaks and valleys, watching with both hope and heartbreak as our nation has struggled to fulfill its destiny. This isn't just another academic analysis or distant critique – it's a deeply personal journey born from lived experience in our beautiful, complex country.

You and I have felt it – that gnawing frustration when electricity fails again during a critical moment. The silent rage watching political figures acquire inexplicable wealth while hospitals lack basic supplies. The weariness of navigating potholed roads that should have been repaired decades ago. This book validates these emotions as legitimate responses to systemic failure. Your anger isn't misplaced – it's the necessary fuel for driving meaningful change.

Nigeria's challenges aren't abstract problems. They manifest as the brilliant student who can't find employment, the entrepreneur whose business collapses under multiple taxation and infrastructure failures, the patient who dies because basic medical supplies weren't available. These daily tragedies reflect a profound systemic crisis demanding urgent, united response.

Why another book on Nigeria's problems? Because this isn't just analysis – it's a call to coordinated action. Nigeria doesn't lack brilliant ideas or capable people. What we've lacked is a unifying framework that channels our collective energy toward strategic transformation.

## Exposing the Crisis

In the chapters ahead, we'll embark on a diagnostic journey that may be uncomfortable but is absolutely necessary. We'll examine how historical legacies – from colonial structures to military rule – continue shaping our present. We'll dissect governance failures that have allowed state institutions to become extraction mechanisms rather than public service providers. We'll confront economic stagnation that has squandered our demographic dividend. And perhaps most painfully, we'll analyze the social fractures that have prevented us from forming the unified front needed for national transformation.

This diagnosis isn't academic – it's practical. Understanding the depth and interconnectedness of our challenges is the essential first step toward addressing them effectively. Without accurate diagnosis, our solutions will remain superficial and ineffective.

## Navigating This Manifesto

Let me be clear about this book's scope. Great Nigeria – Awakening the Giant serves as a Manifesto that focuses specifically on diagnosing the 'Why' – the root causes and systemic nature of our crisis. It aims to build consensus around our shared challenges and create a foundation for collective action.

This volume intentionally focuses on diagnosis rather than detailed prescriptions. The comprehensive 'How' – the detailed operational Masterplan – is reserved for Book 2 (Great Nigeria – The Masterplan). By separating diagnosis from prescription, we ensure a thorough understanding of our challenges before jumping to solutions.

As you read, you'll find not just analysis but invitations to engage. Each chapter includes forum topics for discussion and actionable steps you can take immediately to begin the transformation process. This isn't a book to merely read – it's a manifesto to act upon.

The crisis facing Nigeria demands nothing less than our full, united commitment. Let's begin this journey of awakening the giant together.
`,
                Preface: `
# Foreword: A Nation at the Crossroads – Feel the Raw Urgency

I write these words with an overwhelming sense of urgency. Nigeria stands at its most critical crossroads since independence – a moment demanding honest confrontation with our reality and collective action to alter our trajectory.

The paradox of Nigeria has never been more stark: extraordinary potential alongside devastating underperformance. We are Africa's largest economy yet host to the world's second-largest population of people living in extreme poverty. We have produced world-class minds in every field, yet watch our educational system crumble. We possess abundant natural resources while suffering chronic energy shortages. We demonstrate remarkable individual ingenuity while our systems and institutions fail systematically.

This is not merely unfortunate – it is tragic. It represents the squandering of the potential of 200+ million people and the betrayal of generations yet unborn.

I've traveled extensively throughout our nation, from Lagos to Maiduguri, Port Harcourt to Sokoto. Everywhere, I encounter the same reality: ordinary Nigerians demonstrating extraordinary resilience while yearning for the country they deserve. I've witnessed the market woman in Kano who operates a thriving business despite having no electricity, the teacher in Calabar who purchases supplies from her own meager salary, the young tech entrepreneur in Enugu creating solutions with minimal infrastructure.

These Nigerians aren't asking for handouts – they're asking for systems that work. They're asking for the opportunity to apply their talents in an environment that enables rather than impedes progress. They're asking for leadership that serves rather than extracts.

The purpose of this book is not merely to describe our challenges but to catalyze coordinated action. Nigeria's problems are well-documented, but solutions have remained elusive partly because we've approached them in fragmented, uncoordinated ways. The time has come for a united approach that recognizes the interconnectedness of our challenges and mobilizes Nigerians across all divides.

This is not an academic exercise or a theoretical proposal. It is a call to urgent, practical action. The conditions that have held Nigeria back are not acts of God or immutable fate – they are human-made and can be human-solved.

The question facing each of us is simple: Will we continue accepting the unacceptable, or will we stand up as a united citizenry to demand and create the Nigeria we deserve? History will judge our generation by how we answer this question.

The raw urgency of this moment cannot be overstated. Nigeria's problems will not wait while we debate perfect solutions. Every day of inaction costs lives, squanders potential, and deepens our challenges.

I urge you to read this book not as a passive observer but as a participant in Nigeria's transformation. The analysis and proposals that follow are not the final word but the beginning of a national conversation that must translate quickly into coordinated action.

Nigeria's greatest resource has always been its people. The time has come to unleash our collective power through united, strategic citizen engagement. Nothing less will suffice at this critical moment in our history.

*Dr. Obiageli Ezekwesili*  
Former Minister of Education  
Former Vice President, World Bank Africa Division
`,
                Acknowledgements: `
# Dedication

This book is dedicated to three groups:

First, to the countless Nigerians who struggle daily against systemic obstacles yet continue to demonstrate remarkable resilience and ingenuity. Your persistence in the face of challenges is the foundation upon which a greater Nigeria will be built.

Second, to the emerging generation of young Nigerians who refuse to accept the status quo and are bringing fresh energy and ideas to our national transformation. Your impatience with dysfunction and demand for a better future drive the urgency of this work.

Third, to those who have sacrificed – sometimes with their lives – while standing against corruption, defending democracy, and advocating for a more just Nigeria. Your courage illuminates the path forward.

The Nigeria you deserve is possible, and this book represents one step toward making it reality.
`,
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
        }

        // Save the front matter
        err = s.bookRepo.CreateFrontMatter(frontMatter)
        if err != nil {
                log.Printf("Error creating front matter: %v", err)
                return err
        }

        // Create chapters
        chapters := []models.BookChapter{
                {
                        BookID:      book.ID,
                        Title:       "The Bleeding Giant & Ghosts of the Past",
                        Number:      1,
                        Description: "This chapter vividly illustrates Nigeria's vast, squandered potential and begins unearthing the deep historical roots that continue to bind the nation.",
                        Content:     "This chapter serves as the powerful opening salvo of the diagnosis. It vividly illustrates the 'Bleeding Giant' metaphor by showcasing Nigeria's vast, squandered potential against its current reality, and begins unearthing the deep historical roots ('Ghosts of the Past') – colonialism, military rule, civil conflict, resource curse – that continue to bind the nation.",
                        CreatedAt:   time.Now(),
                        UpdatedAt:   time.Now(),
                },
                {
                        BookID:      book.ID,
                        Title:       "The Rot Within: Diagnosing Our System's Sickness",
                        Number:      2,
                        Description: "This chapter delves into the specific mechanics of Nigeria's contemporary dysfunction.",
                        Content:     "This chapter meticulously dissects the 'System's Sickness' by analyzing the decay in governance structures, the resulting economic hardship and infrastructure collapse, and the consequent erosion of democratic values and trust. The analysis is sharp, evidence-based, and reveals the interconnectedness of these failures.",
                        CreatedAt:   time.Now(),
                        UpdatedAt:   time.Now(),
                },
                {
                        BookID:      book.ID,
                        Title:       "Mirror on Ourselves: Confronting Citizen Complicity",
                        Number:      3,
                        Description: "This chapter holds up an uncomfortable but necessary 'Mirror', shifting the focus inward to examine how ordinary citizens contribute to sustaining the crisis.",
                        Content:     "This chapter shifts the focus inward to examine how ordinary citizens, through action or inaction, contribute to sustaining the crisis. It fosters a sense of shared responsibility without absolving leadership failure.",
                        CreatedAt:   time.Now(),
                        UpdatedAt:   time.Now(),
                },
                {
                        BookID:      book.ID,
                        Title:       "The Roar of the Unheard: Celebrating Emerging Resistance",
                        Number:      4,
                        Description: "This chapter shifts tone to offer hope and inspiration by showcasing existing and emerging resistance efforts.",
                        Content:     "After the heavy diagnosis, this chapter intentionally shifts tone to offer hope and inspiration. It 'celebrates' existing and 'Emerging Resistance' efforts, showcasing the courage and impact of citizens already pushing back ('The Roar of the Unheard'). It highlights key lessons from these actions and emphasizes the transformative 'Power of Unity'.",
                        CreatedAt:   time.Now(),
                        UpdatedAt:   time.Now(),
                },
                {
                        BookID:      book.ID,
                        Title:       "Conclusion & Addendum: Your Manifesto for Unity",
                        Number:      5,
                        Description: "This concluding section serves as the capstone of Book 1's diagnostic manifesto.",
                        Content:     "This concluding section serves as the capstone of Book 1's diagnostic manifesto. It summarizes the crisis concisely, powerfully reiterates the call to awake the giant, invites direct engagement via GreatNigeria.net, and provides a teaser connecting this diagnostic volume explicitly to the solutions detailed in Book 2.",
                        CreatedAt:   time.Now(),
                        UpdatedAt:   time.Now(),
                },
        }

        // Save chapters
        for i, chapter := range chapters {
                err := s.bookRepo.CreateChapter(&chapter)
                if err != nil {
                        log.Printf("Error creating chapter %s: %v", chapter.Title, err)
                        return err
                }
                
                // Create different sections based on chapter number
                var sections []models.BookSection
                
                switch i {
                case 0: // Chapter 1: The Bleeding Giant & Ghosts of the Past
                        sections = []models.BookSection{
                                {
                                        BookID:    book.ID,
                                        ChapterID: chapter.ID,
                                        Title:     "The Bleeding Giant",
                                        Number:    1,
                                        Format:    "markdown",
                                        Content:   `
# 1.1 The Bleeding Giant

I still remember the first time I truly grasped Nigeria's paradox. It was 2016, and I was traveling from Lagos to Ibadan. Along the expressway stood a massive, half-completed factory – steel beams exposed to the elements, weeds growing through the foundation. It had been that way for nearly a decade, I learned from locals. Just a hundred meters away, small businesses operated from makeshift shops, using generators that filled the air with noise and fumes. The juxtaposition was jarring: enormous potential literally rusting away while human ingenuity struggled against impossible odds.

This sight is not exceptional in our Nigeria – it's emblematic. We are the bleeding giant.

Consider our confounding statistics: Africa's largest economy with a GDP of approximately $440 billion, yet more than 40% of our population lives in extreme poverty. We hold Africa's largest oil reserves, yet cannot reliably power homes or businesses. We produce some of the world's most brilliant minds – from Chimamanda Adichie in literature to Dr. Oluyinka Olutoye in medicine – yet watch nearly 20 million children go without education.

The bleeding is not figurative – it manifests in very real human costs. In my work across Nigeria, I've documented countless examples: the university graduate in Kano driving okada because no formal employment exists; the woman in Enugu who lost her pregnancy because the nearby hospital lacked basic supplies; the brilliant entrepreneur in Port Harcourt whose business collapsed under the weight of multiple taxation and infrastructure failure.

What makes this tragedy so acute is that it occurs alongside extraordinary potential. Nigeria possesses everything needed for prosperity:

* **Natural Resources**: Beyond oil and gas, we have vast agricultural land, solid minerals, and water resources. The diversity of our ecological zones provides opportunities for varied agricultural production and tourism development.

* **Geographic Positioning**: Located strategically in West Africa with extensive coastline, Nigeria is naturally positioned as a regional hub for trade, finance, and logistics.

* **Demographic Dividend**: Our population of over 200 million includes a youthful majority that could fuel decades of economic growth – the same demographic advantage that powered Asian economic miracles.

* **Cultural Power**: From Nollywood to Afrobeats, Nigerian cultural production has continental and global influence that could be leveraged economically.

* **Entrepreneurial Spirit**: The Nigerian ability to create opportunity amid chaos is legendary. Our entrepreneurial drive has created innovations that function despite, not because of, our systems.

Yet this potential hemorrhages daily through systems designed for extraction rather than development. The gulf between what is and what could be represents one of the greatest tragedies of modern development.

During a community meeting in Benue State, I listened as a farmer named Emmanuel described harvesting a bumper crop of cassava only to watch much of it rot because roads to market were impassable. "We have the richest soil," he told me, his voice rising with emotion, "but what good is growing food that cannot reach those who need it?" His question haunts me because it applies beyond agriculture – what good is Nigerian potential that cannot reach its expression?

The bleeding of our giant is not natural or inevitable – it is the result of specific historical wounds and ongoing systemic failures that we must diagnose accurately if we hope to stop the hemorrhaging.
`,
                                        Published: true,
                                        CreatedAt: time.Now(),
                                        UpdatedAt: time.Now(),
                                },
                                {
                                        BookID:    book.ID,
                                        ChapterID: chapter.ID,
                                        Title:     "Ghosts of the Past",
                                        Number:    2,
                                        Format:    "markdown",
                                        Content:   `
# 1.2 Ghosts of the Past

The Nigeria we experience today was not created in a vacuum. Our present is haunted by historical legacies that have shaped our institutions, power dynamics, and even our collective psychology. These "ghosts of the past" continue to influence our national development in ways both obvious and subtle.

## The Colonial Blueprint

The most fundamental ghost haunting Nigeria is our colonial foundation. Unlike nations that evolved organically over centuries, modern Nigeria was created through arbitrary colonial boundaries that forced diverse peoples into a single political entity without their consent. This artificial construction created several enduring challenges:

First, the colonial state was designed for extraction, not development. British colonial administration focused on extracting resources and maintaining order, not building effective public institutions or fostering national identity. This extractive blueprint was bequeathed to post-colonial Nigeria, establishing a pattern where state power serves extraction rather than public welfare.

I witnessed this legacy firsthand while researching government archives in Kaduna. Colonial-era documents revealed a consistent prioritization of resource extraction and administrative convenience over community welfare or long-term development – a pattern mirrored in many contemporary governance approaches.

Second, colonial rule deliberately emphasized and institutionalized ethnic differences as a governance strategy. The divide-and-rule approach prevented unified resistance by heightening ethnic distinctions and creating competitive rather than collaborative relationships between ethnic groups. This manipulation of identity planted seeds of distrust that continue bearing bitter fruit.

During my university days at Ahmadu Bello University in the early 1990s, I witnessed how quickly political disagreements transformed into ethnic tensions. What struck me was how these divisions served the interests of those in power – a dynamic straight from the colonial playbook.

## The Military Interregnum

Military rule represents another haunting legacy. Between 1966 and 1999, Nigeria spent nearly 30 years under military dictatorship, with profound consequences for our institutional development:

The militarization of politics centralized power in ways that undermined federalism and local governance. Command-and-control governance replaced deliberative democracy, establishing patterns of authoritarianism that continue influencing our political culture even under civilian rule.

In 2019, I interviewed several state legislators for a research project. Most expressed frustration that their governors operated as "military administrators in civilian clothing," bypassing legislative oversight and treating public funds as personal resources. The ghost of military governance continues walking our corridors of power.

Military rule also normalized corruption at unprecedented scales. The absence of accountability mechanisms allowed grand corruption to become systemic rather than exceptional, creating networks of patronage that persist in our contemporary politics.

During a community forum in Makurdi in 2018, an elderly participant made a poignant observation: "Under colonial rule, corruption was controlled to serve British interests. Under military rule, corruption was unleashed to serve military interests. When will it serve Nigerian interests?"

## The Civil War Trauma

The Nigerian Civil War (1967-1970) represents a national trauma whose effects continue reverberating through our society. Beyond the immediate human toll – over a million deaths – the war left lasting legacies:

It deepened ethnic suspicions and reinforced the perception that politics is a zero-sum game where one group's gain necessitates another's loss. This perception continues undermining efforts at building national consensus on critical issues.

Speaking with community leaders in the Southeast in 2017, I was struck by how fresh civil war wounds remained fifty years later. One elder told me, "The physical reconstruction happened quickly, but the psychological reconstruction never really began." His comment illuminates how unresolved historical traumas continue shaping contemporary political attitudes.

The war also reinforced the primacy of security concerns over development priorities, justifying centralized control of resources and limiting true federalism – patterns that persist in our current governance arrangements.

## The Resource Curse

Nigeria's discovery of oil brought mixed blessings that continue distorting our development. The transformation into a petro-state created several enduring challenges:

Oil wealth facilitated the creation of a rentier state where political power became the primary means of accessing national wealth, intensifying competition for control of the state and disincentivizing productive economic activity.

I remember interviewing a former civil servant who described witnessing the psychological transformation of government in the 1970s. "Before oil," he explained, "government departments focused on generating production. After oil, they focused on distributing consumption." This shift from production to distribution fundamentally altered the state's relationship with citizens.

Oil dependence also deepened regional inequalities and created new environmental challenges, particularly in the Niger Delta, where communities bear the environmental costs of extraction while receiving minimal benefits – a situation that continues generating conflict and instability.

During research visits to oil-producing communities in Bayelsa State, I encountered the stark reality of this ghost: villages surrounded by oil infrastructure yet lacking basic amenities, waters polluted beyond recovery, and communities torn between resistance and accommodation to the industry dominating their lands.

These historical legacies are not deterministic – they don't seal our fate. But they do shape the institutional landscape in which we operate and influence both elite and citizen behavior. Understanding these ghosts allows us to exorcise their negative influences through conscious reform efforts rather than unconsciously reproducing historical patterns.

As the Yoruba proverb wisely observes: "However far a stream flows, it never forgets its source." Nigeria's challenges cannot be understood or addressed without recognizing the historical sources from which they flow.
`,
                                        Published: true,
                                        CreatedAt: time.Now(),
                                        UpdatedAt: time.Now(),
                                },
                                {
                                        BookID:    book.ID,
                                        ChapterID: chapter.ID,
                                        Title:     "Forum Topics: Understanding Potential and History",
                                        Number:    3,
                                        Format:    "markdown",
                                        Content:   `
# 1.3 Forum Topics: Understanding Potential and History

## FT 1.3.1: Personal Resonance with the Paradox
*Which aspect of Nigeria's "bleeding giant" status resonates most with your personal experience? How did it make you feel?*

## FT 1.3.2: Lingering Historical Effects
*How do you see colonial legacies or military rule still affecting governance or social relations in your community today?*

## FT 1.3.3: Beyond the Obvious Historical Factors
*Beyond the "resource curse," what other historical factors do you believe are critical to understanding Nigeria's present challenges?*

## FT 1.3.4: History Towards Solutions
*How can understanding these historical roots move us beyond blame towards constructive solutions?*
`,
                                        Published: true,
                                        CreatedAt: time.Now(),
                                        UpdatedAt: time.Now(),
                                },
                                {
                                        BookID:    book.ID,
                                        ChapterID: chapter.ID,
                                        Title:     "Actionable Steps: Connecting Past to Present",
                                        Number:    4,
                                        Format:    "markdown",
                                        Content:   `
# 1.4 Actionable Steps: Connecting Past to Present

## AS 1.4.1: Reflect & Share Historical Impact
*Write down one personal or family story affected by Nigeria's historical legacies (e.g., impact of civil war, military rule). Share reflections (anonymously if preferred) on the GNN forum for this chapter.*

## AS 1.4.2: Research & Discuss Specific History
*Pick one historical issue (e.g., resource curse, a specific military regime) and find one article or resource (suggested in Appendix/Bibliography) to learn more. Discuss key insights with your Action Cell or online.*

## AS 1.4.3: Observe Local Potential/Waste
*Identify one example of underutilized potential (human or material resource) in your immediate community. Document it (photo/notes) and consider sharing on GNN as evidence of the "Bleeding Giant" locally.*

## AS 1.4.4: Commit to Historical Empowerment
*Verbally affirm to yourself or your group: "Understanding our history empowers us to change the future."*

{{interactive:1}}
`,
                                        Published: true,
                                        CreatedAt: time.Now(),
                                        UpdatedAt: time.Now(),
                                },
                        }
                case 1: // Chapter 2: The Rot Within: Diagnosing Our System's Sickness
                        sections = []models.BookSection{
                                {
                                        BookID:    book.ID,
                                        ChapterID: chapter.ID,
                                        Title:     "Parasitic Governance & State Capture",
                                        Number:    1,
                                        Format:    "markdown",
                                        Content:   `
# 2.1 Parasitic Governance & State Capture

When I served on a state education committee in 2018, I witnessed something disturbing. We had gathered to review a budget allocation of ₦2.5 billion for school infrastructure. The committee chair – a well-connected political appointee – casually explained the unwritten "formula": 40% would actually go to schools, while 60% would be distributed among various officials and contractors. What shocked me wasn't just the brazenness of the theft but how institutionalized it had become – a routine "business practice" rather than an aberration.

This experience illustrates a fundamental reality we must confront: Nigeria doesn't merely have corrupt individuals within its governance system; it has a system of governance designed for extraction rather than service. This distinction is crucial.

## Understanding Parasitic Governance

Parasitic governance occurs when state institutions function primarily to extract resources for the benefit of those controlling them rather than to provide services to citizens. Unlike a symbiotic relationship where government and citizens mutually benefit, parasitic governance drains societal resources while providing minimal returns.

The evidence surrounds us. According to the Nigerian Extractive Industries Transparency Initiative, Nigeria lost over $42 billion to oil theft and mismanagement between 2009 and 2018. The National Bureau of Statistics reports that only 30.1% of children in public primary schools have access to basic sanitation facilities. The stark contrast between resource extraction and service delivery reveals the parasitic nature of our governance.

I've observed this dynamic across sectors. In 2019, I interviewed healthcare workers at a primary healthcare center in Kaduna State. Despite budget allocations for essential medicines, the facility had received no supplies for six months. The funds had vanished into an opaque procurement system designed to enrich officials rather than supply medications.

## The Mechanisms of State Capture

"State capture" describes the process whereby public institutions are redirected to serve private interests rather than public needs. This goes beyond individual corruption to encompass systematic distortion of state functions.

In Nigeria, state capture operates through several mechanisms:

**Control of Key Appointments**: Positions throughout government are allocated based on political loyalty and extraction potential rather than competence. During research interviews with civil servants in 2020, I repeatedly heard variations of the same saying: "Important positions go to those who will 'cooperate' with the sharing formula."

**Budget Manipulation**: Resources are directed toward easily exploitable projects rather than public priorities. The persistent preference for new construction over maintenance exemplifies this tendency – new projects create fresh opportunities for inflated contracts and kickbacks.

**Regulatory Capture**: Regulatory agencies meant to provide oversight are co-opted by the interests they should regulate. In the petroleum sector, for instance, the revolving door between regulatory bodies and industry undermines effective oversight of corporate practices.

**Privatization of Public Power**: Public authority is treated as private property to be leveraged for personal gain. Police checkpoints transformed into toll gates epitomize this dynamic – public security power repurposed for private extraction.

**Informational Asymmetry**: Deliberate opacity in public processes prevents citizen oversight. The resistance to full implementation of the Freedom of Information Act reflects the value of information control to those benefiting from the current system.

## The Consequences of Parasitic Governance

The impacts extend far beyond financial losses:

**Erosion of State Capacity**: When institutions are designed for extraction rather than performance, they lose the ability to fulfill their intended functions. This explains why Nigeria's oil refineries operated at near-zero capacity for years despite billions allocated for maintenance – their true function was wealth extraction, not petroleum refining.

**Public Service Collapse**: Essential services deteriorate as resources are diverted. The state of our public schools, hospitals, and infrastructure isn't a resource problem but a governance problem – funds allocated for these services are systematically diverted.

**Normalization of Corruption**: Systemic corruption creates a moral environment where corrupt practices become expected and even justified as "the way things work." A civil servant I interviewed in Abuja explained, "If you don't participate in the sharing, you're viewed as foolish or naive, not principled."

**Trust Deficit**: Perhaps most devastatingly, parasitic governance destroys the essential trust between citizens and state. A 2019 Nigeria Social Cohesion Survey found that 70% of Nigerians have little or no trust in government institutions. Without this trust, even well-intentioned reforms struggle to gain traction.

## The Elite Consensus

Underlying parasitic governance is what political economists call an "elite consensus" – an informal understanding among powerful groups to maintain systems of extraction regardless of apparent political differences. This explains why similar patterns persist across different administrations and political parties.

The consensus crosses ethnic, religious, and partisan lines. As one retired permanent secretary told me, "When it comes to sharing the national cake, those who appear divided on television find remarkable unity behind closed doors."

This unspoken agreement prevents meaningful reform and explains why Nigeria has abundant anti-corruption agencies but limited anti-corruption progress. The institutions ostensibly fighting corruption operate within a system designed to facilitate it.

## Breaking the Parasitic Cycle

Understanding governance as parasitic rather than merely inefficient fundamentally changes how we approach reform. Instead of technical fixes to improve capacity, we need transformative changes that alter the underlying incentives and power relationships.

This requires disrupting the elite consensus through sustained citizen pressure, creating transparency mechanisms resistant to capture, implementing consistent accountability measures, and cultivating a new leadership ethos focused on service rather than extraction.

Most importantly, it requires that we, as citizens, stop viewing this system as natural or inevitable. The parasitic nature of governance in Nigeria is a human creation that can be uncreated through collective action. As the Niger Delta proverb reminds us, "A disease that has been diagnosed is half cured." By accurately diagnosing our parasitic governance, we take the first crucial step toward transformation.
`,
                                        Published: true,
                                        CreatedAt: time.Now(),
                                        UpdatedAt: time.Now(),
                                },
                                {
                                        BookID:    book.ID,
                                        ChapterID: chapter.ID,
                                        Title:     "Forum Topics: Experiencing Systemic Rot",
                                        Number:    2,
                                        Format:    "markdown",
                                        Content:   `
# 2.4 Forum Topics: Experiencing Systemic Rot

## FT 2.4.1: Governance Impact on Services
*Describe a specific instance where "parasitic governance" or "state capture" directly impacted access to a public service (e.g., healthcare, education, justice) for you or someone you know.*

## FT 2.4.2: Most Impactful Infrastructure Failure
*Which aspect of "economic despair" or "infrastructure failure" (e.g., power, roads, healthcare) impacts your daily life the most severely? How could fixing it unlock potential?*

## FT 2.4.3: Observing Democratic Decay Locally
*What does "democratic integrity" mean to you, and where do you see its most significant erosion in Nigeria today at the grassroots level?*

## FT 2.4.4: Linking Systemic Rot to History
*How does the systemic rot discussed here connect to the historical issues raised in Chapter 1?*
`,
                                        Published: true,
                                        CreatedAt: time.Now(),
                                        UpdatedAt: time.Now(),
                                },
                        }
                default: // Other chapters get a generic section for now
                        sections = []models.BookSection{
                                {
                                        BookID:    book.ID,
                                        ChapterID: chapter.ID,
                                        Title:     "Chapter Overview",
                                        Number:    1,
                                        Format:    "markdown",
                                        Content:   chapter.Content,
                                        Published: true,
                                        CreatedAt: time.Now(),
                                        UpdatedAt: time.Now(),
                                },
                        }
                }

                // Save sections
                for _, section := range sections {
                        err := s.bookRepo.CreateSection(&section)
                        if err != nil {
                                log.Printf("Error creating section %s: %v", section.Title, err)
                                return err
                        }
                        
                        // Create interactive elements for the section
                        err = s.createInteractiveElements(book.ID, section.ID, section.Title)
                        if err != nil {
                                log.Printf("Error creating interactive elements for section %s: %v", section.Title, err)
                                return err
                        }
                }
        }

        log.Printf("Successfully imported Book 1")
        return nil
}