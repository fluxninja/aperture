import grpc from "@grpc/grpc-js";
import express from 'express';
import { apertureClient } from './use_aperture.js';

export const connectedRouter = express.Router();

connectedRouter.get('/', function (req, res) {
    try {
        let clientState = apertureClient.GetState();
        console.log(`Client state ${clientState}`);
        if (clientState != grpc.connectivityState.READY) {
            res.status(503).send('Unavailable\n');
        } else {
            res.status(200).send('Connected\n');
        }
    } catch(e) {
        console.log(e);
        res.send(`Error ocurred: ${e}\n`);
    }
});