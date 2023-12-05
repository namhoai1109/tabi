package dbcore

import "github.com/imdatngo/gowhere"

// ParseCondWithConfig returns standard [sqlString, vars] format for query, powered by gowhere package (configurable version)
func ParseCondWithConfig(cfg gowhere.Config, cond ...interface{}) []interface{} {
	if len(cond) == 1 {
		switch c := cond[0].(type) {
		case map[string]interface{}, []interface{}:
			cond[0] = gowhere.WithConfig(cfg).Where(c)
		}

		if plan, ok := cond[0].(*gowhere.Plan); ok {
			return append([]interface{}{plan.SQL()}, plan.Vars()...)
		}
	}
	return cond
}

// ParseCond returns standard [sqlString, vars] format for query, powered by gowhere package (with default config)
func ParseCond(cond ...interface{}) []interface{} {
	return ParseCondWithConfig(gowhere.DefaultConfig, cond...)
}
