package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MrHenri/go-api/internal/dto"
	"github.com/MrHenri/go-api/internal/entity"
	"github.com/MrHenri/go-api/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{UserDB: db}
}

// Create user godoc
// @Summary      Create user
// @Description  Create user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request     body      dto.CreateUserInput  true  "user request"
// @Success      201
// @Failure      500         {object}  Error
// @Failure      400         {object}  Error
// @Router       /user [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var dtoUser dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&dtoUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := entity.NewUser(dtoUser.Name, dtoUser.Email, dtoUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	err = h.UserDB.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Create user godoc
// @Summary      Generate AccessToken
// @Description  Generate AccessToken
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request     body      dto.GetJWTInput  true  "user request"
// @Success      200		 {object}  dto.GetJWTOutput
// @Failure      400         {object}  Error
// @Failure      401         {object}  Error
// @Failure      404         {object}  Error
// @Router       /user/generate_token [post]
func (h *UserHandler) GetJwt(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExperiesIn := r.Context().Value("jwtExperiesIn").(int)
	var dtoJwt dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&dtoJwt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	user, err := h.UserDB.FindByEmail(dtoJwt.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	if !user.ValidatePassword(dtoJwt.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExperiesIn)).Unix(),
	})

	accessToken := dto.GetJWTOutput{AccessToken: tokenString}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}
