package templates

const MainTemplate = `
import express from "express";

function main() {
    const app = express();
    
    app.listen(%API_PORT%, () => {
		console.log("API started on http://localhost:%API_PORT%"
	});
}
`
