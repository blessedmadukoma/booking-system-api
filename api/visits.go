package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	db "booking-api/db/sqlc"

	"github.com/gin-gonic/gin"
)

// UpdateVisit updates a visit
func (srv *Server) UpdateVisit(ctx *gin.Context) {
	type updateVisitIDParams struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	var reqID updateVisitIDParams
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("visit ID not valid", err))
		return
	}

	type updateVisitParams struct {
		Status     string `json:"status" binding:"required"`
		Reason     string `json:"reason" binding:"required"`
		EmployeeID int64  `json:"employee_id"`
		VisitorID  int64  `json:"visitor_id" binding:"required"`
	}

	var req updateVisitParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("User input not valid", err))
		return
	}

	existingVisit, err := srv.store.GetVisitsByID(ctx, reqID.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse("unable to get visit ID", err))
		return
	}

	if existingVisit.VisitorID != req.VisitorID && existingVisit.EmployeeID.Int64 != req.EmployeeID {
		ctx.JSON(http.StatusForbidden, errorResponse("visitor IDs and Employee IDs do not match", err))
		return
	}

	// get the visit first using the visit ID, visitor ID and employee ID, then update the record and return the record.

	if req.Status != "approved" && req.Status != "pending" && req.Status != "denied" {
		ctx.JSON(http.StatusBadRequest, errorResponse("input not valid", fmt.Errorf("status must be approved, pending or denied")))
		return
	}

	arg := db.UpdateVisitParams{
		ID:         reqID.ID,
		Status:     strings.Trim(req.Status, " "),
		Reason:     strings.Trim(req.Reason, " "),
		VisitorID:  req.VisitorID,
		EmployeeID: sql.NullInt64{Int64: req.EmployeeID, Valid: true},
	}

	visit, err := srv.store.UpdateVisit(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Unable to update visit", err))
		return
	}

	ctx.JSON(http.StatusOK, visit)
}

// CreateVisit creates a new visit
// func (srv *Server) CreateVisit(ctx *gin.Context) {
// 	type createVisitParams struct {
// 		Status     string `json:"status" binding:"required"`
// 		Reason     string `json:"reason" binding:"required"`
// 		EmployeeID int64  `json:"employee_id"`
// 		VisitorID  int64  `json:"visitor_id" binding:"required"`
// 	}

// 	var req createVisitParams
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse("User input not valid", err))
// 		return
// 	}
// 	var arg db.CreateVisitParams

// 	if req.EmployeeID == 0 {
// 		arg = db.CreateVisitParams{
// 			Status:    req.Status,
// 			Reason:    req.Reason,
// 			VisitorID: req.VisitorID,
// 		}
// 	} else {
// 		arg = db.CreateVisitParams{
// 			Status:     req.Status,
// 			Reason:     req.Reason,
// 			EmployeeID: sql.NullInt64{Int64: req.EmployeeID, Valid: true},
// 			VisitorID:  req.VisitorID,
// 		}
// 	}

// 	visit, err := srv.store.CreateVisit(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse("Unable to create account", err))
// 		return
// 	}

// 	srv.CreateSignedVisitor(arg.VisitorID, sql.NullInt64{Int64: visit.ID, Valid: true}, ctx)

// 	ctx.JSON(http.StatusCreated, visit)
// 	return
// }

// GetVisitByID gets a single visit by ID
func (srv *Server) GetVisitByID(ctx *gin.Context) {
	type getVisitByIDParams struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}

	var req getVisitByIDParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("visit ID not valid", err))
		return
	}

	visit, err := srv.store.GetVisitsByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Sprintf("error retrieving single visit %d", req.ID), err))
		return
	}

	ctx.JSON(http.StatusOK, visit)
}

// ListVisits lists all visits
func (srv *Server) ListVisits(ctx *gin.Context) {
	visits, err := srv.store.ListVisits(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("unable to retrieve visits from DB", err))
		return
	}

	if len(visits) == 0 {
		// ctx.JSON(http.StatusNotFound, errorResponse("Unable to get all visits", fmt.Errorf("No visits found")))
		// ctx.JSON(http.StatusNotFound, visits)
		ctx.JSON(http.StatusOK, visits)
		return
	}

	ctx.JSON(http.StatusOK, visits)
}

// ListEmployeeVisits lists all visits for a given employee
func (srv *Server) ListEmployeeVisits(ctx *gin.Context) {
	type listEmployeeVisitsParams struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}

	var req listEmployeeVisitsParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("employee ID for visit not valid", err))
		return
	}

	visits, err := srv.store.ListEmployeeVisits(ctx, sql.NullInt64{Int64: req.ID, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("unable to retrieve employee visits from DB", err))
		return
	}

	if len(visits) == 0 {
		// ctx.JSON(http.StatusNotFound, visits)
		ctx.JSON(http.StatusOK, visits)
		return
	}

	ctx.JSON(http.StatusOK, visits)
}

// ListVisitsByStatus	lists all visits for a given status
func (srv *Server) ListVisitsByStatus(ctx *gin.Context) {
	type listListVisitsByStatusParams struct {
		Status string `json:"status" binding:"required"`
	}

	var req listListVisitsByStatusParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("input not valid", err))
		return
	}

	if req.Status != "approved" && req.Status != "pending" && req.Status != "denied" {
		ctx.JSON(http.StatusBadRequest, errorResponse("input not valid", fmt.Errorf("status must be approved, pending or denied")))
		return
	}

	visits, err := srv.store.ListVisitsByStatus(ctx, req.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("unable to retrieve visits from DB", err))
		return
	}

	if len(visits) == 0 {
		// ctx.JSON(http.StatusNotFound, errorResponse("Unable to get visits", fmt.Errorf("No visits for `%s` found", req.Status)))
		// ctx.JSON(http.StatusNotFound, visits)
		ctx.JSON(http.StatusOK, visits)
		return
	}

	ctx.JSON(http.StatusOK, visits)
}

func (srv *Server) DeleteVisit(ctx *gin.Context) {
	type deleteVisitParams struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}

	var req deleteVisitParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("visit ID not valid", err))
		return
	}

	err := srv.store.DeleteVisit(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Sprintf("error deleting visit %d", req.ID), err))
		return
	}

	ctx.JSON(http.StatusOK, fmt.Sprintf("Visit %d deleted", req.ID))
}
