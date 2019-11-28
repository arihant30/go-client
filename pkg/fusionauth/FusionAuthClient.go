/*
* Copyright (c) 2019, FusionAuth, All Rights Reserved
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*   http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
* either express or implied. See the License for the specific
* language governing permissions and limitations under the License.
*/

package client

import (
  "bytes"
  "encoding/base64"
  "encoding/json"
  "errors"
  "fmt"
  "io"
  "net/http"
  "net/http/httputil"
  "net/url"
  "strconv"
  "strings"
)

var ErrRequestUnsuccessful = errors.New("Response from FusionAuth API was not successful")

// URIWithSegment returns a string with a "/" delimiter between the uri and segment
// If segment is not set (""), just the uri is returned
func URIWithSegment(uri, segment string) string {
	if segment == "" {
		return uri
	}
	return uri + "/" + segment
}

// NewRequest creates a new request for the FusionAuth API call
func (c *FusionAuthClient) NewRequest(method, endpoint string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: endpoint}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if c.APIKey != "" {
		// Send the API Key, but only if it is set
		req.Header.Set("Authorization", c.APIKey)
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// Do makes the request to the FusionAuth API endpoint and decodes the response
func (c *FusionAuthClient) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseDump, _ := httputil.DumpResponse(resp, true)
	fmt.Println(string(responseDump))
	err = json.NewDecoder(resp.Body).Decode(v)
	if err == nil && (resp.StatusCode < 200 || resp.StatusCode > 299) {
		// If everything went well but the API responded with something other than HTTP 2xx, we raise an error
		// That way, consumers can check for ErrRequestUnsuccessful
		return resp, ErrRequestUnsuccessful
	}
	return resp, err
}

// FusionAuthClient describes the Go Client for interacting with FusionAuth's RESTful API
type FusionAuthClient struct {
	BaseURL    *url.URL
	APIKey     string
	HTTPClient *http.Client
}

// ActionUser
// Takes an action on a user. The user being actioned is called the "actionee" and the user taking the action is called the
// "actioner". Both user ids are required. You pass the actionee's user id into the method and the actioner's is put into the
// request object.
//   string actioneeUserId The actionee's user id.
//   interface{} request The action request that includes all of the information about the action being taken including
//   the id of the action, any options and the duration (if applicable).
func (c *FusionAuthClient) ActionUser(actioneeUserId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/action"
    method := http.MethodPost
    uri = URIWithSegment(uri, actioneeUserId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// AddUserToFamily
// Adds a user to an existing family. The family id must be specified.
//   string familyId The id of the family.
//   interface{} request The request object that contains all of the information used to determine which user to add to the family.
func (c *FusionAuthClient) AddUserToFamily(familyId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/family"
    method := http.MethodPut
    uri = URIWithSegment(uri, familyId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CancelAction
// Cancels the user action.
//   string actionId The action id of the action to cancel.
//   interface{} request The action request that contains the information about the cancellation.
func (c *FusionAuthClient) CancelAction(actionId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/action"
    method := http.MethodDelete
    uri = URIWithSegment(uri, actionId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// ChangePassword
// Changes a user's password using the change password Id. This usually occurs after an email has been sent to the user
// and they clicked on a link to reset their password.
//   string changePasswordId The change password Id used to find the user. This value is generated by FusionAuth once the change password workflow has been initiated.
//   interface{} request The change password request that contains all of the information used to change the password.
func (c *FusionAuthClient) ChangePassword(changePasswordId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/change-password"
    method := http.MethodPost
    uri = URIWithSegment(uri, changePasswordId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// ChangePasswordByIdentity
// Changes a user's password using their identity (login id and password). Using a loginId instead of the changePasswordId
// bypasses the email verification and allows a password to be changed directly without first calling the #forgotPassword
// method.
//   interface{} request The change password request that contains all of the information used to change the password.
func (c *FusionAuthClient) ChangePasswordByIdentity(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/change-password"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CommentOnUser
// Adds a comment to the user's account.
//   interface{} request The request object that contains all of the information used to create the user comment.
func (c *FusionAuthClient) CommentOnUser(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/comment"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateApplication
// Creates an application. You can optionally specify an Id for the application, if not provided one will be generated.
//   string applicationId (Optional) The Id to use for the application. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the application.
func (c *FusionAuthClient) CreateApplication(applicationId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/application"
    method := http.MethodPost
    uri = URIWithSegment(uri, applicationId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateApplicationRole
// Creates a new role for an application. You must specify the id of the application you are creating the role for.
// You can optionally specify an Id for the role inside the ApplicationRole object itself, if not provided one will be generated.
//   string applicationId The Id of the application to create the role on.
//   string roleId (Optional) The Id of the role. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the application role.
func (c *FusionAuthClient) CreateApplicationRole(applicationId string, roleId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/application"
    method := http.MethodPost
    uri = URIWithSegment(uri, applicationId)
    uri = URIWithSegment(uri, "role")
    uri = URIWithSegment(uri, roleId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateAuditLog
// Creates an audit log with the message and user name (usually an email). Audit logs should be written anytime you
// make changes to the FusionAuth database. When using the FusionAuth App web interface, any changes are automatically
// written to the audit log. However, if you are accessing the API, you must write the audit logs yourself.
//   interface{} request The request object that contains all of the information used to create the audit log entry.
func (c *FusionAuthClient) CreateAuditLog(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/system/audit-log"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateConsent
// Creates a user consent type. You can optionally specify an Id for the consent type, if not provided one will be generated.
//   string consentId (Optional) The Id for the consent. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the consent.
func (c *FusionAuthClient) CreateConsent(consentId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/consent"
    method := http.MethodPost
    uri = URIWithSegment(uri, consentId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateEmailTemplate
// Creates an email template. You can optionally specify an Id for the template, if not provided one will be generated.
//   string emailTemplateId (Optional) The Id for the template. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the email template.
func (c *FusionAuthClient) CreateEmailTemplate(emailTemplateId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/email/template"
    method := http.MethodPost
    uri = URIWithSegment(uri, emailTemplateId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateFamily
// Creates a family with the user id in the request as the owner and sole member of the family. You can optionally specify an id for the
// family, if not provided one will be generated.
//   string familyId (Optional) The id for the family. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the family.
func (c *FusionAuthClient) CreateFamily(familyId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/family"
    method := http.MethodPost
    uri = URIWithSegment(uri, familyId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateGroup
// Creates a group. You can optionally specify an Id for the group, if not provided one will be generated.
//   string groupId (Optional) The Id for the group. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the group.
func (c *FusionAuthClient) CreateGroup(groupId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/group"
    method := http.MethodPost
    uri = URIWithSegment(uri, groupId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateGroupMembers
// Creates a member in a group.
//   interface{} request The request object that contains all of the information used to create the group member(s).
func (c *FusionAuthClient) CreateGroupMembers(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/group/member"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateIdentityProvider
// Creates an identity provider. You can optionally specify an Id for the identity provider, if not provided one will be generated.
//   string identityProviderId (Optional) The Id of the identity provider. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the identity provider.
func (c *FusionAuthClient) CreateIdentityProvider(identityProviderId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/identity-provider"
    method := http.MethodPost
    uri = URIWithSegment(uri, identityProviderId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateLambda
// Creates a Lambda. You can optionally specify an Id for the lambda, if not provided one will be generated.
//   string lambdaId (Optional) The Id for the lambda. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the lambda.
func (c *FusionAuthClient) CreateLambda(lambdaId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/lambda"
    method := http.MethodPost
    uri = URIWithSegment(uri, lambdaId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateTenant
// Creates a tenant. You can optionally specify an Id for the tenant, if not provided one will be generated.
//   string tenantId (Optional) The Id for the tenant. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the tenant.
func (c *FusionAuthClient) CreateTenant(tenantId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/tenant"
    method := http.MethodPost
    uri = URIWithSegment(uri, tenantId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateTheme
// Creates a Theme. You can optionally specify an Id for the theme, if not provided one will be generated.
//   string themeId (Optional) The Id for the theme. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the theme.
func (c *FusionAuthClient) CreateTheme(themeId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/theme"
    method := http.MethodPost
    uri = URIWithSegment(uri, themeId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateUser
// Creates a user. You can optionally specify an Id for the user, if not provided one will be generated.
//   string userId (Optional) The Id for the user. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the user.
func (c *FusionAuthClient) CreateUser(userId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user"
    method := http.MethodPost
    uri = URIWithSegment(uri, userId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateUserAction
// Creates a user action. This action cannot be taken on a user until this call successfully returns. Anytime after
// that the user action can be applied to any user.
//   string userActionId (Optional) The Id for the user action. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the user action.
func (c *FusionAuthClient) CreateUserAction(userActionId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user-action"
    method := http.MethodPost
    uri = URIWithSegment(uri, userActionId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateUserActionReason
// Creates a user reason. This user action reason cannot be used when actioning a user until this call completes
// successfully. Anytime after that the user action reason can be used.
//   string userActionReasonId (Optional) The Id for the user action reason. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the user action reason.
func (c *FusionAuthClient) CreateUserActionReason(userActionReasonId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user-action-reason"
    method := http.MethodPost
    uri = URIWithSegment(uri, userActionReasonId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateUserConsent
// Creates a single User consent.
//   string userConsentId (Optional) The Id for the User consent. If not provided a secure random UUID will be generated.
//   interface{} request The request that contains the user consent information.
func (c *FusionAuthClient) CreateUserConsent(userConsentId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/consent"
    method := http.MethodPost
    uri = URIWithSegment(uri, userConsentId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// CreateWebhook
// Creates a webhook. You can optionally specify an Id for the webhook, if not provided one will be generated.
//   string webhookId (Optional) The Id for the webhook. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the webhook.
func (c *FusionAuthClient) CreateWebhook(webhookId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/webhook"
    method := http.MethodPost
    uri = URIWithSegment(uri, webhookId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeactivateApplication
// Deactivates the application with the given Id.
//   string applicationId The Id of the application to deactivate.
func (c *FusionAuthClient) DeactivateApplication(applicationId string) (interface{}, error) {
    var body interface{}
    uri := "/api/application"
    method := http.MethodDelete
    uri = URIWithSegment(uri, applicationId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeactivateUser
// Deactivates the user with the given Id.
//   string userId The Id of the user to deactivate.
func (c *FusionAuthClient) DeactivateUser(userId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user"
    method := http.MethodDelete
    uri = URIWithSegment(uri, userId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeactivateUserAction
// Deactivates the user action with the given Id.
//   string userActionId The Id of the user action to deactivate.
func (c *FusionAuthClient) DeactivateUserAction(userActionId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user-action"
    method := http.MethodDelete
    uri = URIWithSegment(uri, userActionId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeactivateUsers
// Deactivates the users with the given ids.
//   []string userIds The ids of the users to deactivate.
func (c *FusionAuthClient) DeactivateUsers(userIds []string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/bulk"
    method := http.MethodDelete
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    for _, userId := range userIds {
 		  q.Add("userId", userId)
 	  }
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteApplication
// Hard deletes an application. This is a dangerous operation and should not be used in most circumstances. This will
// delete the application, any registrations for that application, metrics and reports for the application, all the
// roles for the application, and any other data associated with the application. This operation could take a very
// long time, depending on the amount of data in your database.
//   string applicationId The Id of the application to delete.
func (c *FusionAuthClient) DeleteApplication(applicationId string) (interface{}, error) {
    var body interface{}
    uri := "/api/application"
    method := http.MethodDelete
    uri = URIWithSegment(uri, applicationId)
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("hardDelete", strconv.FormatBool(true))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteApplicationRole
// Hard deletes an application role. This is a dangerous operation and should not be used in most circumstances. This
// permanently removes the given role from all users that had it.
//   string applicationId The Id of the application to deactivate.
//   string roleId The Id of the role to delete.
func (c *FusionAuthClient) DeleteApplicationRole(applicationId string, roleId string) (interface{}, error) {
    var body interface{}
    uri := "/api/application"
    method := http.MethodDelete
    uri = URIWithSegment(uri, applicationId)
    uri = URIWithSegment(uri, "role")
    uri = URIWithSegment(uri, roleId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteConsent
// Deletes the consent for the given Id.
//   string consentId The Id of the consent to delete.
func (c *FusionAuthClient) DeleteConsent(consentId string) (interface{}, error) {
    var body interface{}
    uri := "/api/consent"
    method := http.MethodDelete
    uri = URIWithSegment(uri, consentId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteEmailTemplate
// Deletes the email template for the given Id.
//   string emailTemplateId The Id of the email template to delete.
func (c *FusionAuthClient) DeleteEmailTemplate(emailTemplateId string) (interface{}, error) {
    var body interface{}
    uri := "/api/email/template"
    method := http.MethodDelete
    uri = URIWithSegment(uri, emailTemplateId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteGroup
// Deletes the group for the given Id.
//   string groupId The Id of the group to delete.
func (c *FusionAuthClient) DeleteGroup(groupId string) (interface{}, error) {
    var body interface{}
    uri := "/api/group"
    method := http.MethodDelete
    uri = URIWithSegment(uri, groupId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteGroupMembers
// Removes users as members of a group.
//   interface{} request The member request that contains all of the information used to remove members to the group.
func (c *FusionAuthClient) DeleteGroupMembers(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/group/member"
    method := http.MethodDelete
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteIdentityProvider
// Deletes the identity provider for the given Id.
//   string identityProviderId The Id of the identity provider to delete.
func (c *FusionAuthClient) DeleteIdentityProvider(identityProviderId string) (interface{}, error) {
    var body interface{}
    uri := "/api/identity-provider"
    method := http.MethodDelete
    uri = URIWithSegment(uri, identityProviderId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteKey
// Deletes the key for the given Id.
//   string keyOd The Id of the key to delete.
func (c *FusionAuthClient) DeleteKey(keyOd string) (interface{}, error) {
    var body interface{}
    uri := "/api/key"
    method := http.MethodDelete
    uri = URIWithSegment(uri, keyOd)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteLambda
// Deletes the lambda for the given Id.
//   string lambdaId The Id of the lambda to delete.
func (c *FusionAuthClient) DeleteLambda(lambdaId string) (interface{}, error) {
    var body interface{}
    uri := "/api/lambda"
    method := http.MethodDelete
    uri = URIWithSegment(uri, lambdaId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteRegistration
// Deletes the user registration for the given user and application.
//   string userId The Id of the user whose registration is being deleted.
//   string applicationId The Id of the application to remove the registration for.
func (c *FusionAuthClient) DeleteRegistration(userId string, applicationId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/registration"
    method := http.MethodDelete
    uri = URIWithSegment(uri, userId)
    uri = URIWithSegment(uri, applicationId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteTenant
// Deletes the tenant for the given Id.
//   string tenantId The Id of the tenant to delete.
func (c *FusionAuthClient) DeleteTenant(tenantId string) (interface{}, error) {
    var body interface{}
    uri := "/api/tenant"
    method := http.MethodDelete
    uri = URIWithSegment(uri, tenantId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteTheme
// Deletes the theme for the given Id.
//   string themeId The Id of the theme to delete.
func (c *FusionAuthClient) DeleteTheme(themeId string) (interface{}, error) {
    var body interface{}
    uri := "/api/theme"
    method := http.MethodDelete
    uri = URIWithSegment(uri, themeId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteUser
// Deletes the user for the given Id. This permanently deletes all information, metrics, reports and data associated
// with the user.
//   string userId The Id of the user to delete.
func (c *FusionAuthClient) DeleteUser(userId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user"
    method := http.MethodDelete
    uri = URIWithSegment(uri, userId)
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("hardDelete", strconv.FormatBool(true))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteUserAction
// Deletes the user action for the given Id. This permanently deletes the user action and also any history and logs of
// the action being applied to any users.
//   string userActionId The Id of the user action to delete.
func (c *FusionAuthClient) DeleteUserAction(userActionId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user-action"
    method := http.MethodDelete
    uri = URIWithSegment(uri, userActionId)
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("hardDelete", strconv.FormatBool(true))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteUserActionReason
// Deletes the user action reason for the given Id.
//   string userActionReasonId The Id of the user action reason to delete.
func (c *FusionAuthClient) DeleteUserActionReason(userActionReasonId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user-action-reason"
    method := http.MethodDelete
    uri = URIWithSegment(uri, userActionReasonId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteUsers
// Deletes the users with the given ids.
//   interface{} request The ids of the users to delete.
func (c *FusionAuthClient) DeleteUsers(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/bulk"
    method := http.MethodDelete
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DeleteWebhook
// Deletes the webhook for the given Id.
//   string webhookId The Id of the webhook to delete.
func (c *FusionAuthClient) DeleteWebhook(webhookId string) (interface{}, error) {
    var body interface{}
    uri := "/api/webhook"
    method := http.MethodDelete
    uri = URIWithSegment(uri, webhookId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// DisableTwoFactor
// Disable Two Factor authentication for a user.
//   string userId The Id of the User for which you're disabling Two Factor authentication.
//   string code The Two Factor code used verify the the caller knows the Two Factor secret.
func (c *FusionAuthClient) DisableTwoFactor(userId string, code string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/two-factor"
    method := http.MethodDelete
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("userId", string(userId))
    q.Add("code", string(code))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// EnableTwoFactor
// Enable Two Factor authentication for a user.
//   string userId The Id of the user to enable Two Factor authentication.
//   interface{} request The two factor enable request information.
func (c *FusionAuthClient) EnableTwoFactor(userId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/two-factor"
    method := http.MethodPost
    uri = URIWithSegment(uri, userId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// ExchangeRefreshTokenForJWT
// Exchange a refresh token for a new JWT.
//   interface{} request The refresh request.
func (c *FusionAuthClient) ExchangeRefreshTokenForJWT(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/jwt/refresh"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// ForgotPassword
// Begins the forgot password sequence, which kicks off an email to the user so that they can reset their password.
//   interface{} request The request that contains the information about the user so that they can be emailed.
func (c *FusionAuthClient) ForgotPassword(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/forgot-password"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// GenerateEmailVerificationId
// Generate a new Email Verification Id to be used with the Verify Email API. This API will not attempt to send an
// email to the User. This API may be used to collect the verificationId for use with a third party system.
//   string email The email address of the user that needs a new verification email.
func (c *FusionAuthClient) GenerateEmailVerificationId(email string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/verify-email"
    method := http.MethodPut
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("email", string(email))
    q.Add("sendVerifyEmail", strconv.FormatBool(false))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// GenerateKey
// Generate a new RSA or EC key pair or an HMAC secret.
//   string keyId (Optional) The Id for the key. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the key.
func (c *FusionAuthClient) GenerateKey(keyId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/key/generate"
    method := http.MethodPost
    uri = URIWithSegment(uri, keyId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// GenerateRegistrationVerificationId
// Generate a new Application Registration Verification Id to be used with the Verify Registration API. This API will not attempt to send an
// email to the User. This API may be used to collect the verificationId for use with a third party system.
//   string email The email address of the user that needs a new verification email.
//   string applicationId The Id of the application to be verified.
func (c *FusionAuthClient) GenerateRegistrationVerificationId(email string, applicationId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/verify-registration"
    method := http.MethodPut
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("email", string(email))
    q.Add("sendVerifyPasswordEmail", strconv.FormatBool(false))
    q.Add("applicationId", string(applicationId))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// GenerateTwoFactorSecret
// Generate a Two Factor secret that can be used to enable Two Factor authentication for a User. The response will contain
// both the secret and a Base32 encoded form of the secret which can be shown to a User when using a 2 Step Authentication
// application such as Google Authenticator.
func (c *FusionAuthClient) GenerateTwoFactorSecret() (interface{}, error) {
    var body interface{}
    uri := "/api/two-factor/secret"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// GenerateTwoFactorSecretUsingJWT
// Generate a Two Factor secret that can be used to enable Two Factor authentication for a User. The response will contain
// both the secret and a Base32 encoded form of the secret which can be shown to a User when using a 2 Step Authentication
// application such as Google Authenticator.
//   string encodedJWT The encoded JWT (access token).
func (c *FusionAuthClient) GenerateTwoFactorSecretUsingJWT(encodedJWT string) (interface{}, error) {
    var body interface{}
    uri := "/api/two-factor/secret"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Authorization", "JWT " + encodedJWT)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// IdentityProviderLogin
// Handles login via third-parties including Social login, external OAuth and OpenID Connect, and other
// login systems.
//   interface{} request The third-party login request that contains information from the third-party login
//   providers that FusionAuth uses to reconcile the user's account.
func (c *FusionAuthClient) IdentityProviderLogin(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/identity-provider/login"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// ImportKey
// Import an existing RSA or EC key pair or an HMAC secret.
//   string keyId (Optional) The Id for the key. If not provided a secure random UUID will be generated.
//   interface{} request The request object that contains all of the information used to create the key.
func (c *FusionAuthClient) ImportKey(keyId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/key/import"
    method := http.MethodPost
    uri = URIWithSegment(uri, keyId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// ImportUsers
// Bulk imports multiple users. This does some validation, but then tries to run batch inserts of users. This reduces
// latency when inserting lots of users. Therefore, the error response might contain some information about failures,
// but it will likely be pretty generic.
//   interface{} request The request that contains all of the information about all of the users to import.
func (c *FusionAuthClient) ImportUsers(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/import"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// IssueJWT
// Issue a new access token (JWT) for the requested Application after ensuring the provided JWT is valid. A valid
// access token is properly signed and not expired.
// <p>
// This API may be used in an SSO configuration to issue new tokens for another application after the user has
// obtained a valid token from authentication.
//   string applicationId The Application Id for which you are requesting a new access token be issued.
//   string encodedJWT The encoded JWT (access token).
func (c *FusionAuthClient) IssueJWT(applicationId string, encodedJWT string) (interface{}, error) {
    var body interface{}
    uri := "/api/jwt/issue"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("applicationId", string(applicationId))
    req.Header.Set("Authorization", "JWT " + encodedJWT)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// Login
// Authenticates a user to FusionAuth. 
// 
// This API optionally requires an API key. See <code>Application.loginConfiguration.requireAuthentication</code>.
//   interface{} request The login request that contains the user credentials used to log them in.
func (c *FusionAuthClient) Login(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/login"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// LoginPing
// Sends a ping to FusionAuth indicating that the user was automatically logged into an application. When using
// FusionAuth's SSO or your own, you should call this if the user is already logged in centrally, but accesses an
// application where they no longer have a session. This helps correctly track login counts, times and helps with
// reporting.
//   string userId The Id of the user that was logged in.
//   string applicationId The Id of the application that they logged into.
//   string callerIPAddress (Optional) The IP address of the end-user that is logging in. If a null value is provided
//   the IP address will be that of the client or last proxy that sent the request.
func (c *FusionAuthClient) LoginPing(userId string, applicationId string, callerIPAddress string) (interface{}, error) {
    var body interface{}
    uri := "/api/login"
    method := http.MethodPut
    uri = URIWithSegment(uri, userId)
    uri = URIWithSegment(uri, applicationId)
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("ipAddress", string(callerIPAddress))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// Logout
// The Logout API is intended to be used to remove the refresh token and access token cookies if they exist on the
// client and revoke the refresh token stored. This API does nothing if the request does not contain an access
// token or refresh token cookies.
//   bool global When this value is set to true all of the refresh tokens issued to the owner of the
//   provided token will be revoked.
//   string refreshToken (Optional) The refresh_token as a request parameter instead of coming in via a cookie.
//   If provided this takes precedence over the cookie.
func (c *FusionAuthClient) Logout(global bool, refreshToken string) (interface{}, error) {
    var body interface{}
    uri := "/api/logout"
    method := http.MethodPost
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("global", strconv.FormatBool(global))
    q.Add("refreshToken", string(refreshToken))
    req.Header.Set("Content-Type", "text/plain")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// LookupIdentityProvider
// Retrieves the identity provider for the given domain. A 200 response code indicates the domain is managed
// by a registered identity provider. A 404 indicates the domain is not managed.
//   string domain The domain or email address to lookup.
func (c *FusionAuthClient) LookupIdentityProvider(domain string) (interface{}, error) {
    var body interface{}
    uri := "/api/identity-provider/lookup"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("domain", string(domain))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// ModifyAction
// Modifies a temporal user action by changing the expiration of the action and optionally adding a comment to the
// action.
//   string actionId The Id of the action to modify. This is technically the user action log id.
//   interface{} request The request that contains all of the information about the modification.
func (c *FusionAuthClient) ModifyAction(actionId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/action"
    method := http.MethodPut
    uri = URIWithSegment(uri, actionId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// PasswordlessLogin
// Complete a login request using a passwordless code
//   interface{} request The passwordless login request that contains all of the information used to complete login.
func (c *FusionAuthClient) PasswordlessLogin(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/passwordless/login"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// ReactivateApplication
// Reactivates the application with the given Id.
//   string applicationId The Id of the application to reactivate.
func (c *FusionAuthClient) ReactivateApplication(applicationId string) (interface{}, error) {
    var body interface{}
    uri := "/api/application"
    method := http.MethodPut
    uri = URIWithSegment(uri, applicationId)
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("reactivate", strconv.FormatBool(true))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// ReactivateUser
// Reactivates the user with the given Id.
//   string userId The Id of the user to reactivate.
func (c *FusionAuthClient) ReactivateUser(userId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user"
    method := http.MethodPut
    uri = URIWithSegment(uri, userId)
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("reactivate", strconv.FormatBool(true))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// ReactivateUserAction
// Reactivates the user action with the given Id.
//   string userActionId The Id of the user action to reactivate.
func (c *FusionAuthClient) ReactivateUserAction(userActionId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user-action"
    method := http.MethodPut
    uri = URIWithSegment(uri, userActionId)
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("reactivate", strconv.FormatBool(true))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// ReconcileJWT
// Reconcile a User to FusionAuth using JWT issued from another Identity Provider.
//   interface{} request The reconcile request that contains the data to reconcile the User.
func (c *FusionAuthClient) ReconcileJWT(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/jwt/reconcile"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RefreshUserSearchIndex
// Request a refresh of the User search index. This API is not generally necessary and the search index will become consistent in a
// reasonable amount of time. There may be scenarios where you may wish to manually request an index refresh. One example may be 
// if you are using the Search API or Delete Tenant API immediately following a User Create etc, you may wish to request a refresh to
//  ensure the index immediately current before making a query request to the search index.
func (c *FusionAuthClient) RefreshUserSearchIndex() (interface{}, error) {
    var body interface{}
    uri := "/api/user/search"
    method := http.MethodPut
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// Register
// Registers a user for an application. If you provide the User and the UserRegistration object on this request, it
// will create the user as well as register them for the application. This is called a Full Registration. However, if
// you only provide the UserRegistration object, then the user must already exist and they will be registered for the
// application. The user id can also be provided and it will either be used to look up an existing user or it will be
// used for the newly created User.
//   string userId (Optional) The Id of the user being registered for the application and optionally created.
//   interface{} request The request that optionally contains the User and must contain the UserRegistration.
func (c *FusionAuthClient) Register(userId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/registration"
    method := http.MethodPost
    uri = URIWithSegment(uri, userId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RemoveUserFromFamily
// Removes a user from the family with the given id.
//   string familyId The id of the family to remove the user from.
//   string userId The id of the user to remove from the family.
func (c *FusionAuthClient) RemoveUserFromFamily(familyId string, userId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/family"
    method := http.MethodDelete
    uri = URIWithSegment(uri, familyId)
    uri = URIWithSegment(uri, userId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// ResendEmailVerification
// Re-sends the verification email to the user.
//   string email The email address of the user that needs a new verification email.
func (c *FusionAuthClient) ResendEmailVerification(email string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/verify-email"
    method := http.MethodPut
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("email", string(email))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// ResendRegistrationVerification
// Re-sends the application registration verification email to the user.
//   string email The email address of the user that needs a new verification email.
//   string applicationId The Id of the application to be verified.
func (c *FusionAuthClient) ResendRegistrationVerification(email string, applicationId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/verify-registration"
    method := http.MethodPut
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("email", string(email))
    q.Add("applicationId", string(applicationId))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveAction
// Retrieves a single action log (the log of a user action that was taken on a user previously) for the given Id.
//   string actionId The Id of the action to retrieve.
func (c *FusionAuthClient) RetrieveAction(actionId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/action"
    method := http.MethodGet
    uri = URIWithSegment(uri, actionId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveActions
// Retrieves all of the actions for the user with the given Id. This will return all time based actions that are active,
// and inactive as well as non-time based actions.
//   string userId The Id of the user to fetch the actions for.
func (c *FusionAuthClient) RetrieveActions(userId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/action"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("userId", string(userId))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveActionsPreventingLogin
// Retrieves all of the actions for the user with the given Id that are currently preventing the User from logging in.
//   string userId The Id of the user to fetch the actions for.
func (c *FusionAuthClient) RetrieveActionsPreventingLogin(userId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/action"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("userId", string(userId))
    q.Add("preventingLogin", strconv.FormatBool(true))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveActiveActions
// Retrieves all of the actions for the user with the given Id that are currently active.
// An active action means one that is time based and has not been canceled, and has not ended.
//   string userId The Id of the user to fetch the actions for.
func (c *FusionAuthClient) RetrieveActiveActions(userId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/action"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("userId", string(userId))
    q.Add("active", strconv.FormatBool(true))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveApplication
// Retrieves the application for the given id or all of the applications if the id is null.
//   string applicationId (Optional) The application id.
func (c *FusionAuthClient) RetrieveApplication(applicationId string) (interface{}, error) {
    var body interface{}
    uri := "/api/application"
    method := http.MethodGet
    uri = URIWithSegment(uri, applicationId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveApplications
// Retrieves all of the applications.
func (c *FusionAuthClient) RetrieveApplications() (interface{}, error) {
    var body interface{}
    uri := "/api/application"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveAuditLog
// Retrieves a single audit log for the given Id.
//   int auditLogId The Id of the audit log to retrieve.
func (c *FusionAuthClient) RetrieveAuditLog(auditLogId int) (interface{}, error) {
    var body interface{}
    uri := "/api/system/audit-log"
    method := http.MethodGet
    uri = URIWithSegment(uri, string(auditLogId))
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveConsent
// Retrieves the Consent for the given Id.
//   string consentId The Id of the consent.
func (c *FusionAuthClient) RetrieveConsent(consentId string) (interface{}, error) {
    var body interface{}
    uri := "/api/consent"
    method := http.MethodGet
    uri = URIWithSegment(uri, consentId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveConsents
// Retrieves all of the consent.
func (c *FusionAuthClient) RetrieveConsents() (interface{}, error) {
    var body interface{}
    uri := "/api/consent"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveDailyActiveReport
// Retrieves the daily active user report between the two instants. If you specify an application id, it will only
// return the daily active counts for that application.
//   string applicationId (Optional) The application id.
//   interface{} start The start instant as UTC milliseconds since Epoch.
//   interface{} end The end instant as UTC milliseconds since Epoch.
func (c *FusionAuthClient) RetrieveDailyActiveReport(applicationId string, start interface{}, end interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/report/daily-active-user"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("applicationId", string(applicationId))
    q.Add("start", start.(string))
    q.Add("end", end.(string))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveEmailTemplate
// Retrieves the email template for the given Id. If you don't specify the id, this will return all of the email templates.
//   string emailTemplateId (Optional) The Id of the email template.
func (c *FusionAuthClient) RetrieveEmailTemplate(emailTemplateId string) (interface{}, error) {
    var body interface{}
    uri := "/api/email/template"
    method := http.MethodGet
    uri = URIWithSegment(uri, emailTemplateId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveEmailTemplatePreview
// Creates a preview of the email template provided in the request. This allows you to preview an email template that
// hasn't been saved to the database yet. The entire email template does not need to be provided on the request. This
// will create the preview based on whatever is given.
//   interface{} request The request that contains the email template and optionally a locale to render it in.
func (c *FusionAuthClient) RetrieveEmailTemplatePreview(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/email/template/preview"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveEmailTemplates
// Retrieves all of the email templates.
func (c *FusionAuthClient) RetrieveEmailTemplates() (interface{}, error) {
    var body interface{}
    uri := "/api/email/template"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveEventLog
// Retrieves a single event log for the given Id.
//   int eventLogId The Id of the event log to retrieve.
func (c *FusionAuthClient) RetrieveEventLog(eventLogId int) (interface{}, error) {
    var body interface{}
    uri := "/api/system/event-log"
    method := http.MethodGet
    uri = URIWithSegment(uri, string(eventLogId))
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveFamilies
// Retrieves all of the families that a user belongs to.
//   string userId The User's id
func (c *FusionAuthClient) RetrieveFamilies(userId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/family"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("userId", string(userId))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveFamilyMembersByFamilyId
// Retrieves all of the members of a family by the unique Family Id.
//   string familyId The unique Id of the Family.
func (c *FusionAuthClient) RetrieveFamilyMembersByFamilyId(familyId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/family"
    method := http.MethodGet
    uri = URIWithSegment(uri, familyId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveGroup
// Retrieves the group for the given Id.
//   string groupId The Id of the group.
func (c *FusionAuthClient) RetrieveGroup(groupId string) (interface{}, error) {
    var body interface{}
    uri := "/api/group"
    method := http.MethodGet
    uri = URIWithSegment(uri, groupId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveGroups
// Retrieves all of the groups.
func (c *FusionAuthClient) RetrieveGroups() (interface{}, error) {
    var body interface{}
    uri := "/api/group"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveIdentityProvider
// Retrieves the identity provider for the given id or all of the identity providers if the id is null.
//   string identityProviderId (Optional) The identity provider id.
func (c *FusionAuthClient) RetrieveIdentityProvider(identityProviderId string) (interface{}, error) {
    var body interface{}
    uri := "/api/identity-provider"
    method := http.MethodGet
    uri = URIWithSegment(uri, identityProviderId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveIdentityProviders
// Retrieves all of the identity providers.
func (c *FusionAuthClient) RetrieveIdentityProviders() (interface{}, error) {
    var body interface{}
    uri := "/api/identity-provider"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveInactiveActions
// Retrieves all of the actions for the user with the given Id that are currently inactive.
// An inactive action means one that is time based and has been canceled or has expired, or is not time based.
//   string userId The Id of the user to fetch the actions for.
func (c *FusionAuthClient) RetrieveInactiveActions(userId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/action"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("userId", string(userId))
    q.Add("active", strconv.FormatBool(false))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveInactiveApplications
// Retrieves all of the applications that are currently inactive.
func (c *FusionAuthClient) RetrieveInactiveApplications() (interface{}, error) {
    var body interface{}
    uri := "/api/application"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("inactive", strconv.FormatBool(true))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveInactiveUserActions
// Retrieves all of the user actions that are currently inactive.
func (c *FusionAuthClient) RetrieveInactiveUserActions() (interface{}, error) {
    var body interface{}
    uri := "/api/user-action"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("inactive", strconv.FormatBool(true))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveIntegration
// Retrieves the available integrations.
func (c *FusionAuthClient) RetrieveIntegration() (interface{}, error) {
    var body interface{}
    uri := "/api/integration"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveJWTPublicKey
// Retrieves the Public Key configured for verifying JSON Web Tokens (JWT) by the key Id (kid).
//   string keyId The Id of the public key (kid).
func (c *FusionAuthClient) RetrieveJWTPublicKey(keyId string) (interface{}, error) {
    var body interface{}
    uri := "/api/jwt/public-key"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("kid", string(keyId))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveJWTPublicKeyByApplicationId
// Retrieves the Public Key configured for verifying the JSON Web Tokens (JWT) issued by the Login API by the Application Id.
//   string applicationId The Id of the Application for which this key is used.
func (c *FusionAuthClient) RetrieveJWTPublicKeyByApplicationId(applicationId string) (interface{}, error) {
    var body interface{}
    uri := "/api/jwt/public-key"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("applicationId", string(applicationId))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveJWTPublicKeys
// Retrieves all Public Keys configured for verifying JSON Web Tokens (JWT).
func (c *FusionAuthClient) RetrieveJWTPublicKeys() (interface{}, error) {
    var body interface{}
    uri := "/api/jwt/public-key"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveKey
// Retrieves the key for the given Id.
//   string keyId The Id of the key.
func (c *FusionAuthClient) RetrieveKey(keyId string) (interface{}, error) {
    var body interface{}
    uri := "/api/key"
    method := http.MethodGet
    uri = URIWithSegment(uri, keyId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveKeys
// Retrieves all of the keys.
func (c *FusionAuthClient) RetrieveKeys() (interface{}, error) {
    var body interface{}
    uri := "/api/key"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveLambda
// Retrieves the lambda for the given Id.
//   string lambdaId The Id of the lambda.
func (c *FusionAuthClient) RetrieveLambda(lambdaId string) (interface{}, error) {
    var body interface{}
    uri := "/api/lambda"
    method := http.MethodGet
    uri = URIWithSegment(uri, lambdaId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveLambdas
// Retrieves all of the lambdas.
func (c *FusionAuthClient) RetrieveLambdas() (interface{}, error) {
    var body interface{}
    uri := "/api/lambda"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveLambdasByType
// Retrieves all of the lambdas for the provided type.
//   interface{} type The type of the lambda to return.
func (c *FusionAuthClient) RetrieveLambdasByType(_type interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/lambda"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("type", _type.(string))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveLoginReport
// Retrieves the login report between the two instants. If you specify an application id, it will only return the
// login counts for that application.
//   string applicationId (Optional) The application id.
//   interface{} start The start instant as UTC milliseconds since Epoch.
//   interface{} end The end instant as UTC milliseconds since Epoch.
func (c *FusionAuthClient) RetrieveLoginReport(applicationId string, start interface{}, end interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/report/login"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("applicationId", string(applicationId))
    q.Add("start", start.(string))
    q.Add("end", end.(string))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveMonthlyActiveReport
// Retrieves the monthly active user report between the two instants. If you specify an application id, it will only
// return the monthly active counts for that application.
//   string applicationId (Optional) The application id.
//   interface{} start The start instant as UTC milliseconds since Epoch.
//   interface{} end The end instant as UTC milliseconds since Epoch.
func (c *FusionAuthClient) RetrieveMonthlyActiveReport(applicationId string, start interface{}, end interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/report/monthly-active-user"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("applicationId", string(applicationId))
    q.Add("start", start.(string))
    q.Add("end", end.(string))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveOauthConfiguration
// Retrieves the Oauth2 configuration for the application for the given Application Id.
//   string applicationId The Id of the Application to retrieve OAuth configuration.
func (c *FusionAuthClient) RetrieveOauthConfiguration(applicationId string) (interface{}, error) {
    var body interface{}
    uri := "/api/application"
    method := http.MethodGet
    uri = URIWithSegment(uri, applicationId)
    uri = URIWithSegment(uri, "oauth-configuration")
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrievePasswordValidationRules
// Retrieves the password validation rules for a specific tenant. This method requires a tenantId to be provided 
// through the use of a Tenant scoped API key or an HTTP header X-FusionAuth-TenantId to specify the Tenant Id.
// 
// This API does not require an API key.
func (c *FusionAuthClient) RetrievePasswordValidationRules() (interface{}, error) {
    var body interface{}
    uri := "/api/tenant/password-validation-rules"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrievePasswordValidationRulesWithTenantId
// Retrieves the password validation rules for a specific tenant.
// 
// This API does not require an API key.
//   string tenantId The Id of the tenant.
func (c *FusionAuthClient) RetrievePasswordValidationRulesWithTenantId(tenantId string) (interface{}, error) {
    var body interface{}
    uri := "/api/tenant/password-validation-rules"
    method := http.MethodGet
    uri = URIWithSegment(uri, tenantId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrievePendingChildren
// Retrieves all of the children for the given parent email address.
//   string parentEmail The email of the parent.
func (c *FusionAuthClient) RetrievePendingChildren(parentEmail string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/family/pending"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("parentEmail", string(parentEmail))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveRecentLogins
// Retrieves the last number of login records.
//   interface{} offset The initial record. e.g. 0 is the last login, 100 will be the 100th most recent login.
//   int limit (Optional, defaults to 10) The number of records to retrieve.
func (c *FusionAuthClient) RetrieveRecentLogins(offset interface{}, limit int) (interface{}, error) {
    var body interface{}
    uri := "/api/user/recent-login"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("offset", offset.(string))
    q.Add("limit", string(limit))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveRefreshTokens
// Retrieves the refresh tokens that belong to the user with the given Id.
//   string userId The Id of the user.
func (c *FusionAuthClient) RetrieveRefreshTokens(userId string) (interface{}, error) {
    var body interface{}
    uri := "/api/jwt/refresh"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("userId", string(userId))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveRegistration
// Retrieves the user registration for the user with the given id and the given application id.
//   string userId The Id of the user.
//   string applicationId The Id of the application.
func (c *FusionAuthClient) RetrieveRegistration(userId string, applicationId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/registration"
    method := http.MethodGet
    uri = URIWithSegment(uri, userId)
    uri = URIWithSegment(uri, applicationId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveRegistrationReport
// Retrieves the registration report between the two instants. If you specify an application id, it will only return
// the registration counts for that application.
//   string applicationId (Optional) The application id.
//   interface{} start The start instant as UTC milliseconds since Epoch.
//   interface{} end The end instant as UTC milliseconds since Epoch.
func (c *FusionAuthClient) RetrieveRegistrationReport(applicationId string, start interface{}, end interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/report/registration"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("applicationId", string(applicationId))
    q.Add("start", start.(string))
    q.Add("end", end.(string))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveSystemConfiguration
// Retrieves the system configuration.
func (c *FusionAuthClient) RetrieveSystemConfiguration() (interface{}, error) {
    var body interface{}
    uri := "/api/system-configuration"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveTenant
// Retrieves the tenant for the given Id.
//   string tenantId The Id of the tenant.
func (c *FusionAuthClient) RetrieveTenant(tenantId string) (interface{}, error) {
    var body interface{}
    uri := "/api/tenant"
    method := http.MethodGet
    uri = URIWithSegment(uri, tenantId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveTenants
// Retrieves all of the tenants.
func (c *FusionAuthClient) RetrieveTenants() (interface{}, error) {
    var body interface{}
    uri := "/api/tenant"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveTheme
// Retrieves the theme for the given Id.
//   string themeId The Id of the theme.
func (c *FusionAuthClient) RetrieveTheme(themeId string) (interface{}, error) {
    var body interface{}
    uri := "/api/theme"
    method := http.MethodGet
    uri = URIWithSegment(uri, themeId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveThemes
// Retrieves all of the themes.
func (c *FusionAuthClient) RetrieveThemes() (interface{}, error) {
    var body interface{}
    uri := "/api/theme"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveTotalReport
// Retrieves the totals report. This contains all of the total counts for each application and the global registration
// count.
func (c *FusionAuthClient) RetrieveTotalReport() (interface{}, error) {
    var body interface{}
    uri := "/api/report/totals"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUser
// Retrieves the user for the given Id.
//   string userId The Id of the user.
func (c *FusionAuthClient) RetrieveUser(userId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user"
    method := http.MethodGet
    uri = URIWithSegment(uri, userId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserAction
// Retrieves the user action for the given Id. If you pass in null for the id, this will return all of the user
// actions.
//   string userActionId (Optional) The Id of the user action.
func (c *FusionAuthClient) RetrieveUserAction(userActionId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user-action"
    method := http.MethodGet
    uri = URIWithSegment(uri, userActionId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserActionReason
// Retrieves the user action reason for the given Id. If you pass in null for the id, this will return all of the user
// action reasons.
//   string userActionReasonId (Optional) The Id of the user action reason.
func (c *FusionAuthClient) RetrieveUserActionReason(userActionReasonId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user-action-reason"
    method := http.MethodGet
    uri = URIWithSegment(uri, userActionReasonId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserActionReasons
// Retrieves all the user action reasons.
func (c *FusionAuthClient) RetrieveUserActionReasons() (interface{}, error) {
    var body interface{}
    uri := "/api/user-action-reason"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserActions
// Retrieves all of the user actions.
func (c *FusionAuthClient) RetrieveUserActions() (interface{}, error) {
    var body interface{}
    uri := "/api/user-action"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserByChangePasswordId
// Retrieves the user by a change password Id. The intended use of this API is to retrieve a user after the forgot
// password workflow has been initiated and you may not know the user's email or username.
//   string changePasswordId The unique change password Id that was sent via email or returned by the Forgot Password API.
func (c *FusionAuthClient) RetrieveUserByChangePasswordId(changePasswordId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("changePasswordId", string(changePasswordId))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserByEmail
// Retrieves the user for the given email.
//   string email The email of the user.
func (c *FusionAuthClient) RetrieveUserByEmail(email string) (interface{}, error) {
    var body interface{}
    uri := "/api/user"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("email", string(email))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserByLoginId
// Retrieves the user for the loginId. The loginId can be either the username or the email.
//   string loginId The email or username of the user.
func (c *FusionAuthClient) RetrieveUserByLoginId(loginId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("loginId", string(loginId))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserByUsername
// Retrieves the user for the given username.
//   string username The username of the user.
func (c *FusionAuthClient) RetrieveUserByUsername(username string) (interface{}, error) {
    var body interface{}
    uri := "/api/user"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("username", string(username))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserByVerificationId
// Retrieves the user by a verificationId. The intended use of this API is to retrieve a user after the forgot
// password workflow has been initiated and you may not know the user's email or username.
//   string verificationId The unique verification Id that has been set on the user object.
func (c *FusionAuthClient) RetrieveUserByVerificationId(verificationId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("verificationId", string(verificationId))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserComments
// Retrieves all of the comments for the user with the given Id.
//   string userId The Id of the user.
func (c *FusionAuthClient) RetrieveUserComments(userId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/comment"
    method := http.MethodGet
    uri = URIWithSegment(uri, userId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserConsent
// Retrieve a single User consent by Id.
//   string userConsentId The User consent Id
func (c *FusionAuthClient) RetrieveUserConsent(userConsentId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/consent"
    method := http.MethodGet
    uri = URIWithSegment(uri, userConsentId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserConsents
// Retrieves all of the consents for a User.
//   string userId The User's Id
func (c *FusionAuthClient) RetrieveUserConsents(userId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/consent"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("userId", string(userId))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserLoginReport
// Retrieves the login report between the two instants for a particular user by Id. If you specify an application id, it will only return the
// login counts for that application.
//   string applicationId (Optional) The application id.
//   string userId The userId id.
//   interface{} start The start instant as UTC milliseconds since Epoch.
//   interface{} end The end instant as UTC milliseconds since Epoch.
func (c *FusionAuthClient) RetrieveUserLoginReport(applicationId string, userId string, start interface{}, end interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/report/login"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("applicationId", string(applicationId))
    q.Add("userId", string(userId))
    q.Add("start", start.(string))
    q.Add("end", end.(string))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserLoginReportByLoginId
// Retrieves the login report between the two instants for a particular user by login Id. If you specify an application id, it will only return the
// login counts for that application.
//   string applicationId (Optional) The application id.
//   string loginId The userId id.
//   interface{} start The start instant as UTC milliseconds since Epoch.
//   interface{} end The end instant as UTC milliseconds since Epoch.
func (c *FusionAuthClient) RetrieveUserLoginReportByLoginId(applicationId string, loginId string, start interface{}, end interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/report/login"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("applicationId", string(applicationId))
    q.Add("loginId", string(loginId))
    q.Add("start", start.(string))
    q.Add("end", end.(string))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserRecentLogins
// Retrieves the last number of login records for a user.
//   string userId The Id of the user.
//   interface{} offset The initial record. e.g. 0 is the last login, 100 will be the 100th most recent login.
//   int limit (Optional, defaults to 10) The number of records to retrieve.
func (c *FusionAuthClient) RetrieveUserRecentLogins(userId string, offset interface{}, limit int) (interface{}, error) {
    var body interface{}
    uri := "/api/user/recent-login"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("userId", string(userId))
    q.Add("offset", offset.(string))
    q.Add("limit", string(limit))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveUserUsingJWT
// Retrieves the user for the given Id. This method does not use an API key, instead it uses a JSON Web Token (JWT) for authentication.
//   string encodedJWT The encoded JWT (access token).
func (c *FusionAuthClient) RetrieveUserUsingJWT(encodedJWT string) (interface{}, error) {
    var body interface{}
    uri := "/api/user"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Authorization", "JWT " + encodedJWT)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveWebhook
// Retrieves the webhook for the given Id. If you pass in null for the id, this will return all the webhooks.
//   string webhookId (Optional) The Id of the webhook.
func (c *FusionAuthClient) RetrieveWebhook(webhookId string) (interface{}, error) {
    var body interface{}
    uri := "/api/webhook"
    method := http.MethodGet
    uri = URIWithSegment(uri, webhookId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RetrieveWebhooks
// Retrieves all the webhooks.
func (c *FusionAuthClient) RetrieveWebhooks() (interface{}, error) {
    var body interface{}
    uri := "/api/webhook"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RevokeRefreshToken
// Revokes a single refresh token, all tokens for a user or all tokens for an application. If you provide a user id
// and an application id, this will delete all the refresh tokens for that user for that application.
//   string token (Optional) The refresh token to delete.
//   string userId (Optional) The user id whose tokens to delete.
//   string applicationId (Optional) The application id of the tokens to delete.
func (c *FusionAuthClient) RevokeRefreshToken(token string, userId string, applicationId string) (interface{}, error) {
    var body interface{}
    uri := "/api/jwt/refresh"
    method := http.MethodDelete
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    q.Add("token", string(token))
    q.Add("userId", string(userId))
    q.Add("applicationId", string(applicationId))
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// RevokeUserConsent
// Revokes a single User consent by Id.
//   string userConsentId The User Consent Id
func (c *FusionAuthClient) RevokeUserConsent(userConsentId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/consent"
    method := http.MethodDelete
    uri = URIWithSegment(uri, userConsentId)
    req, err := c.NewRequest(method, uri, body)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// SearchAuditLogs
// Searches the audit logs with the specified criteria and pagination.
//   interface{} request The search criteria and pagination information.
func (c *FusionAuthClient) SearchAuditLogs(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/system/audit-log/search"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// SearchEventLogs
// Searches the event logs with the specified criteria and pagination.
//   interface{} request The search criteria and pagination information.
func (c *FusionAuthClient) SearchEventLogs(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/system/event-log/search"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// SearchLoginRecords
// Searches the login records with the specified criteria and pagination.
//   interface{} request The search criteria and pagination information.
func (c *FusionAuthClient) SearchLoginRecords(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/system/login-record/search"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// SearchUsers
// Retrieves the users for the given ids. If any id is invalid, it is ignored.
//   []string ids The user ids to search for.
func (c *FusionAuthClient) SearchUsers(ids []string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/search"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    q := req.URL.Query()
    for _, ids := range ids {
 		  q.Add("ids", ids)
 	  }
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// SearchUsersByQueryString
// Retrieves the users for the given search criteria and pagination.
//   interface{} request The search criteria and pagination constraints. Fields used: queryString, numberOfResults, startRow,
//   and sort fields.
func (c *FusionAuthClient) SearchUsersByQueryString(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/search"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// SendEmail
// Send an email using an email template id. You can optionally provide <code>requestData</code> to access key value
// pairs in the email template.
//   string emailTemplateId The id for the template.
//   interface{} request The send email request that contains all of the information used to send the email.
func (c *FusionAuthClient) SendEmail(emailTemplateId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/email/send"
    method := http.MethodPost
    uri = URIWithSegment(uri, emailTemplateId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// SendFamilyRequestEmail
// Sends out an email to a parent that they need to register and create a family or need to log in and add a child to their existing family.
//   interface{} request The request object that contains the parent email.
func (c *FusionAuthClient) SendFamilyRequestEmail(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/family/request"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// SendPasswordlessCode
// Send a passwordless authentication code in an email to complete login.
//   interface{} request The passwordless send request that contains all of the information used to send an email containing a code.
func (c *FusionAuthClient) SendPasswordlessCode(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/passwordless/send"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// SendTwoFactorCode
// Send a Two Factor authentication code to assist in setting up Two Factor authentication or disabling.
//   interface{} request The request object that contains all of the information used to send the code.
func (c *FusionAuthClient) SendTwoFactorCode(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/two-factor/send"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// SendTwoFactorCodeForLogin
// Send a Two Factor authentication code to allow the completion of Two Factor authentication.
//   string twoFactorId The Id returned by the Login API necessary to complete Two Factor authentication.
func (c *FusionAuthClient) SendTwoFactorCodeForLogin(twoFactorId string) (interface{}, error) {
    var body interface{}
    uri := "/api/two-factor/send"
    method := http.MethodPost
    uri = URIWithSegment(uri, twoFactorId)
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "text/plain")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// StartIdentityProviderLogin
// Begins a login request for a 3rd party login that requires user interaction such as HYPR.
//   interface{} request The third-party login request that contains information from the third-party login
//   providers that FusionAuth uses to reconcile the user's account.
func (c *FusionAuthClient) StartIdentityProviderLogin(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/identity-provider/start"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// StartPasswordlessLogin
// Start a passwordless login request by generating a passwordless code. This code can be sent to the User using the Send
// Passwordless Code API or using a mechanism outside of FusionAuth. The passwordless login is completed by using the Passwordless Login API with this code.
//   interface{} request The passwordless start request that contains all of the information used to begin the passwordless login request.
func (c *FusionAuthClient) StartPasswordlessLogin(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/passwordless/start"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// TwoFactorLogin
// Complete login using a 2FA challenge
//   interface{} request The login request that contains the user credentials used to log them in.
func (c *FusionAuthClient) TwoFactorLogin(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/two-factor/login"
    method := http.MethodPost
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateApplication
// Updates the application with the given Id.
//   string applicationId The Id of the application to update.
//   interface{} request The request that contains all of the new application information.
func (c *FusionAuthClient) UpdateApplication(applicationId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/application"
    method := http.MethodPut
    uri = URIWithSegment(uri, applicationId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateApplicationRole
// Updates the application role with the given id for the application.
//   string applicationId The Id of the application that the role belongs to.
//   string roleId The Id of the role to update.
//   interface{} request The request that contains all of the new role information.
func (c *FusionAuthClient) UpdateApplicationRole(applicationId string, roleId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/application"
    method := http.MethodPut
    uri = URIWithSegment(uri, applicationId)
    uri = URIWithSegment(uri, "role")
    uri = URIWithSegment(uri, roleId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateConsent
// Updates the consent with the given Id.
//   string consentId The Id of the consent to update.
//   interface{} request The request that contains all of the new consent information.
func (c *FusionAuthClient) UpdateConsent(consentId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/consent"
    method := http.MethodPut
    uri = URIWithSegment(uri, consentId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateEmailTemplate
// Updates the email template with the given Id.
//   string emailTemplateId The Id of the email template to update.
//   interface{} request The request that contains all of the new email template information.
func (c *FusionAuthClient) UpdateEmailTemplate(emailTemplateId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/email/template"
    method := http.MethodPut
    uri = URIWithSegment(uri, emailTemplateId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateGroup
// Updates the group with the given Id.
//   string groupId The Id of the group to update.
//   interface{} request The request that contains all of the new group information.
func (c *FusionAuthClient) UpdateGroup(groupId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/group"
    method := http.MethodPut
    uri = URIWithSegment(uri, groupId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateIdentityProvider
// Updates the identity provider with the given Id.
//   string identityProviderId The Id of the identity provider to update.
//   interface{} request The request object that contains the updated identity provider.
func (c *FusionAuthClient) UpdateIdentityProvider(identityProviderId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/identity-provider"
    method := http.MethodPut
    uri = URIWithSegment(uri, identityProviderId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateIntegrations
// Updates the available integrations.
//   interface{} request The request that contains all of the new integration information.
func (c *FusionAuthClient) UpdateIntegrations(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/integration"
    method := http.MethodPut
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateKey
// Updates the key with the given Id.
//   string keyId The Id of the key to update.
//   interface{} request The request that contains all of the new key information.
func (c *FusionAuthClient) UpdateKey(keyId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/key"
    method := http.MethodPut
    uri = URIWithSegment(uri, keyId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateLambda
// Updates the lambda with the given Id.
//   string lambdaId The Id of the lambda to update.
//   interface{} request The request that contains all of the new lambda information.
func (c *FusionAuthClient) UpdateLambda(lambdaId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/lambda"
    method := http.MethodPut
    uri = URIWithSegment(uri, lambdaId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateRegistration
// Updates the registration for the user with the given id and the application defined in the request.
//   string userId The Id of the user whose registration is going to be updated.
//   interface{} request The request that contains all of the new registration information.
func (c *FusionAuthClient) UpdateRegistration(userId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/registration"
    method := http.MethodPut
    uri = URIWithSegment(uri, userId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateSystemConfiguration
// Updates the system configuration.
//   interface{} request The request that contains all of the new system configuration information.
func (c *FusionAuthClient) UpdateSystemConfiguration(request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/system-configuration"
    method := http.MethodPut
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateTenant
// Updates the tenant with the given Id.
//   string tenantId The Id of the tenant to update.
//   interface{} request The request that contains all of the new tenant information.
func (c *FusionAuthClient) UpdateTenant(tenantId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/tenant"
    method := http.MethodPut
    uri = URIWithSegment(uri, tenantId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateTheme
// Updates the theme with the given Id.
//   string themeId The Id of the theme to update.
//   interface{} request The request that contains all of the new theme information.
func (c *FusionAuthClient) UpdateTheme(themeId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/theme"
    method := http.MethodPut
    uri = URIWithSegment(uri, themeId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateUser
// Updates the user with the given Id.
//   string userId The Id of the user to update.
//   interface{} request The request that contains all of the new user information.
func (c *FusionAuthClient) UpdateUser(userId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user"
    method := http.MethodPut
    uri = URIWithSegment(uri, userId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateUserAction
// Updates the user action with the given Id.
//   string userActionId The Id of the user action to update.
//   interface{} request The request that contains all of the new user action information.
func (c *FusionAuthClient) UpdateUserAction(userActionId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user-action"
    method := http.MethodPut
    uri = URIWithSegment(uri, userActionId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateUserActionReason
// Updates the user action reason with the given Id.
//   string userActionReasonId The Id of the user action reason to update.
//   interface{} request The request that contains all of the new user action reason information.
func (c *FusionAuthClient) UpdateUserActionReason(userActionReasonId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user-action-reason"
    method := http.MethodPut
    uri = URIWithSegment(uri, userActionReasonId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateUserConsent
// Updates a single User consent by Id.
//   string userConsentId The User Consent Id
//   interface{} request The request that contains the user consent information.
func (c *FusionAuthClient) UpdateUserConsent(userConsentId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/user/consent"
    method := http.MethodPut
    uri = URIWithSegment(uri, userConsentId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// UpdateWebhook
// Updates the webhook with the given Id.
//   string webhookId The Id of the webhook to update.
//   interface{} request The request that contains all of the new webhook information.
func (c *FusionAuthClient) UpdateWebhook(webhookId string, request interface{}) (interface{}, error) {
    var body interface{}
    uri := "/api/webhook"
    method := http.MethodPut
    uri = URIWithSegment(uri, webhookId)
    body = request
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "application/json")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// ValidateJWT
// Validates the provided JWT (encoded JWT string) to ensure the token is valid. A valid access token is properly
// signed and not expired.
// <p>
// This API may be used to verify the JWT as well as decode the encoded JWT into human readable identity claims.
//   string encodedJWT The encoded JWT (access token).
func (c *FusionAuthClient) ValidateJWT(encodedJWT string) (interface{}, error) {
    var body interface{}
    uri := "/api/jwt/validate"
    method := http.MethodGet
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Authorization", "JWT " + encodedJWT)
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// VerifyEmail
// Confirms a email verification. The Id given is usually from an email sent to the user.
//   string verificationId The email verification id sent to the user.
func (c *FusionAuthClient) VerifyEmail(verificationId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/verify-email"
    method := http.MethodPost
    uri = URIWithSegment(uri, verificationId)
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "text/plain")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}

// VerifyRegistration
// Confirms an application registration. The Id given is usually from an email sent to the user.
//   string verificationId The registration verification Id sent to the user.
func (c *FusionAuthClient) VerifyRegistration(verificationId string) (interface{}, error) {
    var body interface{}
    uri := "/api/user/verify-registration"
    method := http.MethodPost
    uri = URIWithSegment(uri, verificationId)
    req, err := c.NewRequest(method, uri, body)
    req.Header.Set("Content-Type", "text/plain")
    var resp interface{}
    _, err = c.Do(req, &resp)
    return resp, err
}



// ExchangeOAuthCodeForAccessToken
// Exchanges an OAuth authorization code for an access token.
//   string code The OAuth authorization code.
//   string clientID The OAuth client_id.
//   string clientSecret (Optional: use "" to disregard this parameter) The OAuth client_secret used for Basic Auth.
//   string redirectURI The OAuth redirect_uri.
func (c *FusionAuthClient) ExchangeOAuthCodeForAccessToken(code string, clientID string, clientSecret string, redirectURI string) (interface{}, error) {
  // URL
  rel := &url.URL{Path: "/oauth2/token"}
  u := c.BaseURL.ResolveReference(rel)
  // Body
  body := url.Values{}
  body.Set("code", code)
  body.Set("grant_type", "authorization_code")
  body.Set("client_id", clientID)
  body.Set("redirect_uri", redirectURI)
  encodedBody := strings.NewReader(body.Encode())
  // Request
  method := http.MethodPost
  req, err := http.NewRequest(method, u.String(), encodedBody)
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  // Basic Auth (optional)
  if clientSecret != "" {
    credentials := clientID + ":" + clientSecret
    encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
    req.Header.Set("Authorization", "Basic " + encoded)
  }
  var resp interface{}
  _, err = c.Do(req, &resp)
  return resp, err
}
