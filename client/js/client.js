$(document).ready(function () {

var socket = new WebSocket('ws://localhost:8000/ws/chat?ssdf=qpwoekop');

socket.onopen = function(event) {
    console.log('WebSocket is open now.');
    socket.send(JSON.stringify({'message': 'Hello WebSocket!'}));
};

socket.onmessage = function(event) {
    console.log('Message from server:', event.data);
};

socket.onclose = function(event) {
    console.log('WebSocket is closed now.');
};

socket.onerror = function(error) {
    console.log('WebSocket error:', error);
};



});
