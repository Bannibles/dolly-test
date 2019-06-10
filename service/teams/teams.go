package teams

import (
	"net/http"
	"strconv"

	"github.com/go-phorce/dolly-test/api/v1"
	"github.com/go-phorce/dolly-test/datahub"
	"github.com/go-phorce/dolly/rest"
	"github.com/go-phorce/dolly/xhttp/httperror"
	"github.com/go-phorce/dolly/xhttp/identity"
	"github.com/go-phorce/dolly/xhttp/marshal"
	"github.com/go-phorce/dolly/xlog"
)

// ServiceName provides the Service Name for this package
const ServiceName = "teams"

var logger = xlog.NewPackageLogger("github.com/go-phorce/dolly-test/service", "teams")

// Service defines the Data service
type Service struct {
	server rest.Server
	db     datahub.UsersManager
}

// Factory returns a factory of the service
func Factory(server rest.Server) interface{} {
	if server == nil {
		logger.Panic("teams.Factory: invalid parameter")
	}

	return func(db datahub.UsersManager) {
		svc := &Service{
			server: server,
			db:     db,
		}

		server.AddService(svc)
	}
}

// Name returns the service name
func (s *Service) Name() string {
	return ServiceName
}

// IsReady indicates that the service is ready to serve its end-points
func (s *Service) IsReady() bool {
	return true
}

// Close cleans up background processes of subservices
func (s *Service) Close() {
}

// Register adds the service status endpoints to the overall URL router
func (s *Service) Register(r rest.Router) {
	r.GET(v1.URIForTeams, listTeamsHandler(s))
	r.GET(v1.URIForUsers, listUsersHandler(s))
}

func listTeamsHandler(s *Service) rest.Handle {
	return func(w http.ResponseWriter, r *http.Request, p rest.Params) {
		ctx := identity.ForRequest(r)
		_ = ctx.Identity()

		res, err := s.db.ListTeams(r.Context())
		if err != nil {
			marshal.WriteJSON(w, r, httperror.WithUnexpected("failed to list team").WithCause(err))
			return
		}

		marshal.WritePlainJSON(w, http.StatusOK, res, marshal.PrettyPrint)
	}
}

func listUsersHandler(s *Service) rest.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ rest.Params) {
		ctx := identity.ForRequest(r)
		_ = ctx.Identity()

		params := r.URL.Query()
		minAge := params.Get("min_age")
		maxAge := params.Get("max_age")

		req := &v1.FindUserRequest{
			Name: params.Get("name"),
		}

		if minAge != "" {
			i, _ := strconv.Atoi(minAge)
			req.MinAge = i
		}
		if maxAge != "" {
			i, _ := strconv.Atoi(maxAge)
			req.MaxAge = i
		}

		res, err := s.db.FindUser(r.Context(), req)
		if err != nil {
			marshal.WriteJSON(w, r, httperror.WithUnexpected("failed to list team").WithCause(err))
			return
		}

		marshal.WritePlainJSON(w, http.StatusOK, res, marshal.PrettyPrint)
	}
}

func teamsMembershipHandler(s *Service) rest.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ rest.Params) {
		ctx := identity.ForRequest(r)
		idn := ctx.Identity()
		userID := idn.UserID()

		if userID == "" {
			marshal.WriteJSON(w, r, httperror.WithForbidden("invalid user ID"))
			return
		}

		_, ok := idn.UserInfo().(*v1.UserInfo)
		if !ok {
			marshal.WriteJSON(w, r, httperror.WithForbidden("failed to extract User Info from the token"))
			return
		}

		res := &v1.GetTeamMembershipsResponse{}

		// TODO: response

		marshal.WriteJSON(w, r, res)
	}
}
