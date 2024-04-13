// Copyright (c) 2024 Berk Kirtay

package user

// Action types
const (
	API        string = "API"
	LOGIN      string = "LOGIN"
	LOGOUT     string = "LOGOUT"
	JOIN       string = "JOIN"
	MESSAGE    string = "MESSAGE"
	COMMAND    string = "COMMAND"
	AUTH       string = "AUTH"
	ERROR      string = "ERROR"
	SUSPICIOUS string = "SUSPICIOUS"
)

type Action struct {
	Id         string
	UserId     string
	ActionType string
	Detail     string
	Date       string
}

type ActionOption func(Action) Action

func WithActionId(id string) ActionOption {
	return func(action Action) Action {
		action.Id = id
		return action
	}
}

func WithUserId(userId string) ActionOption {
	return func(action Action) Action {
		action.UserId = userId
		return action
	}
}

func WithActionType(actionType string) ActionOption {
	return func(action Action) Action {
		action.ActionType = actionType
		return action
	}
}

func WithDetail(detail string) ActionOption {
	return func(action Action) Action {
		action.Detail = detail
		return action
	}
}

func WithDate(date string) ActionOption {
	return func(action Action) Action {
		action.Date = date
		return action
	}
}

func CreateDefaultAction() Action {
	return Action{}
}

func CreateAction(options ...ActionOption) Action {
	action := CreateDefaultAction()

	for _, option := range options {
		action = option(action)
	}

	return action
}
