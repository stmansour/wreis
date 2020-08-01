package ws

import (
	"net/http"
	"wreis/session"
	util "wreis/util/lib"
)

// SvcUserProfile handles authentication requests from clients.
//
//  @Title UserProfile
//  @URL /v1/userprofile
//  @Method  POST
//  @Synopsis Get information about the logged in user
//  @Descr Based on the session cookie, this service will return the
//  @Descr information known about the user. As of this writing, that
//  @Descr information includes:  the username, the user's first (or
//  @Descr preferred name), the user's id number, and a url to the
//  @Descr user's image.
//  @Input session.AIRAuthenticateResponse
//  @Response SvcStatus
// wsdoc }
//-----------------------------------------------------------------------------
func SvcUserProfile(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	var funcname = "SvcUserProfile"
	var a session.AIRAuthenticateResponse
	util.Console("Entered: %s\n", funcname)
	a.Status = "success"
	a.UID = d.sess.UID
	a.Name = d.sess.Name
	a.Username = d.sess.Username
	a.ImageURL = d.sess.ImageURL
	a.Expire = util.JSONDateTime(d.sess.Expire)
	a.Token = d.sess.Token

	SvcWriteResponse(&a, w)
}
