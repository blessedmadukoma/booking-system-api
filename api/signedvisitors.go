package api

import (
	db "booking-api/db/sqlc"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateSignedVisitor documents a signed visitor
func (srv *Server) CreateSignedVisitor(data createSignedVisitorData, ctx *gin.Context) (any, error) {

	arg := db.CreateSignedVisitorParams{
		VisitorID: data.VisitorID,
		VisitID:   sql.NullInt64{Int64: data.VisitID, Valid: true},
		SignedOut: data.SignedOut,
	}

	visitor, err := srv.store.CreateSignedVisitor(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("unable to create and document signed visitor: %w", err)
	}

	return visitor, nil
}

// ListSignedVisitors lists all signed in and out visitors
func (srv *Server) ListSignedVisitors(ctx *gin.Context) {
	visitors, err := srv.store.ListSignedVisitors(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("unable to retrieve all signed visitors from DB", err))
		return
	}

	if len(visitors) == 0 {
		// ctx.JSON(http.StatusNotFound, errorResponse("error retrieving documented visitors", fmt.Errorf("No signed visitors found")))
		// ctx.JSON(http.StatusNotFound, visitors)
		ctx.JSON(http.StatusOK, visitors)
		return
	}

	ctx.JSON(http.StatusOK, visitors)
}

// getLatestSignedVisitor gets the latest signed in and out visitor record
// func (srv *Server) getLatestSignedVisitor(ctx *gin.Context, visit_id, visitor_id int64) (int64, error) {
// 	arg := db.ListLatestSignedVisitorParams{
// 		VisitorID: visitor_id,
// 		VisitID:   sql.NullInt64{Int64: visit_id, Valid: true},
// 	}

// 	signedVisit, err := srv.store.ListLatestSignedVisitor(ctx, arg)
// 	if err != nil {
// 		err = fmt.Errorf("unable to retrieve latest signed visitor from DB: %w", err)
// 		ctx.JSON(http.StatusBadRequest, errorResponse("unable to retrieve latest signed visitor from DB", err))
// 		return 0, err
// 	}

// 	return signedVisit.ID, nil

// }

// UpdateSignedVisitors	updates a signed visitor
// func (srv *Server) UpdateSignedVisitors(id int64, visit_id sql.NullInt64, signedOut sql.NullTime) {
// 	var ctx *gin.Context

// 	arg := db.UpdateSignedVisitorsParams{
// 		ID:        id,
// 		VisitID:   visit_id,
// 		SignedOut: signedOut,
// 	}

// 	visitor, err := srv.store.UpdateSignedVisitors(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse("Unable to update and document signed visitor", err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, visitor)
// 	return
// }
