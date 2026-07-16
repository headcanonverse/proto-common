// Package clientevents — единый источник правды для клиентских событий
// discovery-контура (топик client.events). Валидация выполняется на gateway
// (POST /v1/events); consumer в recommendation-service пишет события
// в append-only user_events как есть.
//
// Схема события (JSON в топике):
//
//	{user_id?, work_id?, surface, action, position?, session_id, occurred_at, props?}
package clientevents

import "strings"

// Actions.
const (
	ActionImpression       = "impression"
	ActionCardOpen         = "card_open"
	ActionAddToQueue       = "add_to_queue"
	ActionRemoveFromQueue  = "remove_from_queue"
	ActionNotInterested    = "not_interested"
	ActionShelfScrollDepth = "shelf_scroll_depth"
	ActionAffiliateClick   = "affiliate_click"
	ActionRate             = "rate"
	ActionMarkRead         = "mark_read"
	ActionStatusChange     = "status_change"
	ActionMatchFixed       = "match_fixed"
	ActionSearchQuery      = "search_query"
	ActionSearchClick      = "search_click"
	ActionSearchNoClick    = "search_no_click"
	ActionOnboardingStep   = "onboarding_step"
	ActionReadingStart     = "reading_start"
	ActionProgress         = "progress"
	ActionFinished         = "finished"
	ActionDNF              = "dnf"
	ActionSampleOpen       = "sample_open"
)

// Static surfaces. Полки передаются с типом: "shelf:{type}" — см. IsValidSurface.
const (
	SurfaceSearch     = "search"
	SurfaceWorkCard   = "work_card"
	SurfaceReader     = "reader"
	SurfaceOnboarding = "onboarding"
	SurfaceLibrary    = "library"
	SurfaceImport     = "import"
)

// SurfaceShelfPrefix — префикс динамических полочных поверхностей ("shelf:daily_mix:1").
const SurfaceShelfPrefix = "shelf:"

var validActions = map[string]struct{}{
	ActionImpression: {}, ActionCardOpen: {}, ActionAddToQueue: {},
	ActionRemoveFromQueue: {}, ActionNotInterested: {}, ActionShelfScrollDepth: {},
	ActionAffiliateClick: {}, ActionRate: {}, ActionMarkRead: {},
	ActionStatusChange: {}, ActionMatchFixed: {}, ActionSearchQuery: {},
	ActionSearchClick: {}, ActionSearchNoClick: {}, ActionOnboardingStep: {},
	ActionReadingStart: {}, ActionProgress: {}, ActionFinished: {},
	ActionDNF: {}, ActionSampleOpen: {},
}

var validSurfaces = map[string]struct{}{
	SurfaceSearch: {}, SurfaceWorkCard: {}, SurfaceReader: {},
	SurfaceOnboarding: {}, SurfaceLibrary: {}, SurfaceImport: {},
}

// IsValidAction сообщает, входит ли action в канонический список.
func IsValidAction(action string) bool {
	_, ok := validActions[action]
	return ok
}

// IsValidSurface принимает статические поверхности и любые "shelf:{type}".
func IsValidSurface(surface string) bool {
	if _, ok := validSurfaces[surface]; ok {
		return true
	}
	return strings.HasPrefix(surface, SurfaceShelfPrefix) && len(surface) > len(SurfaceShelfPrefix)
}
