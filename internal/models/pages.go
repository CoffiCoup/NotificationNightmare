package models

//This is just for making a long dictionary of pages for security and http rerouting reasons

//REMINDER FOR SECURITY:
// 0 - admin, 1 - Staff, 2 - TA, 3 - Student, 4 - Guest (everyone can access)

var WEBPAGES = map[string]WebPage{
	"login":             {Security: 4, URL: "internal/pages/login.html"},
	"home":              {Security: 4, URL: "internal/pages/frontPage.html"},
	"about":             {Security: 4, URL: "internal/pages/about.html"},
	"admin":             {Security: 0, URL: "internal/pages/admin.html"},
	"calendar":          {Security: 4, URL: "internal/pages/calendarPage.html"},
	"OHupload":          {Security: 2, URL: "internal/pages/OHupload.html"},
	"tas":               {Security: 2, URL: "internal/pages/TAs.html"},
	"profileupload":     {Security: 2, URL: "internal/pages/ProfileUpload.html"},
	"individualprofile": {Security: 4, URL: "internal/pages/profile.html"},
	"requests":          {Security: 4, URL: "internal/pages/studentreq.html"},
}
