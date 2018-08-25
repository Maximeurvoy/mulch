package main

// AddRoutes defines all API routes for the application
func (app *App) AddRoutes() {
	AddRoute(&Route{
		Methods: []string{"POST"},
		Path:    "/phone",
		Type:    RouteTypeCustom,
		Handler: PhoneController,
	}, app)

	AddRoute(&Route{
		Methods:      []string{"GET"},
		Path:         "/log",
		Type:         RouteTypeStream,
		IsRestricted: true,
		Handler:      LogController,
	}, app)

	AddRoute(&Route{
		Methods:      []string{"PUT"},
		Path:         "/vm",
		Type:         RouteTypeStream,
		IsRestricted: true,
		Handler:      VMController,
	}, app)

	AddRoute(&Route{
		Methods:      []string{"GET"},
		Path:         "/version",
		Type:         RouteTypeCustom,
		IsRestricted: true,
		Handler:      VersionController,
	}, app)

	AddRoute(&Route{
		Methods: []string{"GET"},
		Path:    "/test",
		Type:    RouteTypeStream,
		Handler: TestController,
	}, app)

	AddRoute(&Route{
		Methods: []string{"POST"},
		Path:    "/test2",
		Type:    RouteTypeStream,
		Handler: Test2Controller,
	}, app)
}
