$(document).ready(function() {
  var settings = new Settings();
  var display = new Display(document.getElementById('display'));
  var app = new App(settings, display);
});

