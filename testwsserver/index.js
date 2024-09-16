import * as WebSocket from "ws";

const server = new WebSocket.WebSocketServer({ port: 8081 });

server.on("connection", (ws) => {
    console.log("opened connection");

    ws.send("hey, this is a message from the server");

    ws.on("message", (msg) => {
        console.log(msg);
        ws.send(`you sent: "${msg}"`);
    });
    ws.on("close", () => {
        console.log("closed connection");
    });
});
