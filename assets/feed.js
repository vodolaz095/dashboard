'use strict';

console.log('Feed: making connection...');
const feed = new EventSource('/feed');
feed.onopen = function () {
  console.log('Feed: connection established');
};
feed.onerror = function () {
  console.error(`Feed: error. State: ${this.readyState}`);
  if (this.readyState === EventSource.CONNECTING) {
    console.log('Feed: reconnecting...');
  } else {
    console.error('Feed: fatal error...');
  }
};
feed.onmessage = function (event) {
  console.log('Feed: event received:', event.data);
};

function formatTimestamp(input) {
  const now = new Date(input);
  return [
    now.getHours(),
    (now.getMinutes() > 9) ? now.getMinutes() : '0' + now.getMinutes(),
    (now.getSeconds() > 9) ? now.getSeconds() : '0' + now.getSeconds(),
  ].join(':');
}

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
      itemTimestamp.innerText = formatTimestamp(params.timestamp)
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
      clock.innerText = formatTimestamp(params.timestamp)
    } catch (err) {
      console.error(err);
    }
  });
}
