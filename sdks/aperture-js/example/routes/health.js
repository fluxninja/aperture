import express from 'express';

export const healthRouter = express.Router();

healthRouter.get('/', function (req, res) {
    res.status(200);
    res.send('Healthy');
});
