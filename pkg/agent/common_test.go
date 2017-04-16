package agent

var (
	application  = "notepad"
	instructions = []string{"taskkill", "/f", "/t", "/im", "notepad.exe"}
	config       = `
	{
		"application": "notepad",
		"instructions": [ "taskkill", "/f", "/t", "/im", "notepad.exe" ]	
	}`
)
