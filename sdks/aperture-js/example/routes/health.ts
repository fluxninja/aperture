import express from "express";

export const healthRouter = express.Router();

healthRouter.get("/", function (_: express.Request, res: express.Response) {
  res.status(200);
  res.send("Healthy");
});
