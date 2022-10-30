import express from 'express';
import { ApertureClient } from "../sdk/client.js";

export const apertureRoute = express.Router();

// Create aperture client
export const apertureClient = new ApertureClient();

apertureRoute.get('/', function (req, res) {
    console.log(`Got a request`);

    // do some business logic to collect labels
    var labelsMap = new Map().set('labelKey','labelValue');

    // StartFlow performs a flowcontrolv1.Check call to Aperture Agent. It returns a Flow and an error if any.
    apertureClient.StartFlow("aperture-js", labelsMap).then((response) => {
        console.log('StartFlow Done');
        res.writeHead(200);
        res.end('Hello, World!\n');
    }).then( (rejection) => {
        console.log(`StartFlow failed: ${rejection}`);
        res.status(200).send("StartFlow failed\n");
    }).catch(e => {
        console.log(e);
        res.send(`Error ocurred: ${e}\n`);
    });
});
