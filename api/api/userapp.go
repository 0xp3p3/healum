package api

import (
	"context"
	"net/http"
	account_proto "server/account-srv/proto/account"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	content_proto "server/content-srv/proto/content"
	static_proto "server/static-srv/proto/static"
	track_proto "server/track-srv/proto/track"
	userapp_proto "server/user-app-srv/proto/userapp"
	user_proto "server/user-srv/proto/user"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

type UserAppService struct {
	UserAppClient userapp_proto.UserAppServiceClient
	Auth          Filters
	Audit         AuditFilter
	ServerMetrics metrics.Metrics
}

func (u UserAppService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/user/app")

	audit := &audit_proto.Audit{
		ActionService:  common.UserappSrv,
		ActionResource: common.BASE + common.USER_TYPE,
	}

	ws.Route(ws.GET("/{user_id}").To(u.ReadUser).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Read a user"))

	ws.Route(ws.POST("/bookmark").To(u.CreateBookmark).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Create bookmark"))

	ws.Route(ws.GET("/{user_id}/bookmarks/all").To(u.ReadBookmarkContents).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Create bookmark"))

	ws.Route(ws.GET("/{user_id}/bookmarks/categorys").To(u.ReadBookmarkContentCategorys).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Create bookmark"))

	ws.Route(ws.GET("/{user_id}/{category_id}/bookmarks").To(u.ReadBookmarkByCategory).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Create bookmark"))

	ws.Route(ws.DELETE("/bookmark/{bookmark_id}").To(u.DeleteBookmark).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Delete bookmark"))

	ws.Route(ws.POST("/bookmarks/search").To(u.SearchBookmarks).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Create bookmark"))

	ws.Route(ws.GET("/{user_id}/content/shared").To(u.GetSharedContent).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get shared content"))

	ws.Route(ws.GET("/{user_id}/plan/shared").To(u.GetSharedPlansForUser).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get shared plan"))

	ws.Route(ws.GET("/{user_id}/survey/shared").To(u.GetSharedSurveysForUser).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get shared survey"))

	ws.Route(ws.GET("/{user_id}/goal/shared").To(u.GetSharedGoalsForUser).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get shared goal"))

	ws.Route(ws.GET("/{user_id}/challenge/shared").To(u.GetSharedChallengesForUser).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Doc("Get shared challenge"))

	ws.Route(ws.GET("/{user_id}/habit/shared").To(u.GetSharedHabitsForUser).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get shared habits"))

	ws.Route(ws.POST("/goal/join").To(u.SignupToGoal).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Signup to a goal"))

	ws.Route(ws.GET("/goal/{goal_id}").To(u.GetGoalDetail).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get goal detail"))

	ws.Route(ws.GET("/goals/joined").To(u.GetAllJoinedGoals).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List all user's goal"))

	ws.Route(ws.GET("/goals/current").To(u.GetCurrentJoinedGoals).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List current goal"))

	ws.Route(ws.GET("/goal/current/progress").To(u.GetGoalProgress).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get goal progress"))

	ws.Route(ws.POST("/challenge/join").To(u.SignupToChallenge).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Signup to a challenge"))

	ws.Route(ws.GET("/challenge/{challenge_id}").To(u.GetChallengeDetail).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get challenge detail"))

	ws.Route(ws.GET("/challenges/joined").To(u.GetAllJoinedChallenges).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List all user's challenge"))

	ws.Route(ws.GET("/challenges/current").To(u.GetCurrentJoinedChallenges).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List current challenge"))

	ws.Route(ws.GET("/challenges/current/count").To(u.GetCurrentChallengesWithCount).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List current challenges with count"))

	ws.Route(ws.GET("/current/markers").To(u.ListMarkers).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List markers for the current active goals"))

	ws.Route(ws.GET("/pending").To(u.GetPendingSharedActions).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.Paginate).
		Filter(u.Auth.SortFilter).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get pending shared actions"))

	ws.Route(ws.GET("/marker/default/history").To(u.GetDefaultMarkerHistory).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.Paginate).
		Filter(u.Auth.SortFilter).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get default marker history"))

	ws.Route(ws.GET("/markers/{name_slug}/marker").To(u.MarkerByNameslug).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get marker by name_slug"))

	ws.Route(ws.POST("/habit/join").To(u.SignupToHabit).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Signup to a habit"))

	ws.Route(ws.GET("/habit/{habit_id}").To(u.GetHabitDetail).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get habit detail"))

	ws.Route(ws.GET("/habits/joined").To(u.GetAllJoinedHabits).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List all user's habits joined"))

	ws.Route(ws.GET("/habits/current").To(u.GetCurrentJoinedHabits).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List current habits"))

	ws.Route(ws.GET("/habits/current/count").To(u.GetCurrentHabitsWithCount).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List current habits"))

	ws.Route(ws.GET("/content/categorys/all").To(u.GetContentCategorys).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get content categories"))

	ws.Route(ws.GET("/content/{content_id}").To(u.GetContentDetail).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get content detail"))

	ws.Route(ws.GET("/content/category/{category_id}").To(u.GetContentByCategory).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get content from a category"))

	ws.Route(ws.GET("/content/category/{category_id}/filters").To(u.GetFiltersForCategory).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get filters for a category"))

	ws.Route(ws.POST("/content/category/filters/autocomplete").To(u.FiltersAutocomplete).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Filters autocomplete"))

	ws.Route(ws.POST("/content/category/{category_id}/filter").To(u.FilterContentInParticularCategory).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Filter content in a particular category"))

	ws.Route(ws.POST("/preferences").To(u.SaveUserPreference).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Save a users preferences"))

	ws.Route(ws.GET("/{user_id}/preferences").To(u.GetUserPreference).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get a users preferences"))

	ws.Route(ws.POST("/details").To(u.SaveUserDetails).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Save a users details"))

	ws.Route(ws.GET("/content/recommendations/all").To(u.GetContentRecommendationByUser).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get content recommendations for a particular user"))

	ws.Route(ws.GET("/content/recommendations/category/{category_id}").To(u.GetContentRecommendationByCategory).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get content recommendations by category"))

	ws.Route(ws.POST("/content/{content_id}/rating").To(u.SaveRateForContent).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Save user rating for a particular content object"))

	ws.Route(ws.POST("/content/{content_id}/dislike").To(u.DislikeForContent).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("User not interested in a particular content"))

	ws.Route(ws.POST("/content/{content_id}/dislike/similar").To(u.DislikeForSimilarContent).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("User not interested in similar content"))

	ws.Route(ws.POST("/feedback").To(u.SaveUserFeedback).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Save general user feedback"))

	ws.Route(ws.POST("/plan/join").To(u.JoinUserPlan).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Join/start a plan"))

	ws.Route(ws.POST("/plan/create").To(u.CreateUserPlan).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Create a user's plan"))

	ws.Route(ws.GET("/plan/{user_id}").To(u.GetUserPlan).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get a user's plan"))

	ws.Route(ws.POST("/plan/update").To(u.UpdateUserPlan).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Save general user feedback"))

	ws.Route(ws.GET("/plan/{plan_id}/summary/count").To(u.GetPlanItemsCountByCategory).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get plan items count by category"))

	ws.Route(ws.GET("/plan/{plan_id}/summary/{day_number}/count").To(u.GetPlanItemsCountByDay).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get plan items count by day"))

	ws.Route(ws.GET("/plan/{plan_id}/summary").To(u.GetPlanItemsCountByCategoryAndDay).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get plan items count by category and day"))

	ws.Route(ws.POST("/login").To(u.Login).
		Doc("User login"))

	ws.Route(ws.GET("/logout").To(u.Logout).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Logout the user"))

	ws.Route(ws.GET("/goals/all").To(u.AllGoals).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.Paginate).
		Filter(u.Auth.SortFilter).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List all goals"))

	ws.Route(ws.GET("/challenges/all").To(u.AllChallenges).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.Paginate).
		Filter(u.Auth.SortFilter).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List all challenges"))

	ws.Route(ws.GET("/habits/all").To(u.AllHabits).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.Paginate).
		Filter(u.Auth.SortFilter).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List all habits"))

	ws.Route(ws.POST("/contents/all").To(u.GetShareableContent).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.Paginate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List all habits"))

	ws.Route(ws.POST("/{user_id}/shared").To(u.RecievedItems).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Update shared items status to RECIEVED"))

	ws.Route(ws.POST("/content/category/{nameslug}/items/autocomplete").To(u.AutocompleteContentCategoryItem).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Autocomplete content category item by nameslug"))

	ws.Route(ws.GET("/content/category/{nameslug}/items/all").To(u.AllContentCategoryItemByNameslug).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("All content category item by nameslug"))

	restful.Add(ws)
}

/**
* @api {post} /server/user/app/bookmark?session={session_id} Create Bookmark
* @apiVersion 0.1.0
* @apiName CreateBookmark
* @apiGroup UserApp
*
* @apiDescription create bookmark for a user for a particular content item
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/bookmark?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "content_id": "contentid"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "bookmark_id": "95d31b08-47f6-11e8-b307-20c9d0453b15"
*   },
*   "code": 200,
*   "message": "Created bookmark correctly"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contents were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.CreateBookmark",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) CreateBookmark(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.CreateBookmark API request")
	req_userapp := new(userapp_proto.CreateBookmarkRequest)
	err := req.ReadEntity(req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.CreateBookmark", "BindError")
		return
	}
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.CreateBookmark(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.CreateBookmark", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created bookmark correctly"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/{user_id}/bookmarks/all?session={session_id} List all bookmarks for a user
* @apiVersion 0.1.0
* @apiName ReadBookmarkContents
* @apiGroup UserApp
*
* @apiDescription Get all bookmarks for a particular user across all different categories
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/{user_id}/bookmarks/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "content_id": "contentid"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "bookmarks": [
*       {
*         "id": "95d31b08-47f6-11e8-b307-20c9d0453b15",
*         "content_id": "ContentId",
*         "category_id": "Content.ContentCategory.Id",
*         "image": "Content.Name",
*         "title": "Content.ContentCategory.Title",
*         "category": "Content.ContentCateogry.Name"
*       },
*       ... ...
*     ],
*     "code": 200,
*     "message": "Read bookmarks successfully"
*   }
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The bookmarks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.CreateBookmark",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) ReadBookmarkContents(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.CreateBookmark API request")
	req_userapp := new(userapp_proto.ReadBookmarkContentRequest)
	req_userapp.UserId = req.PathParameter("user_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.ReadBookmarkContents(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.ReadBookmarkContents", "QueryError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read bookmarks successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/{user_id}/bookmarks/categorys?session={session_id} List all bookmark categories
* @apiVersion 0.1.0
* @apiName ReadBookmarkContentCategorys
* @apiGroup UserApp
*
* @apiDescription Get a list of categories for all the bookmarks that a user has created
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/{user_id}/bookmarks/categorys?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "content_id": "contentid"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "categorys": [
*       {
*         "category_id": "Content.ContentCategory.Id",
*         "category": "Content.ContentCateogry.Name"
*       },
*       ... ...
*     ],
*     "code": 200,
*     "message": "Read categories successfully"
*   }
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contents were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.ReadBookmarkContentCategorys",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) ReadBookmarkContentCategorys(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.ReadBookmarkContentCategorys API request")
	req_userapp := new(userapp_proto.ReadBookmarkContentCategorysRequest)
	req_userapp.UserId = req.PathParameter("user_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.ReadBookmarkContentCategorys(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.ReadBookmarkContentCategorys", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read categories successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/{user_id}/{category_id}/bookmarks?session={session_id} List all bookmarks in a category
* @apiVersion 0.1.0
* @apiName ReadBookmarkByCategory
* @apiGroup UserApp
*
* @apiDescription Get a list of all bookmarks that a user has created in a particular category
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/{user_id}/{category_id}/bookmarks?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "bookmarks": [
*       {
*         "id": "95d31b08-47f6-11e8-b307-20c9d0453b15",
*         "content_id": "ContentId",
*         "category_id": "Content.ContentCategory.Id",
*         "image": "Content.Name",
*         "title": "Content.ContentCategory.Title",
*         "category": "Content.ContentCateogry.Name"
*       },
*       ... ...
*     ],
*     "code": 200,
*     "message": "Read bookmarks successfully"
*   }
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The bookmarks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.ReadBookmarkByCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) ReadBookmarkByCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.ReadBookmarkByCategory API request")
	req_userapp := new(userapp_proto.ReadBookmarkByCategoryRequest)
	req_userapp.UserId = req.PathParameter("user_id")
	req_userapp.CategoryId = req.PathParameter("category_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.ReadBookmarkByCategory(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.ReadBookmarkByCategory", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read bookmarks successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/user/app/bookmark/{bookmark_id}?session={session_id} Delete bookmark
* @apiVersion 0.1.0
* @apiName DeleteBookmark
* @apiGroup UserApp
*
* @apiDescription Delete a particular bookmark
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/bookmark/bookmark_id?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "code": 200,
*     "message": "Deleted bookmark successfully"
*   }
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The bookmarks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.DeleteBookmark",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) DeleteBookmark(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.ReadBookmarkByCategory API request")
	req_userapp := new(userapp_proto.DeleteBookmarkRequest)
	req_userapp.BookmarkId = req.PathParameter("bookmark_id")
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.DeleteBookmark(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.ReadBookmarkByCategory", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted bookmark successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/user/app/bookmarks/search?session={session_id} Search bookmark
* @apiVersion 0.1.0
* @apiName SearchBookmarks
* @apiGroup UserApp
*
* @apiDescription Search bookmarks
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/bookmarks/search?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "title": "title",
*   "description": "descript",
*   "summary": "summary"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The bookmarks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.SearchBookmarks",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) SearchBookmarks(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.SearchBookmarks API request")
	req_userapp := new(userapp_proto.SearchBookmarkRequest)
	if err := utils.UnmarshalAny(req, rsp, req_userapp); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.SearchBookmarks", "BindError")
		return
	}
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.SearchBookmarks(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.SearchBookmarks", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted bookmark successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/{user_id}/content/shared?session={session_id} Get shared content
* @apiVersion 0.1.0
* @apiName GetSharedContent
* @apiGroup UserApp
*
* @apiDescription Get a list of content shared with the user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/437c32b2-9dd7-410a-8a97-162e580d8a90/content/shared?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "shared_contents":[
*       {
*         "cotent_id": "f6335c96-0529-4e90-b67b-e2ffd00db160",
*         "image": "Content.Image",
*         "title": "Content.Title",
*         "category": "Content.ContentCategory.Names",
*         "caregory_id": "Content.ContentCategory.Id",
*         "actions": "Content.Actions[]",
*         "type": "Content.ContentType",
*         "source": "Content.Source",
*         "shared_by": "Content.SharedBy.Name",
*         "shared_by_image": "Content.SharedBy.Image"
*       },
*       ... ...
*     ],
*     "code": 200,
*     "message": "Read shared contents successfully"
*   }
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The bookmarks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetSharedContent",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetSharedContent(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetSharedContent API request")
	req_userapp := new(userapp_proto.GetSharedContentRequest)
	req_userapp.UserId = req.PathParameter("user_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetSharedContent(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetSharedContent", "QueryError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read shared contents successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/{user_id}/plan/shared?session={session_id} Get shared plan
* @apiVersion 0.1.0
* @apiName GetSharedPlansForUser
* @apiGroup UserApp
*
* @apiDescription get the plan shared with this user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/437c32b2-9dd7-410a-8a97-162e580d8a90/plan/shared?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "shared_plans":[
*       {
*         "plan_id": "f6335c96-0529-4e90-b67b-e2ffd00db160",
*         "image": "Plan.Image",
*         "title": "Plan.Title",
*         "duration": "P1Y2M3DT4H5M6S",
*         "count": "Plan.ItemsCount",
*         "shared_by": "Plan.SharedBy.Name",
*         "shared_by_image": "Plan.SharedBy.Image"
*       },
*       ... ...
*     ],
*     "code": 200,
*     "message": "Read shared plans successfully"
*   }
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The bookmarks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetSharedPlansForUser",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetSharedPlansForUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetSharedPlansForUser API request")
	req_userapp := new(userapp_proto.GetSharedPlanRequest)
	req_userapp.UserId = req.PathParameter("user_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetSharedPlansForUser(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetSharedPlansForUser", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read shared plans successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/{user_id}/survey/shared?session={session_id} Get shared survey
* @apiVersion 0.1.0
* @apiName GetSharedSurvey
* @apiGroup UserApp
*
* @apiDescription get the survey shared with this user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/437c32b2-9dd7-410a-8a97-162e580d8a90/survey/shared?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "shared_surveys":[
*       {
*         "survey_id": "f6335c96-0529-4e90-b67b-e2ffd00db160",
*         "image": "Survey.Image",
*         "title": "Survey.Title",
*         "count": "Survey.Questions.Count",
*         "shared_by": "Survey.SharedBy.Name",
*         "shared_by_image": "Survey.SharedBy.Image"
*       },
*       ... ...
*     ],
*     "code": 200,
*     "message": "Read shared surveys successfully"
*   }
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The bookmarks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetSharedSurvey",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetSharedSurveysForUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetSharedSurveysForUser API request")
	req_userapp := new(userapp_proto.GetSharedSurveyRequest)
	req_userapp.UserId = req.PathParameter("user_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetSharedSurveysForUser(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetSharedSurveysForUser", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read shared surveys successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/{user_id}/goal/shared?session={session_id} Get shared goal
* @apiVersion 0.1.0
* @apiName GetSharedGoal
* @apiGroup UserApp
*
* @apiDescription get the goal shared with this user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/437c32b2-9dd7-410a-8a97-162e580d8a90/goal/shared?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "shared_goals":[
*       {
*         "goal_id": "f6335c96-0529-4e90-b67b-e2ffd00db160",
*         "image": "Goal.Image",
*         "title": "Goal.Title",
*         "current": "Goal.Target.CurrentValue",
*         "target": "Goal.Users[0].TargetValue",
*         "duration": "P1Y2M3DT4H5M6S",
*         "shared_by": "Survey.SharedBy.Name",
*         "shared_by_image": "Survey.SharedBy.Image"
*       },
*       ... ...
*     ],
*     "code": 200,
*     "message": "Read shared goals successfully"
*   }
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The bookmarks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetSharedGoal",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetSharedGoalsForUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetSharedGoalsForUser API request")
	req_userapp := new(userapp_proto.GetSharedGoalRequest)
	req_userapp.UserId = req.PathParameter("user_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetSharedGoalsForUser(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetSharedGoalsForUser", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read shared goals successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/{user_id}/challenge/shared?session={session_id} Get shared challenge
* @apiVersion 0.1.0
* @apiName GetSharedChallenge
* @apiGroup UserApp
*
* @apiDescription get the challenge shared with this user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/437c32b2-9dd7-410a-8a97-162e580d8a90/challenge/shared?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "shared_challenge":[
*       {
*         "challenge_id": "f6335c96-0529-4e90-b67b-e2ffd00db160",
*         "image": "Challenge.Image",
*         "title": "Challenge.Title",
*         "current": "Challenge.Target.CurrentValue",
*         "target": "Challenge.Users[0].TargetValue",
*         "duration": "P1Y2M3DT4H5M6S",
*         "shared_by": "Challenge.SharedBy.Name",
*         "shared_by_image": "Challenge.SharedBy.Image"
*       },
*       ... ...
*     ],
*     "code": 200,
*     "message": "Read shared challenges successfully"
*   }
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The bookmarks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetSharedChallenge",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetSharedChallengesForUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetSharedChallengesForUser API request")
	req_userapp := new(userapp_proto.GetSharedChallengeRequest)
	req_userapp.UserId = req.PathParameter("user_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetSharedChallengesForUser(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetSharedChallengesForUser", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read shared challenges successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/{user_id}/habit/shared?session={session_id} Get shared habit
* @apiVersion 0.1.0
* @apiName GetSharedHabit
* @apiGroup UserApp
*
* @apiDescription get the habit shared with this user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/437c32b2-9dd7-410a-8a97-162e580d8a90/habit/shared?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "shared_habits":[
*       {
*         "habit_id": "f6335c96-0529-4e90-b67b-e2ffd00db160",
*         "image": "Habit.Image",
*         "title": "Habit.Title",
*         "current": "Habit.Target.CurrentValue",
*         "target": "Habit.Users[0].TargetValue",
*         "shared_by": "Habit.SharedBy.Name",
*         "shared_by_image": "Habit.SharedBy.Image"
*       },
*       ... ...
*     ],
*     "code": 200,
*     "message": "Read shared habits successfully"
*   }
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The bookmarks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetSharedHabit",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetSharedHabitsForUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetSharedHabitsForUser API request")
	req_userapp := new(userapp_proto.GetSharedHabitRequest)
	req_userapp.UserId = req.PathParameter("user_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetSharedHabitsForUser(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetSharedHabitsForUser", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read shared habits successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/user/app/goal/join?session={session_id} Signup to a goal
* @apiVersion 0.1.0
* @apiName SignupToGoal
* @apiGroup UserApp
*
* @apiDescription signup to a shared goal
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/goal/join?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "goal_id": { Share_goal_user.id },
*   "user_id": { User.id },
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "count": 3,
*     "join_goal": {
*       "goal": { Goal },
*       "user": { User },
*       "status": "STARTED",
*       "start": 172435323,
*       "end": 1724335268,
*       "target": { Target },
*     },
*   },
*   "code": 200,
*   "message": "Signup to a goal successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The goal were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.SignupToGoal",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) SignupToGoal(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.SignupToGoal API request")

	req_userapp := new(userapp_proto.SignupToGoalRequest)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)
	err := req.ReadEntity(req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.SignupToGoal", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.SignupToGoal(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.SignupToGoal", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Signup to a goal correctly"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/goals/joined?session={session_id} List all user's goal
* @apiVersion 0.1.0
* @apiName GetAllJoinedGoals
* @apiGroup UserApp
*
* @apiDescription list all the goals that user has joined from the JoinGoal edge table. This will include goals in all different ActionStatus
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/goals/joined?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "join_goals": [
*       {
*         "goal_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "image": "Goal.Image",
*         "title": "Goal.Title",
*         "current": "Goal.Target.CurrentValue",
*         "target": "Goal.Users[0].TargetValue",
*         "duration": "P1Y2M3DT4H5M6S",
*         "shared_by": "Survey.SharedBy.Name",
*         "shared_by_image": "Survey.SharedBy.Image"
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Return all Joined goals correctly"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The goal were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetAllJoinedGoals",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetAllJoinedGoals(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetAllJoinedGoals API request")

	req_userapp := new(userapp_proto.ListGoalRequest)
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetAllJoinedGoals(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetAllJoinedGoals", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Return all Joined goals correctly"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/goals/current?session={session_id} List current goal
* @apiVersion 0.1.0
* @apiName GetCurrentJoinedGoals
* @apiGroup UserApp
*
* @apiDescription list the current INPROGRESS or STARTED goal that user has joined from the JoinGoal edge table.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/goals/current?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "join_goals": [
*       {
*         "goal_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "image": "Goal.Image",
*         "title": "Goal.Title",
*         "current": "Goal.Target.CurrentValue",
*         "target": "Goal.Users[0].TargetValue",
*         "duration": "P1Y2M3DT4H5M6S",
*         "shared_by": "Survey.SharedBy.Name",
*         "shared_by_image": "Survey.SharedBy.Image"
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "List current joined goals correctly"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The goal were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetCurrentJoinedGoals",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetCurrentJoinedGoals(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetCurrentJoinedGoals API request")

	req_userapp := new(userapp_proto.ListGoalRequest)
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetCurrentJoinedGoals(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetCurrentJoinedGoals", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "List current joined goals correctly"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/user/app/challenge/join?session={session_id} Signup to a challenge
* @apiVersion 0.1.0
* @apiName SignupToChallenge
* @apiGroup UserApp
*
* @apiDescription signup to a shared challenge
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/challenge/join?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "challenge_id": { Share_challenge_user.id },
*   "user_id": { User.id },
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "count": 3,
*     "join_challenge": {
*       "challenge": { Challenge },
*       "user": { User },
*       "status": "STARTED",
*       "start": 172435323,
*       "end": 1724335268,
*       "target": { Target },
*     },
*   },
*   "code": 200,
*   "message": "Signup to a challenge successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenge were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.SignupToChallenge",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) SignupToChallenge(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.SignupToChallenge API request")

	req_userapp := new(userapp_proto.SignupToChallengeRequest)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)
	err := req.ReadEntity(req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.SignupToChallenge", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.SignupToChallenge(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.SignupToChallenge", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Signup to a challenge correctly"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/challenges/joined?session={session_id} List all user's challenge
* @apiVersion 0.1.0
* @apiName GetAllJoinedChallenges
* @apiGroup UserApp
*
* @apiDescription list all joined challenges that user has joined. This will include challenges in all different ActionStatus
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/challenges/joined?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "join_challenges": [
*       {
*         "challenge_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "image": "Challenge.Image",
*         "title": "Challenge.Title",
*         "current": "Challenge.Target.CurrentValue",
*         "target": "Challenge.Users[0].TargetValue",
*         "duration": "P1Y2M3DT4H5M6S",
*         "shared_by": "Survey.SharedBy.Name",
*         "shared_by_image": "Survey.SharedBy.Image"
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Returned all Joined challenges successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenge were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetAllJoinedChallenges",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetAllJoinedChallenges(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetAllJoinedChallenges API request")

	req_userapp := new(userapp_proto.ListChallengeRequest)
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetAllJoinedChallenges(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetAllJoinedChallenges", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Returned all Joined challenges successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/user/app/habit/join?session={session_id} Signup to a habit
* @apiVersion 0.1.0
* @apiName SignupToHabit
* @apiGroup UserApp
*
* @apiDescription signup to a shared habit
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/habit/join?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "habit_id": { Share_habit_user.id },
*   "user_id": { User.id },
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "count": 3,
*     "join_habit": {
*       "habit": { Habit },
*       "user": { User },
*       "status": "STARTED",
*       "start": 172435323,
*       "end": 1724335268,
*       "target": { Target },
*     },
*   },
*   "code": 200,
*   "message": "Signup to a habit successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The habit were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.SignupToHabit",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) SignupToHabit(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.SignupToHabit API request")

	req_userapp := new(userapp_proto.SignupToHabitRequest)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)
	err := req.ReadEntity(req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.SignupToHabit", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.SignupToHabit(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.SignupToHabit", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Signup to a habit correctly"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/habits/joined?session={session_id} List all user's habit
* @apiVersion 0.1.0
* @apiName GetAllJoinedHabits
* @apiGroup UserApp
*
* @apiDescription list all the habits that user has joined. This will include habits in all different ActionStatus
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/habits/joined?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "join_habits": [
*       {
*         "habit_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "image": "Habit.Image",
*         "title": "Habit.Title",
*         "current": "Habit.Target.CurrentValue",
*         "target": "Habit.Users[0].TargetValue",
*         "duration": "P1Y2M3DT4H5M6S",
*         "shared_by": "Survey.SharedBy.Name",
*         "shared_by_image": "Survey.SharedBy.Image"
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Returned all Joined habits successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The habit were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetAllJoinedHabits",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetAllJoinedHabits(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetAllJoinedHabits API request")

	req_userapp := new(userapp_proto.ListHabitRequest)
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetAllJoinedHabits(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetAllJoinedHabits", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Returned all Joined habits successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/habits/current?session={session_id} List current habit
* @apiVersion 0.1.0
* @apiName GetCurrentJoinedHabits
* @apiGroup UserApp
*
* @apiDescription list the current INPROGRESS or STARTED habit that user has joined from the JoinHabit edge table.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/habits/current?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "join_habits": [
*       {
*         "habit_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "image": "Habit.Image",
*         "title": "Habit.Title",
*         "current": "Habit.Target.CurrentValue",
*         "target": "Habit.Users[0].TargetValue",
*         "duration": "P1Y2M3DT4H5M6S",
*         "shared_by": "Survey.SharedBy.Name",
*         "shared_by_image": "Survey.SharedBy.Image"
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read current joined habits successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The habit were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.ListCurrentHabit",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetCurrentJoinedHabits(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetCurrentJoinedHabits API request")

	req_userapp := new(userapp_proto.ListHabitRequest)
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetCurrentJoinedHabits(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetCurrentJoinedHabits", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read current joined habits successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/habits/current/count?session={session_id} List current habits with count
* @apiVersion 0.1.0
* @apiName GetCurrentHabitsWithCount
* @apiGroup UserApp
*
* @apiDescription list the current INPROGRESS habits that user has joined with count.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/habits/current/count?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "habit_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "title": "Title",
*		  "image": "Image",
*         "current": 100,
*         "targe": 100
*         "duration": "P1Y2M3DT4H5M6S"
*         "count": 1
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "List current habits with count correctly"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The habits were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetCurrentHabitsWithCount",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *UserAppService) GetCurrentHabitsWithCount(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetCurrentHabitsWithCount API request")

	req_userapp := new(userapp_proto.GetCurrentHabitsWithCountRequest)
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetCurrentHabitsWithCount(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetCurrentHabitsWithCount", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "List current habits with count correctly"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/challenges/current?session={session_id} List current challenge
* @apiVersion 0.1.0
* @apiName GetCurrentJoinedChallenges
* @apiGroup UserApp
*
* @apiDescription list the current INPROGRESS or STARTED challenge that user has joined.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/challenges/current?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "join_challenges": [
*       {
*         "challenge_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "image": "Challenge.Image",
*         "title": "Challenge.Title",
*         "current": "Challenge.Target.CurrentValue",
*         "target": "Challenge.Users[0].TargetValue",
*         "duration": "P1Y2M3DT4H5M6S",
*         "shared_by": "Survey.SharedBy.Name",
*         "shared_by_image": "Survey.SharedBy.Image"
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read current joined challenges successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenge were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.ListCurrentChallenge",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetCurrentJoinedChallenges(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetCurrentJoinedChallenges API request")

	req_userapp := new(userapp_proto.ListChallengeRequest)
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetCurrentJoinedChallenges(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetCurrentJoinedChallenges", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read current joined challenges successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/current/markers?session={session_id} List current markers
* @apiVersion 0.1.0
* @apiName ListMarkers
* @apiGroup UserApp
*
* @apiDescription Return a list of unique markers for the current goals and challenges
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/current/markers?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "markers": [
*       {
*         "is_default": true,
*         "marker_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "name": "Title",
*         "icon_slug": "ICON-SLUG",
*		  "unit":["unit1","unit2"],
*		  "tracker_methods":[{TracketMethod1},{TrackerMethod2}]
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "List markers for the current correctly"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenge were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.ListMarkers",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) ListMarkers(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.ListMarkers API request")

	req_userapp := new(userapp_proto.ListMarkersRequest)
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.ListMarkers(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.ListMarkers", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "List markers for the current correctly"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/pending?session={session_id} Get pending shared actions
* @apiVersion 0.1.0
* @apiName GetPendingSharedActions
* @apiGroup UserApp
*
* @apiDescription Return a list of all pending items that are in the Pending collection where ordered by date
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/pending?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "pendings": [
*       {
*         "id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "title": "Shared_X_Object",
*         "shared_by": { User },
*         "created": 14546365745
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get pending shared actions correctly"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenge were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetPendingSharedActions",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetPendingSharedActions(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetPendingSharedActions API request")

	req_userapp := new(userapp_proto.GetPendingSharedActionsRequest)
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_userapp.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_userapp.Offset = req.Attribute(PaginateOffsetParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetPendingSharedActions(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetPendingSharedActions", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get pending shared actions correctly"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/goal/current/progress?session={session_id} Get current goal progress
* @apiVersion 0.1.0
* @apiName GetGoalProgress
* @apiGroup UserApp
*
* @apiDescription Return progress on the current goal
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/goal/current/progress?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "goal": { Goal },
*         "user": { User },
*         "latestValue": 4,
*         "target": 100,
*         "unit": Kgs
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get goal progress successully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenge were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetGoalProgress",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetGoalProgress(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetGoalProgress API request")

	req_userapp := new(userapp_proto.GetGoalProgressRequest)
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetGoalProgress(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetGoalProgress", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get goal progress successully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/marker/default/history?session={session_id} Get default marker history
* @apiVersion 0.1.0
* @apiName GetDefaultMarkerHistory
* @apiGroup UserApp
*
* @apiDescription Return tracking history of the default marker for the goal using from, to and limit and offset
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/marker/default/history?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "user": { User },
*         "org_id": "orgid",
*         "marker": { Marekr }
*         "created": 1524324325
*         "value": String|Bool|Byte|Number,
*         "unit": "unit string"
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get default marker hsitory correctly"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenge were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetDefaultMarkerHistory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *UserAppService) GetDefaultMarkerHistory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetDefaultMarkerHistory API request")

	req_userapp := new(track_proto.GetDefaultMarkerHistoryRequest)
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)
	req_userapp.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_userapp.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_userapp.SortParameter = req.Attribute(SortParameter).(string)
	req_userapp.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetDefaultMarkerHistory(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetDefaultMarkerHistory", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get default marker hsitory correctly"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/challenges/current/count?session={session_id} List current challenges with count
* @apiVersion 0.1.0
* @apiName GetCurrentChallengesWithCount
* @apiGroup UserApp
*
* @apiDescription list the current INPROGRESS challenge that user has joined from the JoinChallenge edge table.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/challenges/current/count?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "challenge_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "title": "Title",
*		  "image": "Image",
*         "current": 100,
*         "targe": 100
*         "duration": "P1Y2M3DT4H5M6S"
*         "count": 1
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "List current challenges with count correctly"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenge were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetCurrentChallengesWithCount",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *UserAppService) GetCurrentChallengesWithCount(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetCurrentChallengesWithCount API request")

	req_userapp := new(userapp_proto.GetCurrentChallengesWithCountRequest)
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetCurrentChallengesWithCount(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetCurrentChallengesWithCount", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "List current challenges with count correctly"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/content/categorys/all?session={session_id} Get content categories
* @apiVersion 0.1.0
* @apiName GetContentCategorys
* @apiGroup UserApp
*
* @apiDescription Get content categories returns a subset of values from a list of content.categorys
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/content/categorys/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "categorys": [
*       {
*         "name": "Category.Name",
*         "icon_slug": "IconSlug",
*         "category_id": "437c32b2-9dd7-410a-8a97-162e580d8a90"
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get content categories successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenge were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetContentCategorys",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *UserAppService) GetContentCategorys(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetContentCategorys API request")

	req_userapp := new(content_proto.GetContentCategorysRequest)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetContentCategorys(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetContentCategorys", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get content categories successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/content/{content_id}?session={session_id} Get content detail
* @apiVersion 0.1.0
* @apiName GetContentDetail
* @apiGroup UserApp
*
* @apiDescription Get content detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/content/437c32b2-9dd7-410a-8a97-162e580d8a90?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "user_id": "856c32b2-sdf7-9dfa-8a97-987asdfa90",
*     "content": {Content},
*     "rating": 10,
*     "bookmarked": true,
* 	"actions": [action, action ],
* 	"bookmarked": true,
* 	"category": "Relationships",
* 	"category_id": "b76be63a-58bc-11e8-b913-42010a9a0002",
* 	"id": "79a49660-91c2-11e8-9446-00155d4b0101",
* 	"item": {
* 		"@type": "healum.com/proto/go.micro.srv.content.Recipe"
* 	},
* 	"shared_by": {user},
* 	"title": "fsadfsd"
*   },
*   "code": 200,
*   "message": "Get content detail successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenge were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetContentDetail",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *UserAppService) GetContentDetail(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetContentDetail API request")

	req_userapp := new(content_proto.ReadContentRequest)
	req_userapp.Id = req.PathParameter("content_id")
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetContentDetail(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetContentDetail", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get content detail successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/content/category/{category_id}?session={session_id} Get content from a category
* @apiVersion 0.1.0
* @apiName GetContentByCategory
* @apiGroup UserApp
*
* @apiDescription Get content from a category
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/content/category/437c32b2-9dd7-410a-8a97-162e580d8a90?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contents": [
*       {
*         "image": "Content.Image",
*         "title": "Content.Title",
*         "author": "Content.Author",
*         "source": { Content.Source },
*         "content_id": "Content.Id",
*         "category_id": "Content.ContentCategory.Id",
*         "icon_lsug": "Content.ContentCategory.IconSlug",
*         "category_name": "Content.ContentCategory.Name"
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get content from a category successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The content were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetContentByCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *UserAppService) GetContentByCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetContentByCategory API request")

	req_userapp := new(content_proto.GetContentByCategoryRequest)
	req_userapp.CategoryId = req.PathParameter("category_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetContentByCategory(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetContentByCategory", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get content from a category successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/content/category/{category_id}/filters?session={session_id} Get filters for a category
* @apiVersion 0.1.0
* @apiName GetFiltersForCategory
* @apiGroup UserApp
*
* @apiDescription Get filters for a category
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/content/category/437c32b2-9dd7-410a-8a97-162e580d8a90/filters?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentCategoryItems": [
*       {
*         "id": "111",
*         "name": "title",
*         "name_slug": "nameSlug",
*         "icon_slug": "iconSlug",
*         "summary": "summary",
*         "description": "description",
*         "org_id": "orgid",
*         "tags": ["tag1", "tag2"],
*         "taxonomy": { Taxonomy },
*         "weight": 100,
*         "priority": 1,
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get filters for a category successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The content were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetFiltersForCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *UserAppService) GetFiltersForCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetFiltersForCategory API request")

	req_userapp := new(content_proto.GetFiltersForCategoryRequest)
	req_userapp.CategoryId = req.PathParameter("category_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetFiltersForCategory(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetFiltersForCategory", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get filters for a category successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/user/app/content/category/filters/autocomplete?session={session_id} Filters autocomplete
* @apiVersion 0.1.0
* @apiName FiltersAutocomplete
* @apiGroup UserApp
*
* @apiDescription Filters autocomplete
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/content/category/filters/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "name": "sub_string"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "categorys": [
*       {
*         "name": "Category.Name",
*         "icon_slug": "IconSlug",
*         "category_id": "437c32b2-9dd7-410a-8a97-162e580d8a90"
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Filters autocomplete successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The habit were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.FiltersAutocomplete",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) FiltersAutocomplete(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.FiltersAutocomplete API request")

	req_userapp := new(content_proto.FiltersAutocompleteRequest)
	err := req.ReadEntity(req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.FiltersAutocomplete", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.FiltersAutocomplete(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.FiltersAutocomplete", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Filters autocomplete successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/user/app/content/category/{category_id}/filter?session={session_id} Filter content in a particular category
* @apiVersion 0.1.0
* @apiName FilterContentInParticularCategory
* @apiGroup UserApp
*
* @apiDescription Filter content in a particular category as per the filter. Return content where contentcategoryitems match content.tags
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/content/category/{category_id}/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "contentcategoryItems": ["_id1","_id2"]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contents": [
*       {
*         "image": "Content.Image",
*         "title": "Content.Title",
*         "author": "Content.Author",
*         "source": { Content.Source },
*         "content_id": "Content.Id",
*         "category_id": "Content.ContentCategory.Id",
*         "icon_lsug": "Content.ContentCategory.IconSlug",
*         "category_name": "Content.ContentCategory.Name"
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Filter content in a particular category successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The content were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.FilterContentInParticularCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) FilterContentInParticularCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.FilterContentInParticularCategory API request")

	req_userapp := new(content_proto.FilterContentInParticularCategoryRequest)
	err := req.ReadEntity(req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.FilterContentInParticularCategory", "BindError")
		return
	}
	req_userapp.CategoryId = req.PathParameter("category_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.FilterContentInParticularCategory(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.FilterContentInParticularCategory", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Filter content in a particular category successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/{user_id}/preferences?session={session_id} Get user preferences
* @apiVersion 0.1.0
* @apiName GetUserPreference
* @apiGroup UserApp
*
* @apiDescription API allows a user to get their preferences
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=/preferences?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "preference": {
*       "allergies": [
*         {
*           "category": null,
*           "created": "0",
*           "description": "",
*           "icon_slug": "",
*           "id": "",
*           "name": "name_2",
*           "name_slug": "name_slug_2",
*           "org_id": "",
*           "priority": "0",
*           "summary": "",
*           "tags": [],
*           "taxonomy": null,
*           "updated": "0",
*           "weight": "0"
*         }
*       ],
*       "conditions": [
*         {
*           "category": null,
*           "created": "0",
*           "description": "",
*           "icon_slug": "",
*           "id": "",
*           "name": "name_1",
*           "name_slug": "name_slug_1",
*           "org_id": "",
*           "priority": "0",
*           "summary": "",
*           "tags": [],
*           "taxonomy": null,
*           "updated": "0",
*           "weight": "0"
*         }
*       ],
*       "created": "1523701578",
*       "cuisines": [
*         {
*           "category": null,
*           "created": "0",
*           "description": "",
*           "icon_slug": "",
*           "id": "",
*           "name": "name_4",
*           "name_slug": "name_slug_4",
*           "org_id": "",
*           "priority": "0",
*           "summary": "",
*           "tags": [],
*           "taxonomy": null,
*           "updated": "0",
*           "weight": "0"
*         }
*       ],
*       "currentMeasurements": [
*         {
*           "created": "0",
*           "id": "measure_id",
*           "marker": null,
*           "measuredBy": null,
*           "method": null,
*           "org_id": "orgid",
*           "unit": "",
*           "updated": "0",
*           "userId": "userid",
*           "value": null
*         }
*       ],
*       "ethinicties": [
*         {
*           "category": null,
*           "created": "0",
*           "description": "",
*           "icon_slug": "",
*           "id": "",
*           "name": "name_5",
*           "name_slug": "name_slug_5",
*           "org_id": "",
*           "priority": "0",
*           "summary": "",
*           "tags": [],
*           "taxonomy": null,
*           "updated": "0",
*           "weight": "0"
*         }
*       ],
*       "food": [
*         {
*           "category": null,
*           "created": "0",
*           "description": "",
*           "icon_slug": "",
*           "id": "",
*           "name": "name_3",
*           "name_slug": "name_slug_3",
*           "org_id": "",
*           "priority": "0",
*           "summary": "",
*           "tags": [],
*           "taxonomy": null,
*           "updated": "0",
*           "weight": "0"
*         }
*       ],
*       "id": "449c6ec2-3fce-11e8-8085-20c9d0453b15",
*       "org_id": "orgid",
*       "updated": "1523701578",
*       "userId": "userid"
*     }
*   },
*   "code": 200,
*   "message": "Read user successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetUserPreference",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetUserPreference",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserAppService) GetUserPreference(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.GetUserPreference API request")
	req_user := new(user_proto.ReadUserPreferenceRequest)
	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetUserPreference(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetUserPreference", "ReadError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read user preferences successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/user/app/preferences?session={session_id} Save a users preferences
* @apiVersion 0.1.0
* @apiName SaveUserPreference
* @apiGroup UserApp
*
* @apiDescription The API endpoint for this will be in user-app-srv. This should call the userClient to update preferences.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/preferences?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "preference": {
*     "org_id": "orgid",
*     "currentMeasurements": [
*       {
*         "user_id": "UserId",
*         "org_id": "org_id",
*         "marker": { Marker },
*         "method": { TrackerMethod} ,
*         "measuredBy": { User },
*         "value": {
*           "@type": ""
*         },
*         "unit": "unit",
*       },
*       ... ...
*     ],
*     "conditions": [
*       {
*         "id": "111",
*         "name": "title",
*         "name_slug": "nameSlug",
*         "icon_slug": "iconSlug",
*         "summary": "summary",
*         "description": "description",
*         "org_id": "orgid",
*         "tags": ["tag1", "tag2"],
*         "taxonomy": { Taxonomy },
*         "weight": 100,
*         "priority": 1,
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ],
*     "allergies": [ {ContentCategoryItem} ],
*     "food": [ {ContentCategoryItem} ],
*     "cuisines": [ {ContentCategoryItem} ],
*     "ethinicties": [ {ContentCategoryItem} ]
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "preference": [
*       {
*         "image": "Content.Image",
*         "title": "Content.Title",
*         "author": "Content.Author",
*         "source": { Content.Source },
*         "content_id": "Content.Id",
*         "category_id": "Content.ContentCategory.Id",
*         "icon_lsug": "Content.ContentCategory.IconSlug",
*         "category_name": "Content.ContentCategory.Name"
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Save a users preferences successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.SaveUserPreference",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) SaveUserPreference(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.SaveUserPreference API request")

	req_userapp := new(user_proto.SaveUserPreferenceRequest)
	if err := utils.UnmarshalAny(req, rsp, req_userapp); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.SaveUserPreference", "BindError")
		return
	}
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.SaveUserPreference(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.SaveUserPreference", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Save a users preferences successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/user/app/details?session={session_id} Save or update user details
* @apiVersion 0.1.0
* @apiName SaveUserDetails
* @apiGroup UserApp
*
* @apiDescription The API endpoint for this will be in user-app-srv. This should call the userClient to update user details.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/details?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*     "firstname": "david",
*     "lastname": "john",
*     "avatar_url": "http://example.com",
*     "gender": "MALE"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "user": {
*       "id": "userid",
*       "org_id": "orgid",
*       "firstname": "david",
*       "lastname": "john",
*       "avatar_url": "http://example.com",
*       "gender": "MALE",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Updated user details successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.SaveUserDetails",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) SaveUserDetails(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.SaveUserDetails API request")

	req_userapp := new(userapp_proto.SaveUserDetailsRequest)
	if err := utils.UnmarshalAny(req, rsp, req_userapp); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.SaveUserDetails", "BindError")
		return
	}
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.SaveUserDetails(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.SaveUserDetails", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Updated user details successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/content/recommendations/all?session={session_id} Get content recommendations for a particular user
* @apiVersion 0.1.0
* @apiName GetContentRecommendationByUser
* @apiGroup UserApp
*
* @apiDescription This API endpoint will use contentClient to fetch the data from content_recommendation edge table
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/content/recommendations/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "recommendations": [
*       {
*         "content_id": "content.id",
*         "content_title": "content.title",
*         "content_author": "content.author",
*         "content_source": { content.source },
*         "category_id": "content.category.id",
*         "category_icon_slug": "content.category.icon_slug",
*         "category_name": "content.category.name",
*         "user_id": "user_id",
*         "tags": [ {ContentCategoryItem}, ... ]
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get content recommendations for a particular user successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The content recommendation were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetContentRecommendationByUser",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *UserAppService) GetContentRecommendationByUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetContentRecommendationByUser API request")

	req_userapp := new(content_proto.GetContentRecommendationByUserRequest)
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetContentRecommendationByUser(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetContentRecommendationByUser", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get content recommendations for a particular user successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/content/recommendations/category/{category_id}?session={session_id} Get content recommendations by category
* @apiVersion 0.1.0
* @apiName GetContentRecommendationByCategory
* @apiGroup UserApp
*
* @apiDescription This API endpoint will use contentClient to fetch the data from content_recommendation edge table based on a particular
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/content/recommendations/category/437c32b2-9dd7-410a-8a97-162e580d8a90?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "recommendations": [
*       {
*         "content_id": "content.id",
*         "content_title": "content.title",
*         "content_author": "content.author",
*         "content_source": { content.source },
*         "category_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "category_icon_slug": "content.category.icon_slug",
*         "category_name": "content.category.name",
*         "user_id": "user_id",
*         "tags": [ {ContentCategoryItem}, ... ]
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get content recommendations by category successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The content recommendation were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetContentRecommendationByCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *UserAppService) GetContentRecommendationByCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetContentRecommendationByCategory API request")

	req_userapp := new(content_proto.GetContentRecommendationByCategoryRequest)
	req_userapp.CategoryId = req.PathParameter("category_id")
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetContentRecommendationByCategory(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetContentRecommendationByCategory", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get content recommendations by category successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/user/app/content/{content_id}/rating?session={session_id} Filter content in a particular category
* @apiVersion 0.1.0
* @apiName SaveRateForContent
* @apiGroup UserApp
*
* @apiDescription Allow user to rate a content item
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/content/437c32b2-9dd7-410a-8a97-162e580d8a90/rating?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "rating": 9
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "content_rating": {
*       "org_id": "orgid",
*       "user_id": "userid",
*       "content_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*       "rating": 9,
*       "created": 1420890823,
*       "updated": 1420890823
*     }
*   },
*   "code": 200,
*   "message": "Save user rating for a particular content object successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The content were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.SaveRateForContent",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) SaveRateForContent(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.SaveRateForContent API request")

	req_userapp := new(userapp_proto.SaveRateForContentRequest)
	err := req.ReadEntity(req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.SaveRateForContent", "BindError")
		return
	}
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_userapp.ContentId = req.PathParameter("content_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.SaveRateForContent(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.SaveRateForContent", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Save user rating for a particular content object successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/user/app/content/{content_id}/dislike?session={session_id} User not interested in a particular content
* @apiVersion 0.1.0
* @apiName DislikeForContent
* @apiGroup UserApp
*
* @apiDescription Allow user to dislike a content item
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/content/437c32b2-9dd7-410a-8a97-162e580d8a90/dislike?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "content_dislike": {
*       "org_id": "orgid",
*       "user_id": "userid",
*       "content_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*       "created": 1420890823,
*       "updated": 1420890823
*     }
*   },
*   "code": 200,
*   "message": "Allowed user to dislike a content item successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The content were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.DislikeForContent",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) DislikeForContent(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.DislikeForContent API request")

	req_userapp := new(userapp_proto.DislikeForContentRequest)
	err := req.ReadEntity(req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.DislikeForContent", "BindError")
		return
	}
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_userapp.ContentId = req.PathParameter("content_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.DislikeForContent(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.DislikeForContent", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Allowed user to dislike a content item successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/user/app/content/{content_id}/dislike/similar?session={session_id} User not interested in similar content
* @apiVersion 0.1.0
* @apiName DislikeForSimilarContent
* @apiGroup UserApp
*
* @apiDescription Allow user to dislike content that's simillar to this one
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/content/437c32b2-9dd7-410a-8a97-162e580d8a90/dislike/similar?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "tags": [ {ContentCategoryItem}, ... ]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "content_dislike_similar": {
*       "org_id": "orgid",
*       "user_id": "userid",
*       "content_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*       "tags": [ {ContentCategoryItem}, ... ],
*       "created": 1420890823,
*       "updated": 1420890823
*     }
*   },
*   "code": 200,
*   "message": "Allowed user to dislike content that's simillar to this one successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The content were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.DislikeForSimilarContent",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) DislikeForSimilarContent(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.DislikeForSimilarContent API request")

	req_userapp := new(userapp_proto.DislikeForSimilarContentRequest)
	err := req.ReadEntity(req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.DislikeForSimilarContent", "BindError")
		return
	}
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_userapp.ContentId = req.PathParameter("content_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.DislikeForSimilarContent(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.DislikeForSimilarContent", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Allowed user to dislike content that's simillar to this one successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/user/app/feedback?session={session_id} Save general user feedback
* @apiVersion 0.1.0
* @apiName SaveUserFeedback
* @apiGroup UserApp
*
* @apiDescription Save user comment / feedback
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/feedback?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "feedback": "Very Nice!!!",
*   "rating": 9
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "feedback": {
*       "org_id": "orgid",
*       "user_id": "userid",
*       "feedback": "Very Nice!!!",
*       "rating": 9,
*       "created": 1420890823,
*       "updated": 1420890823
*     }
*   },
*   "code": 200,
*   "message": "Save user comment feedback successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The content were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.SaveUserFeedback",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) SaveUserFeedback(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.SaveUserFeedback API request")

	req_userapp := new(userapp_proto.SaveUserFeedbackRequest)
	err := req.ReadEntity(req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.SaveUserFeedback", "BindError")
		return
	}
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.SaveUserFeedback(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.SaveUserFeedback", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Save user comment feedback successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/user/app/plan/join?session={session_id} Join/start a plan
* @apiVersion 0.1.0
* @apiName JoinUserPlan
* @apiGroup UserApp
*
* @apiDescription Create a plan for the user when the user choses to join / start a plan shared with the user.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/plan/join?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "user_id": "userid",
*   "plan_id": "planid"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "user_plan": {
*       "id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*       "name": "plan_name",
*       "org_id": "orgid",
*       "pic": "http://example.com",
*       "description": "description",
*       "created": 1420890823,
*       "updated": 1420890823,
*       "targetUser": "userid",
*       "plan": { Plan },
*       "goals": [ {Goal}, {Goal} ],
*       "duration": "PDT3",
*       "start": 1420890823,
*       "end": 1540890823,
*       "creator": { User },
*       "days": {
*         "1":{"items":[ {DayItem}, ...] },
*         "2":{"items":[ {DayItem}, ...] }
*       },
*       "items_count": 10,
*     }
*   },
*   "code": 200,
*   "message": "Join/start a plan successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plan were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.JoinUserPlan",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) JoinUserPlan(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.JoinUserPlan API request")

	req_userapp := new(userapp_proto.JoinUserPlanRequest)
	err := req.ReadEntity(req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.JoinUserPlan", "BindError")
		return
	}
	req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.JoinUserPlan(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.JoinUserPlan", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Join/start a plan successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/user/app/plan/create?session={session_id} Create a user's plan
* @apiVersion 0.1.0
* @apiName CreateUserPlan
* @apiGroup UserApp
*
* @apiDescription Create a plan for the user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/plan/create?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "user_id": "userid",
*   "goal_id": "goalid",
*   "days": 3,
*   "itemsPerDay": 7
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "user_plan": {
*       "id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*       "name": "plan_name",
*       "org_id": "orgid",
*       "pic": "http://example.com",
*       "description": "description",
*       "created": 1420890823,
*       "updated": 1420890823,
*       "targetUser": "userid",
*       "plan": {
*         "collaborators": [],
*         "creator": {
*           "addresses": [],
*           "contactDetails": [],
*           "created": "0",
*           "dob": "0",
*           "firstname": "",
*           "gender": "Gender_NONE",
*           "id": "222",
*           "image": "",
*           "lastname": "",
*           "org_id": "",
*           "tokens": []
*         },
*         "days": {
*           "1": {
*             "items": [
*               {
*                 "categoryIconSlug": "",
*                 "categoryId": "",
*                 "categoryName": "",
*                 "contentId": "",
*                 "contentPicUrl": "",
*                 "contentTitle": "",
*                 "id": "day_item_001",
*                 "options": [],
*                 "post": {
*                   "created": "0",
*                   "creator": {
*                     "addresses": [],
*                     "contactDetails": [],
*                     "created": "0",
*                     "dob": "0",
*                     "firstname": "",
*                     "gender": "Gender_NONE",
*                     "id": "userid",
*                     "image": "",
*                     "lastname": "",
*                     "org_id": "",
*                     "tokens": [],
*                     "updated": "0"
*                   },
*                   "id": "111",
*                   "items": [],
*                   "name": "todo1",
*                   "org_id": "orgid",
*                   "updated": "0"
*                 },
*                 "pre": {
*                   "created": "0",
*                   "creator": {
*                     "addresses": [],
*                     "contactDetails": [],
*                     "created": "0",
*                     "dob": "0",
*                     "firstname": "",
*                     "gender": "Gender_NONE",
*                     "id": "userid",
*                     "image": "",
*                     "lastname": "",
*                     "org_id": "",
*                     "tokens": [],
*                     "updated": "0"
*                   },
*                   "id": "111",
*                   "items": [],
*                   "name": "todo1",
*                   "org_id": "orgid",
*                   "updated": "0"
*                 },
*                 "primary": false,
*                 "time": ""
*               }
*             ]
*           },
*           "2": {
*             "items": [
*               {
*                 "categoryIconSlug": "",
*                 "categoryId": "",
*                 "categoryName": "",
*                 "contentId": "",
*                 "contentPicUrl": "",
*                 "contentTitle": "",
*                 "id": "day_item_002",
*                 "options": [],
*                 "post": {
*                   "created": "0",
*                   "creator": {
*                     "addresses": [],
*                     "contactDetails": [],
*                     "created": "0",
*                     "dob": "0",
*                     "firstname": "",
*                     "gender": "Gender_NONE",
*                     "id": "userid",
*                     "image": "",
*                     "lastname": "",
*                     "org_id": "",
*                     "tokens": [],
*                     "updated": "0"
*                   },
*                   "id": "111",
*                   "items": [],
*                   "name": "todo1",
*                   "org_id": "orgid",
*                   "updated": "0"
*                 },
*                 "pre": {
*                   "created": "0",
*                   "creator": {
*                     "addresses": [],
*                     "contactDetails": [],
*                     "created": "0",
*                     "dob": "0",
*                     "firstname": "",
*                     "gender": "Gender_NONE",
*                     "id": "userid",
*                     "image": "",
*                     "lastname": "",
*                     "org_id": "",
*                     "tokens": [],
*                     "updated": "0"
*                   },
*                   "id": "111",
*                   "items": [],
*                   "name": "todo1",
*                   "org_id": "orgid",
*                   "updated": "0"
*                 },
*                 "primary": false,
*                 "time": ""
*               }
*             ]
*           }
*         },
*         "description": "hello world",
*         "duration": "P1Y2M3DT4H5M6S",
*         "end": "0",
*         "endTimeUnspecified": false,
*         "goals": [
*            {
*              "articles": [],
*              "category": null,
*              "challenges": [],
*              "completionApprovalRequired": false,
*              "created": "0",
*              "createdBy": null,
*              "description": "",
*              "duration": "",
*              "habits": [],
*              "id": "1",
*              "image": "",
*              "notifications": [],
*              "org_id": "",
*              "setbacks": [],
*              "social": [],
*              "source": "",
*              "status": "Status_NONE",
*              "successCriterias": [],
*              "summary": "",
*              "tags": [],
*              "target": null,
*              "title": "",
*              "trackers": [],
*              "triggers": [],
*              "updated": "0",
*              "users": [],
*              "visibility": "Visibility_NONE"
*            },
*            {
*              "articles": [],
*              "category": null,
*              "challenges": [],
*              "completionApprovalRequired": false,
*              "created": "0",
*              "createdBy": null,
*              "description": "",
*              "duration": "",
*              "habits": [],
*              "id": "2",
*              "image": "",
*              "notifications": [],
*              "org_id": "",
*              "setbacks": [],
*              "social": [],
*              "source": "",
*              "status": "Status_NONE",
*              "successCriterias": [],
*              "summary": "",
*              "tags": [],
*              "target": null,
*              "title": "",
*              "trackers": [],
*              "triggers": [],
*              "updated": "0",
*              "users": [],
*              "visibility": "Visibility_NONE"
*            }
*          ],
*          "id": "1212",
*          "isTemplate": true,
*          "itemsCount": "2",
*          "linkSharingEnabled": false,
*          "name": "plan1",
*          "org_id": "",
*          "pic": "",
*          "recurrence": [],
*          "shares": [
*            {
*              "addresses": [],
*              "contactDetails": [],
*              "created": "0",
*              "dob": "0",
*              "firstname": "",
*              "gender": "Gender_NONE",
*              "id": "userid",
*              "image": "",
*              "lastname": "",
*              "org_id": "",
*              "tokens": [],
*              "updated": "0"
*            }
*          ],
*          "start": "0",
*          "status": "DRAFT",
*          "templateId": "template1",
*          "updated": "1523428756",
*          "users": [
*             {
*               "addresses": [],
*               "contactDetails": [],
*               "created": "0",
*               "dob": "0",
*               "firstname": "",
*               "gender": "Gender_NONE",
*               "id": "userid",
*               "image": "",
*               "lastname": "",
*               "org_id": "",
*               "tokens": [],
*               "updated": "0"
*             }
*           ],
*          "setting": {
*            "embeddingEnabled": false,
*            "linkSharingEnabled": false,
*            "notifications": [],
*            "shareableLink": "",
*            "social": [],
*            "visibility": "PUBLIC"
*          },
*         }
*       }
*       "goals": [ {Goal}, {Goal} ],
*       "duration": "PDT3",
*       "start": 1420890823,
*       "end": 1540890823,
*       "creator": { User },
*       "days": {
*         "1":{"items":[ {DayItem}, ...] },
*         "2":{"items":[ {DayItem}, ...] }
*       },
*       "items_count": 10,
*     }
*   },
*   "code": 200,
*   "message": "Create a user's plan successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plan were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.CreateUserPlan",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) CreateUserPlan(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.CreateUserPlan API request")

	req_userapp := new(userapp_proto.CreateUserPlanRequest)
	// err := req.ReadEntity(req_userapp)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.CreateUserPlan", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_userapp); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.CreateUserPlan", "BindError")
		return
	}
	// req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.CreateUserPlan(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.CreateUserPlan", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Create a user's plan successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/plan/{user_id}?session={session_id} Get a user's plan
* @apiVersion 0.1.0
* @apiName GetUserPlan
* @apiGroup UserApp
*
* @apiDescription Get a user's plan
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/plan/{user_id}?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "user_plan": {
*       "id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*       "name": "plan_name",
*       "org_id": "orgid",
*       "pic": "http://example.com",
*       "description": "description",
*       "created": 1420890823,
*       "updated": 1420890823,
*       "targetUser": "userid",
*       "plan": { Plan },
*       "goals": [ {Goal}, {Goal} ],
*       "duration": "PDT3",
*       "start": 1420890823,
*       "end": 1540890823,
*       "creator": { User },
*       "days": {
*         "1":{"items":[ {DayItem}, ...] },
*         "2":{"items":[ {DayItem}, ...] }
*       },
*       "items_count": 10,
*     }
*   },
*   "code": 200,
*   "message": "Get a user's plan successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plan were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetUserPlan",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetUserPlan(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetUserPlan API request")

	req_userapp := new(userapp_proto.GetUserPlanRequest)
	req_userapp.UserId = req.PathParameter("user_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetUserPlan(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetUserPlan", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Join/start a plan successfully successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/user/app/plan/update?session={session_id} Update user's plan
* @apiVersion 0.1.0
* @apiName UpdateUserPlan
* @apiGroup UserApp
*
* @apiDescription Update user's plan
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/plan/update?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "id": "id",
*   "org_id": "orgid",
*   "goals": [ {Goal}, ... ],
*   "days": {
*     "1":{"items":[ {DayItem}, ...] },
*     "2":{"items":[ {DayItem}, ...] }
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "user_plan": {
*       "id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*       "name": "plan_name",
*       "org_id": "orgid",
*       "pic": "http://example.com",
*       "description": "description",
*       "created": 1420890823,
*       "updated": 1420890823,
*       "targetUser": "userid",
*       "plan": { Plan },
*       "goals": [ {Goal}, {Goal} ],
*       "duration": "PDT3",
*       "start": 1420890823,
*       "end": 1540890823,
*       "creator": { User },
*       "days": {
*         "1":{"items":[ {DayItem}, ...] },
*         "2":{"items":[ {DayItem}, ...] }
*       },
*       "items_count": 10,
*     }
*   },
*   "code": 200,
*   "message": "Update user's plan successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plan were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.UpdateUserPlan",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) UpdateUserPlan(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.UpdateUserPlan API request")

	req_userapp := new(userapp_proto.UpdateUserPlanRequest)
	// err := req.ReadEntity(req_userapp)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.UpdateUserPlan", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_userapp); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.UpdateUserPlan", "BindError")
		return
	}
	// req_userapp.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.UpdateUserPlan(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.UpdateUserPlan", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Update user's plan successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/plan/{plan_id}/summary/count?session={session_id} Get plan items count by category
* @apiVersion 0.1.0
* @apiName GetPlanItemsCountByCategory
* @apiGroup UserApp
*
* @apiDescription Get the current count of plan items by category across the entire plan
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/plan/{plan_id}/summary/count?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "content_count": [
*       {
*         "category_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "icon_slug": "cateogry_icon_slug",
*         "category_name": "category_name",
*         "item_count": 4
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get plan items count by category successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plan were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetPlanItemsCountByCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetPlanItemsCountByCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetPlanItemsCountByCategory API request")

	req_userapp := new(userapp_proto.GetPlanItemsCountByCategoryRequest)
	req_userapp.PlanId = req.PathParameter("plan_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetPlanItemsCountByCategory(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetPlanItemsCountByCategory", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get plan items count by category successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/plan/{plan_id}/summary/{day_number}/count?session={session_id} Get plan items count by day
* @apiVersion 0.1.0
* @apiName GetPlanItemsCountByDay
* @apiGroup UserApp
*
* @apiDescription Get the current count of plan items by category across a particular day
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/plan/{plan_id}/summary/{day_number}/count?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "content_count": [
*       {
*         "day_number": "1",
*         "category_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "icon_slug": "cateogry_icon_slug",
*         "category_name": "category_name",
*         "item_count": 4
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get plan items count by day successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plan were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetPlanItemsCountByDay",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetPlanItemsCountByDay(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetPlanItemsCountByDay API request")

	req_userapp := new(userapp_proto.GetPlanItemsCountByDayRequest)
	req_userapp.PlanId = req.PathParameter("plan_id")
	req_userapp.DayNumber = req.PathParameter("day_number")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetPlanItemsCountByDay(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetPlanItemsCountByDay", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get plan items count by day successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/plan/{plan_id}/summary/{day_number}/count?session={session_id} Get plan items count by category and day
* @apiVersion 0.1.0
* @apiName GetPlanItemsCountByCategoryAndDay
* @apiGroup UserApp
*
* @apiDescription Get the current count of plan items by category
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/plan/{plan_id}/summary/{day_number}/count?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "content_count": [
*       {
*         "day_number": "1",
*         "category_id": "437c32b2-9dd7-410a-8a97-162e580d8a90",
*         "icon_slug": "cateogry_icon_slug",
*         "category_name": "category_name",
*         "item_count": 4
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get plan items count by category and day successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plan were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetPlanItemsCountByCategoryAndDay",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) GetPlanItemsCountByCategoryAndDay(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetPlanItemsCountByCategoryAndDay API request")

	req_userapp := new(userapp_proto.GetPlanItemsCountByCategoryAndDayRequest)
	req_userapp.PlanId = req.PathParameter("plan_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetPlanItemsCountByCategoryAndDay(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetPlanItemsCountByCategoryAndDay", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get plan items count by category and day successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/user/app/login User login
* @apiVersion 0.1.0
* @apiName Login
* @apiGroup UserApp
*
* @apiDescription The functionality is to be able to login using password or code. The account needs to exist in the database for this to be a valid request.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/login
*
*
* @apiParamExample {json} Request-Phone:
* {
*   "phone": "123-4567-890",
*   "passcode": "12345",
*   "device_token": "token_string",
*   "unique_identifier": "unique_id",
*   "app_identifier": "app.bundle.id",
*   "platform": 1
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "session": {
*       "id": "f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=",
*       "expires_at": 153252466
*     },
*     "user": { User }
*   },
*   "code": 200,
*   "message": "Login successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, QueryError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "BindError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.account.Login",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *UserAppService) Login(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.Login API request")

	req_login := new(account_proto.LoginRequest)
	err := req.ReadEntity(req_login)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.Login", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.Login(ctx, req_login)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.Login", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Login successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/logout Logout a user
* @apiVersion 0.1.0
* @apiName Logout
* @apiGroup UserApp
*
* @apiDescription The functionality is to be able to login using password or code. The account needs to exist in the database for this to be a valid request.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/logout?session={session_id}
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Logout successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, QueryError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "BindError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.account.Logout",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *UserAppService) Logout(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.Logout API request")

	req_logout := new(account_proto.LogoutRequest)
	req_logout.SessionId = req.QueryParameter("session")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.Logout(ctx, req_logout)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.Logout", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Logout successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/goals/all?session={session_id}&offset={offset}&limit={limit} List all challenges List all goals
* @apiVersion 0.1.0
* @apiName AllGoals
* @apiGroup UserApp
*
* @apiDescription List all goals
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/goals/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "goals": [
*       {
*         "id": "g111",
*         "title": "g_title",
*         "image": "http://image.com",
*         "summary": "summary",
*         "createdby": "david john"
*         "createdby_pic": "http://image.com/image",
*         "target": { Target }
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all goals successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The goals were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.AllGoals",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.AllGoals",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserAppService) AllGoals(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.AllGoals API request")
	req_goal := new(user_proto.AllGoalResponseRequest)
	req_goal.UserId = req.Attribute(UserIdAttrName).(string)
	req_goal.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_goal.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_goal.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_goal.SortParameter = req.Attribute(SortParameter).(string)
	req_goal.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.AllGoalResponse(ctx, req_goal)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.AllGoals", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all goals successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/challenges/all?session={session_id}&offset={offset}&limit={limit} List all challenges
* @apiVersion 0.1.0
* @apiName AllChallenges
* @apiGroup UserApp
*
* @apiDescription List all challenges
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/challenges/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "challenges": [
*       {
*         "id": "c111",
*         "title": "c_title",
*         "image": "http://image.com",
*         "summary": "summary",
*         "createdby": "david john"
*         "createdby_pic": "http://image.com/image",
*         "target": { Target }
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all challenges successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenges were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.AllChallenges",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.AllChallenges",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserAppService) AllChallenges(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.AllChallenges API request")
	req_challenge := new(user_proto.AllChallengeResponseRequest)
	req_challenge.UserId = req.Attribute(UserIdAttrName).(string)
	req_challenge.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_challenge.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_challenge.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_challenge.SortParameter = req.Attribute(SortParameter).(string)
	req_challenge.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.AllChallengeResponse(ctx, req_challenge)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.AllChallenges", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all challenges successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/habits/all?session={session_id}&offset={offset}&limit={limit} List all habits
* @apiVersion 0.1.0
* @apiName AllHabits
* @apiGroup UserApp
*
* @apiDescription List all habits
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/habits/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "habits": [
*       {
*         "id": "h111",
*         "title": "h_title",
*         "image": "http://image.com",
*         "summary": "summary",
*         "createdby": "david john"
*         "createdby_pic": "http://image.com/image",
*         "target": { Target }
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all habits successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The habits were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.AllHabits",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.AllHabits",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserAppService) AllHabits(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.AllHabits API request")
	req_habit := new(user_proto.AllHabitResponseRequest)
	req_habit.UserId = req.Attribute(UserIdAttrName).(string)
	req_habit.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_habit.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_habit.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_habit.SortParameter = req.Attribute(SortParameter).(string)
	req_habit.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.AllHabitResponse(ctx, req_habit)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.AllHabits", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all habits successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/contents/all?session={session_id}&offset={offset}&limit={limit} List all contents
* @apiVersion 0.1.0
* @apiName GetShareableContent
* @apiGroup UserApp
*
* @apiDescription List all content
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/contents/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*	"type":["healum.com/proto/go.micro.srv.content.Recipe","healum.com/proto/go.micro.srv.content.Article","healum.com/proto/go.micro.srv.content.Video"],
*	"created_by":["user_id","user_id"]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "resources": [
*       {
*         "id": "g111",
*         "title": "g_title",
*         "image": "http://image.com",
*         "summary": "summary",
*         "createdby": "david john"
*         "createdby_pic": "http://image.com/image",
*         "target": { Target }
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all shareable content successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contents were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.GetShareableContent",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.GetShareableContent",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserAppService) GetShareableContent(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetShareableContent API request")
	req_habit := new(user_proto.GetShareableContentRequest)
	req_habit.UserId = req.Attribute(UserIdAttrName).(string)
	req_habit.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_habit.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_habit.Offset = req.Attribute(PaginateOffsetParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetShareableContent(ctx, req_habit)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetShareableContent", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all contents successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/user/app/{user_id}/shared?session={session_id} Update shared items status to RECIEVED
* @apiVersion 0.1.0
* @apiName ReceivedItems
* @apiGroup UserApp
*
* @apiDescription Update shared items status to RECIEVED
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/849ba14c-4a31-11e8-b4d2-20c9d0453b15/shared?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "shared": [
*     {
*       "type", "go.micro.srv.behaviour.Goal"
*       "id", "sdf7987df-45d1-90df-457s-gs09d8fds78"
*     },
*     ... ...
*   ]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Received all items successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The items were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.RecievedItems",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.RecievedItems",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserAppService) RecievedItems(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.RecievedItems API request")
	req_userapp := new(userapp_proto.ReceivedItemsRequest)
	err := req.ReadEntity(req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.RecievedItems", "BindError")
		return
	}
	req_userapp.UserId = req.Attribute(UserIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.ReceivedItems(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.RecievedItems", "QueryError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Received all items successfully"

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/{user_id}?session={session_id} Get user details
* @apiVersion 0.1.0
* @apiName ReadUser
* @apiGroup UserApp
*
* @apiDescription Get details of a particular user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/849ba14c-4a31-11e8-b4d2-20c9d0453b15?session=qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ=
*
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*     "code": "200",
*     "data": {
*         "user": {
*             "contact_details": [
*                 {
*                     "type": "PHONE",
*                     "value": "0123456789"
*                 }
*             ],
*             "created": "1524500000",
*             "firstname": "Test User",
*             "gender": "MALE",
*             "id": "849ba14c-4a31-11e8-b4d2-20c9d0453b15",
*             "lastname": "Test",
*             "org_id": "sdf7987df-45d1-90df-457s-gs09d8fds78",
*             "updated": "1524500000"
*         }
*     },
*     "message": "Read user successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The items were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.ReadUser",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.ReadUser",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserAppService) ReadUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.ReadUser API request")
	req_user := new(user_proto.ReadRequest)
	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.ReadUser(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.ReadUser", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read user successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/user/app/content/category/{nameslug}/items/autocomplete?session={session_id} Autocomplete content category item
* @apiVersion 0.1.0
* @apiName AutocompleteContentCategoryItem
* @apiGroup UserApp
*
* @apiDescription Autocomplete content category item
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/content/category/{nameslug}/items/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "name": "ti",
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "category_id": "111",
*         "category_nameslug": "slug",
*         "categoryitem_id": "222",
*         "categoryitem_name": name
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Autocomplete content category item successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, SearchError.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.AutocompleteContentCategoryItem",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) AutocompleteContentCategoryItem(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.AutocompleteContentCategoryItem API request")

	req_search := new(content_proto.AutocompleteContentCategoryItemRequest)
	err := req.ReadEntity(req_search)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.AutocompleteContentCategoryItem", "BindError")
		return
	}
	req_search.NameSlug = req.PathParameter("nameslug")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.AutocompleteContentCategoryItem(ctx, req_search)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.AutocompleteContentCategoryItem", "SearchError")
		return
	} else if resp.Data.Response == nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.AutocompleteContentCategoryItem", "NotFound")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Autocomplete content category item successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/content/category/{nameslug}/items/all?session={session_id} All content category items by name slug
* @apiVersion 0.1.0
* @apiName AllContentCategoryItemByNameslug
* @apiGroup UserApp
*
* @apiDescription All content category items by name slug
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/content/category/{nameslug}/items/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "category_id": "111",
*         "category_nameslug": "slug",
*         "categoryitem_id": "222",
*         "categoryitem_name": name
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "List all content category items by name slug successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, SearchError.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.AllContentCategoryItemByNameslug",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) AllContentCategoryItemByNameslug(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.AllContentCategoryItemByNameslug API request")

	req_search := new(content_proto.AllContentCategoryItemByNameslugRequest)
	err := req.ReadEntity(req_search)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.AllContentCategoryItemByNameslug", "BindError")
		return
	}
	req_search.NameSlug = req.PathParameter("nameslug")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.AllContentCategoryItemByNameslug(ctx, req_search)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.AllContentCategoryItemByNameslug", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "List all content category items by name slug"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/markers/{nameslug}/marker?session={session_id} All content category items by name slug
* @apiVersion 0.1.0
* @apiName MarkerByNameslug
* @apiGroup UserApp
*
* @apiDescription Return marker details by name slug
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/markers/{nameslug}/marker?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "category_id": "111",
*         "category_nameslug": "slug",
*         "categoryitem_id": "222",
*         "categoryitem_name": name
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "List all content category items by name slug successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, SearchError.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.MarkerByNameslug",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserAppService) MarkerByNameslug(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.MarkerByNameslug API request")

	req_marker := new(static_proto.ReadByNameslugRequest)

	req_marker.NameSlug = req.PathParameter("name_slug")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.ReadMarkerByNameslug(ctx, req_marker)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.MarkerByNameslug", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Get marker details by name slug"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/user/app/goal/{goal_id}?session={session_id} View goal detail
* @apiVersion 0.1.0
* @apiName GetGoalDetail
* @apiGroup UserApp
*
* @apiDescription View goal detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/goal/g111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "detail": {
*       "goal_id": "g111",
*       "title": "g_title",
*       "summary": "summary",
*       "description": "description",
*       "shared_by": "someon one",
*       "category": { "id" : "category111", ...},
*		"tags":["tags","tags"],
*		"image":"url",
*		"shared_by_image":"url",
*		"targer":0,
*		"duration":"P1D",
*		"source":"Some source",
*		"challenges":[{Challenge Object 1},{Challenge Object 2}],
*		"habits":[{Habit Object 1},{Habit Object 2}],
*		"todos":{"items":[]},
*		"category":{Category Object}
*     }
*   },
*   "code": 200,
*   "message": "Read goal successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The goal was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetGoalDetail",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetGoalDetail",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserAppService) GetGoalDetail(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetGoalDetail API request")
	req_goal := new(behaviour_proto.ReadGoalRequest)
	req_goal.GoalId = req.PathParameter("goal_id")
	req_goal.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetGoalDetail(ctx, req_goal)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetGoalDetail", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read goal successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/challenge/{challenge_id}?session={session_id} View challenge detail
* @apiVersion 0.1.0
* @apiName GetChallengeDetail
* @apiGroup UserApp
*
* @apiDescription View challenge detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/challenge/c111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "detail": {
*       "challenge_id": "g111",
*       "title": "g_title",
*       "summary": "summary",
*       "description": "description",
*       "shared_by": "someon one",
*       "category": { "id" : "category111", ...},
*		"tags":["tags","tags"],
*		"image":"url",
*		"shared_by_image":"url",
*		"targer":0,
*		"duration":"P1D",
*		"source":"Some source",
*		"habits":[{Habit Object 1},{Habit Object 2}],
*		"todos":{"items":[]},
*		"category":{Category Object}
*     },
*	}
*   "code": 200,
*   "message": "Read challenge successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenge was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetChallengeDetail",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetChallengeDetail",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserAppService) GetChallengeDetail(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetChallengeDetail API request")
	req_challenge := new(behaviour_proto.ReadChallengeRequest)
	req_challenge.ChallengeId = req.PathParameter("challenge_id")
	req_challenge.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetChallengeDetail(ctx, req_challenge)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetChallengeDetail", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read challenge successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/user/app/habit/{habit_id}?session={session_id} View habit detail
* @apiVersion 0.1.0
* @apiName GetHabitDetail
* @apiGroup UserApp
*
* @apiDescription View habit detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/user/app/habit/h111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*      "detail": {
*       "habit_id": "g111",
*       "title": "g_title",
*       "summary": "summary",
*       "description": "description",
*       "shared_by": "someon one",
*       "category": { "id" : "category111", ...},
*		"tags":["tags","tags"],
*		"image":"url",
*		"shared_by_image":"url",
*		"targer":0,
*		"duration":"P1D",
*		"source":"Some source",
*		"todos":{"items":[]},
*		"category":{Category Object}
*     }
*   },
*   "code": 200,
*   "message": "Read habit successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The habit was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetHabitDetail",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.userapp.GetHabitDetail",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserAppService) GetHabitDetail(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UserApp.GetHabitDetail API request")
	req_habit := new(behaviour_proto.ReadHabitRequest)
	req_habit.HabitId = req.PathParameter("habit_id")
	req_habit.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetHabitDetail(ctx, req_habit)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.userapp.GetHabitDetail", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read habit successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}
