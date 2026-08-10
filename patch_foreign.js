const fs = require('fs');
let code = fs.readFileSync('src/Foreign.go', 'utf8');

code = code.replace(`
	if mapRaw, ok := v.(map[string]interface{}); ok {
`, `
	if mapGopurs, ok := v.(map[string]gopurs_runtime.Value); ok {
		res := make(map[string]interface{})
		for k, x := range mapGopurs {
			res[k] = deepUnbox(x)
		}
		return res
	}
	if mapRaw, ok := v.(map[string]interface{}); ok {
`);
fs.writeFileSync('src/Foreign.go', code);
