package api

func (s *APIServer) registerRoutes() {
	api := s.echo.Group("/api")

	// Public routes
	api.GET("/health", s.handleHealth)
	api.GET("/status", s.handleStatus)

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/signup", s.handleSignup)
	auth.POST("/login", s.handleLogin)
	auth.POST("/logout", s.handleLogout)
	auth.GET("/me", s.handleMe, s.requireAuth)

	// Protected routes
	protected := api.Group("", s.requireAuth)
	_ = protected // placeholder until protected routes are added
}
