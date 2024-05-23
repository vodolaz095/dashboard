'use strict';

console.log('Feed: making connection...');
const feed = new EventSource('/feed');
feed.onopen = function () {
  console.log('Feed: connection established');
};
feed.onerror = function () {
  console.log(`Feed: error. State: ${this.readyState}`);
  if (this.readyState === EventSource.CONNECTING) {
    console.log('Feed: reconnecting...');
  } else {
    console.log('Feed: fatal error...');
  }
};
feed.onmessage = function (event) {
  console.log('Feed: event received:', event.data);
};

function doSubscribeOn(eventTypeName) {
  const itemValue = document.getElementById('value_' + eventTypeName);
  const itemError = document.getElementById('error_' + eventTypeName);
  const itemTimestamp = document.getElementById('timestamp_' + eventTypeName);

  feed.addEventListener(eventTypeName, function (e) {
    console.log('Feed: %s event received:', eventTypeName, e.type, e.data);
    try {
      const params = JSON.parse(e.data);
      itemValue.innerText = params.value;
      itemError.innerText = params.error;
      const now = new Date(params.timestamp);
      itemTimestamp.innerText = now.getHours()+':'+now.getMinutes()+':'+now.getSeconds();
    } catch (err) {
      console.error(err);
    }
  });
}

function doClock() {
  const clock = document.getElementById('clock');
  feed.addEventListener('clock', function (e) {
    try {
      const params = JSON.parse(e.data);
      const now = new Date(params.timestamp);
      clock.innerText = now.getHours()+':'+now.getMinutes()+':'+now.getSeconds();
    } catch (err) {
      console.error(err);
    }
  });
}
