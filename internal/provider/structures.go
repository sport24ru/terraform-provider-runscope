package provider

import (
	"time"

	"github.com/adRise/tubi-terraform-provider-runscope/internal/runscope"
)

func expandStringSlice(s []interface{}) []string {
	result := make([]string, len(s), len(s))
	for k, v := range s {
		result[k] = v.(string)
	}
	return result
}

func flattenStepVariables(variables []runscope.StepVariable) []interface{} {
	result := make([]interface{}, len(variables))
	for i, v := range variables {
		result[i] = map[string]interface{}{
			"name":     v.Name,
			"property": v.Property,
			"source":   v.Source,
		}
	}
	return result
}

func flattenStepAssertions(assertions []runscope.StepAssertion) []interface{} {
	result := make([]interface{}, len(assertions))
	for i, a := range assertions {
		result[i] = map[string]interface{}{
			"source":     a.Source,
			"property":   a.Property,
			"comparison": a.Comparison,
			"value":      a.Value,
		}
	}
	return result
}

func flattenStepHeaders(headers map[string][]string) []interface{} {
	result := []interface{}{}
	for header, values := range headers {
		for _, value := range values {
			result = append(result, map[string]interface{}{
				"header": header,
				"value":  value,
			})
		}
	}
	return result
}

func flattenFormParameters(form map[string][]string) []interface{} {
	result := []interface{}{}
	for name, values := range form {
		for _, value := range values {
			result = append(result, map[string]interface{}{
				"name":  name,
				"value": value,
			})
		}
	}
	return result
}

func flattenStepAuth(auth runscope.StepAuth) []map[string]interface{} {
	return []map[string]interface{}{{
		"username":  auth.Username,
		"auth_type": auth.AuthType,
		"password":  auth.Password,
	}}
}

func expandStepVariables(variables []interface{}) []runscope.StepVariable {
	result := make([]runscope.StepVariable, len(variables))
	for i, variable := range variables {
		v := variable.(map[string]interface{})
		result[i] = runscope.StepVariable{
			Name:     v["name"].(string),
			Property: v["property"].(string),
			Source:   v["source"].(string),
		}
	}
	return result
}

func expandStepAssertions(assertions []interface{}) []runscope.StepAssertion {
	result := make([]runscope.StepAssertion, len(assertions))
	for i, assertion := range assertions {
		a := assertion.(map[string]interface{})
		result[i] = runscope.StepAssertion{
			Source:     a["source"].(string),
			Property:   a["property"].(string),
			Comparison: a["comparison"].(string),
			Value:      a["value"].(string),
		}
	}
	return result
}

func expandStepHeaders(headers []interface{}) map[string][]string {
	result := map[string][]string{}
	for _, h := range headers {
		header := h.(map[string]interface{})
		name := header["header"].(string)
		value := header["value"].(string)
		if _, ok := result[name]; !ok {
			result[name] = []string{}
		}
		result[name] = append(result[name], value)
	}
	return result
}

func expandStepForm(formParameters []interface{}) map[string][]string {
	result := map[string][]string{}
	for _, fp := range formParameters {
		formParameter := fp.(map[string]interface{})
		name := formParameter["name"].(string)
		value := formParameter["value"].(string)
		if _, ok := result[name]; !ok {
			result[name] = []string{}
		}
		result[name] = append(result[name], value)
	}
	return result
}

func expandStepAuth(auth []interface{}) runscope.StepAuth {
	result := runscope.StepAuth{}
	if len(auth) > 0 {
		a := auth[0].(map[string]interface{})
		result.Username = a["username"].(string)
		result.Password = a["password"].(string)
		result.AuthType = a["auth_type"].(string)
	}
	return result
}

func flattenTime(t time.Time) string {
	if t.Unix() == 0 {
		return ""
	}

	return t.Format(time.RFC1123)
}

func flattenCreatedBy(c *runscope.CreatedBy) []map[string]interface{} {
	return []map[string]interface{}{{
		"id":    c.Id,
		"name":  c.Name,
		"email": c.Email,
	}}
}

func flattenRemoteAgents(ra []*runscope.RemoteAgent) []map[string]interface{} {
	remoteAgents := make([]map[string]interface{}, len(ra))
	for i, r := range ra {
		remoteAgents[i] = map[string]interface{}{
			"id":      r.Id,
			"name":    r.Name,
			"version": r.Version,
		}
	}
	return remoteAgents
}
