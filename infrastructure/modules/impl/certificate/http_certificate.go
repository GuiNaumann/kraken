package ertificate

import (
	"database/sql"
	"encoding/json"
	"kraken/domain/entities"
	"kraken/domain/usecases"
	//setup "kraken/infrastructure"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"io"
	au "kraken/infrastructure/modules/impl/auth"
	"kraken/infrastructure/modules/impl/http_error"
	"log"
	"net/http"
	"strconv"
)

type CertificateModule struct {
	Db                 *sql.DB
	Cookie             *securecookie.SecureCookie
	CertificateUseCase usecases.CertificateUseCase
}

func (c *CertificateModule) Path() string {
	return "/certificate"
}

func (c *CertificateModule) Setup(router *mux.Router) {
	//privateRoutes := router.PathPrefix("/private").Subrouter()
	//privateRoutes.Use(setup.AuthorizationMiddleware)

	privateRoutes := router.PathPrefix(c.Path()).Subrouter()

	privateRoutes.HandleFunc("/create", c.createCertificate).Methods(http.MethodPost)
	privateRoutes.HandleFunc("/list", c.listCertificate).Methods(http.MethodGet)
	privateRoutes.HandleFunc("/get/{certificateID}", c.getCertificateById).Methods(http.MethodGet)
	privateRoutes.HandleFunc("/update/{certificateID}", c.updateCertificate).Methods(http.MethodPost)
	privateRoutes.HandleFunc("/delete/{certificateID}", c.deleteCertificate).Methods(http.MethodDelete)
}

func (c *CertificateModule) createCertificate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[createCertificate] Error ReadAll", err)
		http_error.HandleError(w, http_error.NewBadRequestError(http_error.InvalidRequestBody))
		return
	}

	var certificate entities.Certificate
	err = json.Unmarshal(body, &certificate)
	if err != nil {
		log.Println("[createCertificate] Error Unmarshal", err)
		http_error.HandleError(w, http_error.NewBadRequestError(http_error.InvalidRequestBody))
		return
	}

	ctx := r.Context()
	user := ctx.Value(au.CtxUserKey).(*entities.User)
	CertificateID, err := c.CertificateUseCase.CreateCertificateUseCase(ctx, *user, certificate)
	if err != nil {
		log.Println("[createCertificate] Error CreateCertificateUseCase", err)
		http_error.HandleError(w, err)
		return
	}

	b, err := json.Marshal(CertificateID)
	if err != nil {
		log.Println("[createCertificate] Error Marshal", err)
		http_error.HandleError(w, http_error.NewUnexpectedError(http_error.Unexpected))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println("[createCertificate] Error Write", err)
		http_error.HandleError(w, http_error.NewUnexpectedError(http_error.Unexpected))
		return
	}
}

func (c *CertificateModule) listCertificate(w http.ResponseWriter, r *http.Request) {
	var filter entities.GeneralFilter
	var err error

	page := r.URL.Query().Get("page")
	if page != "" {
		filter.Page, err = strconv.ParseInt(page, 10, 64)
		if err != nil {
			log.Println("[listCertificate] Error ParseInt page", err)
			http_error.HandleError(w, http_error.NewBadRequestError(http_error.InvalidParameter))
			return
		}
	}

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		filter.Limit, err = strconv.ParseInt(limit, 10, 64)
		if err != nil {
			http_error.HandleError(w, http_error.NewBadRequestError(http_error.InvalidParameter))
			log.Println("[listCertificate] Error ParseInt limit", err)
			return
		}
	}

	filter.Column = r.URL.Query().Get("orderBy")

	ordinationAsc := r.URL.Query().Get("ordinationAsc")
	if ordinationAsc == "true" {
		filter.OrdinationAsc = true
	}

	filter.Search = r.URL.Query().Get("search")

	if filter.Limit == 0 && filter.Page != 0 {
		log.Println("[listCertificate] Error invalidParameter", err)
		http_error.HandleError(w, http_error.NewBadRequestError(http_error.InvalidParameter))
		return
	}

	ctx := r.Context()
	user := ctx.Value(au.CtxUserKey).(*entities.User)
	response, err := c.CertificateUseCase.ListCertificateUseCase(ctx, *user, filter)
	if err != nil {
		log.Println("[listCertificate] Error ListCertificateUseCase", err)
		http_error.HandleError(w, err)
		return
	}

	b, err := json.Marshal(response)
	if err != nil {
		log.Println("[listCertificate] Error Marshal", err)
		http_error.HandleError(w, http_error.NewUnexpectedError(http_error.Unexpected))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println("[listCertificate] Error Write", err)
		http_error.HandleError(w, http_error.NewUnexpectedError(http_error.Unexpected))
		return
	}
}

func (c *CertificateModule) getCertificateById(w http.ResponseWriter, r *http.Request) {
	certificateID, err := strconv.Atoi(mux.Vars(r)["certificateID"])
	if err != nil {
		log.Println("[getCertificateById] Error Atoi CertificateID", err)
		http_error.HandleError(w, http_error.NewUnexpectedError(http_error.Unexpected))
		return
	}

	ctx := r.Context()
	user := ctx.Value(au.CtxUserKey).(*entities.User)
	certificate, err := c.CertificateUseCase.GetCertificateByIdUseCase(ctx, *user, int64(certificateID))
	if err != nil {
		log.Println("[getCertificateById] Error GetCertificateByIdUseCase", err)
		http_error.HandleError(w, err)
		return
	}

	b, err := json.Marshal(certificate)
	if err != nil {
		log.Println("[getCertificateById] Error Marshal", err)
		http_error.HandleError(w, http_error.NewUnexpectedError(http_error.Unexpected))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println("[getCertificateById] Error Write", err)
		http_error.HandleError(w, http_error.NewUnexpectedError(http_error.Unexpected))
		return
	}
}

func (c *CertificateModule) updateCertificate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[updateCertificate] Error ReadAll", err)
		http_error.HandleError(w, http_error.NewBadRequestError(http_error.InvalidCertificate))
		return
	}

	var certificate entities.Certificate
	err = json.Unmarshal(body, &certificate)
	if err != nil {
		log.Println("[updateCertificate] Error Unmarshal", err)
		http_error.HandleError(w, http_error.NewBadRequestError(http_error.InvalidCertificate))
		return
	}

	certificateID, err := strconv.Atoi(mux.Vars(r)["certificateID"])
	if err != nil {
		log.Println("[updateCertificate] Error Atoi certificateID", err)
		http_error.HandleError(w, http_error.NewBadRequestError(http_error.InvalidParameter))
		return
	}

	certificate.Id = int64(certificateID)
	ctx := r.Context()
	user := ctx.Value(au.CtxUserKey).(*entities.User)

	err = c.CertificateUseCase.EditCertificateUseCase(ctx, *user, certificate)
	if err != nil {
		log.Println("[updateCertificate] Error EditCertificateUseCase", err)
		http_error.HandleError(w, err)
		return
	}

	b, err := json.Marshal(entities.NewSuccessfulRequest())
	if err != nil {
		log.Println("[updateCertificate] Error Marshal", err)
		http_error.HandleError(w, http_error.NewUnexpectedError(http_error.Unexpected))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println("[updateCertificate] Error Write", err)
		http_error.HandleError(w, http_error.NewUnexpectedError(http_error.Unexpected))
		return
	}
}

func (c *CertificateModule) deleteCertificate(w http.ResponseWriter, r *http.Request) {
	certificateID, err := strconv.ParseInt(mux.Vars(r)["certificateID"], 10, 64)
	if err != nil {
		log.Println("[deleteCertificate] Error Atoi certificateID", err)
		http_error.HandleError(w, http_error.NewUnexpectedError(http_error.Unexpected))
		return
	}

	ctx := r.Context()
	user := ctx.Value(au.CtxUserKey).(*entities.User)
	err = c.CertificateUseCase.DeleteCertificateUseCase(ctx, *user, certificateID)
	if err != nil {
		log.Println("[deleteCertificate] Error DeleteCertificateUseCase", err)
		http_error.HandleError(w, err)
		return
	}

	b, err := json.Marshal(entities.NewSuccessfulRequest())
	if err != nil {
		log.Println("[deleteCertificate] Error Marshal", err)
		http_error.HandleError(w, http_error.NewUnexpectedError(http_error.Unexpected))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println("[deleteCertificate] Error Write", err)
		http_error.HandleError(w, http_error.NewUnexpectedError(http_error.Unexpected))
		return
	}
}
