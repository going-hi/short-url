package link

import (
	"net/http"
	"short-url/pkg/middleware"
	"short-url/pkg/utils"
	"strconv"
)

type LinkController struct {
	Repository *LinkRepository
}

type LinkControllerParams struct {
	Repository *LinkRepository
}

func NewLinkController(params *LinkControllerParams) *LinkController {
	return &LinkController{
		Repository: params.Repository,
	}
}

func (c *LinkController) create(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.GetBody[CreateLinkRequest](r.Body)

	if err != nil {
		utils.SendJson(w, http.StatusBadRequest, err.Error())
		return
	}

	userId := r.Context().Value(middleware.ContextIdKey).(int)
	code := utils.GenerateCode()
	linkData, err := c.Repository.Create(payload.Url, code, userId)

	if err != nil {
		utils.SendJson(w, http.StatusBadRequest, err.Error())
		return
	}

	response := &CreateLinkResponse{
		Url:  linkData.Url,
		Code: linkData.Code,
		Id:   linkData.Id,
	}

	utils.SendJson(w, http.StatusCreated, response)
}

func (c *LinkController) findById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")

		id, err := strconv.ParseInt(idString, 10, 0)

		if err != nil {
			utils.SendJson(w, http.StatusBadRequest, err.Error())
			return
		}

		linkData, err := c.Repository.FindById(int(id))

		if err != nil {
			utils.SendJson(w, http.StatusNotFound, err.Error())
			return
		}

		response := &FindLinkResponse{
			Id:     linkData.Id,
			Url:    linkData.Url,
			Code:   linkData.Code,
			Clicks: linkData.Clicks,
		}

		utils.SendJson(w, http.StatusOK, response)
	}
}

func (c *LinkController) findByCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		if code == "" {
			utils.SendJson(w, http.StatusBadGateway, "Валидация")
			return
		}

		linkData, err := c.Repository.FindByCode(code)

		if err != nil {
			utils.SendJson(w, http.StatusNotFound, err.Error())
			return
		}

		response := &FindLinkResponse{
			Id:     linkData.Id,
			Url:    linkData.Url,
			Code:   linkData.Code,
			Clicks: linkData.Clicks,
		}

		utils.SendJson(w, http.StatusOK, response)
	}
}

func (c *LinkController) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")

		if code == "" {
			utils.SendJson(w, http.StatusBadRequest, "Валидация")
			return
		}

		linkData, err := c.Repository.FindByCode(code)

		if err != nil {
			utils.SendJson(w, http.StatusNotFound, err.Error())
			return
		}

		err = c.Repository.UpdateClick(linkData.Id)

		if err != nil {
			utils.SendJson(w, http.StatusBadRequest, err.Error())
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
			utils.SendJson(w, http.StatusBadRequest, err.Error())
			return
		}

		userId := r.Context().Value(middleware.ContextIdKey).(int)

		linkData, err := c.Repository.FindById(int(id))

		if err != nil {
			utils.SendJson(w, http.StatusNotFound, err.Error())
			return
		}

		if userId != linkData.UserId {
			utils.SendJson(w, http.StatusForbidden, "Нет прав")
			return
		}

		err = c.Repository.Delete(linkData.Id)

		if err != nil {
			utils.SendJson(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.SendJson(w, http.StatusNoContent, "Success")
	}
}

func (c *LinkController) getList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(middleware.ContextIdKey).(int)

		links, err := c.Repository.FindAllByUserId(userId)

		if err != nil {
			utils.SendJson(w, http.StatusBadRequest, links)
			return
		}

		utils.SendJson(w, http.StatusOK, links)
	}
}
