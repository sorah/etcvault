package engine

import (
	"reflect"
	"testing"
)

func TestTransformEtcdJsonResponse(t *testing.T) {
	tests := []struct {
		Name   string
		Case   []byte
		Expect []byte
	}{
		{
			Name:   "plain (node.value)",
			Case:   []byte(`{"node": {"value": "ETCVAULT::asis:plain::ETCVAULT"}}`),
			Expect: []byte(`{"node":{"value":"plain"}}`),
		},
		{
			Name:   "plain (prevNode.value)",
			Case:   []byte(`{"prevNode": {"value": "ETCVAULT::asis:plain::ETCVAULT"}}`),
			Expect: []byte(`{"prevNode":{"value":"plain"}}`),
		},
		{
			Name:   "both (node.value, prevNode.value)",
			Case:   []byte(`{"node": {"value": "ETCVAULT::asis:plain::ETCVAULT"}, "prevNode": {"value": "ETCVAULT::asis:plain::ETCVAULT"}}`),
			Expect: []byte(`{"node":{"value":"plain"},"prevNode":{"value":"plain"}}`),
		},
		{
			Name:   "inside directory (node.nodes[0].value)",
			Case:   []byte(`{"node": {"nodes": [{"value": "ETCVAULT::asis:plain::ETCVAULT"}]}}`),
			Expect: []byte(`{"node":{"nodes":[{"value":"plain"}]}}`),
		},
		{
			Name:   "inside directory, multiple (node.nodes[0].value, node.nodes[1].value)",
			Case:   []byte(`{"node": {"nodes": [{"value": "ETCVAULT::asis:plain::ETCVAULT"}, {"value": "ETCVAULT::asis:plain::ETCVAULT"}]}}`),
			Expect: []byte(`{"node":{"nodes":[{"value":"plain"},{"value":"plain"}]}}`),
		},
		{
			Name:   "nested, inside directory (node.nodes[0].nodes[0].value)",
			Case:   []byte(`{"node": {"nodes": [{"nodes": [{"value": "ETCVAULT::asis:plain::ETCVAULT"}]}]}}`),
			Expect: []byte(`{"node":{"nodes":[{"nodes":[{"value":"plain"}]}]}}`),
		},
	}

	engine := NewEngine(testKeychain)

	for _, test := range tests {
		transformedJson, err := engine.TransformEtcdJsonResponse(test.Case)

		if err != nil {
			t.Errorf("%s:\n\tunexpected err: %s", test.Name, err.Error())
		}

		if !reflect.DeepEqual(transformedJson, test.Expect) {
			t.Errorf("%s:\n\tunexpected result: %s", test.Name, transformedJson)
		}
	}
}

func TestTransformEtcdJsonResponseFailures(t *testing.T) {
	tests := []struct {
		Name   string
		Case   []byte
		Expect []byte
	}{
		{
			Name:   "plain (node.value)",
			Case:   []byte(`{"node": {"value": "ETCVAULT::plain1::ETCVAULT"}}`),
			Expect: []byte(`{"node":{"_etcvault_error":"couldn't parse","value":"ETCVAULT::plain1::ETCVAULT"}}`),
		},
		{
			Name:   "plain (prevNode.value)",
			Case:   []byte(`{"prevNode": {"value": "ETCVAULT::plain1::ETCVAULT"}}`),
			Expect: []byte(`{"prevNode":{"_etcvault_error":"couldn't parse","value":"ETCVAULT::plain1::ETCVAULT"}}`),
		},
		{
			Name:   "both (node.value, prevNode.value)",
			Case:   []byte(`{"node": {"value": "ETCVAULT::plain1::ETCVAULT"}, "prevNode": {"value": "ETCVAULT::plain1::ETCVAULT"}}`),
			Expect: []byte(`{"node":{"_etcvault_error":"couldn't parse","value":"ETCVAULT::plain1::ETCVAULT"},"prevNode":{"_etcvault_error":"couldn't parse","value":"ETCVAULT::plain1::ETCVAULT"}}`),
		},
		{
			Name:   "inside directory (node.nodes[0].value)",
			Case:   []byte(`{"node": {"nodes": [{"value": "ETCVAULT::plain1::ETCVAULT"}]}}`),
			Expect: []byte(`{"node":{"nodes":[{"_etcvault_error":"couldn't parse","value":"ETCVAULT::plain1::ETCVAULT"}]}}`),
		},
		{
			Name:   "inside directory, multiple (node.nodes[0].value, node.nodes[1].value)",
			Case:   []byte(`{"node": {"nodes": [{"value": "ETCVAULT::plain1::ETCVAULT"}, {"value": "ETCVAULT::plain1::ETCVAULT"}]}}`),
			Expect: []byte(`{"node":{"nodes":[{"_etcvault_error":"couldn't parse","value":"ETCVAULT::plain1::ETCVAULT"},{"_etcvault_error":"couldn't parse","value":"ETCVAULT::plain1::ETCVAULT"}]}}`),
		},
		{
			Name:   "nested, inside directory (node.nodes[0].nodes[0].value)",
			Case:   []byte(`{"node": {"nodes": [{"nodes": [{"value": "ETCVAULT::plain1::ETCVAULT"}]}]}}`),
			Expect: []byte(`{"node":{"nodes":[{"nodes":[{"_etcvault_error":"couldn't parse","value":"ETCVAULT::plain1::ETCVAULT"}]}]}}`),
		},
	}

	engine := NewEngine(testKeychain)

	for _, test := range tests {
		transformedJson, err := engine.TransformEtcdJsonResponse(test.Case)

		if err != nil {
			t.Errorf("%s:\n\tunexpected err: %s", test.Name, err.Error())
		}

		if !reflect.DeepEqual(transformedJson, test.Expect) {
			t.Errorf("%s:\n\t  expected result: %s\n\tunexpected result: %s", test.Name, test.Expect, transformedJson)
		}
	}
}
