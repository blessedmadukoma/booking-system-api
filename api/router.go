package api

import "github.com/gin-gonic/gin"

func Routes(router *gin.Engine, srv *Server) {
	// serverrr := srv.getVisitorsViaSocketIO(&gin.Context{})

	routes := router.Group("/api")
	{

		routes.GET("/health", healthy)

		// web socket
		// routes.GET("/socketio/*any", gin.WrapH(serverrr))
		// routes.POST("/socketio/*any", gin.WrapH(serverrr))
		// routes.GET("/ws", srv.getVisitorsViaWebSocket)
		// routes.GET("/ws/:id", srv.getEmployeeVisitorsViaWebSocket)

		// employee routes
		empRoutes := routes.Group("/employees")
		{
			empRoutes.GET("/", srv.ListEmployees)
			empRoutes.POST("/", srv.CreateEmployee)
			empRoutes.POST("/login", srv.LoginEmployee)
			empRoutes.GET("/:id", srv.GetEmployeeByID)
			empRoutes.DELETE("/:id", srv.DeleteEmployee)
		}

		// visitors routes
		visitorRoutes := routes.Group("/visitors")
		{
			visitorRoutes.POST("/", srv.CreateVisitor)
			visitorRoutes.GET("/:id", srv.GetvisitorByID)
			visitorRoutes.PUT("/:id", srv.UpdateVisitor)
			// visitorRoutes.DELETE("/:id", srv.DeleteVisitor)
			visitorRoutes.GET("/", srv.ListVisitors)
			visitorRoutes.GET("/all", srv.ListAllVisitors)
		}

		// signed visitors (sv) routes
		svRoutes := routes.Group("/signedvisitors")
		{
			svRoutes.GET("/", srv.ListSignedVisitors)
		}

		// visit routes
		visitRoutes := routes.Group("/visits")
		{
			// visitRoutes.POST("/", srv.CreateVisit)
			visitRoutes.GET("/", srv.ListVisits)
			visitRoutes.PUT("/:id", srv.UpdateVisit)
			visitRoutes.GET("/:id", srv.GetVisitByID)
			visitRoutes.DELETE("/:id", srv.DeleteVisit)
			visitRoutes.GET("/status", srv.ListVisitsByStatus)
			visitRoutes.GET("/employee/:id", srv.ListEmployeeVisits)
		}

		// admin routes
		adminRoutes := routes.Group("/admins")
		{
			adminRoutes.POST("/", srv.CreateAdmin)
			adminRoutes.GET("/", srv.ListAdmins)
			adminRoutes.POST("/login", srv.LoginAdmin)
			adminRoutes.POST("/logout", srv.LogoutAdmin)
			adminRoutes.DELETE("/:id", srv.DeleteAdmin)
		}

		// adminlogs routes
		admlogRoutes := routes.Group("/adminlogs")
		{
			admlogRoutes.GET("/", srv.ListAllAdminLogs)
			admlogRoutes.GET("/:id", srv.ListUserAdminLog)

		}
	}
}
