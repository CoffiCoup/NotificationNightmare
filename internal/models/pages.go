package models

//This is just for making a long dictionary of pages for security and http rerouting reasons

//REMINDER FOR SECURITY:
// 0 - admin, 1 - Staff, 2 - TA, 3 - Student, 4 - Guest (everyone can access)

var WEBPAGES = map[string]WebPage{
	"login":        {Security: 4, URL: "internal/pages/login.html"},
	"home":         {Security: 4, URL: "internal/pages/frontPage.html"},
	"about":        {Security: 4, URL: "internal/pages/calendarPage.html"},
	"admin":        {Security: 0, URL: "internal/pages/admin.html"},
	"calendarPage": {Security: 4, URL: "internal/pages/calendarPage.html"},
	"frontPage":    {Security: 4, URL: "internal/pages/frontPage.html"},
	"OHupload":     {Security: 2, URL: "internal/pages/OHupload.html"},
	"profiles":     {Security: 4, URL: "internal/pages/profiles.html"},
	"TAs":          {Security: 2, URL: "internal/pages/TAs.html"},
}
