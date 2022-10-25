import { client } from "./client.js";
import http from "http";
import { host, port } from "./consts.js";


const requestListener = function (req, res) {
  console.log(`Got a request`);
  var labelsMap = new Map().set('labelKey','labelValue');
  client.StartFlow("aperture-js", labelsMap).then(() => {
    console.log('StartFlow Done');
    res.writeHead(200);
    res.end('Hello, World!\n');
  });
};

const server = http.createServer(requestListener);
server.listen(port, host, () => {
    console.log(`Server is running on http://${host}:${port}`);
});
