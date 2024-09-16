// Package untis - Code for getting information from the Untis API
//
//	TODO: Deserialize schools to School struct instead of reading values manually
package untis

import (
	"fmt"
	"guntis/jsonrpc"
	"net/http"
)

// GetAPIString - This retrieves the actual API endpoint that should be used for requests
func (s *School) GetAPIString() string {
	return "https://" + s.Server + "/WebUntis/" + "jsonrpc.do" + "?school=" + s.LoginName
}

func (c *Client) isLoggedIn() error {
	if c.loggedIn {
		return nil
	}

	return fmt.Errorf("This client is logged out and no longer usable")
}

// SearchSchools - Searches for schools matching the search query
func SearchSchools(search string) ([]School, error) {
	client, err := jsonrpc.NewClient("https://mobile.webuntis.com/ms/schoolquery2")

	if err != nil {
		return nil, err
	}

	res, err := client.SendRequest(jsonrpc.NewRequestRaw("searchSchool", fmt.Sprintf("[{\"search\": \"%s\"}]", search)))
	if err != nil {
		return nil, err
	}

	con, conerr := res.Content()
	if conerr != nil {
		return nil, fmt.Errorf("Failed getting schools: %s", conerr.Format())
	}

	schoolsObj, ok := con.(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("Response result has unexpected format")
	}

	schoolsData, ok := schoolsObj["schools"]
	if !ok {
		return nil, fmt.Errorf("Schools field not found in response")
	}

	schoolsSlice, ok := schoolsData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected format for schools field")
	}

	var schools []School
	for _, schoolItem := range schoolsSlice {
		schoolMap, ok := schoolItem.(map[string]interface{})

		if !ok {
			return nil, fmt.Errorf("Unexpected format for school item")
		}

		school := School{
			Server:      schoolMap["server"].(string),
			Address:     schoolMap["address"].(string),
			DisplayName: schoolMap["displayName"].(string),
			LoginName:   schoolMap["loginName"].(string),
			SchoolID:    int64(schoolMap["schoolId"].(float64)),
			ServerURL:   schoolMap["serverUrl"].(string),
		}

		schools = append(schools, school)
	}

	return schools, nil
}

// Login - Log into this school with credentials and get a client
func (s *School) Login(username, password string) (*Client, error) {
	client, err := jsonrpc.NewClient(s.GetAPIString())

	if err != nil {
		return nil, err
	}

	params := jsonrpc.NewParams().Add("user", username).Add("password", password).Add("client", clientName)

	res, err := client.SendRequest(jsonrpc.NewRequest("authenticate", params))

	if err != nil {
		return nil, err
	}

	con, conerr := res.Content()

	if conerr != nil {
		return nil, fmt.Errorf("Failed logging in with credentials: %s", conerr.Format())
	}

	loginObj, ok := con.(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("Response result has unexpected format")
	}

	sessionID, ok := loginObj["sessionId"].(string)

	if !ok {
		return nil, fmt.Errorf("Failed logging in with credentials: API didn't return sessionID")
	}

	client.AddCookie(http.Cookie{
		Name:  "JSESSIONID",
		Value: sessionID,
	})

	return &Client{
		loggedIn:   true,
		SessionID:  sessionID,
		PersonType: int64(loginObj["personType"].(float64)),
		PersonID:   int64(loginObj["personId"].(float64)),
		client:     client,
	}, nil
}

// Logout - Logout of the client sesssion, this makes the client no longer usable
func (c *Client) Logout() error {
	if err := c.isLoggedIn(); err != nil {
		return err
	}

	res, err := c.client.SendRequest(jsonrpc.NewRequest("logout", jsonrpc.NewParams()))

	if err != nil {
		return err
	}

	_, conerr := res.Content()

	if conerr != nil {
		return fmt.Errorf("Failed logging out of session: %s", conerr.Format())
	}

	c.loggedIn = false

	return nil
}

// GetTeachers - Get a list of teachers, may require special client permissions
func (c *Client) GetTeachers() (interface{}, error) {
	if err := c.isLoggedIn(); err != nil {
		return nil, err
	}

	res, err := c.client.SendRequest(jsonrpc.NewRequest("getTeachers", jsonrpc.NewParams()))

	if err != nil {
		return nil, err
	}

	con, conerr := res.Content()

	if conerr != nil {
		return nil, fmt.Errorf("Failed getting teachers: %s", conerr.Format())
	}

	return con, nil
}

// GetSubjects - Get a list of subjects
func (c *Client) GetSubjects() (interface{}, error) {
	if err := c.isLoggedIn(); err != nil {
		return nil, err
	}

	res, err := c.client.SendRequest(jsonrpc.NewRequest("getSubjects", jsonrpc.NewParams()))

	if err != nil {
		return nil, err
	}

	con, conerr := res.Content()

	if conerr != nil {
		return nil, fmt.Errorf("Failed getting subjects: %s", conerr.Format())
	}

	return con, nil
}
