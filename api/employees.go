package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "booking-api/db/sqlc"
	"booking-api/util"

	"github.com/gin-gonic/gin"
)

type employeeResponse struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"fullname"`
	Email     string    `json:"email"`
	Mobile    string    `json:"mobile"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

// newEmployeeResponse creates a new employee response
func newEmployeeResponse(employee db.Employee) employeeResponse {
	return employeeResponse{
		ID:        employee.ID,
		FullName:  employee.Fullname,
		Email:     employee.Email,
		Mobile:    employee.Mobile,
		Token:     employee.Token,
		CreatedAt: employee.CreatedAt,
	}
}

// CreateEmployee creates a new employee
func (srv *Server) CreateEmployee(ctx *gin.Context) {
	type createEmployeeParams struct {
		Fullname string `json:"fullname" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Mobile   string `json:"mobile" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req createEmployeeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("User input not valid", err))
		return
	}

	password, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Error with password", err))
		return
	}

	arg := db.CreateEmployeeParams{
		Fullname: req.Fullname,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: password,
	}

	employee, err := srv.store.CreateEmployee(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("Unable to create employee", err))
		return
	}

	ctx.JSON(http.StatusCreated, employee)
}

// DeleteEmployee deletes an exisiting employee
func (srv *Server) DeleteEmployee(ctx *gin.Context) {
	type deleteEmployeeParams struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}

	var req deleteEmployeeParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("employee ID not valid", err))
		return
	}

	err := srv.store.DeleteEmployee(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Sprintf("error deleting employee %d", req.ID), err))
		return
	}

	ctx.JSON(http.StatusOK, fmt.Sprintf("Employee %d deleted!", req.ID))
}

type getEmployeeByIDParams struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// GetEmployeeByID gets a single employee by ID
func (srv *Server) GetEmployeeByID(ctx *gin.Context) {
	var req getEmployeeByIDParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("employee ID not valid", err))
		return
	}

	employee, err := srv.store.GetEmployeeByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Sprintf("error retrieving single employee %d", req.ID), err))
		return
	}

	ctx.JSON(http.StatusOK, employee)
}

// ListEmployees lists all employees
func (srv *Server) ListEmployees(ctx *gin.Context) {

	employees, err := srv.store.ListEmployees(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("unable to retrieve employees from DB", err))
		return
	}

	ctx.JSON(http.StatusOK, employees)
}

// LoginEmployee signs out a employee
func (srv *Server) LoginEmployee(ctx *gin.Context) {
	type loginEmployeeParams struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Token    string `json:"token"`
	}

	type loginEmployeeResponse struct {
		AccessToken string           `json:"access_token"`
		Employee    employeeResponse `json:"user"`
	}

	var req loginEmployeeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("employee input not valid", err))
		return
	}

	employee, err := srv.store.GetEmployeeByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse("employee not found", err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse("error communicating with db", err))
		return
	}

	err = util.CheckPassword(req.Password, employee.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse("password does not match", err))
		return
	}

	accessToken, err := srv.tokenMaker.CreateToken(employee.Email, srv.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("error creating access token", err))
		return
	}

	arg := db.LoginEmployeeParams{
		Email: req.Email,
		Token: req.Token,
	}

	_, err = srv.store.LoginEmployee(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("error storing mobile token", err))
		return
	}

	employee.Token = arg.Token

	response := loginEmployeeResponse{
		AccessToken: accessToken,
		Employee:    newEmployeeResponse(employee),
	}

	ctx.JSON(http.StatusOK, response)
}
