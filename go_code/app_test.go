package main

import (
	"reflect"
	"testing"
)

func TestFileSort(t *testing.T) {
	tests := []struct {
		search_filter   string
		column_select   string
		input_data      [][]string
		expected_data   [][]string
		expected_target int
	}{
		// Happy Path Case
		{
			"picture",
			"id",
			[][]string{
				{"org_id", "id", "agent_id"},
				{"12345", "picture", "7777"},
				{"12346", "flipper", "7778"},
				{"12347", "picture", "7779"},
			},
			[][]string{
				{"org_id", "id", "agent_id"},
				{"12345", "picture", "7777"},
				{"12347", "picture", "7779"},
			},
			1,
		},
		// Case where all rows match the search filter
		{
			"picture",
			"id",
			[][]string{
				{"org_id", "id", "agent_id"},
				{"12345", "picture", "7777"},
				{"12346", "picture", "7778"},
				{"12347", "picture", "7779"},
			},
			[][]string{
				{"org_id", "id", "agent_id"},
				{"12345", "picture", "7777"},
				{"12346", "picture", "7778"},
				{"12347", "picture", "7779"},
			},
			1,
		},
		// Case where none of the rows match the search filter
		{
			"picture",
			"id",
			[][]string{
				{"org_id", "id", "agent_id"},
				{"12345", "financial", "7777"},
				{"12346", "flipper", "7778"},
				{"12347", "flipper", "7779"},
			},
			[][]string{
				{"org_id", "id", "agent_id"},
			},
			1,
		},
		// Case where no rows in input
		{
			"picture",
			"id",
			[][]string{},
			nil,
			-1,
		},
		// Case where search column does not exist in input
		{
			"picture",
			"id",
			[][]string{
				{"org_id", "flopper", "agent_id"},
				{"12345", "picture", "7777"},
				{"12346", "flipper", "7778"},
				{"12347", "picture", "7779"},
			},
			[][]string{
				{"org_id", "flopper", "agent_id"},
			},
			-1,
		},
	}

	for _, tt := range tests {
		result_data, result_target := fileSort(tt.search_filter, tt.column_select, tt.input_data)
		if !reflect.DeepEqual(result_data, tt.expected_data) {
			t.Errorf("fileSort(%v, %v, %v) --> = %v; want %v", tt.search_filter, tt.column_select, tt.input_data, result_data, tt.expected_data)
		}

		if result_target != tt.expected_target {
			t.Errorf("fileSort(%v, %v, %v,) --> = %d; want %d", tt.search_filter, tt.column_select, tt.input_data, result_target, tt.expected_target)
		}
	}
}

func TestColumnRemoval(t *testing.T) {
	tests := []struct {
		filtered_data [][]string
		target        int
		expected_data [][]string
	}{
		// Happy Path Case
		{
			[][]string{
				{"org_id", "id", "agent_id"},
				{"12345", "picture", "7777"},
				{"12347", "picture", "7779"},
			},
			1,
			[][]string{
				{"org_id", "agent_id"},
				{"12345", "7777"},
				{"12347", "7779"},
			},
		},
		// Case where all rows match the search filter
		{
			[][]string{
				{"org_id", "id", "agent_id"},
				{"12345", "picture", "7777"},
				{"12346", "picture", "7778"},
				{"12347", "picture", "7779"},
			},
			1,
			[][]string{
				{"org_id", "agent_id"},
				{"12345", "7777"},
				{"12346", "7778"},
				{"12347", "7779"},
			},
		},
		// Case where none of the rows match the search filter
		{
			[][]string{
				{"org_id", "id", "agent_id"},
			},
			1,
			[][]string{
				{"org_id", "agent_id"},
			},
		},
		// Case where no rows in input
		{
			nil,
			-1,
			nil,
		},
		// Case where search column does not exist in input
		{
			[][]string{
				{"org_id", "flopper", "agent_id"},
			},
			-1,
			[][]string{
				{"org_id", "flopper", "agent_id"},
			},
		},
	}

	for _, tt := range tests {
		result_data := columnRemoval(tt.filtered_data, tt.target)
		if !reflect.DeepEqual(result_data, tt.expected_data) {
			t.Errorf("columnRemoval(%v, %d) --> = %v; want %v", tt.filtered_data, tt.target, result_data, tt.expected_data)
		}
	}

}
