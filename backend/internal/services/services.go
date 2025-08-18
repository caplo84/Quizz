package services

// Services holds all service dependencies
type Services struct {
    TopicService   TopicService
    QuizService    QuizService
    AttemptService AttemptService
    AdminService   AdminService
}

// NewServices creates a new Services instance
func NewServices(
    topicService TopicService,
    quizService QuizService,
    attemptService AttemptService,
    adminService AdminService,
) *Services {
    return &Services{
        TopicService:   topicService,
        QuizService:    quizService,
        AttemptService: attemptService,
        AdminService:   adminService,
    }
}