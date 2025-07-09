package migration

import (
	"github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/content/models"
	"gorm.io/gorm"
)

// RunCitationMigrations sets up citation-related database tables
func RunCitationMigrations(db *gorm.DB) error {
	// Enable auto-migration for the citation models
	if err := db.AutoMigrate(
		&models.Citation{},
		&models.CitationUsage{},
		&models.BibliographyMetadata{},
		&models.CitationCategory{},
	); err != nil {
		return err
	}

	// Insert default citation categories if they don't exist
	categories := []models.CitationCategory{
		{Name: "book", Description: "Academic books and monographs", DisplayOrder: 1},
		{Name: "journal", Description: "Academic journal articles", DisplayOrder: 2},
		{Name: "report", Description: "Research reports and working papers", DisplayOrder: 3},
		{Name: "government", Description: "Government and institutional publications", DisplayOrder: 4},
		{Name: "interview", Description: "Field interviews and focus groups", DisplayOrder: 5},
		{Name: "survey", Description: "Surveys and statistical data", DisplayOrder: 6},
		{Name: "media", Description: "Media and online resources", DisplayOrder: 7},
	}

	for _, category := range categories {
		var count int64
		db.Model(&models.CitationCategory{}).Where("name = ?", category.Name).Count(&count)
		if count == 0 {
			db.Create(&category)
		}
	}

	return nil
}

// ImportBookCitations imports sample citation data for testing
func ImportBookCitations(db *gorm.DB, bookID uint, citations []models.Citation) error {
	// Clear existing citations for this book to avoid duplicates
	db.Where("book_id = ?", bookID).Delete(&models.Citation{})

	// Process each citation
	for i, c := range citations {
		c.BookID = bookID
		c.RefNumber = i + 1
		if err := db.Create(&c).Error; err != nil {
			return err
		}
	}

	return nil
}

// ImportBook3Citations imports sample citations for Book 3
func ImportBook3Citations(db *gorm.DB) error {
	book3Citations := []models.Citation{
		{CitationKey: "achebe1983", Author: "Achebe, Chinua", Year: "1983", Title: "The Trouble with Nigeria", Source: "Heinemann Educational Publishers", Type: "book", CitedCount: 4},
		{CitationKey: "collier2007", Author: "Collier, Paul", Year: "2007", Title: "The Bottom Billion: Why the Poorest Countries are Failing and What Can Be Done About It", Source: "Oxford University Press", Type: "book", CitedCount: 3},
		{CitationKey: "diamond2010", Author: "Diamond, Larry", Year: "2010", Title: "The Spirit of Democracy: The Struggle to Build Free Societies Throughout the World", Source: "St. Martin's Griffin", Type: "book", CitedCount: 2},
		{CitationKey: "easterly2006", Author: "Easterly, William", Year: "2006", Title: "The White Man's Burden: Why the West's Efforts to Aid the Rest Have Done So Much Ill and So Little Good", Source: "Penguin Books", Type: "book", CitedCount: 2},
		{CitationKey: "falola2008", Author: "Falola, Toyin & Heaton, Matthew M.", Year: "2008", Title: "A History of Nigeria", Source: "Cambridge University Press", Type: "book", CitedCount: 3},
		{CitationKey: "acemoglu2012", Author: "Acemoglu, Daron & Robinson, James A.", Year: "2012", Title: "Why Nations Fail: The Origins of Power, Prosperity, and Poverty", Source: "Crown Business", Type: "book", CitedCount: 2},
		{CitationKey: "elrufai2013", Author: "El-Rufai, Nasir A.", Year: "2013", Title: "The Accidental Public Servant", Source: "Safari Books", Type: "book", CitedCount: 2},
		{CitationKey: "okonjo2018", Author: "Okonjo-Iweala, Ngozi", Year: "2018", Title: "Fighting Corruption Is Dangerous: The Story Behind the Headlines", Source: "MIT Press", Type: "book", CitedCount: 2},
		{CitationKey: "adebanwi2016", Author: "Adebanwi, Wale", Year: "2016", Title: "Nation as Grand Narrative: The Nigerian Press and the Politics of Meaning", Source: "University of Rochester Press", Type: "book", CitedCount: 2},
		{CitationKey: "olukoshi2015", Author: "Olukoshi, Adebayo", Year: "2015", Title: "State, Conflict, and Democracy in Africa", Source: "Lynne Rienner Publishers", Type: "book", CitedCount: 2},
		{CitationKey: "lewis2018", Author: "Lewis, Peter", Year: "2018", Title: "Nigeria's Democracy and the Crisis of Political Instability: An Overview", Source: "Journal of Modern African Studies, 56(1), 141-166", Type: "journal", CitedCount: 2},
		{CitationKey: "onyeukwu2023", Author: "Onyeukwu, C.E. & Adegboye, F.B.", Year: "2023", Title: "Digital Transformation of Public Services in Nigeria: Challenges and Opportunities", Source: "African Journal of Science, Technology, Innovation and Development, 15(2), 209-221", Type: "journal", CitedCount: 2},
		{CitationKey: "suberu2021", Author: "Suberu, Rotimi", Year: "2021", Title: "Federalism, Power Sharing, and the COVID-19 Challenge in Nigeria", Source: "Commonwealth & Comparative Politics, 59(1), 101-125", Type: "journal", CitedCount: 2},
		{CitationKey: "cdd2023", Author: "Centre for Democracy and Development", Year: "2023", Title: "The State of Democracy in Nigeria", Source: "CDD West Africa", Type: "report", CitedCount: 2},
		{CitationKey: "icg2023", Author: "International Crisis Group", Year: "2023", Title: "Managing Vigilantism in Nigeria: A Near-term Necessity", Source: "Africa Report No. 308", Type: "report", CitedCount: 2},
		{CitationKey: "cppa2023", Author: "Centre for Public Policy Alternatives", Year: "2023", Title: "Nigeria Governance Reform Roadmap", Source: "CPPA Policy Papers", Type: "report", CitedCount: 2},
		{CitationKey: "mgi2022", Author: "McKinsey Global Institute", Year: "2022", Title: "Nigeria's Renewal: Delivering Inclusive Growth", Source: "MGI Reports", Type: "report", CitedCount: 2},
		{CitationKey: "worldbank2023", Author: "World Bank", Year: "2023", Title: "Nigeria Public Finance Review: Fiscal Management for Resilient and Inclusive Growth", Source: "World Bank Group", Type: "report", CitedCount: 2},
		{CitationKey: "cbn2022", Author: "Central Bank of Nigeria", Year: "2022", Title: "Annual Economic Report", Source: "CBN Publications", Type: "government", CitedCount: 2},
		{CitationKey: "nbs2023", Author: "National Bureau of Statistics", Year: "2023", Title: "Nigeria Poverty Assessment Report", Source: "NBS Publications", Type: "government", CitedCount: 2},
		{CitationKey: "fmf2023", Author: "Federal Ministry of Finance", Year: "2023", Title: "Medium Term Expenditure Framework and Fiscal Strategy Paper", Source: "Federal Government of Nigeria", Type: "government", CitedCount: 2},
		{CitationKey: "interviews2023", Author: "Series of interviews with civil servants in federal ministries, Abuja (January-March 2023)", Year: "2023", Title: "Insights on bureaucratic processes and reform challenges", Source: "Field Research (Names changed to protect privacy)", Type: "interview", CitedCount: 2},
		{CitationKey: "focusgroups2023", Author: "Focus group discussions with youth entrepreneurs across six geopolitical zones (April-June 2023)", Year: "2023", Title: "Entrepreneurial obstacles and opportunities", Source: "Field Research (Names changed to protect privacy)", Type: "interview", CitedCount: 2},
		{CitationKey: "tradrulers2023", Author: "Roundtable with traditional rulers on community governance systems (August 2023)", Year: "2023", Title: "Traditional governance mechanisms in modern Nigeria", Source: "Field Research (Names changed to protect privacy)", Type: "interview", CitedCount: 2},
		{CitationKey: "natsurvey2023", Author: "National survey on citizen perceptions of governance (n=3,600), conducted across 36 states (May-July 2023)", Year: "2023", Title: "Citizen Perceptions of Government Performance", Source: "Field Research Data", Type: "survey", CitedCount: 2},
		{CitationKey: "urbansurvey2023", Author: "Urban infrastructure needs assessment survey in Lagos, Kano, Port Harcourt, and Enugu (August-September 2023)", Year: "2023", Title: "Urban Infrastructure Priorities", Source: "Field Research Data", Type: "survey", CitedCount: 2},
	}

	return ImportBookCitations(db, 3, book3Citations)
}