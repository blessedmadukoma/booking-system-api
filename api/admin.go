package api

import (
	db "booking-api/db/sqlc"
	"booking-api/util"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type adminResponse struct {
	ID        int64          `json:"id"`
	Fullname  string         `json:"fullname"`
	Email     string         `json:"email"`
	Role      sql.NullString `json:"role"`
	LoggedIn  sql.NullTime   `json:"logged_in"`
	LoggedOut sql.NullTime   `json:"logged_out"`
	// Password  string         `json:"password"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

// newAdminResponse creates a new admin response
func newAdminResponse(admin db.Admin) adminResponse {
	return adminResponse{
		ID:        admin.ID,
		Fullname:  admin.Fullname,
		Email:     admin.Email,
		Role:      admin.Role,
		LoggedIn:  admin.LoggedIn,
		LoggedOut: admin.LoggedOut,
		// Password:  admin.Password,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}
}

func (srv *Server) CreateAdmin(ctx *gin.Context) {
	type createAdminParams struct {
		Fullname string `json:"fullname" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Role     string `json:"role"`
		Password string `json:"password" binding:"required"`
	}

	var params createAdminParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("User input not valid", err))
		return
	}

	role := sql.NullString{
		String: params.Role,
		Valid:  true,
	}

	password, err := util.HashPassword(params.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Error with password", err))
		return
	}

	if role.String != "admin" && role.String != "superadmin" {
		role.String = "admin"
	}

	arg := db.CreateAdminParams{
		Fullname: params.Fullname,
		Email:    params.Email,
		Role:     role,
		Password: password,
	}

	admin, err := srv.store.CreateAdmin(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Error creating admin", err))
		return
	}
	ctx.JSON(http.StatusOK, admin)
}

func (srv *Server) LoginAdmin(ctx *gin.Context) {
	type loginAdminParams struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	type loginAdminResponse struct {
		AccessToken string        `json:"access_token"`
		Admin       adminResponse `json:"admin"`
	}

	var params loginAdminParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("Error with login params", err))
		return
	}

	userAdmin, err := srv.store.GetAdminByEmail(ctx, params.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse("admin not found", err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse("error communicating with db", err))
		return
	}

	err = util.CheckPassword(params.Password, userAdmin.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse("password does not match", err))
		return
	}

	admin, err := srv.store.LoginAdmin(ctx, params.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Error logging in admin", err))
		return
	}

	srv.store.CreateAdminLog(ctx, userAdmin.ID)

	token, err := srv.tokenMaker.CreateToken(userAdmin.Email, srv.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Error creating token", err))
		return
	}

	response := loginAdminResponse{
		AccessToken: token,
		Admin:       newAdminResponse(admin),
	}

	ctx.JSON(http.StatusOK, response)
}

func (srv *Server) LogoutAdmin(ctx *gin.Context) {
	type logoutAdminParams struct {
		Email string `json:"email" binding:"required"`
	}

	var params logoutAdminParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("Error with logout params", err))
		return
	}

	userAdmin, err := srv.store.GetAdminByEmail(ctx, params.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse("admin not found", err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse("error communicating with db", err))
		return
	}

	admin, err := srv.store.LogoutAdmin(ctx, userAdmin.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Error logging out admin", err))
		return
	}

	lastestAdminDetail, _ := srv.store.GetAdminLogsByAdminID(ctx, userAdmin.ID)

	srv.store.UpdateAdminLogs(ctx, lastestAdminDetail.AdminID)

	ctx.JSON(http.StatusOK, admin)
}

func (srv *Server) ListAdmins(ctx *gin.Context) {
	admins, err := srv.store.ListAdmins(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, admins)
}

func (srv *Server) DeleteAdmin(ctx *gin.Context) {
	type deleteAdminParams struct {
		ID int64 `json:"id" binding:"required"`
	}
	var params deleteAdminParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := srv.store.DeleteAdmin(ctx, params.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Admin deleted successfully"})
}
