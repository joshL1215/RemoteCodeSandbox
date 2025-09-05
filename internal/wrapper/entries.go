package wrapper

const PythonEntry = `
import json
from submission import user_solution

with open("cases.json") as f:
	cases = json.load(f)

results = []
for i, case in enumerate(cases):
	raw_in = case["input"]
	raw_out = case["expectedOutput"]

	try:
		actual_input = json.loads(raw_in)
	except json.JSONDecodeError:
		actual_input = raw_in
	
	try:
		actual_expected_output = json.loads(raw_out)
	except json.JSONDecodeError:
		actual_expected_output = raw_out
	
	user_out = user_solution(actual_input)
	results.append({"returned" : user_out, "res" : "pass" if actual_expected_output == user_out else "fail"})

print(json.dumps(results))
`

const NodeEntry = ``
