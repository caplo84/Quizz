package services

// Technology configuration centralized in a separate file for easier management

// TechConfigData holds all technology configurations
var TechConfigData = map[string]TechConfig{
	// Core Programming Languages
	"javascript": {
		DisplayName: "JavaScript",
		Aliases:     []string{"js"},
		Category:    "Programming Language",
	},
	"typescript": {
		DisplayName: "TypeScript",
		Aliases:     []string{"ts"},
		Category:    "Programming Language",
	},
	"python": {
		DisplayName: "Python",
		Aliases:     []string{},
		Category:    "Programming Language",
	},
	"java": {
		DisplayName: "Java",
		Aliases:     []string{},
		Category:    "Programming Language",
	},
	"go": {
		DisplayName: "Go Programming",
		Aliases:     []string{"golang"},
		Category:    "Programming Language",
	},
	"c++": {
		DisplayName: "C++",
		Aliases:     []string{"cpp"},
		Category:    "Programming Language",
	},
	"c-sharp": {
		DisplayName: "C#",
		Aliases:     []string{"csharp", "c#"},
		Category:    "Programming Language",
	},
	"c-(programming-language)": {
		DisplayName: "C Programming",
		Aliases:     []string{"c(programming-language)"},
		Category:    "Programming Language",
	},
	"php": {
		DisplayName: "PHP",
		Aliases:     []string{},
		Category:    "Programming Language",
	},
	"ruby": {
		DisplayName: "Ruby",
		Aliases:     []string{},
		Category:    "Programming Language",
	},
	"swift": {
		DisplayName: "Swift",
		Aliases:     []string{},
		Category:    "Programming Language",
	},
	"kotlin": {
		DisplayName: "Kotlin",
		Aliases:     []string{},
		Category:    "Programming Language",
	},
	"scala": {
		DisplayName: "Scala",
		Aliases:     []string{},
		Category:    "Programming Language",
	},
	"rust": {
		DisplayName: "Rust",
		Aliases:     []string{},
		Category:    "Programming Language",
	},
	"bash": {
		DisplayName: "Bash Scripting",
		Aliases:     []string{},
		Category:    "Scripting",
	},
	"objective-c": {
		DisplayName: "Objective-C",
		Aliases:     []string{},
		Category:    "Programming Language",
	},
	"r": {
		DisplayName: "R Programming",
		Aliases:     []string{},
		Category:    "Programming Language",
	},
	"matlab": {
		DisplayName: "MATLAB",
		Aliases:     []string{},
		Category:    "Programming Language",
	},

	// Web Technologies
	"html": {
		DisplayName: "HTML",
		Aliases:     []string{},
		Category:    "Web Technology",
	},
	"css": {
		DisplayName: "CSS",
		Aliases:     []string{},
		Category:    "Web Technology",
	},
	"reactjs": {
		DisplayName: "React",
		Aliases:     []string{"react"},
		Category:    "Web Framework",
	},
	"angular": {
		DisplayName: "Angular",
		Aliases:     []string{},
		Category:    "Web Framework",
	},
	"node.js": {
		DisplayName: "Node.js",
		Aliases:     []string{"nodejs"},
		Category:    "Runtime",
	},
	"jquery": {
		DisplayName: "jQuery",
		Aliases:     []string{},
		Category:    "Web Library",
	},
	"django": {
		DisplayName: "Django",
		Aliases:     []string{},
		Category:    "Web Framework",
	},
	"spring-framework": {
		DisplayName: "Spring Framework",
		Aliases:     []string{"spring"},
		Category:    "Web Framework",
	},
	"dotnet-framework": {
		DisplayName: ".NET Framework",
		Aliases:     []string{".net", "asp.net"},
		Category:    "Framework",
	},
	"ruby-on-rails": {
		DisplayName: "Ruby on Rails",
		Aliases:     []string{},
		Category:    "Web Framework",
	},
	"front-end-development": {
		DisplayName: "Frontend Development",
		Aliases:     []string{"frontend"},
		Category:    "Development Area",
	},
	"rest-api": {
		DisplayName: "REST API",
		Aliases:     []string{},
		Category:    "API Technology",
	},
	"json": {
		DisplayName: "JSON",
		Aliases:     []string{},
		Category:    "Data Format",
	},
	"xml": {
		DisplayName: "XML",
		Aliases:     []string{},
		Category:    "Data Format",
	},

	// Databases
	"sql": {
		DisplayName: "SQL",
		Aliases:     []string{},
		Category:    "Database",
	},
	"t-sql": {
		DisplayName: "T-SQL",
		Aliases:     []string{"tsql"},
		Category:    "Database",
	},
	"mysql": {
		DisplayName: "MySQL",
		Aliases:     []string{},
		Category:    "Database",
	},
	"mongodb": {
		DisplayName: "MongoDB",
		Aliases:     []string{},
		Category:    "Database",
	},
	"nosql": {
		DisplayName: "NoSQL",
		Aliases:     []string{},
		Category:    "Database",
	},
	"machine-learning": {
		DisplayName: "Machine Learning",
		Aliases:     []string{"ml"},
		Category:    "Technology Area",
	},

	// Mobile Development
	"android": {
		DisplayName: "Android Development",
		Aliases:     []string{},
		Category:    "Mobile Development",
	},
	"ios": {
		DisplayName: "iOS Development",
		Aliases:     []string{},
		Category:    "Mobile Development",
	},
	"unity": {
		DisplayName: "Unity",
		Aliases:     []string{},
		Category:    "Game Development",
	},
	"xamarin": {
		DisplayName: "Xamarin",
		Aliases:     []string{},
		Category:    "Mobile Development",
	},

	// Cloud & DevOps
	"git": {
		DisplayName: "Git Version Control",
		Aliases:     []string{},
		Category:    "Version Control",
	},
	"aws": {
		DisplayName: "Amazon Web Services",
		Aliases:     []string{},
		Category:    "Cloud Platform",
	},
	"aws-lambda": {
		DisplayName: "AWS Lambda",
		Aliases:     []string{},
		Category:    "Cloud Service",
	},
	"microsoft-azure": {
		DisplayName: "Microsoft Azure",
		Aliases:     []string{"azure"},
		Category:    "Cloud Platform",
	},
	"google-cloud-platform": {
		DisplayName: "Google Cloud Platform",
		Aliases:     []string{"gcp"},
		Category:    "Cloud Platform",
	},
	"docker": {
		DisplayName: "Docker",
		Aliases:     []string{},
		Category:    "Containerization",
	},
	"kubernetes": {
		DisplayName: "Kubernetes",
		Aliases:     []string{},
		Category:    "Container Orchestration",
	},
	"linux": {
		DisplayName: "Linux",
		Aliases:     []string{},
		Category:    "Operating System",
	},
	"windows-server": {
		DisplayName: "Windows Server",
		Aliases:     []string{},
		Category:    "Operating System",
	},
	"maven": {
		DisplayName: "Apache Maven",
		Aliases:     []string{},
		Category:    "Build Tool",
	},
}

// GetTechnologyCategories returns all available technology categories
func GetTechnologyCategories() []string {
	categories := make(map[string]bool)
	for _, config := range TechConfigData {
		categories[config.Category] = true
	}

	result := make([]string, 0, len(categories))
	for category := range categories {
		result = append(result, category)
	}

	return result
}

// GetTechnologiesByCategory returns technologies in a specific category
func GetTechnologiesByCategory(category string) map[string]TechConfig {
	filtered := make(map[string]TechConfig)

	for key, config := range TechConfigData {
		if config.Category == category {
			filtered[key] = config
		}
	}

	return filtered
}
