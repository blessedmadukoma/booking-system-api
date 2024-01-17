package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (srv *Server) ListAllAdminLogs(ctx *gin.Context) {
	logs, err := srv.store.ListAllAdminLogs(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Error with database", err))
		return
	}

	ctx.JSON(http.StatusOK, logs)
}

func (srv *Server) ListUserAdminLog(ctx *gin.Context) {
	type listUserAdminLogParams struct {
		AdminID int64 `uri:"id" binding:"required"`
	}

	var params listUserAdminLogParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("Error with params", err))
		return
	}

	logs, err := srv.store.GetAdminLogByAdminID(ctx, params.AdminID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Error with database", err))
		return
	}

	ctx.JSON(http.StatusOK, logs)
}
