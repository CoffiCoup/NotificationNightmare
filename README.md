# NotificationNightmare
ENGR Foundations 2 Project

In order to test the project in its current form, run in the terminal:
go run ./cmd/server/main.go

make sure you're running this command from the project folder root NotificationNightmare

and in your local browser navigate to http://localhost:8080/view/home

For this product to be effective, it will need to be hosted not on a localhost, and the code in main.go will need to be edited to reflect this. Specifcially refer to:
http.ListenAndServe(":8080", router)

If you want to include netbadge integration, you will need to try to connect the uid found in the SAML response from netbadge and connect it to roleCache and sessionCache. Only use sessionCache if the SAML response is discarded after login is completed, or anytime in the near future after, as it is simply a way of connecting the current sessions to users in the roleCache.

All backend to frontend bridging is handled in main.go and the handlers package in cmd/server/handling. Sorry it's all crammed together there... hopefully it's not too messy to dig through.

Editing UI can be accomplished in the html files in internal/pages. Be mindful editing the profile pages as they use the html/template setup to properly integrate with the golang code in internal/profiles.

## KNOWN BUGS:

There's currently an issue with storing admin roles where the roleList.json makes a bunch of duplicate files and breaks, so that will need to be fixed in order to make the admin page function correctly.

Cannot delete or modify office hours currently; not working. But the UI exists for it. Might be pulling from the wrong file, so look into the filepaths being opened.

Images for profiles are currently being recieved in the edit profiles page, but the images are not properly displaying in the profile grid or individual profile pages.

There's an issue with switching between student and admin logins; can produce 404 and other errors. This is likely a sessionCache or roleCache issue,.


