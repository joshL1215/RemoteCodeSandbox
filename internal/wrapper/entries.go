package wrapper

const PythonEntry = `
import json
from submission import user_solution

with open("cases.json") as f:
	cases = json.load(f)
for case in cases:
	raw = case["input"]

	try:
		converted = json.loads(raw)
	except json.JSONDecodeError:
		converted = raw
	
	print(user_solution(case))
`

const NodeEntry = ``
