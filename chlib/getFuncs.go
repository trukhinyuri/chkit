package chlib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
)

func LoadGenericJsonFromFile(path string) (b []GenericJson, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	err = json.NewDecoder(file).Decode(&b)
	return
}

func GetCmdRequestJson(client *Client, kind, name string) (ret []GenericJson, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("can`t extract field: %s", r)
		}
	}()
	var apiResult TcpApiResult
	switch kind {
	case KindNamespaces:
		apiResult, err = client.Get(KindNamespaces, name, "")
		if err != nil {
			return ret, err
		}
		items := apiResult["results"].([]interface{})
		for _, itemI := range items {
			item := itemI.(map[string]interface{})
			_, hasNs := item["data"].(map[string]interface{})["metadata"].(map[string]interface{})["namespace"]
			if hasNs {
				ret = append(ret, GenericJson(item))
			}
		}
	default:
		apiResult, err := client.Get(kind, name, client.userConfig.Namespace)
		if err != nil {
			return ret, err
		}
		items := apiResult["results"].([]interface{})
		for _, itemI := range items {
			ret = append(ret, itemI.(map[string]interface{}))
		}
	}
	return
}

type PrettyPrintConfig struct {
	Columns []string
	Data    [][]string
}

func CreatePrettyPrintConfig(kind string, jsonContent []GenericJson) (ppc PrettyPrintConfig, err error) {
	switch kind {
	case KindNamespaces:
		var nsResults []nsResult
		nsResults, err = extractNsResults(jsonContent)
		if err != nil {
			break
		}
		ppc = formatNamespacePrettyPrint(nsResults)
	case KindPods:
		var podResults []podResult
		podResults, err = extractPodResults(jsonContent)
		if err != nil {
			break
		}
		ppc = formatPodPrettyPrint(podResults)
	case KindDeployments:
		var deployResults []deployResult
		deployResults, err = extractDeployResult(jsonContent)
		if err != nil {
			break
		}
		ppc, err = formatDeployPrettyPrint(deployResults)
	case KindService:
		var serviceResults []serivceResult
		serviceResults, err = extractServiceResult(jsonContent)
		if err != nil {
			break
		}
		ppc = formatServicePrettyPrint(serviceResults)
	}
	return
}

func PrettyPrint(ppc PrettyPrintConfig, writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(ppc.Columns)
	table.AppendBulk(ppc.Data)
	table.Render()
}
