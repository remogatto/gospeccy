var App = function (settings, display) {
  var app = this;

  this.ws = new WebSocket(settings.socketURL);
  this.hasConnection = false;

  this.ws.onopen = function(evt) {
    this.hasConnection = true;
  };
  this.ws.onmessage = function(evt) {
    app.ws.send("RECEIVED");
    display.render(JSON.parse(evt.data));
  };
  this.ws.onerror = function(evt) {
    console.log("Error: ", evt);
    ws.close();
  };
};
