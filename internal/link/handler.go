package link

import (
	"net/http"
	"short-url/pkg/jwt"
	"short-url/pkg/middleware"
	"short-url/pkg/utils"
	"strconv"
)

type LinkController struct {
	Repository *LinkRepository
}

type LinkHandlerParams struct {
	Repository *LinkRepository
	JwtService *jwt.JwtService
}

func NewLinkHandler(router *http.ServeMux, params LinkHandlerParams) {
	linkController := &LinkController{
		Repository: params.Repository,
	}

	router.Handle("POST /link", middleware.IsAuthMiddleware(http.HandlerFunc(linkController.create), params.JwtService))

	router.Handle("GET /link/{id}", middleware.IsAuthMiddleware(linkController.findById(), params.JwtService))
	router.Handle("GET /link", middleware.IsAuthMiddleware(linkController.getList(), params.JwtService))
	router.Handle("DELETE /link/{id}", middleware.IsAuthMiddleware(linkController.delete(), params.JwtService))
	router.Handle("GET /{code}", middleware.IsAuthMiddleware(linkController.GoTo(), params.JwtService))
}

func (c *LinkController) create(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.GetBody[CreateLinkRequest](r.Body)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	userId := r.Context().Value(middleware.ContextIdKey).(int)
	code := utils.GenerateCode()
	linkData, err := c.Repository.Create(payload.Url, code, userId)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	response := &CreateLinkResponse{
		Url:  linkData.Url,
		Code: linkData.Code,
		Id:   linkData.Id,
	}

	utils.SendJson(w, 201, response)
}

func (c *LinkController) findById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")

		id, err := strconv.ParseInt(idString, 10, 0)

		if err != nil {
			utils.SendJson(w, 400, err.Error())
			return
		}

		linkData, err := c.Repository.FindById(int(id))

		if err != nil {
			utils.SendJson(w, 404, err.Error())
			return
		}

		response := &FindLinkResponse{
			Id:     linkData.Id,
			Url:    linkData.Url,
			Code:   linkData.Code,
			Clicks: linkData.Clicks,
		}

		utils.SendJson(w, 200, response)
	}
}

func (c *LinkController) findByCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		if code == "" {
			utils.SendJson(w, 400, "Валидация")
			return
		}

		linkData, err := c.Repository.FindByCode(code)

		if err != nil {
			utils.SendJson(w, 404, err.Error())
			return
		}

		response := &FindLinkResponse{
			Id:     linkData.Id,
			Url:    linkData.Url,
			Code:   linkData.Code,
			Clicks: linkData.Clicks,
		}

		utils.SendJson(w, 200, response)
	}
}

func (c *LinkController) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		if code == "" {
			utils.SendJson(w, 400, "Валидация")
			return
		}

		linkData, err := c.Repository.FindByCode(code)

		if err != nil {
			utils.SendJson(w, 404, err.Error())
			return
		}

		err = c.Repository.UpdateClick(linkData.Id)

		if err != nil {
			utils.SendJson(w, 400, err.Error())
			return
		}

		http.Redirect(w, r, linkData.Url, http.StatusFound)
	}
}

func (c *LinkController) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")

		id, err := strconv.ParseInt(idString, 10, 0)

		if err != nil {
			utils.SendJson(w, 400, err.Error())
			return
		}

		userId := r.Context().Value(middleware.ContextIdKey).(int)

		linkData, err := c.Repository.FindById(int(id))

		if err != nil {
			utils.SendJson(w, 404, err.Error())
			return
		}

		if userId != linkData.UserId {
			utils.SendJson(w, 403, "Нет прав")
			return
		}

		err = c.Repository.Delete(linkData.Id)

		if err != nil {
			utils.SendJson(w, 400, err.Error())
			return
		}

		utils.SendJson(w, 403, "Success")
	}
}

func (c *LinkController) getList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(middleware.ContextIdKey).(int)

		links, err := c.Repository.FindAllByUserId(userId)

		if err != nil {
			utils.SendJson(w, 400, links)
			return
		}

		utils.SendJson(w, 200, links)
	}
}
