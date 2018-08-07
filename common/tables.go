package common

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// all the Db tables.
var (
	DbHealumName   = Name("healum")
	DbHealumDriver = "arangodb"
	ErrNotFound    = errors.New("not found")

	DbAccountTable                   = Name("account")
	DbActivityConfigTable            = Name("activity_config")
	DbGoalTable                      = Name("goal")
	DbChallengeTable                 = Name("challenge")
	DbHabitTable                     = Name("habit")
	DbShareGoalUserEdgeTable         = Name("share_goal_user")            //edge
	DbShareGoalUserGraph             = Name("share_goal_user_graph")      //graph
	DbShareChallengeUserEdgeTable    = Name("share_challenge_user")       //edge
	DbShareChallengeUserGraph        = Name("share_challenge_user_graph") //graph
	DbShareHabitUserEdgeTable        = Name("share_habit_user")           //edge
	DbShareHabitUserGraph            = Name("share_habit_user_graph")     //graph
	DbUserTable                      = Name("user")
	DbPendingTable                   = Name("pending")
	DbSourceTable                    = Name("source")
	DbTaxonomyTable                  = Name("taxonomy")
	DbContentCategoryItemTable       = Name("content_category_item")
	DbContentTable                   = Name("content")
	DbContentRuleTable               = Name("content_rule")
	DbContentTagEdgeTable            = Name("content_tag_edge")         //edge
	DbContentTagGraph                = Name("content_tag_graph")        //graph
	DbShareContentUserEdgeTable      = Name("share_content_user")       //edge
	DbShareContentUserGraph          = Name("share_content_user_graph") //graph
	DbContentCategoryTable           = Name("content_category")
	DbBookmarkEdgeTable              = Name("bookmark_edge")  //edge
	DbBookmarkGraph                  = Name("bookmark_graph") //graph
	DbNoteTable                      = Name("note")
	DbUserAccountEdgeTable           = Name("user_account_edge")     //edge
	DbUserAccountGraph               = Name("user_account_graph")    //graph
	DbEmployeeTable                  = Name("employee")              //edge
	DbEmployeeGraph                  = Name("employee_graph")        //graph
	DbEmployeeModuleEdgeTable        = Name("employee_module_edge")  //edge
	DbEmployeeModuleGraph            = Name("employee_module_graph") //graph
	DbOrgProfileTable                = Name("organisation_profile")
	DbOrgSettingTable                = Name("organisation_setting")
	DbOrgModuleEdgeTable             = Name("organisation_module_edge")  //edge
	DbOrgModuleGraph                 = Name("organisation_module_graph") //graph
	DbModuleTable                    = Name("module")
	DbRoleTable                      = Name("role")
	DbPlanTable                      = Name("plan")
	DbPlanItemTable                  = Name("plan_item")
	DbPlanTodoTable                  = Name("plan_todo")
	DbPlanGoalTable                  = Name("plan_goal")
	DbPlanFilterTable                = Name("plan_filter")
	DbSharePlanUserEdgeTable         = Name("share_plan_user")       //edge
	DbSharePlanUserGraph             = Name("share_plan_user_graph") //graph
	DbSurveyTable                    = Name("survey")
	DbSurveyResponseEdgeTable        = Name("survey_response_edge")  //edge
	DbResponseGraph                  = Name("survey_response_graph") //graph
	DbResponseTable                  = Name("response")
	DbShareSurveyUserEdgeTable       = Name("share_survey_user")       //edge
	DbShareSurveyUserGraph           = Name("share_survey_user_graph") //graph
	DbAppTable                       = Name("app")
	DbPlatformTable                  = Name("platform")
	DbWearableTable                  = Name("wearable")
	DbDeviceTable                    = Name("device")
	DbMarkerTable                    = Name("marker")
	DbBehaviourCategoryTable         = Name("behaviour_category")
	DbSocialTypeTable                = Name("social_type")
	DbNotificationTable              = Name("notification")
	DbTrackerMethodTable             = Name("tracker_method")
	DbBehaviourCategoryAimTable      = Name("behaviour_category_aim")
	DbContentParentCategoryTable     = Name("content_parent_category")
	DbContentTypeTable               = Name("content_type")
	DbContentSourceTypeTable         = Name("content_source_type")
	DbMarkerTrackerEdgeTable         = Name("marker_tracker_method_edge")  //edge
	DbMarkerTrackerGraph             = Name("marker_tracker_method_graph") //graph
	DbModuleTriggerTable             = Name("module_trigger")
	DbTriggerContentTypeTable        = Name("trigger_content_type")
	DbSurveyQuestionTable            = Name("survey_question")
	DbSurveyHashTable                = Name("survey_hash")
	DbSurveyQuestionEdgeTable        = Name("survey_question_edge")  //edge
	DbSurveyQuestionGraph            = Name("survey_question_graph") //graph
	DbTaskTable                      = Name("task")
	DbTeamTable                      = Name("team")
	DbEmployeeProfileTable           = Name("employee_profile")
	DbOrganisationTable              = Name("organisation")
	DbTeamMembershipTable            = Name("team_membership")       //edge
	DbTeamMembershipGraph            = Name("team_membership_graph") //graph
	DbTodoTable                      = Name("todo")
	DbTrackGoalTable                 = Name("track_goal")
	DbTrackGoalEdgeTable             = Name("track_goal_edge")  //edge
	DbTrackGoalGraph                 = Name("track_goal_graph") //graph
	DbTrackChallengeTable            = Name("track_challenge")
	DbTrackChallengeEdgeTable        = Name("track_challenge_edge")  //edge
	DbTrackChallengeGraph            = Name("track_challenge_graph") //graph
	DbTrackHabitTable                = Name("track_habit")
	DbTrackHabitEdgeTable            = Name("track_habit_edge")  //edge
	DbTrackHabitGraph                = Name("track_habit_graph") //graph
	DbTrackContentTable              = Name("track_content")
	DbTrackContentEdgeTable          = Name("track_content_edge")  //edge
	DbTrackContentGraph              = Name("track_content_graph") //graph
	DbTrackMarkerTable               = Name("track_marker")
	DbJoinGoalEdgeTable              = Name("join_goal_edge")       //edge
	DbJoinGoalGraph                  = Name("join_goal_graph")      //graph
	DbJoinChallengeEdgeTable         = Name("join_challenge_edge")  //edge
	DbJoinChallengeGraph             = Name("join_challenge_graph") //graph
	DbJoinHabitEdgeTable             = Name("join_habit_edge")      //edge
	DbJoinHabitGraph                 = Name("join_habit_graph")     //graph
	DbUserOrgEdgeTable               = Name("user_org_edge")        //edge
	DbUserOrgGraph                   = Name("user_org_graph")       //graph
	DbPreferenceTable                = Name("user_preference")
	DbContentRecommendationEdgeTable = Name("content_recommendation_edge")   //edge
	DbContentRecommendationGraph     = Name("content_recommendation_graph")  //graph
	DbContentRatingEdgeTable         = Name("content_rating_edge")           //edge
	DbContentRatingGraph             = Name("content_rating_graph")          //graph
	DbContentDislikeEdgeTable        = Name("content_dislike_edge")          //edge
	DbContentDislikeGraph            = Name("content_dislike_graph")         //graph
	DbContentDislikeSimilarEdgeTable = Name("content_dislike_similar_edge")  //edge
	DbContentDislikeSimilarGraph     = Name("content_dislike_similar_graph") //graph
	DbUserFeedbackTable              = Name("user_feedback")
	DbUserPlanTable                  = Name("user_plan")
	DbUserPlanEdgeTable              = Name("user_plan_edge")  //edge
	DbUserPlanGraph                  = Name("user_plan_graph") //graph
	DbSetbackTable                   = Name("setback")
	DbUserMeasurementEdgeTable       = Name("user_measurement_edge")  //edge
	DbUserMeasurementGraph           = Name("user_measurement_graph") //graph
	DbProductTable                   = Name("product")
	DbTeamProductEdgeTable           = Name("team_product_edge")  //edge
	DbTeamProductGraph               = Name("team_product_graph") //graph
	DbServiceTable                   = Name("service")
	DbTeamServiceEdgeTable           = Name("team_service_edge")  //edge
	DbTeamServiceGraph               = Name("team_service_graph") //graph
	DbBatchTable                     = Name("batch")
	DbUserBatchEdgeTable             = Name("user_batch_edge")  //edge
	DbUserBatchGraph                 = Name("user_batch_graph") //graph
	DbActionTable                    = Name("action")
	DbTriggerEventTable              = Name("trigger_event")
	DbTodoItemTable                  = Name("todo_item")
	DbEmailSubAccountTable           = Name("email_subaccount")
	DbEmailTransmissionTable         = Name("email_transmission")
	DbEmailConfigTable               = Name("email_config")
	DbSmsSubAccountTable             = Name("sms_subaccount")
	DbAuditTable                     = Name("audit_log")

	DbHealum = [][]string{
		// table
		{DbAccountTable},
		{DbActivityConfigTable},
		{DbGoalTable},
		{DbChallengeTable},
		{DbHabitTable},
		{DbUserTable},
		{DbPendingTable},
		{DbSourceTable},
		{DbTaxonomyTable},
		{DbContentCategoryItemTable},
		{DbContentTable},
		{DbContentRuleTable},
		{DbContentCategoryTable},
		{DbNoteTable},
		{DbOrgProfileTable},
		{DbOrgSettingTable},
		{DbModuleTable},
		{DbRoleTable},
		{DbPlanTable},
		{DbPlanItemTable},
		{DbPlanTodoTable},
		{DbPlanGoalTable},
		{DbPlanFilterTable},
		{DbSurveyTable},
		{DbResponseTable},
		{DbAppTable},
		{DbPlatformTable},
		{DbWearableTable},
		{DbDeviceTable},
		{DbMarkerTable},
		{DbBehaviourCategoryTable},
		{DbSocialTypeTable},
		{DbNotificationTable},
		{DbTrackerMethodTable},
		{DbBehaviourCategoryAimTable},
		{DbContentParentCategoryTable},
		{DbContentTypeTable},
		{DbContentSourceTypeTable},
		{DbModuleTriggerTable},
		{DbTriggerContentTypeTable},
		{DbSurveyQuestionTable},
		{DbSurveyHashTable},
		{DbTaskTable},
		{DbTeamTable},
		{DbEmployeeProfileTable},
		{DbOrganisationTable},
		{DbTodoTable},
		{DbTrackGoalTable},
		{DbTrackChallengeTable},
		{DbTrackHabitTable},
		{DbTrackContentTable},
		{DbTrackMarkerTable},
		{DbPreferenceTable},
		{DbUserFeedbackTable},
		{DbUserPlanTable},
		{DbSetbackTable},
		{DbProductTable},
		{DbServiceTable},
		{DbBatchTable},
		{DbActionTable},
		{DbTriggerEventTable},
		{DbTodoItemTable},
		{DbEmailSubAccountTable},
		{DbEmailSubAccountTable},
		{DbEmailConfigTable},
		{DbSmsSubAccountTable},
		{DbAuditTable},
		{},
		// egde & graph
		{DbShareGoalUserEdgeTable, DbShareGoalUserGraph, DbGoalTable, DbUserTable},
		{DbShareChallengeUserEdgeTable, DbShareChallengeUserGraph, DbChallengeTable, DbUserTable},
		{DbShareHabitUserEdgeTable, DbShareHabitUserGraph, DbHabitTable, DbUserTable},
		{DbContentTagEdgeTable, DbContentTagGraph, DbContentTable, DbContentCategoryItemTable},
		{DbShareContentUserEdgeTable, DbShareContentUserGraph, DbContentTable, DbUserTable},
		{DbBookmarkEdgeTable, DbBookmarkGraph, DbUserTable, DbContentTable},
		{DbUserAccountEdgeTable, DbUserAccountGraph, DbUserTable, DbAccountTable},
		{DbEmployeeTable, DbEmployeeGraph, DbUserTable, DbOrganisationTable},
		{DbEmployeeModuleEdgeTable, DbEmployeeModuleGraph, DbEmployeeTable, DbModuleTable},
		{DbOrgModuleEdgeTable, DbOrgModuleGraph, DbOrganisationTable, DbModuleTable},
		{DbSharePlanUserEdgeTable, DbSharePlanUserGraph, DbPlanTable, DbUserTable},
		{DbSurveyResponseEdgeTable, DbResponseGraph, DbSurveyTable, DbResponseTable},
		{DbShareSurveyUserEdgeTable, DbShareSurveyUserGraph, DbSurveyTable, DbUserTable},
		{DbMarkerTrackerEdgeTable, DbMarkerTrackerGraph, DbMarkerTable, DbTrackerMethodTable},
		{DbSurveyQuestionEdgeTable, DbSurveyQuestionGraph, DbSurveyTable, DbSurveyQuestionTable},
		{DbTeamMembershipTable, DbTeamMembershipGraph, DbTeamTable, DbUserTable},
		{DbTrackGoalEdgeTable, DbTrackGoalGraph, DbTrackGoalTable, DbGoalTable},
		{DbTrackChallengeEdgeTable, DbTrackChallengeGraph, DbTrackChallengeTable, DbChallengeTable},
		{DbTrackHabitEdgeTable, DbTrackHabitGraph, DbTrackHabitTable, DbHabitTable},
		{DbTrackContentEdgeTable, DbTrackContentGraph, DbTrackContentTable, DbContentTable},
		{DbJoinGoalEdgeTable, DbJoinGoalGraph, DbUserTable, DbGoalTable},
		{DbJoinChallengeEdgeTable, DbJoinChallengeGraph, DbUserTable, DbChallengeTable},
		{DbJoinHabitEdgeTable, DbJoinHabitGraph, DbUserTable, DbHabitTable},
		{DbUserOrgEdgeTable, DbUserOrgGraph, DbUserTable, DbOrganisationTable},
		{DbContentRecommendationEdgeTable, DbContentRecommendationGraph, DbUserTable, DbContentTable},
		{DbContentRatingEdgeTable, DbContentRatingGraph, DbUserTable, DbContentTable},
		{DbContentDislikeEdgeTable, DbContentDislikeGraph, DbUserTable, DbContentTable},
		{DbContentDislikeSimilarEdgeTable, DbContentDislikeSimilarGraph, DbUserTable, DbContentTable},
		{DbUserPlanEdgeTable, DbUserPlanGraph, DbUserTable, DbPlanTable},
		{DbUserMeasurementEdgeTable, DbUserMeasurementGraph, DbUserTable, DbTrackMarkerTable},
		{DbTeamProductEdgeTable, DbTeamProductGraph, DbTeamTable, DbProductTable},
		{DbTeamServiceEdgeTable, DbTeamServiceGraph, DbTeamTable, DbServiceTable},
		{DbUserBatchEdgeTable, DbUserBatchGraph, DbUserTable, DbBatchTable},
		{},
	}
)

func Name(table string) string {
	err := godotenv.Load(os.Getenv("GOPATH") + "/src/server/.env")
	if err == nil {
		env := os.Getenv("HEALUM_ENV")
		if env == "test" {
			return TestingName(table)
		}
	}
	return table
}
