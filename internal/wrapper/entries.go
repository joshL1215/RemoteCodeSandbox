package wrapper

const PythonEntry = `
import json
from submission import user_solution

with open("cases.json") as f:
	cases = json.load(f)
for case in cases:
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
	
	print(user_solution(actual_input))
	print("Expected:", raw_out)
`

const NodeEntry = ``
