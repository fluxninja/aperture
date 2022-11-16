import http from "http";
import express from "express";
import { createHttpTerminator } from 'http-terminator';
import { apertureRoute, apertureClient } from "./routes/use_aperture.js";
import { healthRouter } from "./routes/health.js";
import { connectedRouter } from "./routes/connected.js";

const host = 'localhost';
const port = process.env.FN_APP_PORT ? process.env.FN_APP_PORT : "8080";

// Create server
const router = express()
const server = http.createServer(router);
const httpTerminator = createHttpTerminator({
  server,
});

// Add routes
router.use("/health", healthRouter);
router.use("/connected", connectedRouter);
router.use("/super", apertureRoute);

// Start listening
server.listen(port, host, () => {
  console.log(`Server is running on http://${host}:${port}`);
  process.on('SIGTERM', startGracefulShutdown);
  process.on('SIGINT', startGracefulShutdown);
});

// Handle graceful shutdown
const startGracefulShutdown = () => {
  apertureClient.Shutdown();
  httpTerminator.terminate().then(() => {
    console.log("Finished shutting down server")
  });
}
