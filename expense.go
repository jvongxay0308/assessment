package main

import (
	"errors"
	"fmt"
	"html"
	"strings"
)

// ErrInvalidExpense is returned when expense is invalid
var ErrInvalidExpense = errors.New("invalid expense")

// Expense represents an expense in tracking system
type Expense struct {
	ID     int64    `json:"id"`
	Amount float64  `json:"amount"`
	Title  string   `json:"title"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

// Validate validates the expense
func (e *Expense) Validate() error {
	if e.Title == "" {
		return fmt.Errorf("%w: title must not be empty", ErrInvalidExpense)
	}
	if e.Amount <= 0 {
		return fmt.Errorf("%w: amount must be greater than zero", ErrInvalidExpense)
	}
	return nil
}

// Sanitize sanitizes the expense and returns a new expense
func (e *Expense) Sanitize() *Expense {
	tags := make([]string, 0, len(e.Tags))
	for _, tag := range e.Tags {
		tag = html.EscapeString(strings.TrimSpace(tag))
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return &Expense{
		ID:     e.ID,
		Amount: e.Amount,
		Title:  html.EscapeString(strings.TrimSpace(e.Title)),
		Note:   html.EscapeString(strings.TrimSpace(e.Note)),
		Tags:   tags,
	}
}
