var WebSocketService = function (url) {
  this.ws = new WebSocket(setting.socketURL);
  this.hasConnection = false;

  this.ws.onopen = function(evt) {
    this.hasConnection = true;
  };

  this.ws.onmessage = function(evt) {
    render(JSON.parse(evt.data));
    sendRcv();
  };

  this.ws.onerror = function(evt) {
    console.log("Error: ", evt);
    this.ws.close();
  };

};