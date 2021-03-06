package expert

import (
	"github.com/graarh/golang/expert/condition"
	"github.com/graarh/golang/expert/data"
	"github.com/graarh/golang/expert/modifier"
	"testing"
)

var conditions = [...]Condition{
	&condition.Map{
		Name: "first",
		Values: map[data.Parameter]struct{}{
			"one":   struct{}{},
			"two":   struct{}{},
			"three": struct{}{},
		},
	},
	&condition.Map{
		Name: "second",
		Values: map[data.Parameter]struct{}{
			"four": struct{}{},
			"five": struct{}{},
			"six":  struct{}{},
		},
	},
	&condition.Map{
		Name: "third",
		Values: map[data.Parameter]struct{}{
			"seven": struct{}{},
			"eight": struct{}{},
			"nine":  struct{}{},
		},
	},
}

var mods = []Modifier{
	&modifier.Static{&data.SingleWeight{9}},
	&modifier.Static{&data.SingleWeight{10}},
}
var weight = data.SingleWeight{5}
var paramsMap = map[string]data.Parameter{
	"zero":   "hour",
	"first":  "two",
	"second": "four",
	"third":  "nine",
	"last":   "unused",
	"Type":   uint(2),
	"Time":   int64(3),
}
var params = data.CreateParameters(paramsMap)

// conditions passed, modifier applied -> mod weight
func TestApplied(t *testing.T) {
	rule := RuleRecord{"rule", conditions[0:3], mods}
	calculated, affected := rule.Calculate(&weight, params)

	if !affected {
		t.Error("Rule affects weight cause conditions should be passed, but false given")
	}
	result := calculated.(*data.SingleWeight)

	if result.Value.(int) != 10 {
		t.Error("Modifier of rule was set to return 10, but given ", result.Value)
	}
}

// middle condition failed, modifier not applied -> initial weight
func TestNotApplied(t *testing.T) {
	failedParamsMap := make(map[string]data.Parameter)
	for key, value := range paramsMap {
		failedParamsMap[key] = value
	}
	failedParamsMap["second"] = "failed"

	rule := RuleRecord{"rule", conditions[0:3], mods}
	_, affected := rule.Calculate(&weight, data.CreateParameters(failedParamsMap))

	if affected {
		t.Error("Second condition of rule should fail, so the affected should be false")
	}
}

// no conditions, modifier applied -> mod weight
func TestNoConditions(t *testing.T) {
	rule := RuleRecord{"rule", conditions[0:0], mods}
	calculated, affected := rule.Calculate(&weight, params)
	result := calculated.(*data.SingleWeight)

	if !affected {
		t.Error("No conditions == always passed, affected should be true")
	}

	if result.Value.(int) != 10 {
		t.Error("No condtions == always passed, and modifier of promotion was set to return 10, "+
			"it should return 10, but given ", result.Value)
	}
}

func BenchmarkCommonRule(b *testing.B) {
	rule := RuleRecord{"rule", conditions[0:3], mods}
	for i := 0; i < b.N; i++ {
		rule.Calculate(&weight, params)
	}
}
