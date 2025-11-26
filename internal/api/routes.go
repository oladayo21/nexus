package api

func (s *APIServer) registerRoutes() {
	api := s.echo.Group("/api")

	// Public routes
	api.GET("/health", s.handleHealth)

	// Protected routes
	protected := api.Group("", s.requireAuth)
	_ = protected // placeholder until protected routes are added
}
