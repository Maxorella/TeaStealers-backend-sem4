package delivery

import (
	"github.com/TeaStealers-backend-sem4/internal/auth"
	"github.com/TeaStealers-backend-sem4/internal/models"
	"github.com/TeaStealers-backend-sem4/pkg/jwt"
	"github.com/TeaStealers-backend-sem4/pkg/middleware"
	"github.com/TeaStealers-backend-sem4/pkg/utils"
	"github.com/satori/uuid"
	"net/http"
)

// AuthHandler handles HTTP requests for user authentication.
type AuthHandler struct {
	// uc represents the usecase interface for authentication.
	uc auth.AuthUsecase
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(uc auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// @Summary Register a new user
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.UserLoginData true "User data"
// @Success 201 {object} models.User
// @Failure 400 {string} string "Incorrect data format"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/signup [post]
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	data := models.UserSignUpData{}

	if err := utils.ReadRequestData(r, &data); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	newUser, token, exp, err := h.uc.SignUp(r.Context(), &data)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "data already is used")
		return
	}

	newUser.Token = token
	http.SetCookie(w, jwt.TokenCookie(middleware.CookieName, token, exp))

	if err = utils.WriteResponse(w, http.StatusCreated, newUser); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
	}
}

// @Summary User login
// @Description User login
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.UserLoginData true "User login data"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Incorrect password or login"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	data := models.UserLoginData{}
	if err := utils.ReadRequestData(r, &data); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, token, exp, err := h.uc.Login(r.Context(), &data)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "incorrect password or login")
		return
	}

	http.SetCookie(w, jwt.TokenCookie(middleware.CookieName, token, exp))
	user.Token = token
	if err := utils.WriteResponse(w, http.StatusOK, user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
	}
}

// @Summary User logout
// @Description User logout
// @Tags auth
// @Success 200 {string} string "Logged out"
// @Router /auth/logout [get]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  middleware.CookieName,
		Value: "",
		Path:  "/",
	})
	if err := utils.WriteResponse(w, http.StatusOK, "success logout"); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
	}
}

func (h *AuthHandler) MeHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(middleware.CookieName)
	if id == nil {
		utils.WriteError(w, http.StatusUnauthorized, "token cookie not found")
		return
	}
	uID := id.(uuid.UUID)

	user, err := h.uc.GetUserByID(r.Context(), uID)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "jwt token is invalid")
		return
	}
	user.ID = uuid.Nil
	if err := utils.WriteResponse(w, http.StatusOK, user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

/*func (h *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(middleware.CookieName)
	if id == nil {
		utils.WriteError(w, http.StatusUnauthorized, "token cookie not found")
		return
	}
	uId := id.(string)

	gotid, err := h.uc.CheckAuth(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = utils.WriteResponse(w, http.StatusOK, gotid); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
	}
}

*/

func (h *AuthHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(middleware.CookieName)
	if id == nil {
		utils.WriteError(w, http.StatusUnauthorized, "token cookie not found")
		return
	}
	uID := id.(uuid.UUID)

	dat := &models.UserUpdatePassword{
		ID: uID,
	}

	if err := utils.ReadRequestData(r, &dat); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "incorrect data format")
		return
	}

	token, exp, err := h.uc.UpdateUserPassword(dat)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	http.SetCookie(w, jwt.TokenCookie(middleware.CookieName, token, exp))

	if err := utils.WriteResponse(w, http.StatusOK, "success update password"); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "error write response")
	}
}
