"use strict";

/**
* @type WebSocket | null
*/
let ws = null;

/**
* @param {MessageEvent<any>} event
*/
function handleMessage(event) {
    console.log(event.data);
}

/**
* @param {MessageEvent<any>} event
*/
function handleOpen(event) {
    console.log(`connected with ${event.data}`);
}

function connectToGo() {
    console.log("connecting to go");

    ws?.close();
    ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = handleOpen;
    ws.onmessage = handleMessage;
}

function connectToJS() {
    console.log("connecting to js");

    ws?.close();
    ws = new WebSocket("ws://localhost:8081");

    ws.onmessage = handleMessage;
}

function sendHello() {
    ws?.send("hello there");
}
