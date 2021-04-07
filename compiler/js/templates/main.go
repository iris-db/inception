package templates

// MainTemplate is the entrypoint for the api.
//
// Parameters:
// 	- API_PORT : The port to run the API on
const MainTemplate = `
import express from "express";

function main() {
    const app = express();
    
    app.listen(%API_PORT%, () => {
		console.log("API started on http://localhost:%API_PORT%");
    });
}
`
