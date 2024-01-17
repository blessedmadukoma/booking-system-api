package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	db "booking-api/db/sqlc"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	// socketio "github.com/googollee/go-socket.io"
	// "github.com/gorilla/websocket"
)

// socket.IO to get created visitor data from postgres handler
// func (srv *Server) getVisitorsViaSocketIO(ctx *gin.Context) *socketio.Server {
// 	fmt.Println("socketIO")

// 	server := socketio.NewServer(nil)

// 	server.OnConnect("/", func(s socketio.Conn) error {
// 		s.SetContext("")
// 		log.Println("connected:", s.ID())
// 		return nil
// 	})

// 	fmt.Println("socketIO 2")

// 	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
// 		visitors, _ := srv.store.ListVisitors(ctx)
// 		log.Println("notice:", msg)
// 		s.Emit("reply", visitors)
// 	})

// 	// server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
// 	// 	s.SetContext(msg)
// 	// 	return "recv " + msg
// 	// })

// 	// server.OnEvent("/", "bye", func(s socketio.Conn) string {
// 	// 	last := s.Context().(string)
// 	// 	s.Emit("bye", last)
// 	// 	s.Close()
// 	// 	return last
// 	// })

// 	server.OnError("/", func(s socketio.Conn, e error) {
// 		log.Println("meet error:", e)
// 	})

// 	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
// 		log.Println("closed", reason)
// 	})

// 	return server
// }

// var wsupgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// // web socket to get created visitor data from postgres handler
// func (srv *Server) getEmployeeVisitorsViaWebSocket(ctx *gin.Context) {
// 	wsupgrader.CheckOrigin = func(r *http.Request) bool {
// 		return true
// 	}

// 	conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
// 	if err != nil {
// 		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
// 		return
// 	}

// 	defer conn.Close()

// 	var req getEmployeeByIDParams
// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse("employee ID not valid", err))
// 		return
// 	}

// 	for {
// 		visitors, err := srv.store.GetEmployeeVisitors(ctx, req.ID)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, errorResponse("unable to retrieve 24-hour span visitors from DB", err))
// 			return
// 		}

// 		conn.WriteJSON(visitors)
// 	}
// }

// // web socket to get created visitor data from postgres handler
// func (srv *Server) getVisitorsViaWebSocket(ctx *gin.Context) {
// 	wsupgrader.CheckOrigin = func(r *http.Request) bool {
// 		return true
// 	}

// 	conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
// 	if err != nil {
// 		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
// 		return
// 	}

// 	defer conn.Close()

// 	for {
// 		visitors, err := srv.store.ListVisitors(ctx)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, errorResponse("unable to retrieve 24-hour span visitors from DB", err))
// 			return
// 		}

// 		conn.WriteJSON(visitors)
// 	}
// }

type createSignedVisitorData struct {
	VisitID   int64
	SignedOut sql.NullTime
	VisitorID int64
}

// CreateVisitor creates a new visitor
func (srv *Server) CreateVisitor(ctx *gin.Context) {
	type createVisitorParams struct {
		Fullname    string `json:"fullname" binding:"required"`
		Email       string `json:"email" binding:"required"`
		Mobile      string `json:"mobile" binding:"required"`
		CompanyName string `json:"company_name" binding:"required"`
		Picture     string `json:"picture"`
		EmployeeID  int64  `json:"employee_id"`
	}

	var req createVisitorParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("User input not valid", err))
		return
	}

	var arg db.CreateVisitorParams

	picture := srv.uploadFile(req.Picture, fmt.Sprintf("visitors/%s", req.Fullname))

	if req.EmployeeID == 0 {
		req.EmployeeID = -1

		arg = db.CreateVisitorParams{
			Fullname:    strings.Trim(req.Fullname, " "),
			Email:       strings.Trim(req.Email, " "),
			Mobile:      strings.Trim(req.Mobile, " "),
			CompanyName: strings.Trim(req.CompanyName, " "),
			Picture:     picture,
		}
	} else {
		arg = db.CreateVisitorParams{
			Fullname:    strings.Trim(req.Fullname, " "),
			Email:       strings.Trim(req.Email, " "),
			Mobile:      strings.Trim(req.Mobile, " "),
			CompanyName: strings.Trim(req.CompanyName, " "),
			Picture:     picture,
			EmployeeID:  sql.NullInt64{Int64: req.EmployeeID, Valid: true},
		}
	}

	visitor, err := srv.store.CreateVisitor(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Unable to create visitor account", err))
		return
	}

	cvArgs := db.CreateVisitParams{
		VisitorID:  visitor.ID,
		Status:     "pending",
		Reason:     "yet to be approved",
		EmployeeID: arg.EmployeeID,
	}
	newVisit, err := srv.store.CreateVisit(ctx, cvArgs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Unable to create visit", err))
		return
	}

	employee, _ := srv.store.GetEmployeeByID(ctx, req.EmployeeID)

	// send mail to employee
	mailinfo := mailStructure{
		visitorName:    visitor.Fullname,
		employeeName:   employee.Fullname,
		employeeEmail:  employee.Email,
		visitorPicture: picture,
	}

	err = srv.sendMailEmployee(mailinfo)
	if err != nil {
		log.Println("Error sending mail to employee", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse("Unable to send mail to employee", err))
		return
	}

	visitor.Picture = picture

	response := struct {
		Visitor     db.Visitor `json:"visitor"`
		NewVisit    db.Visit   `json:"newVisit"`
		MailMessage string     `json:"mailMessage"`
	}{
		visitor,
		newVisit,
		"Mail sent successfully!",
	}

	ctx.JSON(http.StatusCreated, response)
}

// UpdateVisitor signs out a visitor
func (srv *Server) UpdateVisitor(ctx *gin.Context) {
	type updateVisitorParams struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}

	var req updateVisitorParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("visitor ID not valid", err))
		return
	}

	visitor, err := srv.store.UpdateVisitor(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Sprintf("error signing out visitor %d", req.ID), err))
		return
	}

	visit, err := srv.store.GetVisitByVisitorID(ctx, visitor.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("error getting the visit ID of this visitor", err))
		return
	}

	// log.Printf("Visit ID: %d => Visitor ID: %d\n", visit.ID, visitor.ID)

	// id, err := srv.getLatestSignedVisitor(ctx, visit.ID, visitor.ID)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, errorResponse("error getting latest signed visitor ID", err))
	// 	return
	// }

	// visit_id := sql.NullInt64{Int64: visit.ID, Valid: true}
	visitor_id := sql.NullInt64{Int64: visitor.ID, Valid: true}

	signedOut := sql.NullTime{Time: visitor.SignOut.Time, Valid: true}

	data := createSignedVisitorData{
		VisitID:   visit.ID,
		VisitorID: visitor_id.Int64,
		SignedOut: signedOut,
	}

	srv.CreateSignedVisitor(data, ctx)

	// send mail
	mailinfo := mailStructure{
		visitorName:  visitor.Fullname,
		visitorEmail: visitor.Email,
	}
	srv.sendMailVisitor(mailinfo)

	ctx.JSON(http.StatusOK, fmt.Sprintf("Visitor %d signed out by %v !", visitor.ID, visitor.SignOut.Time))
}

// GetvisitorByID retrieves a single visitor
func (srv *Server) GetvisitorByID(ctx *gin.Context) {
	type getVisitorByIDParams struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}

	var req getVisitorByIDParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("visitor ID not valid", err))
		return
	}

	visitor, err := srv.store.GetVisitorByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Sprintf("error retrieving single visitor %d", req.ID), err))
		return
	}

	ctx.JSON(http.StatusOK, visitor)
}

// ListVisitors retrieves visitors within a 24-hour span
func (srv *Server) ListVisitors(ctx *gin.Context) {
	visitors, err := srv.store.ListVisitors(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("unable to retrieve 24-hour span visitors from DB", err))
		return
	}

	ctx.JSON(http.StatusOK, visitors)
}

// ListAllVisitors retrieves all visitors
func (srv *Server) ListAllVisitors(ctx *gin.Context) {
	visitors, err := srv.store.ListAllVisitors(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("unable to retrieve all visitors from DB", err))
		return
	}

	ctx.JSON(http.StatusOK, visitors)
}

// // DeleteVisitor deletes a visitor
// func (srv *Server) DeleteVisitor(ctx *gin.Context) {
// 	type deleteVisitorParams struct {
// 		ID int64 `uri:"id" binding:"required,min=1"`
// 	}

// 	var req deleteVisitorParams
// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse("visitor ID not valid", err))
// 		return
// 	}

// 	err := srv.store.DeleteVisitor(ctx, req.ID)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Sprintf("error deleting visitor %d", req.ID), err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, fmt.Sprintf("Visitor %d deleted!", req.ID))
// 	return
// }

// uploadFile uploads the file to cloudinary and retrieves a compressed link
// func (srv *Server) uploadFile(file *multipart.FileHeader, filename string) string {
func (srv *Server) uploadFile(file string, filename string) string {
	CLOUDI_NAME := srv.config.Cloudinary.CloudiName
	CLOUDI_API_KEY := srv.config.Cloudinary.CloudiAPIKey
	CLOUDI_API_SECRET := srv.config.Cloudinary.CloudiAPISecret

	var cld, _ = cloudinary.NewFromParams(CLOUDI_NAME, CLOUDI_API_KEY, CLOUDI_API_SECRET)

	var ctx = context.Background()

	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{PublicID: "visitorsapi/Visitor_Images/" + filename})
	if err != nil {
		log.Fatalf("Failed to upload file  to cloudinary, %v\n", err)
	}

	// log.Println("Response url:", resp.SecureURL)
	return resp.SecureURL
}
