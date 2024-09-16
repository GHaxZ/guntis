// Package untis - Code for getting information from the Untis API
//
//	TODO: Finish subject struct
package untis

import (
	"guntis/jsonrpc"
)

const clientName string = "guntis"

// School - Represents a school
type School struct {
	Server      string `json:"server"`
	Address     string `json:"address"`
	DisplayName string `json:"displayName"`
	LoginName   string `json:"loginName"`
	SchoolID    int64  `json:"schoolId"`
	ServerURL   string `json:"serverUrl"`
}

// Client - A client to a user login
type Client struct {
	loggedIn   bool
	SessionID  string
	PersonType int64
	PersonID   int64
	client     *jsonrpc.Client
}

type Subject struct {
	active bool
}
