package engine

import (
	"encoding/json"
)

// transform node.value, node.**.nodes[].value, prevNode.value, prevNode.**.nodes[].value.
func (engine *Engine) TransformEtcdJsonResponse(jsonData []byte) ([]byte, error) {
	var data interface{}
	json.Unmarshal(jsonData, &data)

	root, ok := data.(map[string]interface{})
	if !ok {
		return jsonData, nil
	}

	if nodeRaw, ok := root["node"]; ok {
		if node, ok := nodeRaw.(map[string]interface{}); ok {
			engine.transformEtcdJsonResponse0(&node, 0)
		}
	}

	if nodeRaw, ok := root["prevNode"]; ok {
		if node, ok := nodeRaw.(map[string]interface{}); ok {
			engine.transformEtcdJsonResponse0(&node, 0)
		}
	}

	return json.Marshal(data)
}

func (engine *Engine) transformEtcdJsonResponse0(nodePtr *map[string]interface{}, depth int) {
	if depth > 100 {
		return
	}

	node := *nodePtr

	if value, ok := node["value"]; ok {
		if str, ok := value.(string); ok {
			newValue, container, err := engine.TransformAndParse(str)
			if err == nil {
				node["value"] = newValue
				if container != nil {
					node["_etcvault"] = map[string]interface{}{
						"version":   container.Version(),
						"container": container,
					}
				}
			} else {
				node["_etcvault_error"] = err.Error()
			}
		}
	}

	if nodesRaw, ok := node["nodes"]; ok {
		if nodes, ok := nodesRaw.([]interface{}); ok {
			for _, subNodeRaw := range nodes {
				subNode, ok := subNodeRaw.(map[string]interface{})
				if !ok {
					continue
				}

				engine.transformEtcdJsonResponse0(&subNode, depth+1)
			}
		}
	}

	return
}
