package templates

// Main is the entrypoint for the api.
//
// Parameters:
// 	- API_PORT : The port to run the API on
const Main = `
import express from "express";
import router from "./routes";

async function main() {
    const app = express();

	app.use("/%API_PREFIX%", router);

    app.listen(%API_PORT%, () => {
		console.log("API started on http://localhost:%API_PORT%");
    });
}

main().catch(err => console.log(err));
`
