package models

//This is just for making a long dictionary of pages for security and http rerouting reasons

//REMINDER FOR SECURITY:
// 0 - admin, 1 - Staff, 2 - TA, 3 - Student, 4 - Guest (everyone can access)

var WEBPAGES = map[string]WebPage{
	"login": {Security: 4, URL: "internal/pages/login.html"},
}
