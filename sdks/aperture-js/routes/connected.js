import express from 'express';
import { apertureClient } from './use_aperture.js';

export const connectedRouter = express.Router();

connectedRouter.get('/', function (req, res) {
    //a.grpcClient.Connect()
	//state := a.grpcClient.GetState()
	//if state != connectivity.Ready {
	//	w.WriteHeader(http.StatusServiceUnavailable)
	//}
	//_, _ = w.Write([]byte(state.String()))

    apertureClient.GetState().then((stateResult) => {
        console.log(`GRPC state result ${stateResult}`)
        res.status(200).send('Connected');
    }).catch(e => {
        console.log(e);
        res.send(`Error ocurred: ${e}\n`);
    });
});