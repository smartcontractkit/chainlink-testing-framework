package seth

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/awalterschulze/gographviz"
)

func findShortestPath(calls []*DecodedCall) []string {
	callMap := make(map[string]*DecodedCall)
	for _, call := range calls {
		callMap[call.CommonData.Signature] = call
	}

	var root *DecodedCall
	for _, call := range calls {
		if call.CommonData.ParentSignature == "" {
			root = call
			break
		}
	}

	if root == nil {
		return nil // No root found
	}

	var end *DecodedCall
	for i := len(calls) - 1; i >= 0; i-- {
		if calls[i].CommonData.Error != "" {
			end = calls[i]
			break
		}
	}
	if end == nil {
		end = calls[len(calls)-1]
	}

	type node struct {
		call *DecodedCall
		path []string
	}

	queue := []node{{call: root, path: []string{root.CommonData.Signature}}}
	visited := make(map[string]bool)

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		currentCall := currentNode.call
		currentPath := currentNode.path

		if currentCall.CommonData.Signature == end.CommonData.Signature {
			return currentPath
		}

		visited[currentCall.CommonData.Signature] = true

		for _, call := range calls {
			if call.CommonData.ParentSignature == currentCall.CommonData.Signature && !visited[call.CommonData.Signature] {
				newPath := append([]string{}, currentPath...)
				newPath = append(newPath, call.CommonData.Signature)
				queue = append(queue, node{call: call, path: newPath})
			}
		}
	}

	return nil // No path found
}

var defaultTruncateTo = 20

func (t *Tracer) generateDotGraph(txHash string, calls []*DecodedCall, revertErr error) error {
	if !t.Cfg.hasOutput(TraceOutput_DOT) {
		return nil
	}

	shortestPath := findShortestPath(calls)

	callHashToID := make(map[string]int)
	nextID := 1

	g := gographviz.NewGraph()
	if err := g.SetName("G"); err != nil {
		return fmt.Errorf("failed to set graph name: %w", err)
	}
	if err := g.SetDir(true); err != nil {
		return fmt.Errorf("failed to set graph direction: %w", err)
	}

	nodesAtLevel := make(map[int][]string)
	revertedCallIdx := -1

	if len(calls) > 0 {
		for i, dc := range calls {
			if dc.Error != "" {
				revertedCallIdx = i
			}
		}
	}

	if err := g.AddNode("G", "start", map[string]string{"label": "\"Start\n\"", "shape": "circle", "style": "filled", "fillcolor": "darkseagreen3", "color": "darkslategray", "fontcolor": "darkslategray"}); err != nil {
		return fmt.Errorf("failed to add start node: %w", err)
	}

	for idx, call := range calls {
		hash := hashCall(call)

		var callID int
		_, exists := callHashToID[hash]
		if !exists {
			callID = nextID
			nextID++
			callHashToID[hash] = callID

			basicNodeID := "node" + strconv.Itoa(callID) + "_basic"
			extraNodeID := "node" + strconv.Itoa(callID) + "_extra"

			var from, to string
			if call.From != "" && call.From != UNKNOWN {
				from = call.From
			} else {
				from = call.FromAddress
			}

			if call.To != "" && call.To != UNKNOWN {
				to = call.To
			} else {
				to = call.ToAddress
			}

			basicLabel := fmt.Sprintf("\"%s -> %s\n %s\"", from, to, call.CommonData.Method)
			extraLabel := fmt.Sprintf("\"Inputs: %s\nOutputs: %s\"", formatMapForLabel(call.CommonData.Input, defaultTruncateTo), formatMapForLabel(call.CommonData.Output, defaultTruncateTo))

			isMajorNode := false
			for _, path := range shortestPath {
				if path == call.Signature {
					isMajorNode = true
					break
				}
			}

			style := "filled"
			nodeColor := "darkslategray"
			fontSize := "14.0"
			if !isMajorNode {
				style = "dashed"
				fontSize = "9.0"
				nodeColor = "lightslategray"
			}

			var subgraphAttrs map[string]string
			subgraphAttrs = map[string]string{"color": "darkslategray"}
			if call.Error != "" {
				subgraphAttrs = map[string]string{"color": "lightcoral"}
				nodeColor = "lightcoral"
			}

			if err := g.AddNode("G", basicNodeID, map[string]string{"label": basicLabel, "shape": "box", "style": style, "fillcolor": "ivory", "color": nodeColor, "fontcolor": "darkslategray", "fontsize": fontSize, "tooltip": formatTooltip(call)}); err != nil {
				return fmt.Errorf("failed to add node: %w", err)
			}
			if err := g.AddNode("G", extraNodeID, map[string]string{"label": extraLabel, "shape": "box", "style": style, "fillcolor": "gainsboro", "color": nodeColor, "fontcolor": "darkslategray", "fontsize": fontSize}); err != nil {
				return fmt.Errorf("failed to add node: %w", err)
			}

			subGraphName := "cluster_" + strconv.Itoa(callID)
			if err := g.AddSubGraph("G", subGraphName, subgraphAttrs); err != nil {
				return fmt.Errorf("failed to add node: %w", err)
			}

			if err := g.AddNode(subGraphName, basicNodeID, nil); err != nil {
				return fmt.Errorf("failed to add node: %w", err)
			}
			if err := g.AddNode(subGraphName, extraNodeID, map[string]string{"rank": "same"}); err != nil {
				return fmt.Errorf("failed to add node: %w", err)
			}

			if idx == 0 {
				if err := g.AddEdge("start", basicNodeID, true, nil); err != nil {
					return fmt.Errorf("failed to add edge: %w", err)
				}
			}

			if call.CommonData.ParentSignature != "" {
				for _, parentCall := range calls {
					if parentCall.CommonData.Signature == call.CommonData.ParentSignature {
						parentHash := hashCall(parentCall)
						parentID := callHashToID[parentHash]
						parentBasicNodeID := "node" + strconv.Itoa(parentID) + "_basic"
						attrs := map[string]string{"fontsize": fontSize, "label": fmt.Sprintf(" \"(%s)\"", ordinalNumber(idx))}
						if call.Error != "" {
							attrs["color"] = "lightcoral"
							attrs["fontcolor"] = "lightcoral"
						} else {
							attrs["color"] = "darkslategray"
							attrs["fontcolor"] = "darkslategray"
						}

						if err := g.AddEdge(parentBasicNodeID, basicNodeID, true, attrs); err != nil {
							return fmt.Errorf("failed to add edge: %w", err)
						}
						break
					}
				}
			}
		} else {
			// This could be also valid if the same call is present twice in the trace, but in typical scenarios that should not happen
			L.Warn().Msg("The same call was present twice. This should not happen and might indicate a bug in the tracer. Check debug log for details")
			marshalled, err := json.Marshal(call)
			if err == nil {
				L.Debug().Msgf("Call: %v", marshalled)
			}
			continue
		}
	}

	// Create dummy nodes to adjust vertical positions
	for level, nodes := range nodesAtLevel {
		for i, node := range nodes {
			if i > 0 {
				dummyNode := fmt.Sprintf("dummy_%d_%d", level, i)
				if err := g.AddNode("G", dummyNode, map[string]string{"label": "\"\"", "shape": "none", "height": "0.1", "width": "0.1"}); err != nil {
					return fmt.Errorf("failed to add node: %w", err)
				}
				if err := g.AddEdge(nodes[i-1], dummyNode, true, map[string]string{"style": "invis"}); err != nil {
					return fmt.Errorf("failed to add node: %w", err)
				}
				if err := g.AddEdge(dummyNode, node, true, map[string]string{"style": "invis"}); err != nil {
					return fmt.Errorf("failed to add node: %w", err)
				}
			}
		}
	}

	if revertErr != nil {
		revertNode := fmt.Sprintf("revert_node_%d", nextID-1)

		if err := g.AddNode("G", revertNode, map[string]string{"label": fmt.Sprintf("\"%s\"", revertErr.Error()), "shape": "rectangle", "style": "filled", "color": "lightcoral", "fillcolor": "lightcoral", "fontcolor": "darkslategray"}); err != nil {
			return fmt.Errorf("failed to add node: %w", err)
		}

		hash := hashCall(calls[revertedCallIdx])
		revertParentNodeId, ok := callHashToID[hash]
		if !ok {
			return fmt.Errorf("failed to find parent node for revert node. This should never happen and likely indicates a bug in code")
		}

		parentBasicNodeID := "node" + strconv.Itoa(revertParentNodeId) + "_basic"

		if err := g.AddEdge(revertNode, parentBasicNodeID, true, map[string]string{"style": "filled", "fillcolor": "lightcoral", "color": "lightcoral", "fontcolor": "darkslategray"}); err != nil {
			return fmt.Errorf("failed to add node: %w", err)
		}
	} else {
		if err := g.AddNode("G", "end", map[string]string{"label": "\"End\n\"", "shape": "circle", "style": "filled", "fillcolor": "darkseagreen3", "color": "darkslategray", "fontcolor": "darkslategray"}); err != nil {
			return fmt.Errorf("failed to add end node: %w", err)
		}

		hash := hashCall(calls[len(calls)-1])
		lastNodeId, ok := callHashToID[hash]
		if !ok {
			return fmt.Errorf("failed to find parent node for revert node. This should never happen and likely indicates a bug in code")
		}

		parentBasicNodeID := "node" + strconv.Itoa(lastNodeId) + "_basic"
		if err := g.AddEdge(parentBasicNodeID, "end", true, nil); err != nil {
			return fmt.Errorf("failed to add edge: %w", err)
		}
	}

	dirPath := filepath.Join(t.Cfg.ArtifactsDir, "dot_graphs")
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(dirPath, fmt.Sprintf("%s.dot", txHash))

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Error creating file: %v\n", err)
	}
	defer func() { _ = f.Close() }()

	if _, err := f.WriteString(g.String()); err != nil {
		return fmt.Errorf("Error writing to file: %v\n", err)
	}

	L.Debug().Msgf("DOT graph saved to %s", filePath)
	L.Debug().Msgf("To view run: xdot %s", filePath)

	return nil
}

func formatTooltip(call *DecodedCall) string {
	basicTooltip := fmt.Sprintf("\"BASIC\nFrom: %s\nTo: %s\nType: %s\nGas Used/Limit: %s\nValue: %d\n\nINPUTS%s\n\nOUTPUTS%s\n\nEVENTS%s\n\"",
		call.FromAddress, call.ToAddress, call.CommonData.CallType, fmt.Sprintf("%d/%d", call.GasUsed, call.GasLimit), call.Value, formatMapForTooltip(call.CommonData.Input), formatMapForTooltip(call.CommonData.Output), formatEvent(call.Events))

	if call.Comment == "" {
		return basicTooltip
	}

	return fmt.Sprintf("%s\nCOMMENT\n%s", basicTooltip, call.Comment)
}

func formatEvent(events []DecodedCommonLog) string {
	if len(events) == 0 {
		return "\n{}"
	}
	parts := make([]string, 0, len(events))
	for _, event := range events {
		parts = append(parts, fmt.Sprintf("\n%s    %s", event.Signature, formatMapForTooltip(event.EventData)))
	}
	return strings.Join(parts, "\n")
}

func prepareMapParts(m map[string]interface{}, truncateTo int) []string {
	if len(m) == 0 {
		return []string{}
	}
	parts := make([]string, 0, len(m))
	for k, v := range m {
		value := fmt.Sprint(v)
		if truncateTo != -1 && len(value) > truncateTo {
			value = value[:truncateTo] + "..."
		}
		parts = append(parts, fmt.Sprintf("%s: %v", k, value))
	}

	return parts
}

func formatMapForTooltip(m map[string]interface{}) string {
	if len(m) == 0 {
		return "\n{}"
	}
	parts := prepareMapParts(m, -1)
	return "\n" + strings.Join(parts, "\n")
}

func formatMapForLabel(m map[string]interface{}, truncateTo int) string {
	if len(m) == 0 {
		return "{}"
	}
	parts := prepareMapParts(m, truncateTo)
	return "\n" + strings.Join(parts, "\\l") + "\\l"
}

func hashCall(call *DecodedCall) string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%v", call)))
	return hex.EncodeToString(h.Sum(nil))
}

func ordinalNumber(n int) string {
	if n <= 0 {
		return strconv.Itoa(n)
	}

	var suffix string
	switch n % 10 {
	case 1:
		if n%100 == 11 {
			suffix = "th"
		} else {
			suffix = "st"
		}
	case 2:
		if n%100 == 12 {
			suffix = "th"
		} else {
			suffix = "nd"
		}
	case 3:
		if n%100 == 13 {
			suffix = "th"
		} else {
			suffix = "rd"
		}
	default:
		suffix = "th"
	}

	return strconv.Itoa(n) + suffix
}
