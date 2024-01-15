package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		command := r.FormValue("command")
		result, err := executeCommand(command)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error executing command: %s", err), http.StatusInternalServerError)
			return
		}

		renderResult(w, result)
		return
	}

	renderForm(w)
}

func renderForm(w http.ResponseWriter) {
	form := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Command Executor</title>
		</head>
		<body>
			<h2>Введите команду:</h2>
			<form method="post">
				<input type="text" name="command" required>
				<input type="submit" value="Выполнить">
			</form>
		</body>
		</html>
	`

	w.Write([]byte(form))
}

func renderResult(w http.ResponseWriter, result string) {
	output := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Результат выполнения команды</title>
		</head>
		<body>
			<h2>Результат выполнения команды:</h2>
			<pre>%s</pre>
		</body>
		</html>
	`

	fmt.Fprintf(w, output, result)
}

func executeCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
