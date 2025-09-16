package services

// Technology configuration centralized in a separate file for easier management

// TechConfigData holds all technology configurations
var TechConfigData = map[string]TechConfig{
	// Core Programming Languages
	"javascript": {
		DisplayName: "JavaScript",
		Aliases:     []string{"js"},
		Category:    "Programming Language",
		IconURL:     "/icon-js.svg",
	},
	"typescript": {
		DisplayName: "TypeScript",
		Aliases:     []string{"ts"},
		Category:    "Programming Language",
		IconURL:     "/icon-ts.svg",
	},
	"python": {
		DisplayName: "Python",
		Aliases:     []string{},
		Category:    "Programming Language",
		IconURL:     "/icon-python.svg",
	},
	"java": {
		DisplayName: "Java",
		Aliases:     []string{},
		Category:    "Programming Language",
		IconURL:     "/icon-java.svg",
	},
	"go": {
		DisplayName: "Go Programming",
		Aliases:     []string{"golang"},
		Category:    "Programming Language",
		IconURL:     "/icon-go.svg",
	},
	"c++": {
		DisplayName: "C++",
		Aliases:     []string{"cpp"},
		Category:    "Programming Language",
		IconURL:     "/icon-cpp.svg",
	},
	"c-sharp": {
		DisplayName: "C#",
		Aliases:     []string{"csharp", "c#"},
		Category:    "Programming Language",
		IconURL:     "/icon-csharp.svg",
	},
	"c-(programming-language)": {
		DisplayName: "C Programming",
		Aliases:     []string{"c(programming-language)"},
		Category:    "Programming Language",
		IconURL:     "/icon-c.svg",
	},
	"php": {
		DisplayName: "PHP",
		Aliases:     []string{},
		Category:    "Programming Language",
		IconURL:     "/icon-php.svg",
	},
	"ruby": {
		DisplayName: "Ruby",
		Aliases:     []string{},
		Category:    "Programming Language",
		IconURL:     "/icon-ruby.svg",
	},
	"swift": {
		DisplayName: "Swift",
		Aliases:     []string{},
		Category:    "Programming Language",
		IconURL:     "/icon-swift.svg",
	},
	"kotlin": {
		DisplayName: "Kotlin",
		Aliases:     []string{},
		Category:    "Programming Language",
		IconURL:     "/icon-kotlin.svg",
	},
	"scala": {
		DisplayName: "Scala",
		Aliases:     []string{},
		Category:    "Programming Language",
		IconURL:     "/icon-scala.svg",
	},
	"rust": {
		DisplayName: "Rust",
		Aliases:     []string{},
		Category:    "Programming Language",
		IconURL:     "/icon-rust.svg",
	},
	"bash": {
		DisplayName: "Bash Scripting",
		Aliases:     []string{},
		Category:    "Scripting",
		IconURL:     "/icon-bash.svg",
	},
	"objective-c": {
		DisplayName: "Objective-C",
		Aliases:     []string{},
		Category:    "Programming Language",
		IconURL:     "/icon-objc.svg",
	},
	"r": {
		DisplayName: "R Programming",
		Aliases:     []string{},
		Category:    "Programming Language",
		IconURL:     "/icon-r.svg",
	},
	"matlab": {
		DisplayName: "MATLAB",
		Aliases:     []string{},
		Category:    "Programming Language",
		IconURL:     "/icon-matlab.svg",
	},

	// Web Technologies
	"html": {
		DisplayName: "HTML",
		Aliases:     []string{},
		Category:    "Web Technology",
		IconURL:     "/icon-html.svg",
	},
	"css": {
		DisplayName: "CSS",
		Aliases:     []string{},
		Category:    "Web Technology",
		IconURL:     "/icon-css.svg",
	},
	"reactjs": {
		DisplayName: "React",
		Aliases:     []string{"react"},
		Category:    "Web Framework",
		IconURL:     "/icon-react.svg",
	},
	"angular": {
		DisplayName: "Angular",
		Aliases:     []string{},
		Category:    "Web Framework",
		IconURL:     "/icon-angular.svg",
	},
	"node.js": {
		DisplayName: "Node.js",
		Aliases:     []string{"nodejs"},
		Category:    "Runtime",
		IconURL:     "/icon-nodejs.svg",
	},
	"jquery": {
		DisplayName: "jQuery",
		Aliases:     []string{},
		Category:    "Web Library",
		IconURL:     "/icon-jquery.svg",
	},
	"django": {
		DisplayName: "Django",
		Aliases:     []string{},
		Category:    "Web Framework",
		IconURL:     "/icon-django.svg",
	},
	"spring-framework": {
		DisplayName: "Spring Framework",
		Aliases:     []string{"spring"},
		Category:    "Web Framework",
		IconURL:     "/icon-spring.svg",
	},
	"dotnet-framework": {
		DisplayName: ".NET Framework",
		Aliases:     []string{".net", "asp.net"},
		Category:    "Framework",
		IconURL:     "/icon-dotnet.svg",
	},
	"ruby-on-rails": {
		DisplayName: "Ruby on Rails",
		Aliases:     []string{},
		Category:    "Web Framework",
		IconURL:     "/icon-ruby.svg",
	},
	"front-end-development": {
		DisplayName: "Frontend Development",
		Aliases:     []string{"frontend"},
		Category:    "Development Area",
		IconURL:     "/icon-frontend.svg",
	},
	"rest-api": {
		DisplayName: "REST API",
		Aliases:     []string{},
		Category:    "API Technology",
		IconURL:     "/icon-api.svg",
	},
	"json": {
		DisplayName: "JSON",
		Aliases:     []string{},
		Category:    "Data Format",
		IconURL:     "/icon-json.svg",
	},
	"xml": {
		DisplayName: "XML",
		Aliases:     []string{},
		Category:    "Data Format",
		IconURL:     "/icon-xml.svg",
	},

	// Databases
	"sql": {
		DisplayName: "SQL",
		Aliases:     []string{},
		Category:    "Database",
		IconURL:     "/icon-sql.svg",
	},
	"t-sql": {
		DisplayName: "T-SQL",
		Aliases:     []string{"tsql"},
		Category:    "Database",
		IconURL:     "/icon-sql.svg",
	},
	"mysql": {
		DisplayName: "MySQL",
		Aliases:     []string{},
		Category:    "Database",
		IconURL:     "/icon-mysql.svg",
	},
	"mongodb": {
		DisplayName: "MongoDB",
		Aliases:     []string{},
		Category:    "Database",
		IconURL:     "/icon-mongodb.svg",
	},
	"nosql": {
		DisplayName: "NoSQL",
		Aliases:     []string{},
		Category:    "Database",
		IconURL:     "/icon-nosql.svg",
	},
	"machine-learning": {
		DisplayName: "Machine Learning",
		Aliases:     []string{"ml"},
		Category:    "Technology Area",
		IconURL:     "/icon-ai.svg",
	},

	// Mobile Development
	"android": {
		DisplayName: "Android Development",
		Aliases:     []string{},
		Category:    "Mobile Development",
		IconURL:     "/icon-android.svg",
	},
	"ios": {
		DisplayName: "iOS Development",
		Aliases:     []string{},
		Category:    "Mobile Development",
		IconURL:     "/icon-ios.svg",
	},
	"unity": {
		DisplayName: "Unity",
		Aliases:     []string{},
		Category:    "Game Development",
		IconURL:     "/icon-unity.svg",
	},
	"xamarin": {
		DisplayName: "Xamarin",
		Aliases:     []string{},
		Category:    "Mobile Development",
		IconURL:     "/icon-xamarin.svg",
	},

	// Cloud & DevOps
	"git": {
		DisplayName: "Git Version Control",
		Aliases:     []string{},
		Category:    "Version Control",
		IconURL:     "/icon-git.svg",
	},
	"aws": {
		DisplayName: "Amazon Web Services",
		Aliases:     []string{},
		Category:    "Cloud Platform",
		IconURL:     "/icon-aws.svg",
	},
	"aws-lambda": {
		DisplayName: "AWS Lambda",
		Aliases:     []string{},
		Category:    "Cloud Service",
		IconURL:     "/icon-aws.svg",
	},
	"microsoft-azure": {
		DisplayName: "Microsoft Azure",
		Aliases:     []string{"azure"},
		Category:    "Cloud Platform",
		IconURL:     "/icon-azure.svg",
	},
	"google-cloud-platform": {
		DisplayName: "Google Cloud Platform",
		Aliases:     []string{"gcp"},
		Category:    "Cloud Platform",
		IconURL:     "/icon-gcp.svg",
	},
	"docker": {
		DisplayName: "Docker",
		Aliases:     []string{},
		Category:    "Containerization",
		IconURL:     "/icon-docker.svg",
	},
	"kubernetes": {
		DisplayName: "Kubernetes",
		Aliases:     []string{},
		Category:    "Container Orchestration",
		IconURL:     "/icon-k8s.svg",
	},
	"linux": {
		DisplayName: "Linux",
		Aliases:     []string{},
		Category:    "Operating System",
		IconURL:     "/icon-linux.svg",
	},
	"windows-server": {
		DisplayName: "Windows Server",
		Aliases:     []string{},
		Category:    "Operating System",
		IconURL:     "/icon-windows.svg",
	},
	"maven": {
		DisplayName: "Apache Maven",
		Aliases:     []string{},
		Category:    "Build Tool",
		IconURL:     "/icon-maven.svg",
	},
	"accessibility": {
		DisplayName: "Accessibility",
		Aliases:     []string{},
		Category:    "Development Area",
		IconURL:     "/icon-accessibility.svg",
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

// GetIconForSlug returns the icon URL for a given slug
func GetIconForSlug(slug string) (string, bool) {
	// Check direct match
	if config, exists := TechConfigData[slug]; exists {
		return config.IconURL, true
	}

	// Check aliases
	for _, config := range TechConfigData {
		for _, alias := range config.Aliases {
			if alias == slug {
				return config.IconURL, true
			}
		}
	}

	return "", false
}
