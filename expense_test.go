package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestExpense(t *testing.T) {
	cases := []struct {
		name    string
		expense *Expense
		wantErr error
	}{
		{
			name: "empty title",
			expense: &Expense{
				Title:  "",
				Amount: 100,
			},
			wantErr: ErrInvalidExpense,
		},
		{
			name: "negative amount",
			expense: &Expense{
				Title:  "negative amount",
				Amount: -100,
			},
			wantErr: ErrInvalidExpense,
		},
		{
			name: "zero amount",
			expense: &Expense{
				Title:  "zero amount",
				Amount: 0,
			},
			wantErr: ErrInvalidExpense,
		},
		{
			name: "good expense",
			expense: &Expense{
				Title:  "good expense",
				Amount: 100,
			},
			wantErr: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotErr := c.expense.Validate()
			if !errors.Is(gotErr, c.wantErr) {
				t.Errorf("expected error %v, got %v", c.wantErr, gotErr)
			}
		})
	}
}

func TestSanitize(t *testing.T) {
	cases := []struct {
		name    string
		expense *Expense
		want    *Expense
	}{
		{
			name: "space and html unsafe characters",
			expense: &Expense{
				Title:  "  <script>alert('hello')</script>  ",
				Amount: 100,
				Note:   "  <script>alert('hello')</script>  ",
				Tags:   []string{"  <script>alert('hello')</script>  "},
			},
			want: &Expense{
				Title:  "&lt;script&gt;alert(&#39;hello&#39;)&lt;/script&gt;",
				Amount: 100,
				Note:   "&lt;script&gt;alert(&#39;hello&#39;)&lt;/script&gt;",
				Tags:   []string{"&lt;script&gt;alert(&#39;hello&#39;)&lt;/script&gt;"},
			},
		},
		{
			name: "empty space title, note and tags",
			expense: &Expense{
				Title:  " ",
				Amount: 100,
				Note:   " ",
				Tags:   []string{" "},
			},
			want: &Expense{
				Title:  "",
				Amount: 100,
				Note:   "",
				Tags:   []string{},
			},
		},
		{
			name: "good expense",
			expense: &Expense{
				Title:  "good expense",
				Amount: 100,
				Note:   "good note",
				Tags:   []string{"good", "tags"},
			},
			want: &Expense{
				Title:  "good expense",
				Amount: 100,
				Note:   "good note",
				Tags:   []string{"good", "tags"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.expense.Sanitize()
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("expected %v, got %v", c.want, c.expense)
			}
		})
	}
}
